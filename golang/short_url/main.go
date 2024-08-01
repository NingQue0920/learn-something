package main

import (
	"github.com/BurntSushi/toml"
	"log"
	"net/http"
)

type Config struct {
	GenType string `toml:"gen_type"`
	Host    string `toml:"host"`
	Port    string `toml:"port"`
}

var config *Config

func main() {

	config = ReadConfig()
	switch config.GenType {
	case "redis":
		redisGen = NewRedisGenerator(config.Host, config.Port)
	case "murmurhash":
		murmurGen = NewMurmurGenerator()
	}
	http.HandleFunc("/shorten", ShortUrlHandler)
	http.ListenAndServe(":8080", nil)
}

func ReadConfig() *Config {
	var config Config
	_, err := toml.DecodeFile("./app.conf", &config)
	if err != nil {
		log.Fatal(err)
	}
	return &config
}
