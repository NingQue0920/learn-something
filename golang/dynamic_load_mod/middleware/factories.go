package middleware

import "fmt"

type MiddlewareFactory func() (Middleware, error)

var middlewareFactories = make(map[string]MiddlewareFactory)

func getMiddlewareFactory(name string) (MiddlewareFactory, error) {
	factory, ok := middlewareFactories[name]
	if !ok {
		return nil, fmt.Errorf("middleware %s not supported or not compiled in", name)
	}
	return factory, nil
}
