package middleware

import (
	"fmt"
	"sync"
)

type MiddlewareFactory func() (Middleware, error)

var (
	middlewareFactories = make(map[string]MiddlewareFactory)
	mutex               = sync.Mutex{}
)

func getMiddlewareFactory(name string) (MiddlewareFactory, error) {
	factory, ok := middlewareFactories[name]
	if !ok {
		return nil, fmt.Errorf("middleware %s not supported or not compiled in", name)
	}
	return factory, nil
}

func RegisterMiddleware(name string, factory MiddlewareFactory) {
	mutex.Lock()
	defer mutex.Unlock()
	middlewareFactories[name] = factory
}
