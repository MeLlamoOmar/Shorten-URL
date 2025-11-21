package main

import (
	"fmt"
	"time"
)

type Service interface {
	CreateShortURL(originalURL string) (*ShortenStatResponse, error)
	GetShortURL(shortCode string) (*ShortenStatResponse, error)
	UpdateShortURL(shortCode string, newURL string) (*ShortenStatResponse, error)
	DeleteShortURL(shortCode string) error
	IncrementAccessCount(shortCode string) error
}

type URLService struct {
	store Store
}

func NewURLService(store Store) Service {
	return &URLService{
		store: store,
	}
}

func (s *URLService) CreateShortURL(originalURL string) (*ShortenStatResponse, error) {
	shortCode := GenerateChortCode(originalURL)
	if s.store.Exists(shortCode) {
		return nil, fmt.Errorf("%s", ErrShortCodeAlreadyExists{})
	}

	shorten := &ShortenStatResponse{
		ID: s.store.NextID(),
		URL: originalURL,
		ShortCode: shortCode,
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
		AccessCount: 0,
	}
	
	if err := s.store.Create(shorten); err != nil {
		return nil, err
	}

	return shorten, nil
}

func (s *URLService) GetShortURL(shortCode string) (*ShortenStatResponse, error) {
	return s.store.Get(shortCode)
}

func (s *URLService) UpdateShortURL(shortCode string, newURL string) (*ShortenStatResponse, error) {
	updatedShorten, err := s.store.Update(shortCode, newURL)
	if err != nil {
		return nil, err
	}

	updatedShorten.UpdatedAt = time.Now().String()
	return updatedShorten, nil
}

func (s *URLService) DeleteShortURL(shortCode string) error {
	return s.store.Delete(shortCode)
}

func (s *URLService) IncrementAccessCount(shortCode string) error {
	return s.store.IncrementAccessCount(shortCode)
}