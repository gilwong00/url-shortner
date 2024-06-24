package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gilwong00/url-shortner/pkg/models"
	"github.com/redis/go-redis/v9"
)

func (h *Handler) UpdateURL(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	defer r.Body.Close()
	shortName := r.PathValue("shortName")
	if shortName == "" {
		GenericErrorResponse(w, http.StatusBadRequest, "missing param")
		return
	}
	payload, err := ParseJSONBody[models.UpdateUrlRequest](r)
	if err != nil {
		GenericErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if payload.Expiry == 0 {
		GenericErrorResponse(w, http.StatusBadRequest, "expiration cannot be empty")
		return
	}
	// look up url
	result, err := h.store.Get(ctx, shortName).Result()
	if err == redis.Nil || result == "" {
		GenericErrorResponse(w, http.StatusNotFound, "could not find url for short name")
		return
	} else if err != nil {
		GenericErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err = h.store.Set(ctx, shortName, result, payload.Expiry*3600*time.Second).Err(); err != nil {
		GenericErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set(contentType, appJSON)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Expiry has been updated"))
}
