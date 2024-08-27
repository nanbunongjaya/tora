package main

import (
	"tora"
	"tora/server"
)

func main() {
	tora.ServeWithGRPC(
		server.WithAsSlave(),
		server.WithInCloud(),
		server.WithGRPC("tcp", ":50051"),
	)
}
