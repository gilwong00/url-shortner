package redisapp

import (
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(host string, port string, password string) (*redis.Client, error) {
	if host == "" {
		return nil, errors.New("missing redis host")
	}
	if port == "" {
		return nil, errors.New("missing redis port")
	}
	address := fmt.Sprintf("%s:%s", host, port)
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		// Maybe config this?
		DB: 1,
	})
	return client, nil
}
