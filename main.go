package main

import (
	"log"
	"plugin"

	"tora/compiler"
	"tora/component"
	plug "tora/plugin"
)

func main() {
	Run()
	//Compile()
}

func Compile() {
	comps := &component.Components{}
	comps.Register(&plug.Plugin{})

	program, err := compiler.Compile(comps)
	if err != nil {
		log.Fatal(err)
	}

	err = compiler.Output("controller/controller.go", program)
	if err != nil {
		log.Fatal(err)
	}

	err = compiler.Build("controller/controller.go")
	if err != nil {
		log.Fatal(err)
	}
}

func Run() {
	p, err := plugin.Open("controller.so")
	if err != nil {
		log.Fatal(err)
	}

	s, err := p.Lookup("Setup")
	if err != nil {
		log.Fatal(err)
	}
	s.(func())()

	s, err = p.Lookup("List")
	if err != nil {
		log.Fatal(err)
	}
	s.(func())()

	s, err = p.Lookup("Handle")
	if err != nil {
		log.Fatal(err)
	}
	err = s.(func(string, []byte) error)("Plugin.Test", []byte("123"))
	if err != nil {
		log.Fatal(err)
	}
}
