package main

type ShortenRequest struct {
	OriginalURL string `json:"url" validate:"required,url"`
}

type ShortenResponse struct {
	ID          int64  `json:"id"`
	OriginalURL string `json:"url"`
	ShortCode   string `json:"shortCode"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}

type ShortenStatResponse struct {
	ID          int64  `json:"id"`
	OriginalURL string `json:"url"`
	ShortCode   string `json:"shortCode"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
	AccessCount int    `json:"accessCount"`
}

type Store interface {
	Create(shortCode, originalURL string) (*ShortenStatResponse, error)
	Get(shortCode string) (*ShortenStatResponse, error)
	Update(shortCode string, url string) (*ShortenStatResponse, error)
	Delete(shortCode string) error
	Exists(shortCode string) bool
	IncrementAccessCount(shortCode string, currentCount int) error
}
