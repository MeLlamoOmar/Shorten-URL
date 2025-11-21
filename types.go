package main

type ErrShortCodeAlreadyExists struct{}

func (e ErrShortCodeAlreadyExists) Error() string {
	return "short code already exists"
}

type ShortenRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type ShortenResponse struct {
	ID int64  `json:"id"`
	URL string `json:"url"`
	ShortCode string `json:"shortCode"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ShortenStatResponse struct {
	ID int64  `json:"id"`
	URL string `json:"url"`
	ShortCode string `json:"shortCode"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	AccessCount int64  `json:"accessCount"`
}

type Store interface {
	Create(shortenURL *ShortenStatResponse) error
	Get(shortCode string) (*ShortenStatResponse, error)
	Update(shortCode string, url string) (*ShortenStatResponse, error)
	Delete(shortCode string) error
	Exists(shortCode string) bool
	IncrementAccessCount(shortCode string) error
	NextID() int64
}