package component

type (
	// Component interface with lifecycle functions
	Component interface {
		New() interface{}
		Init()
		Shutdown()
	}

	Components struct {
		comps []Component
	}
)

// Register registers a component with options to Components
func (cs *Components) Register(c Component) {
	cs.comps = append(cs.comps, c.New().(Component))
}

// List returns all components with options
func (cs *Components) List() []Component {
	return cs.comps
}
