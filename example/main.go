package main

import (
	"log"

	"tora"
	"tora/component"
	"tora/example/greeter"
	"tora/example/plugin"
	"tora/server"
)

func main() {
	Compile()
}

func Compile() {
	comps := &component.Components{}
	comps.Register(&plugin.Plugin{})
	comps.Register(&greeter.Greeter{})

	s, err := tora.Serve(
		server.WithComponents(comps),
		server.WithInCloud(),
	)
	if err != nil {
		log.Print(err)
	}

	s.Handle("Greeter.Increase", nil)
	s.Handle("Greeter.Increase", nil)
	s.Handle("Greeter.Increase", nil)
	s.Handle("Greeter.Increase", nil)
	s.Handle("Greeter.Increase", nil)
	s.Handle("Greeter.Increase", nil)
	s.Handle("Greeter.Increase", nil)
	s.Handle("Greeter.Increase", nil)
}
