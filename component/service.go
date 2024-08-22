package component

import (
	"errors"
	"log"
	"reflect"
	"runtime/debug"
	"strings"
)

type (
	service struct {
		Name          string              // name of service
		Type          reflect.Type        // type of the receiver
		Receiver      reflect.Value       // receiver of methods for the service
		Handlers      map[string]*Handler // registered methods
		SchedulerName string              // name of scheduler variable in session data
	}

	Services struct {
		services map[string]*service
	}
)

func (ss *Services) Setup(comps *Components) {
	if ss.services == nil {
		ss.services = make(map[string]*service, 0)
	}

	for _, comp := range comps.comps {
		s := newService(comp)
		ss.services[s.Name] = s
	}
}

func (ss *Services) List() {
	for serviceName, service := range ss.services {
		for handlerName := range service.Handlers {
			log.Printf("CMD: %s.%s", serviceName, handlerName)
		}
	}
}

func (ss *Services) Handle(cmd string, data []byte) error {
	parts := strings.Split(cmd, ".")
	if len(parts) != 2 {
		return errors.New("invalid cmd")
	}

	service, exist := ss.services[parts[0]]
	if !exist {
		return errors.New("service not exist")
	}

	handler, exist := service.Handlers[parts[1]]
	if !exist {
		return errors.New("handler not exist")
	}

	args := make([]reflect.Value, handler.Type.NumIn())
	args[0] = handler.Receiver
	args[1] = reflect.ValueOf(data)

	task := func() {
		// Call the handler
		results := handler.Method.Func.Call(args)

		// Check results
		if len(results) > 0 {
			if err := results[0].Interface(); err != nil {
				log.Printf("Handle request on service error: %+v", err)
			}
		}
	}

	safecall(task)

	return nil
}

func safecall(task func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered: %v", r)
			log.Printf("Stack trace: %s", debug.Stack())
		}
	}()
	task()
}

func newService(comp Component) *service {
	s := &service{
		Type:     reflect.TypeOf(comp),
		Receiver: reflect.ValueOf(comp),
		Handlers: make(map[string]*Handler),
	}

	// Register handlers
	s.registerHandlers()

	// Store service name
	s.Name = reflect.Indirect(s.Receiver).Type().Name()

	return s
}

func (s *service) registerHandlers() {
	for i := 0; i < s.Type.NumMethod(); i++ {
		m := s.Type.Method(i)

		if isHandlerMethod(m) {
			// Check raw data or not
			raw := false
			if m.Type.In(1) == typeBytes {
				raw = true
			}

			s.Handlers[m.Name] = &Handler{
				Type:     m.Type,
				Receiver: s.Receiver,
				Method:   m,
				IsRawArg: raw,
			}
		}
	}
}
