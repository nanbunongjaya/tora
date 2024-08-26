package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"plugin"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "tora/proto/servicespb"
)

type services struct {
	pb.UnimplementedServicesServer
	plugin *plugin.Plugin
}

func (ss *services) Handle(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	if err := ss.handle(req.CMD, req.Data); err != nil {
		return nil, err
	}

	fmt.Printf("Received CMD: %s, Data: %s\n", req.CMD, string(req.Data))

	return &pb.Response{Data: req.Data}, nil
}

func (ss *services) setup() error {
	p, err := plugin.Open("../example/tora_slave_services_plugin.so")
	if err != nil {
		return err
	}

	ss.plugin = p

	// Call setup function
	f, err := p.Lookup("Setup")
	if err != nil {
		return err
	}
	f.(func())()

	// Call list function
	f, err = p.Lookup("List")
	if err != nil {
		return err
	}
	f.(func())()

	return nil
}

func (ss *services) handle(cmd string, data []byte) error {
	f, err := ss.plugin.Lookup("Handle")
	if err != nil {
		return err
	}

	err = f.(func(string, []byte) error)(cmd, data)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Listen on a port
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	s := grpc.NewServer()

	// Setup services
	services := &services{}
	if err := services.setup(); err != nil {
		log.Fatal(err)
	}

	// Register the server with the gRPC server
	pb.RegisterServicesServer(s, services)

	// Enable server reflection for easier debugging
	reflection.Register(s)

	fmt.Println("gRPC server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
