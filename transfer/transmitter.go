package transfer

import (
	"context"
	"errors"

	"github.com/juntaki/transparent"
	"github.com/juntaki/transparent/simple"
	pb "github.com/juntaki/transparent/transfer/pb"
	"google.golang.org/grpc"
)

func (t *Transmitter) convertSendMessage(m *transparent.Message) (*pb.Message, error) {
	var converted pb.Message
	if m.Key != nil {
		keyStr, err := t.ValidateKey(m.Key)
		if err != nil {
			return nil, err
		}
		converted.Key = keyStr
	}
	if m.Value != nil {
		valueBytes, err := t.ValidateValue(m.Value)
		if err != nil {
			return nil, err
		}
		converted.Value = valueBytes
	}
	switch m.Message {
	case transparent.MessageSet:
		converted.MessageType = pb.MessageType_Set
	case transparent.MessageGet:
		converted.MessageType = pb.MessageType_Get
	case transparent.MessageRemove:
		converted.MessageType = pb.MessageType_Remove
	case transparent.MessageSync:
		converted.MessageType = pb.MessageType_Sync
	default:
		return nil, errors.New("Unknown type")
	}

	return &converted, nil
}

func (t *Transmitter) convertReceiveMessage(m *pb.Message) (*transparent.Message, error) {
	if m == nil {
		return nil, errors.New("nil")
	}
	var converted transparent.Message
	converted.Key = m.Key
	converted.Value = m.Value
	switch m.MessageType {
	case pb.MessageType_Set:
		converted.Message = transparent.MessageSet
	case pb.MessageType_Get:
		converted.Message = transparent.MessageGet
	case pb.MessageType_Remove:
		converted.Message = transparent.MessageRemove
	case pb.MessageType_Sync:
		converted.Message = transparent.MessageSync
	default:
		return nil, errors.New("Unknown type")
	}

	return &converted, nil
}

// Transmitter is simple Transmitter
type Transmitter struct {
	simple.Validator
	client     pb.TransferClient
	serverAddr string
	conn       *grpc.ClientConn
}

// NewSimpleTransmitter returns Sender
func NewSimpleTransmitter(serverAddr string) *Transmitter {
	return &Transmitter{
		Validator:  simple.Validator{},
		serverAddr: serverAddr,
	}
}

func (t *Transmitter) Request(m *transparent.Message) (*transparent.Message, error) {
	message, err := t.convertSendMessage(m)
	if err != nil {
		return nil, err
	}
	r, err := t.client.Request(context.Background(), message)
	response, err := t.convertReceiveMessage(r)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (t *Transmitter) Start() error {
	conn, err := grpc.Dial(t.serverAddr, grpc.WithInsecure())
	t.conn = conn
	if err != nil {
		return err
	}
	t.client = pb.NewTransferClient(t.conn)
	return nil
}

func (t *Transmitter) Stop() error {
	t.conn.Close()
	return nil
}

func (t *Transmitter) SetCallback(m func(*transparent.Message) error) error {
	return nil
}
