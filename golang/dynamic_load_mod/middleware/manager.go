package middleware

type Middleware interface {
	Initialize() error
}

type Manager struct {
	middlewares map[string]Middleware
}

// Common interfaces for all middleware types
type Setter interface {
	Set(key, value string) error
}

type Producer interface {
	Produce(topic, message string) error
}

type Selector interface {
	Select(key string) (any, error)
}

func (m *Manager) InitMiddleware(name string) error {
	factory, err := getMiddlewareFactory(name)
	if err != nil {
		return err
	}
	middleware, err := factory()
	if err != nil {
		return err
	}
	if err = middleware.Initialize(); err != nil {
		return err
	}
	m.middlewares[name] = middleware
	return nil
}

func (m *Manager) GetMiddleware(name string) (Middleware, bool) {
	mw, ok := m.middlewares[name]
	return mw, ok
}
func NewManager() *Manager {
	return &Manager{
		middlewares: make(map[string]Middleware),
	}
}

// Type assertion helpers
func (m *Manager) GetSetter(name string) (Setter, bool) {
	if mw, ok := m.GetMiddleware(name); ok {
		if setter, ok := mw.(Setter); ok {
			return setter, true
		}
	}
	return nil, false
}

func (m *Manager) GetProducer(name string) (Producer, bool) {
	if mw, ok := m.GetMiddleware(name); ok {
		if producer, ok := mw.(Producer); ok {
			return producer, true
		}
	}
	return nil, false
}

func (m *Manager) GetSelector(name string) (Selector, bool) {
	if mw, ok := m.GetMiddleware(name); ok {
		if selector, ok := mw.(Selector); ok {
			return selector, true
		}
	}
	return nil, false
}
