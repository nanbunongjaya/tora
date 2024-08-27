package main

import (
	"github.com/nanbunongjaya/tora"
	"github.com/nanbunongjaya/tora/server"
)

func main() {
	tora.ServeWithGRPC(
		server.WithAsSlave(),
		server.WithInCloud(),
		server.WithGRPC("tcp", ":50051"),
	)
}
