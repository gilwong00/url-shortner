package handlers

import "context"

type Handler struct{}

func (h *Handler) ShortenURL(ctx context.Context) (string, error) {
	return "", nil
}
