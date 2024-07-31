package main

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"sync"
)

type RedisGenerator struct {
	client *redis.Client
	once   sync.Once
}

func NewRedisGenerator(host, port string) *RedisGenerator {

	return &RedisGenerator{
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: "Yjc@123456", // no password set
			DB:       0,            // use default DB
		}),
	}
}
func (rg *RedisGenerator) Generate(input string) string {
	// 确保Redis中有初始的计数器
	rg.once.Do(func() {
		err := rg.client.SetNX("short_url:counter", 1000000000, 0).Err()
		if err != nil {
			log.Println("SET NX ERROR : ", err.Error())
		}
	})
	// load from cache
	shorten, err := rg.client.Get("origin_url:" + input).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Println("Get origin_url error: ", err.Error())
	} else if errors.Is(err, redis.Nil) {
		//	cache no data , generate short code
	} else {
		return shorten
	}
	id := rg.client.Incr("short_url:counter").Val()
	shortCode := FormatInt62(uint64(id))
	rg.Store(input, shortCode)
	return shortCode
}

func (rg *RedisGenerator) Store(input, shorten string) {
	// 缓存短链，避免重复生成
	err := rg.client.SetNX("origin_url:"+input, shorten, 0).Err()
	if err != nil {
		log.Println("SET NX ERROR : ", err.Error())
	}
}
