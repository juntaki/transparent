//Package twopc is two phase commit implements for key-value store
package twopc

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/juntaki/transparent"
	pb "github.com/juntaki/transparent/twopc/pb"
)

// DebugLevel determine the amount of debug output
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
	stateReady       // Participant only
	stateAbort
	stateCommit
)

type member string

// Coodinator distribute vote request
type Coodinator struct {
	lock    sync.RWMutex
	in      chan *pb.Message
	out     map[uint64]chan *pb.Message
	request chan *pb.SetRequest
	summary map[uint64]*pb.Message
	ack     map[uint64]*pb.Message
	timeout time.Duration
	status  state
	current uint64
}

// NewCoodinator returns started Coodinator
func NewCoodinator(serverAddr string) (*Coodinator, error) {
	c := &Coodinator{
		timeout: 1000,
		in:      make(chan *pb.Message, 1),
		lock:    sync.RWMutex{},
		out:     make(map[uint64]chan *pb.Message),
		request: make(chan *pb.SetRequest, 10),
		status:  stateInit,
	}
	started := make(chan error)
	go c.start(serverAddr, started)
	err := <-started
	return c, err
}

// StartServ Starts cluster coodinator
func (c *Coodinator) start(address string, started chan error) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		started <- err
		return
	}
	grpcServer := grpc.NewServer()
	pb.RegisterClusterServer(grpcServer, c)

	go c.run()
	started <- nil
	grpcServer.Serve(lis)
}

// SetTimeout change timeout default is 1000 milliseconds
func (c *Coodinator) SetTimeout(millisecond time.Duration) {
	c.timeout = millisecond
}

// Set accepts request from any client
func (c *Coodinator) Set(ctx context.Context, req *pb.SetRequest) (*pb.EmptyMessage, error) {
	c.request <- req
	return &pb.EmptyMessage{}, nil
}

// Connection start and keep connection for each client
func (c *Coodinator) Connection(stream pb.Cluster_ConnectionServer) error {
	// Assign clientID and tell current request ID
	clientID := uint64(len(c.out))
	c.lock.Lock()
	m := &pb.Message{
		ClientID:    clientID,
		MessageType: pb.MessageType_ACK,
		RequestID:   c.current,
	}
	c.out[clientID] = make(chan *pb.Message, 1)
	c.lock.Unlock()
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
		c.lock.RLock()
		sendChannel := c.out[clientID]
		c.lock.RUnlock()
		for {
			m := <-sendChannel
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

func (c *Coodinator) run() {
	for {
		c.initialize()
		debugPrintln(1, "ServerStatus:", c.status)
		select {
		case r := <-c.request:
			commit := c.voteRequest(r)
			debugPrintln(1, "ServerStatus:", c.status)
			if commit {
				c.globalcommit()
			} else {
				c.globalAbort()
			}
			debugPrintln(1, "ServerStatus:", c.status)
			c.waitsendACK()
		}
	}
}

func (c *Coodinator) initialize() {
	c.lock.Lock()
	c.status = stateInit
	c.summary = make(map[uint64]*pb.Message)
	c.ack = make(map[uint64]*pb.Message)
	c.current++
	c.lock.Unlock()
}

func (c *Coodinator) broadcast(m *pb.Message) {
	for _, channel := range c.out {
		channel <- m
	}
}

func (c *Coodinator) globalcommit() {
	m := &pb.Message{
		MessageType: pb.MessageType_GlobalCommit,
		RequestID:   c.current,
	}
	c.status = stateCommit
	c.broadcast(m)
}

func (c *Coodinator) globalAbort() {
	m := &pb.Message{
		MessageType: pb.MessageType_GlobalAbort,
		RequestID:   c.current,
	}
	c.status = stateAbort
	c.broadcast(m)
}

func (c *Coodinator) waitsendACK() (ok bool) {
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

func (c *Coodinator) voteRequest(r *pb.SetRequest) (commit bool) {
	c.status = stateWait

	message := &pb.Message{
		MessageType: pb.MessageType_VoteRequest,
		RequestID:   c.current,
		Payload:     r.Payload,
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

// NewParticipant returns started Participant
func NewParticipant(serverAddr string) *Participant {
	p := &Participant{
		timeout:    1000, //millisecond
		serverAddr: serverAddr,
	}
	return p
}

// Participant manage its resource
type Participant struct {
	in             chan *pb.Message
	out            chan *pb.Message
	timeout        time.Duration
	status         state
	current        uint64
	clientID       uint64
	currentRequest *pb.Message
	client         pb.ClusterClient
	committer      func(m *transparent.Message) (*transparent.Message, error)
	serverAddr     string
}

// SetTimeout change timeout default is 1000 milliseconds
func (a *Participant) SetTimeout(millisecond time.Duration) {
	a.timeout = millisecond
}

// start start participant service
func (a *Participant) start(serverAddr string, started chan error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		started <- err
		return
	}
	defer conn.Close()
	client := pb.NewClusterClient(conn)
	a.client = client
	stream, err := client.Connection(context.Background())
	if err != nil {
		started <- err
		return
	}

	// Get ID from server
	in, err := stream.Recv()
	debugPrintln(5, "Client:Recv", a.clientID, in)
	if err == io.EOF {
		started <- err
		return
	}
	if err != nil {
		started <- err
		return
	}
	if in.MessageType != pb.MessageType_ACK {
		started <- err
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
				started <- err
				return
			}
			if err != nil {
				started <- err
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

	go a.mainLoop()

	started <- nil
	<-recv
	<-send
	stream.CloseSend()
}

func (a *Participant) Start() error {
	started := make(chan error)
	go a.start(a.serverAddr, started)
	err := <-started
	return err
}
func (a *Participant) Stop() error {
	return nil
}

func (a *Participant) SetCallback(committer func(m *transparent.Message) (*transparent.Message, error)) error {
	a.committer = committer
	return nil
}

// Request send request to Coodinator
func (a *Participant) Request(operation *transparent.Message) (*transparent.Message, error) {
	request, err := a.encode(operation)
	if err != nil {
		debugPrintln(1, "Encode error", err)
		return nil, err
	}
	debugPrintln(1, "Client Set", operation)
	_, err = a.client.Set(context.Background(), request)
	return nil, err
}

func (a *Participant) encode(operation *transparent.Message) (*pb.SetRequest, error) {
	gob.Register(operation)
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(&operation)
	if err != nil {
		return nil, err
	}
	request := &pb.SetRequest{
		Payload: buf.Bytes(),
	}
	return request, nil
}

func (a *Participant) commit() (*transparent.Message, error) {
	operation, err := a.decode(a.currentRequest.Payload)
	if err != nil {
		debugPrintln(1, "Decode error", err)
		return nil, err
	}
	debugPrintln(1, "Client Commit", a.clientID, operation)
	return a.committer(operation)
}

func (a *Participant) decode(encoded []byte) (*transparent.Message, error) {
	var operation transparent.Message
	buf := bytes.NewBuffer(encoded)
	encoder := gob.NewDecoder(buf)
	err := encoder.Decode(&operation)
	if err != nil {
		return nil, err
	}
	return &operation, nil
}

func (a *Participant) mainLoop() {
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
					a.voteAbort(m.RequestID)
					break
				}
				a.currentRequest = m
				a.status = stateReady
				a.votecommit(m.RequestID)
			case pb.MessageType_GlobalCommit:
				if a.status != stateReady ||
					m.RequestID != a.current {
					debugPrintln(5, "Ignore globalCommit", a.clientID, a.status, m.RequestID, a.current)
					// ignore
					break
				}
				a.status = stateCommit
				a.commit()
				a.sendACK(m.RequestID)
			case pb.MessageType_GlobalAbort:
				if a.status != stateReady ||
					m.RequestID != a.current {
					debugPrintln(5, "Ignore globalAbort", a.clientID, a.status, m.RequestID, a.current)
					// ignore
					break
				}
				a.status = stateAbort
				a.sendACK(m.RequestID)
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

func (a *Participant) votecommit(requestID uint64) {
	a.out <- &pb.Message{
		MessageType: pb.MessageType_VoteCommit,
		Payload:     nil,
		ClientID:    a.clientID,
		RequestID:   requestID,
	}
}

func (a *Participant) voteAbort(requestID uint64) {
	a.out <- &pb.Message{
		MessageType: pb.MessageType_VoteAbort,
		Payload:     nil,
		ClientID:    a.clientID,
		RequestID:   requestID,
	}
}

func (a *Participant) sendACK(requestID uint64) {
	a.out <- &pb.Message{
		MessageType: pb.MessageType_ACK,
		Payload:     nil,
		ClientID:    a.clientID,
		RequestID:   requestID,
	}

	a.current++
	a.status = stateInit
}
