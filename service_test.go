package main

import (
	"fmt"
	"testing"
)

type MockStore struct {
	Response   ShortenStatResponse
	Err        error
	ExistsResp bool
}

func (s *MockStore) Create(shortCode, originalURL string) (*ShortenStatResponse, error) {
	return &s.Response, s.Err
}

func (s *MockStore) Get(shortCode string) (*ShortenStatResponse, error) {
	return &s.Response, s.Err
}

func (s *MockStore) Update(shortCode string, url string) (*ShortenStatResponse, error) {
	s.Response.OriginalURL = url
	return &s.Response, s.Err
}

func (s *MockStore) Delete(shortCode string) error {
	return s.Err
}

func (s *MockStore) Exists(shortCode string) bool {
	return s.ExistsResp
}

func (s *MockStore) IncrementAccessCount(shortCode string, currentCount int) error {
	if s.Err != nil {
		return s.Err
	}
	s.Response.AccessCount = currentCount + 1
	return nil
}

func TestCreateShortUrlSuccess(t *testing.T) {
	store := MockStore{
		Response: ShortenStatResponse{
			OriginalURL: "mock.com",
			ShortCode: "1234",
		},
	}

	s := NewURLService(&store)

	res, err := s.CreateShortURL("mock.com")
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "mock.com"
	if res.OriginalURL != expected {
		t.Errorf("Expected: %s, Got: %s", expected, res.OriginalURL)
	}
}

func TestGetShortURL(t *testing.T) {
	store := MockStore{
		Response: ShortenStatResponse{
			OriginalURL: "mock.com",
			ShortCode:    "1234",
		},
	}

	s := NewURLService(&store)

	res, err := s.GetShortURL("1234")
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "mock.com"
	if res.OriginalURL != expected {
		t.Errorf("Expected: %s, Got: %s", expected, res.OriginalURL)
	}
}

func TestCreateShortURLCollision(t *testing.T) {
	store := MockStore{
		Response: ShortenStatResponse{OriginalURL: "mock.com", ShortCode: "abcd"},
		ExistsResp: true,
	}

	s := NewURLService(&store)

	_, err := s.CreateShortURL("mock.com")
	if err == nil {
		t.Fatal("Expected error due to short code collision, got nil")
	}
}

func TestUpdateShortURLSuccess(t *testing.T) {
	store := MockStore{
		Response: ShortenStatResponse{OriginalURL: "old.com", ShortCode: "u1"},
	}

	s := NewURLService(&store)

	updated, err := s.UpdateShortURL("u1", "new.com")
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	if updated.OriginalURL != "new.com" {
		t.Fatalf("Expected updated URL 'new.com', got: %s", updated.OriginalURL)
	}
}

func TestDeleteShortURLSuccess(t *testing.T) {
	store := MockStore{}
	s := NewURLService(&store)

	if err := s.DeleteShortURL("any"); err != nil {
		t.Fatalf("Expected nil error on delete, got: %v", err)
	}
}

func TestIncrementAccessCountSuccess(t *testing.T) {
	store := MockStore{
		Response: ShortenStatResponse{OriginalURL: "site", ShortCode: "c1", AccessCount: 5},
	}
	s := NewURLService(&store)

	if err := s.IncrementAccessCount("c1", 5); err != nil {
		t.Fatalf("Unexpected error incrementing access count: %v", err)
	}

	if store.Response.AccessCount != 6 {
		t.Fatalf("Expected access count 6, got %d", store.Response.AccessCount)
	}
}

func TestErrorPropagation(t *testing.T) {
	e := fmt.Errorf("db failure")
	store := MockStore{Err: e}
	s := NewURLService(&store)

	if _, err := s.CreateShortURL("x"); err == nil {
		t.Fatalf("Expected error from store, got nil")
	}
}

func TestGetShortURLError(t *testing.T) {
	e := fmt.Errorf("get failure")
	store := MockStore{Err: e}
	s := NewURLService(&store)

	if _, err := s.GetShortURL("nope"); err == nil {
		t.Fatalf("Expected error from Get, got nil")
	}
}

func TestUpdateShortURLError(t *testing.T) {
	e := fmt.Errorf("update failure")
	store := MockStore{Err: e}
	s := NewURLService(&store)

	if _, err := s.UpdateShortURL("u1", "new"); err == nil {
		t.Fatalf("Expected error from Update, got nil")
	}
}

func TestDeleteShortURLError(t *testing.T) {
	e := fmt.Errorf("delete failure")
	store := MockStore{Err: e}
	s := NewURLService(&store)

	if err := s.DeleteShortURL("d1"); err == nil {
		t.Fatalf("Expected error from Delete, got nil")
	}
}

func TestIncrementAccessCountError(t *testing.T) {
	e := fmt.Errorf("increment failure")
	store := MockStore{Err: e}
	s := NewURLService(&store)

	if err := s.IncrementAccessCount("c1", 1); err == nil {
		t.Fatalf("Expected error from IncrementAccessCount, got nil")
	}
}