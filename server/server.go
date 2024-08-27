package server

import (
	"context"
	"log"
	"net"

	"tora/cluster"
	"tora/compiler"
	"tora/component"
	"tora/config"

	pb "tora/proto/servicespb"

	grpclib "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type (
	Server struct {
		pb.UnimplementedServicesServer

		comps     *component.Components
		clientset *cluster.ClientSet
		services  services
		grpc      *grpc
		httpx     *httpx
		incloud   bool
		ismaster  bool
	}

	grpc struct {
		enable   bool
		network  string
		address  string
		listener net.Listener
		server   *grpclib.Server
	}

	// TODO: implement httpx server
	httpx struct {
		enable bool
	}
)

func New(opts ...Option) *Server {
	s := &Server{
		comps:    &component.Components{},
		grpc:     &grpc{enable: false},  // default as disable
		httpx:    &httpx{enable: false}, // default as disable
		incloud:  false,                 // default as local
		ismaster: true,                  // default as master
	}

	for i := range opts {
		opts[i](s)
	}

	// Setup in local
	if !s.incloud {
		s.services = newMasterServices(s.comps)
		return s
	}

	switch s.ismaster {
	case true:
		// Initialize as master
		if err := s.compilePlugin(); err != nil {
			log.Fatal(err)
		}

		if err := s.initializeK8sResources(); err != nil {
			log.Fatal(err)
		}

	case false:
		// Initialize as slave
		services, err := newSlaveServices()
		if err != nil {
			log.Fatal(err)
		}

		s.services = services
		s.grpc.server = grpclib.NewServer()

		// Enable server reflection for easier debugging
		reflection.Register(s.grpc.server)

		// Register the server with the gRPC server
		pb.RegisterServicesServer(s.grpc.server, s)
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

func (s *Server) initializeK8sResources() error {
	clientset, err := cluster.Init()
	if err != nil {
		return err
	}

	s.clientset = clientset

	if err := s.clientset.UpsertNamespace(); err != nil {
		return err
	}

	if err := s.clientset.UpsertConfigMap(config.SO_FILE); err != nil {
		return err
	}

	if err := s.clientset.UpsertClusterRole(); err != nil {
		return err
	}

	if err := s.clientset.UpsertClusterRoleBinding(); err != nil {
		return err
	}

	for service := range s.services.List() {
		if err := s.clientset.UpsertService(service); err != nil {
			return err
		}

		if err := s.clientset.UpsertDeployment(service); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) Handle(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	if err := s.services.Handle(req.CMD, req.Data); err != nil {
		return nil, err
	}

	log.Printf("Received CMD: %s, Data: %s", req.CMD, string(req.Data))

	return &pb.Response{Data: req.Data}, nil
}

func (s *Server) Info() {
	for service, handlers := range s.services.List() {
		for _, handler := range handlers {
			log.Printf("Registered handler: %v.%v", service, handler)
		}
	}
}

func (s *Server) ServeGRPC() error {
	if err := s.grpc.Listen(); err != nil {
		return err
	}

	if err := s.grpc.Serve(); err != nil {
		return err
	}

	return nil
}

func (g *grpc) Listen() error {
	listener, err := net.Listen(g.network, g.address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	g.listener = listener

	return nil
}

func (g *grpc) Serve() error {
	log.Printf("gRPC server is listening %v on: (%v)", g.network, g.address)

	if err := g.server.Serve(g.listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	return nil
}
