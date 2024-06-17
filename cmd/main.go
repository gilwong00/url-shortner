package main

import (
	"fmt"
	"os"

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
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	store, err := redisapp.NewRedisClient(redisHost, redisPort, redisPassword)
	// fmt.Println("redis ping", store.Ping(context.Background()))
	if err != nil {
		fmt.Println("failed to created redis store")
		panic(err)
	}
	s := server.NewServer(port, store)
	s.StartServer()
}
