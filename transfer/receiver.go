package transfer

import (
	"net"

	"golang.org/x/net/context"

	"github.com/juntaki/transparent"

	pb "github.com/juntaki/transparent/transfer/pb"
	"google.golang.org/grpc"
)

type receiver struct {
	serverAddr     string
	grpcServer     *grpc.Server
	transferServer *server
}

func NewSimpleLayerReceiver(serverAddr string) transparent.Layer {
	r := NewSimpleReceiver(serverAddr)
	return transparent.NewLayerReceiver(r)
}

// NewSimpleReceiver returns simple Receiver
func NewSimpleReceiver(serverAddr string) transparent.BackendReceiver {
	return &receiver{
		serverAddr:     serverAddr,
		transferServer: &server{converter: converter{}},
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

func (r *receiver) SetCallback(cb func(m *transparent.Message) (*transparent.Message, error)) error {
	r.transferServer.callback = cb
	return nil
}

type server struct {
	converter
	callback func(m *transparent.Message) (*transparent.Message, error)
}

func (t *server) Request(c context.Context, m *pb.Message) (*pb.Message, error) {
	decoded, err := t.convertReceiveMessage(m)
	if err != nil {
		return nil, err
	}
	res, err := t.callback(decoded)
	if err != nil {
		return nil, err
	}
	message, err := t.convertSendMessage(res)
	if err != nil {
		return nil, err
	}
	return message, nil
}
