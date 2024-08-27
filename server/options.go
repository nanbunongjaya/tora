package server

import (
	"tora/component"
)

type Option func(s *Server)

func WithComponents(comps *component.Components) Option {
	return func(s *Server) {
		for _, comp := range comps.List() {
			s.comps.Register(comp)
		}
	}
}

func WithInCloud() Option {
	return func(s *Server) {
		s.incloud = true
	}
}

func WithAsSlave() Option {
	return func(s *Server) {
		s.ismaster = false
	}
}

func WithGRPC(network, address string) Option {
	return func(s *Server) {
		s.grpc.enable = true
		s.grpc.network = network
		s.grpc.address = address
	}
}
