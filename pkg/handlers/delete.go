package handlers

import (
	"context"
	"net/http"
)

func (h *Handler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	shortName := r.PathValue("shortName")
	if shortName == "" {
		GenericErrorResponse(w, http.StatusBadRequest, "missing param")
		return
	}
	if err := h.store.Del(ctx, shortName).Err(); err != nil {
		GenericErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set(contentType, appJSON)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Url been deleted"))
}
