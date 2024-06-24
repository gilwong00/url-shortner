package models

import "time"

type NewShortenUrlRequest struct {
	URL         string        `json:"url" validate:"required"`
	ShortenName string        `json:"shortenName"`
	Expiry      time.Duration `json:"expiry"`
}

type NewShortenUrlResponse struct {
	ShortURL string        `json:"shortURL"`
	Expiry   time.Duration `json:"expiry"`
}

type UpdateUrlRequest struct {
	Expiry time.Duration `json:"expiry" validate:"required"`
}
