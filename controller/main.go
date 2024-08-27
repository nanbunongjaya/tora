package controller

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/nanbunongjaya/tora/proto/servicespb"
)

type (
	Controller struct {
		name    string
		address string
		conn    *grpc.ClientConn
		client  pb.ServicesClient
	}
)

func New(name, addr string) (*Controller, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Fail to connect grpc slave server: %v, name: %v, address: %v", err, name, addr)
		return nil, err
	}

	return &Controller{
		name:    name,
		address: addr,
		conn:    conn,
		client:  pb.NewServicesClient(conn),
	}, nil
}

func (c *Controller) HandleRequest(ctx context.Context, cmd string, data []byte) ([]byte, error) {
	req := &pb.Request{
		CMD:  cmd,
		Data: data,
	}

	res, err := c.client.Handle(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (c *Controller) Close() error {
	return c.conn.Close()
}
