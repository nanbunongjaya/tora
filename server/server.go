package server

import (
	"context"
	"log"
	"plugin"

	"tora/compiler"
	"tora/component"
	"tora/config"

	pb "tora/proto/servicespb"
)

type (
	Server struct {
		pb.UnimplementedServicesServer

		comps    *component.Components
		services *component.Services
		plugin   *plugin.Plugin
		incloud  bool
		ismaster bool
	}
)

func New(opts ...Option) *Server {
	s := &Server{
		comps:    &component.Components{},
		services: &component.Services{},
		plugin:   &plugin.Plugin{},
		incloud:  false, // default as local
		ismaster: true,  // default as master
	}

	for i := range opts {
		opts[i](s)
	}

	// Setup the local services
	if !s.incloud {
		s.services.Setup(s.comps)
		return s
	}

	if err := s.compilePlugin(); err != nil {
		log.Fatal(err)
	}

	switch s.ismaster {
	case true:
		// Initialize as master
		if err := s.initializeK8sResources(); err != nil {
			log.Fatal(err)
		}

	case false:
		// Initialize as slave
		if err := s.loadPlugin(); err != nil {
			log.Fatal(err)
		}
	}

	return s
}

func (s *Server) compilePlugin() error {
	program, err := compiler.Compile(s.comps)
	if err != nil {
		return err
	}

	err = compiler.Output(config.GO_FILE, program)
	if err != nil {
		return err
	}

	err = compiler.Build(config.SO_FILE, config.GO_FILE)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) loadPlugin() error {
	p, err := plugin.Open(config.SO_FILE)
	if err != nil {
		return err
	}

	s.plugin = p

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

func (s *Server) initializeK8sResources() error {
	// TODO:
	// Setup ConfigMap
	// Create RBAC resources
	// Create Pods
	return nil
}

func (s *Server) Handle(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	if err := s.handle(req.CMD, req.Data); err != nil {
		return nil, err
	}

	log.Printf("Received CMD: %s, Data: %s", req.CMD, string(req.Data))

	return &pb.Response{Data: req.Data}, nil
}

func (s *Server) handle(cmd string, data []byte) error {
	if !s.incloud {
		return s.services.Handle(cmd, data)
	}

	f, err := s.plugin.Lookup("Handle")
	if err != nil {
		return err
	}

	err = f.(func(string, []byte) error)(cmd, data)
	if err != nil {
		return err
	}

	return nil
}
