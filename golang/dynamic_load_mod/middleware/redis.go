//go:build redis
// +build redis

package middleware

import (
	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
}

func (r *Redis) Initialize() error {
	// init redis
	r.client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return nil
}

func (r *Redis) Set(key, value string) error {
	return r.client.Set(key, value, 0).Err()
}
func (r *Redis) Get(key string) string {
	return r.client.Get(key).Val()
}
