package server

import (
	"log"
	"plugin"

	"tora/compiler"
	"tora/component"
)

const (
	sofile = "tora_slave_services_plugin.so"
	gofile = "tora_slave_services/tora_slave_services.go"
)

type (
	Server struct {
		comps    *component.Components
		services *component.Services
		plugin   *plugin.Plugin
		incloud  bool
	}
)

func New(opts ...Option) *Server {
	s := &Server{
		comps:    &component.Components{},
		services: &component.Services{},
		plugin:   &plugin.Plugin{},
		incloud:  false,
	}

	for i := range opts {
		opts[i](s)
	}

	// Setup the local services
	if !s.incloud {
		s.services.Setup(s.comps)
		return s
	}

	if err := s.compile(); err != nil {
		log.Fatal(err)
	}

	if err := s.setup(); err != nil {
		log.Fatal(err)
	}

	return s
}

func (s *Server) setup() error {
	p, err := plugin.Open(sofile)
	if err != nil {
		return err
	}

	s.plugin = p

	// Call setup function
	f, err := p.Lookup("Setup")
	if err != nil {
		return err
	}
	f.(func())()

	// Call list function
	f, err = p.Lookup("List")
	if err != nil {
		return err
	}
	f.(func())()

	return nil
}

func (s *Server) compile() error {
	program, err := compiler.Compile(s.comps)
	if err != nil {
		return err
	}

	err = compiler.Output(gofile, program)
	if err != nil {
		return err
	}

	err = compiler.Build(sofile, gofile)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Handle(cmd string, data []byte) error {
	if !s.incloud {
		return s.services.Handle(cmd, data)
	}

	f, err := s.plugin.Lookup("Handle")
	if err != nil {
		return err
	}

	err = f.(func(string, []byte) error)(cmd, data)
	if err != nil {
		return err
	}

	return nil
}
