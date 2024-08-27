package main

import (
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

	tora.NewServer(
		server.WithComponents(comps),
	)
}
