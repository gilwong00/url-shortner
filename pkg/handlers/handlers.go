package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type Handler struct {
	Store *redis.Client
}

const (
	contentType = "content-type"
	appJSON     = "application/json"
)

func NewHandler(store *redis.Client) *Handler {
	return &Handler{
		Store: store,
	}
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}

func GenericErrorResponse(w http.ResponseWriter, statusCode int, errMessage string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(errMessage))
}

// ParseJSONBody unbox post request to type, returns ErrBadRequest if unable to parse
func ParseJSONBody[T any](r *http.Request) (T, error) {
	var parsed T
	err := json.NewDecoder(r.Body).Decode(&parsed)
	if err != nil {
		return parsed, fmt.Errorf("unable to parse JSON: `%s`", err)
	}
	return parsed, nil
}
