package tora

import (
	"tora/server"
)

func NewServer(opts ...server.Option) *server.Server {
	return server.New(opts...)
}

func ServeWithGRPC(opts ...server.Option) error {
	s := server.New(opts...)

	s.Info()

	return s.ServeGRPC()
}
