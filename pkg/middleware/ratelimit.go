package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

// ErrorResponse format response (this is public facing)
type ErrorResponse struct {
	Error   string `json:"error"`
	Message any    `json:"message"`
}

func RateLimiter(next http.Handler, ctx context.Context, store *redis.Client, maxLimit int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipAddress := GetIPAddress(r)
		counter, err := store.Get(ctx, ipAddress).Int64()
		// could not find the IP in the DB
		// we need to set it. This means that this IP is making a request for the first time
		if err == redis.Nil {
			// this sets the IP in redis and gives it a TTL of 1 hour
			err = store.Set(ctx, ipAddress, 1, 1*time.Hour).Err()
			if err != nil {
				// return err
				writeErrResponse(w, 500, err.Error())
			}
		} else if err != nil {
			writeErrResponse(w, 500, err.Error())
		} else {
			// Check if rate limit is exceeded
			if counter > int64(maxLimit) {
				limit, err := store.TTL(ctx, ipAddress).Result()
				if err != nil {
					writeErrResponse(w, 400, err.Error())
				}
				writeErrResponse(w, 500, fmt.Sprintf("rate limit exceeded, will reset in: %v", limit/time.Nanosecond/time.Minute))
			}
			_, err = store.Incr(ctx, ipAddress).Result()
			if err != nil {
				writeErrResponse(w, 500, err.Error())
			}
		}
		next.ServeHTTP(w, r)
	})
}

func GetIPAddress(r *http.Request) string {
	ipAddress := r.Header.Get("X-Real-Ip")
	if ipAddress == "" {
		ipAddress = r.Header.Get("X-FORWARDED-FOR")
	}
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}
	return ipAddress
}

func writeErrResponse(w http.ResponseWriter, statusCode int, msg string) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: msg,
	})
}
