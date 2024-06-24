package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gilwong00/url-shortner/pkg/models"
	"github.com/google/uuid"
)

func (h *Handler) CreateShortenURL(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	defer r.Body.Close()
	var payload models.NewShortenUrlRequest
	parsed, err := ParseJSONBody[models.NewShortenUrlRequest](r)
	if err != nil {
		GenericErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	// validation
	// check for valid url
	url, err := url.ParseRequestURI(parsed.URL)
	if err != nil {
		GenericErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	isSameDomain, err := isDifferentDomain(url.Host, "")
	if err != nil {
		GenericErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if !isSameDomain {
		GenericErrorResponse(w, http.StatusBadRequest, "domain cannot be different")
		return
	}
	// logic
	if payload.Expiry == 0 {
		payload.Expiry = 24
	}
	// make sure url has http or https prefix
	payload.URL = appendPrefix(payload.URL)
	if payload.ShortenName == "" {
		payload.ShortenName = uuid.New().String()[:8]
	}
	// check if shortname already exist in redis
	val, err := h.Store.Get(ctx, payload.ShortenName).Result()
	if !isSameDomain {
		GenericErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if val != "" {
		GenericErrorResponse(w, http.StatusForbidden, "shorten url already exists")
		return
	}
	// saving
	expiration := parsed.Expiry * 3600 * time.Second
	if err = h.Store.Set(ctx, payload.ShortenName, payload.URL, expiration).Err(); err != nil {
		GenericErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	// TODO replace with env var
	domain := ""
	resp := models.NewShortenUrlResponse{
		ShortURL: fmt.Sprintf("%s/%s", domain, payload.ShortenName),
		Expiry:   expiration,
	}
	w.Header().Set(contentType, appJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func isDifferentDomain(u string, domain string) (bool, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}
	hostname := strings.TrimPrefix(parsed.Hostname(), "www.")
	fmt.Println(hostname)
	return hostname == domain, nil
}

func appendPrefix(url string) string {
	if strings.HasPrefix("http://", url) || strings.HasPrefix("https://", url) {
		return url
	}
	return fmt.Sprintf("https://%s", url)
}
