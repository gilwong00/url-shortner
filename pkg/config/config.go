package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort      int
	RedisHost       string
	RedisPort       int
	RedisPassword   string
	Domain          string
	MaxRequestLimit int
}

const (
	defaultMaxRequestLimit = int(10)
)

func NewConfig() (*Config, error) {
	port := os.Getenv("PORT")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	domain := os.Getenv("DOMAIN")
	maxRequestLimit := os.Getenv("MAX_REQUEST_LIMIT")
	rPort, err := strconv.Atoi(redisPort)
	if err != nil {
		return nil, err
	}
	serverPort, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	requestLimit, err := strconv.Atoi(maxRequestLimit)
	if err != nil {
		return nil, err
	}
	if requestLimit == 0 {
		requestLimit = defaultMaxRequestLimit
	}
	return &Config{
		ServerPort:      serverPort,
		RedisHost:       redisHost,
		RedisPort:       rPort,
		RedisPassword:   redisPassword,
		Domain:          domain,
		MaxRequestLimit: requestLimit,
	}, nil
}
