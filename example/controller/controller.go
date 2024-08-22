package main

import (
	"tora/component"
	"tora/example/greeter"
	"tora/example/plugin"
)

var (
	comps    = &component.Components{}
	services = &component.Services{}
)

func Setup() {
	// Registrations of components
	comps.Register(&plugin.Plugin{})
	comps.Register(&greeter.Greeter{})

	// Setup services
	services.Setup(comps)
}

func List() {
	services.List()
}

func Handle(cmd string, data []byte) error {
	return services.Handle(cmd, data)
}
