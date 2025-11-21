package main

import "fmt"

type inMemoryStore struct {
	data map[string]*ShortenStatResponse
}

func NewInMemoryStore() Store {
	return &inMemoryStore{
		data: make(map[string]*ShortenStatResponse),
	}
}

func (s *inMemoryStore) Create(shortenURL *ShortenStatResponse) error {
	if s.data[shortenURL.URL] != nil {
		return fmt.Errorf("esta url ya existe")
	}

	s.data[shortenURL.ShortCode] = shortenURL
	return nil
}

func (s *inMemoryStore) Get(shortCode string) (*ShortenStatResponse, error) {
	if shortenURL, exists := s.data[shortCode]; exists {
		return shortenURL, nil
	}

	return nil, fmt.Errorf("url no encontrada")
}

func (s *inMemoryStore) Update(shortCode string, url string) (*ShortenStatResponse, error) {
	if shortenURL, exists := s.data[shortCode]; exists {
		shortenURL.URL = url
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

func (s *inMemoryStore) IncrementAccessCount(shortCode string) error {
	if shortenURL, exists := s.data[shortCode]; exists {
		shortenURL.AccessCount++
		return nil
	}

	return fmt.Errorf("url no encontrada")
}

func (s *inMemoryStore) NextID() int64 {
	return int64(len(s.data) + 1)
}