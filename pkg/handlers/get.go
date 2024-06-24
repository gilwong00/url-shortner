package handlers

import (
	"context"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func (h *Handler) GetURL(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	shortName := r.PathValue("shortName")
	if shortName == "" {
		GenericErrorResponse(w, http.StatusBadRequest, "missing param")
		return
	}
	result, err := h.store.Get(ctx, shortName).Result()
	if err == redis.Nil {
		GenericErrorResponse(w, http.StatusNotFound, "could not find url for short name")
	} else if err != nil {
		GenericErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set(contentType, appJSON)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}
