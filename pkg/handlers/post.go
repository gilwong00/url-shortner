package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gilwong00/url-shortner/pkg/models"
	"github.com/google/uuid"
)

const (
	defaultExpiration = 24
)

func (h *Handler) CreateShortenURL(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	defer r.Body.Close()
	payload, err := ParseJSONBody[models.NewShortenUrlRequest](r)
	if err != nil {
		GenericErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	// validation
	domain := h.config.Domain
	// ensures url has http or https protocol prefix
	payload.URL = appendProtocol(payload.URL)
	url, err := url.ParseRequestURI(payload.URL)
	if err != nil {
		GenericErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	isSameDomain, err := isDifferentDomain(url, domain)
	if err != nil {
		GenericErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if !isSameDomain {
		GenericErrorResponse(w, http.StatusBadRequest, "domain cannot be different")
		return
	}
	// logic
	if payload.ShortenName == "" {
		payload.ShortenName = uuid.New().String()[:8]
	}
	// check if shortname already exist in redis
	val, err := h.store.Get(ctx, payload.ShortenName).Result()
	if !isSameDomain {
		GenericErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if val != "" {
		GenericErrorResponse(w, http.StatusForbidden, "shorten url already exists")
		return
	}
	if payload.Expiry == 0 {
		payload.Expiry = defaultExpiration
	}
	// creating
	expiration := payload.Expiry * 3600 * time.Second
	if err = h.store.Set(ctx, payload.ShortenName, payload.URL, expiration).Err(); err != nil {
		GenericErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp := models.NewShortenUrlResponse{
		ShortURL: fmt.Sprintf("%s/%s", domain, payload.ShortenName),
		Expiry:   expiration,
	}
	w.Header().Set(contentType, appJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func isDifferentDomain(u *url.URL, domain string) (bool, error) {
	hostname := strings.TrimPrefix(u.Host, "www.")
	fmt.Println(hostname)
	serverHost := domain
	port := u.Port()
	if port != "" {
		serverHost = fmt.Sprintf("%s:%s", domain, port)
	}
	return hostname == serverHost, nil
}

func appendProtocol(url string) string {
	if strings.HasPrefix("http://", url) || strings.HasPrefix("https://", url) {
		return url
	}
	return fmt.Sprintf("https://%s", url)
}
