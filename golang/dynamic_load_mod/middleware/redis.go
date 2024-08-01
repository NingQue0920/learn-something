//go:build redis
// +build redis

package middleware

import (
	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
}

func init() {
	RegisterMiddleware("redis", NewRedis)
}

func NewRedis() (Middleware, error) {
	return &Redis{}, nil
}

func (r *Redis) Initialize() error {
	// init redis
	r.client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return nil
}

func (r *Redis) Write(key, value string) error {
	return r.client.Set(key, value, 0).Err()
}
func (r *Redis) Read(key string) (any, error) {
	return r.client.Get(key).Val(), nil
}
