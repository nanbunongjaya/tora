package plugin

import (
	"fmt"
)

const (
	tag = "Plugin"
)

type Plugin struct {
	s string
	n int
	f func()
}

func (p *Plugin) New() interface{} {
	return &Plugin{
		s: "",
		n: 0,
		f: func() {},
	}
}

func (p *Plugin) Test(data []byte) error {
	fmt.Println(tag + " Test...")
	return nil
}

func (p *Plugin) Init() {
	fmt.Println(tag + " init...")
}

func (p *Plugin) Shutdown() {
	fmt.Println(tag + " shutdown...")
}
