package main

import (
	"fmt"

	"github.com/gilwong00/url-shortner/pkg/config"
	"github.com/gilwong00/url-shortner/pkg/redisapp"
	"github.com/gilwong00/url-shortner/pkg/server"
	"github.com/joho/godotenv"
)

const (
	defaultPort = "8080"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("failed to load env vars", err)
	}
	config, err := config.NewConfig()
	store, err := redisapp.NewRedisClient(config.RedisHost, config.RedisPort, config.RedisPassword)
	// fmt.Println("redis ping", store.Ping(context.Background()))
	if err != nil {
		fmt.Println("failed to created redis store")
		panic(err)
	}
	s := server.NewServer(config.ServerPort, store)
	s.StartServer(config)
}
