package twopc

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	pb "github.com/juntaki/transparent/2pc/pb"
)

var DebugLevel int

func init() {
	DebugLevel = 0
}

func debugPrintln(level int, a ...interface{}) (n int, err error) {
	if DebugLevel >= level {
		return fmt.Println(a)
	}
	return 0, nil
}

type request struct {
	key         interface{}
	value       interface{}
	requestType pb.RequestType
}

type state int

func (s state) String() string {
	switch s {
	case stateInit:
		return "Init"
	case stateWait:
		return "Wait"
	case stateReady:
		return "Ready"
	case stateAbort:
		return "Abort"
	case stateCommit:
		return "Commit"
	}
	panic("Unknown value")
}

const (
	stateInit  state = iota + 1
	stateWait        // Coodinator only
	stateReady       // Attendee only
	stateAbort
	stateCommit
)

type member string

// Coodinator distribute vote request
type Coodinator struct {
	in      chan *pb.Message
	out     map[uint64]chan *pb.Message
	request chan *pb.SetRequest
	summary map[uint64]*pb.Message
	ack     map[uint64]*pb.Message
	timeout time.Duration
	status  state
	current uint64
}

// StartServ Starts cluster coodinator
func (c *Coodinator) StartServ(timeoutMillisecond time.Duration) {
	c.timeout = timeoutMillisecond
	c.in = make(chan *pb.Message, 1)
	c.out = make(map[uint64]chan *pb.Message)
	c.request = make(chan *pb.SetRequest, 10)
	c.status = stateInit
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)

	}
	grpcServer := grpc.NewServer()
	pb.RegisterClusterServer(grpcServer, c)

	go c.Run()
	grpcServer.Serve(lis)
}

// SetRequest accepts request from any client
func (c *Coodinator) Set(ctx context.Context, req *pb.SetRequest) (*pb.EmptyMessage, error) {
	c.request <- req
	return &pb.EmptyMessage{}, nil
}

// Connection for each client
func (c *Coodinator) Connection(stream pb.Cluster_ConnectionServer) error {
	// Assign clientID and tell current request ID
	clientID := uint64(len(c.out))
	m := &pb.Message{
		ClientID:    clientID,
		MessageType: pb.MessageType_ACK,
		RequestID:   c.current,
	}
	c.out[clientID] = make(chan *pb.Message, 1)
	debugPrintln(5, "Server:Send", m)
	if err := stream.Send(m); err != nil {
		debugPrintln(5, err)
		return err
	}

	// Receiver
	recv := make(chan bool)
	go func(stream pb.Cluster_ConnectionServer, finish chan bool) {
		for {
			in, err := stream.Recv()
			debugPrintln(5, "Server:Recv", in)
			if err == io.EOF {
				break
			}
			if err != nil {
				debugPrintln(5, err)
				break
			}
			c.in <- in
		}
		finish <- true
		return
	}(stream, recv)

	// Sender
	send := make(chan bool)
	go func(stream pb.Cluster_ConnectionServer, clientID uint64, finish chan bool) {
		for {
			m := <-c.out[clientID]
			debugPrintln(5, "Server:Send", m)
			if err := stream.Send(m); err != nil {
				debugPrintln(5, err)
				break
			}
		}
		finish <- true
		return
	}(stream, clientID, send)

	<-recv
	<-send
	return nil
}

func (c *Coodinator) Run() {
	for {
		c.Initialize()
		debugPrintln(1, "ServerStatus:", c.status)
		select {
		case r := <-c.request:
			commit := c.VoteRequest(r)
			debugPrintln(1, "ServerStatus:", c.status)
			if commit {
				c.GlobalCommit()
			} else {
				c.GlobalAbort()
			}
			debugPrintln(1, "ServerStatus:", c.status)
			c.WaitACK()
		}
	}
}

func (c *Coodinator) Initialize() {
	c.status = stateInit
	c.summary = make(map[uint64]*pb.Message)
	c.ack = make(map[uint64]*pb.Message)
	c.current++
}

func (c *Coodinator) broadcast(m *pb.Message) {
	for _, channel := range c.out {
		channel <- m
	}
}

func (c *Coodinator) GlobalCommit() {
	m := &pb.Message{
		MessageType: pb.MessageType_GlobalCommit,
		RequestID:   c.current,
	}
	c.status = stateCommit
	c.broadcast(m)
}

func (c *Coodinator) GlobalAbort() {
	m := &pb.Message{
		MessageType: pb.MessageType_GlobalAbort,
		RequestID:   c.current,
	}
	c.status = stateAbort
	c.broadcast(m)
}

func (c *Coodinator) WaitACK() (ok bool) {
	ok = true
	for {
		select {
		case v := <-c.in:
			if v.MessageType != pb.MessageType_ACK &&
				v.RequestID != c.current {
				break
			}
			c.ack[v.ClientID] = v
			debugPrintln(5, "ack total", len(c.out), "current", len(c.ack))
			if len(c.out) == len(c.ack) {
				return
			}
		case <-time.After(time.Millisecond * c.timeout):
			debugPrintln(5, "Timeout")
			ok = false
			return
		}
	}
}

func (c *Coodinator) VoteRequest(r *pb.SetRequest) (commit bool) {
	c.status = stateWait

	message := &pb.Message{
		MessageType: pb.MessageType_VoteRequest,
		RequestID:   c.current,
		RequestType: r.RequestType,
		Key:         r.Key,
		Value:       r.Value,
	}
	c.broadcast(message)
	commit = true
	for {
		select {
		case v := <-c.in:
			if v.MessageType != pb.MessageType_VoteCommit &&
				v.MessageType != pb.MessageType_VoteAbort &&
				v.RequestID != c.current {
				break
			} else if v.MessageType == pb.MessageType_VoteAbort {
				debugPrintln(5, "Server:Get VoteAbort")
				commit = false
			}
			c.summary[v.ClientID] = v
			debugPrintln(5, "vote total", len(c.out), "current", len(c.summary))
			if len(c.out) == len(c.summary) {
				return
			}
		case <-time.After(time.Millisecond * c.timeout):
			debugPrintln(5, "Server:Timeout")
			commit = false
			return
		}
	}
}

type Attendee struct {
	in             chan *pb.Message
	out            chan *pb.Message
	timeout        time.Duration
	status         state
	current        uint64
	clientID       uint64
	currentRequest *pb.Message
	client         pb.ClusterClient
}

func (a *Attendee) StartClient(timeoutMillisecond time.Duration) {
	a.timeout = timeoutMillisecond
	serverAddr := "127.0.0.1:8080"

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		debugPrintln(5, err)
	}
	defer conn.Close()
	client := pb.NewClusterClient(conn)
	a.client = client
	stream, err := client.Connection(context.Background())

	// Get ID from server
	in, err := stream.Recv()
	debugPrintln(5, "Client:Recv", a.clientID, in)
	if err == io.EOF {
		return
	}
	if err != nil {
		return
	}
	if in.MessageType != pb.MessageType_ACK {
		return
	}
	a.current = in.RequestID
	a.clientID = in.ClientID

	a.in = make(chan *pb.Message, 1)
	a.out = make(chan *pb.Message, 1)

	// Receiver
	recv := make(chan bool)
	go func(stream pb.Cluster_ConnectionClient, finish chan bool) {
		for {
			in, err := stream.Recv()
			debugPrintln(5, "Client:Recv", a.clientID, in)
			if err == io.EOF {
				return
			}
			if err != nil {
				return
			}
			a.in <- in
		}
	}(stream, recv)

	// Sender
	send := make(chan bool)
	go func(stream pb.Cluster_ConnectionClient, finish chan bool) {
		for {
			m := <-a.out
			debugPrintln(5, "Client:Send", a.clientID, m)
			if err := stream.Send(m); err != nil {
				break
			}
		}
		finish <- true
		return
	}(stream, recv)

	go a.Run()
	<-recv
	<-send
	stream.CloseSend()
}

func (a *Attendee) Set() {
	a.client.Set(context.Background(), &pb.SetRequest{})
}

func (a *Attendee) Run() {
	a.status = stateInit
	for {
		select {
		case m := <-a.in:
			switch m.MessageType {
			case pb.MessageType_VoteRequest:
				if a.status != stateInit {
					debugPrintln(5, "Ignore VoteRequest", a.clientID, a.status)
					// ignore
					break
				}
				if m.RequestID != a.current {
					// attendee may miss the last request
					a.status = stateAbort
					a.VoteAbort(m)
					break
				}
				a.currentRequest = m
				a.status = stateReady
				a.VoteCommit(m)
			case pb.MessageType_GlobalCommit:
				if a.status != stateReady ||
					m.RequestID != a.current {
					debugPrintln(5, "Ignore GlobalCommit", a.clientID, a.status, m.RequestID, a.current)
					// ignore
					break
				}
				a.status = stateCommit
				a.Commit()
				a.ACK(m)
			case pb.MessageType_GlobalAbort:
				if a.status != stateReady ||
					m.RequestID != a.current {
					debugPrintln(5, "Ignore GlobalAbort", a.clientID, a.status, m.RequestID, a.current)
					// ignore
					break
				}
				a.status = stateAbort
				a.ACK(m)
			}
		case <-time.After(time.Millisecond * a.timeout):
			debugPrintln(5, "Client:Timeout", a.clientID)
			if a.status == stateReady {
				a.status = stateInit
				a.current++
			}
		}
	}
}

func (a *Attendee) VoteCommit(v *pb.Message) {
	// send votePayload to coodinator
	v.MessageType = pb.MessageType_VoteCommit
	v.Key = nil
	v.Value = nil
	v.ClientID = a.clientID
	a.out <- v
}

func (a *Attendee) VoteAbort(v *pb.Message) {
	// send votePayload to coodinator
	v.MessageType = pb.MessageType_VoteAbort
	v.Key = nil
	v.Value = nil
	v.ClientID = a.clientID
	a.out <- v
}

func (a *Attendee) Commit() {
	debugPrintln(1, "ClientCommit", a.clientID, a.currentRequest)
	// Set value
}

func (a *Attendee) ACK(v *pb.Message) {
	// send votePayload to coodinator
	v.MessageType = pb.MessageType_ACK
	v.Key = nil
	v.Value = nil
	v.ClientID = a.clientID
	a.out <- v

	a.current++
	a.status = stateInit
}
