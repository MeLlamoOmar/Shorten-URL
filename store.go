package main

import "fmt"

type Store interface {
	Create(shortenURL *ShortenStatResponse) error
	Get(shortCode string) (*ShortenStatResponse, error)
	Update(shortCode string, url string) (*ShortenStatResponse, error)
	Delete(shortCode string) error
	Exists(shortCode string) bool
	IncrementAccessCount(shortCode string) error
	NextID() int64
}

type InMemoryStore struct {
	data map[string]*ShortenStatResponse
}

func NewInMemoryStore() Store {
	return &InMemoryStore{
		data: make(map[string]*ShortenStatResponse),
	}
}

func (s *InMemoryStore) Create(shortenURL *ShortenStatResponse) error {
	if s.data[shortenURL.URL] != nil {
		return fmt.Errorf("esta url ya existe")
	}

	s.data[shortenURL.ShortCode] = shortenURL
	return nil
}

func (s *InMemoryStore) Get(shortCode string) (*ShortenStatResponse, error) {
	if shortenURL, exists := s.data[shortCode]; exists {
		return shortenURL, nil
	}

	return nil, fmt.Errorf("url no encontrada")
}

func (s *InMemoryStore) Update(shortCode string, url string) (*ShortenStatResponse, error) {
	if shortenURL, exists := s.data[shortCode]; exists {
		shortenURL.URL = url
		return shortenURL, nil
	}
	
	return nil, fmt.Errorf("url no encontrada")
}

func (s *InMemoryStore) Delete(shortCode string) error {
	if _, exists := s.data[shortCode]; exists {
		delete(s.data, shortCode)
		return nil
	}

	return fmt.Errorf("url no encontrada")
}

func (s *InMemoryStore) Exists(shortCode string) bool {
	_, exists := s.data[shortCode]
	return exists
}

func (s *InMemoryStore) IncrementAccessCount(shortCode string) error {
	if shortenURL, exists := s.data[shortCode]; exists {
		shortenURL.AccessCount++
		return nil
	}

	return fmt.Errorf("url no encontrada")
}

func (s *InMemoryStore) NextID() int64 {
	return int64(len(s.data) + 1)
}