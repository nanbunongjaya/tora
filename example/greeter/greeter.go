package greeter

import (
	"fmt"
	"strconv"
)

const (
	tag = "Greeter"
)

type Greeter struct {
	s string
	n int
	f func(*Greeter)
}

func (g *Greeter) New() interface{} {
	return &Greeter{
		s: "",
		n: 0,
		f: func(g *Greeter) { g.n++ },
	}
}

func (g *Greeter) Increase(date []byte) error {
	g.f(g)
	fmt.Println(tag + "  Increase..." + strconv.Itoa(g.n))
	return nil
}

func (g *Greeter) Test(data []byte) error {
	fmt.Println(tag + " Test...")
	return nil
}

func (g *Greeter) Init() {
	fmt.Println(tag + " init...")
}

func (g *Greeter) Shutdown() {
	fmt.Println(tag + " shutdown...")
}
