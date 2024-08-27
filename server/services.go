package server

import (
	"log"
	"plugin"

	"tora/component"
	"tora/config"
)

type (
	services interface {
		List()
		Handle(cmd string, data []byte) error
	}

	master struct {
		services *component.Services
	}

	slave struct {
		plugin *plugin.Plugin
	}
)

func newMasterServices(comps *component.Components) services {
	s := &component.Services{}
	s.Setup(comps)

	return &master{
		services: s,
	}
}

func (m *master) Handle(cmd string, data []byte) error {
	return m.Handle(cmd, data)
}

func (m *master) List() {
	m.services.List()
}

func newSlaveServices() (services, error) {
	p, err := plugin.Open(config.SO_FILE)
	if err != nil {
		return nil, err
	}

	s := &slave{plugin: p}

	// Call "Setup" function
	f, err := p.Lookup("Setup")
	if err != nil {
		return nil, err
	}
	f.(func())()

	return s, nil
}

func (s *slave) Handle(cmd string, data []byte) error {
	f, err := s.plugin.Lookup("Handle")
	if err != nil {
		return err
	}

	return f.(func(string, []byte) error)(cmd, data)
}

func (s *slave) List() {
	// Call "List" function
	f, err := s.plugin.Lookup("List")
	if err != nil {
		log.Fatal(err)
	}
	f.(func())()
}
