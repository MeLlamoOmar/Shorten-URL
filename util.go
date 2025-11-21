package main

import "github.com/google/uuid"

func GenerateChortCode(originalURL string) string {
	uuid := uuid.New()
	return uuid.String()[:6]
}

func TransformToShortenResponse(stat *ShortenStatResponse) *ShortenResponse {
	return &ShortenResponse{
		ID:        stat.ID,
		URL:       stat.URL,
		ShortCode: stat.ShortCode,
		CreatedAt: stat.CreatedAt,
		UpdatedAt: stat.UpdatedAt,
	}
}