package component

type (
	// Component interface with lifecycle functions
	Component interface {
		New() interface{}
		Init()
		Shutdown()
	}

	Components struct {
		comps []CompWithOptions
	}

	CompWithOptions struct {
		Comp Component
		Opts []Option
	}
)

// Register registers a component with options to Components
func (cs *Components) Register(c Component, options ...Option) {
	cs.comps = append(cs.comps, CompWithOptions{c, options})
}

// List returns all components with options
func (cs *Components) List() []CompWithOptions {
	return cs.comps
}
