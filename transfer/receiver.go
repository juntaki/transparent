package transfer

import (
	"errors"
	"net"

	"golang.org/x/net/context"

	"github.com/juntaki/transparent"

	"github.com/juntaki/transparent/simple"
	pb "github.com/juntaki/transparent/transfer/pb"
	"google.golang.org/grpc"
)

type receiver struct {
	serverAddr     string
	grpcServer     *grpc.Server
	transferServer *server
}

// NewReceiver returns simple Receiver
func NewReceiver(serverAddr string) transparent.BackendReceiver {
	return &receiver{
		serverAddr:     serverAddr,
		transferServer: &server{},
	}
}

func (r *receiver) Start() error {
	lis, err := net.Listen("tcp", r.serverAddr)
	if err != nil {
		return err
	}
	r.grpcServer = grpc.NewServer()
	pb.RegisterTransferServer(r.grpcServer, r.transferServer)

	go r.grpcServer.Serve(lis)
	return nil
}

func (r *receiver) Stop() error {
	return nil
}

func (r *receiver) SetNext(l transparent.Layer) error {
	r.transferServer.next = l
	return nil
}

type server struct {
	next transparent.Layer
	simple.Validator
}

func (t *server) Request(c context.Context, m *pb.Message) (*pb.Message, error) {
	var err error
	var message pb.Message
	message.MessageType = m.MessageType

	switch m.MessageType {
	case pb.MessageType_Set:
		err = t.next.Set(m.Key, m.Value)
	case pb.MessageType_Get:
		value, err := t.next.Get(m.Key)
		valueStr, err := t.ValidateValue(value)
		if err != nil {
			return nil, err
		}
		message.Key = m.Key
		message.Value = valueStr
	case pb.MessageType_Remove:
		err = t.next.Remove(m.Key)
	case pb.MessageType_Sync:
		err = t.next.Sync()
	default:
		err = errors.New("unknown message")
	}
	if err != nil {
		return nil, err
	}
	return &message, nil
}
