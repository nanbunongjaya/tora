package main

import (
	"github.com/nanbunongjaya/tora"
	"github.com/nanbunongjaya/tora/component"
	"github.com/nanbunongjaya/tora/example/greeter"
	"github.com/nanbunongjaya/tora/example/plugin"
	"github.com/nanbunongjaya/tora/server"
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
