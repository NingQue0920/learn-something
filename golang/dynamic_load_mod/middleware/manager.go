package middleware

type Middleware interface {
	Initialize() error
}

type Manager struct {
	middlewares map[string]Middleware
}

// Common interfaces for all middleware types
type Reader interface {
	Read(key string) (any, error)
}

type Writer interface {
	Write(topic, message string) error
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
func (m *Manager) GetReader(name string) (Reader, bool) {
	if mw, ok := m.GetMiddleware(name); ok {
		if setter, ok := mw.(Reader); ok {
			return setter, true
		}
	}
	return nil, false
}

func (m *Manager) GetWriter(name string) (Writer, bool) {
	if mw, ok := m.GetMiddleware(name); ok {
		if producer, ok := mw.(Writer); ok {
			return producer, true
		}
	}
	return nil, false
}
