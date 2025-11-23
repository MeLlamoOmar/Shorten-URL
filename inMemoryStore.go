package main

import (
	"fmt"
	"time"
)

type inMemoryStore struct {
	data map[string]*ShortenStatResponse
}

func NewInMemoryStore() Store {
	return &inMemoryStore{
		data: make(map[string]*ShortenStatResponse),
	}
}

func (s *inMemoryStore) Create(shortCode, originalURL string) (*ShortenStatResponse, error) {
	shorten := ShortenStatResponse{
		ID:          s.NextID(),
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		AccessCount: 0,
	}

	if s.data[shorten.OriginalURL] != nil {
		return nil, fmt.Errorf("esta url ya existe")
	}

	s.data[shorten.ShortCode] = &shorten
	return &shorten, nil
}

func (s *inMemoryStore) Get(shortCode string) (*ShortenStatResponse, error) {
	if shortenURL, exists := s.data[shortCode]; exists {
		return shortenURL, nil
	}

	return nil, fmt.Errorf("url no encontrada")
}

func (s *inMemoryStore) Update(shortCode string, url string) (*ShortenStatResponse, error) {
	if shortenURL, exists := s.data[shortCode]; exists {
		shortenURL.OriginalURL = url
		shortenURL.UpdatedAt = time.Now().Unix()
		return shortenURL, nil
	}

	return nil, fmt.Errorf("url no encontrada")
}

func (s *inMemoryStore) Delete(shortCode string) error {
	if _, exists := s.data[shortCode]; exists {
		delete(s.data, shortCode)
		return nil
	}

	return fmt.Errorf("url no encontrada")
}

func (s *inMemoryStore) Exists(shortCode string) bool {
	_, exists := s.data[shortCode]
	return exists
}

func (s *inMemoryStore) IncrementAccessCount(shortCode string, currentCount int) error {
	if shortenURL, exists := s.data[shortCode]; exists {
		shortenURL.AccessCount = currentCount + 1
		return nil
	}

	return fmt.Errorf("url no encontrada")
}

func (s *inMemoryStore) NextID() int64 {
	return int64(len(s.data) + 1)
}
