package tora

import (
	"tora/server"
)

func Serve(opts ...server.Option) (*server.Server, error) {
	return server.New(opts...), nil
}
