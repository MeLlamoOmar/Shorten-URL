package main

import (
	"database/sql"
	"time"
)

type SQLStore struct {
	db *sql.DB
}

func NewSQLStore(db *sql.DB) Store {
	return &SQLStore{
		db: db,
	}
}

func (s *SQLStore) Create(shortCode, originalURL string) (*ShortenStatResponse, error) {
	q := "INSERT INTO urls (original_url, short_code) VALUES (?, ?)"

	var u ShortenStatResponse
	err := s.db.QueryRow(q, shortCode, originalURL).Scan(&u.ID, &u.OriginalURL, &u.ShortCode, &u.CreatedAt, &u.UpdatedAt, &u.AccessCount)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *SQLStore) Get(shortCode string) (*ShortenStatResponse, error) {
	q := `SELECT id, original_url, short_code, created_at, updated_at, access_count FROM urls WHERE short_code = ?`

	var u ShortenStatResponse
	err := s.db.QueryRow(q, shortCode).Scan(&u.ID, &u.OriginalURL, &u.ShortCode, &u.CreatedAt, &u.UpdatedAt, &u.AccessCount)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *SQLStore) Update(shortCode string, url string) (*ShortenStatResponse, error) {
	q := `UPDATE urls SET original_url = ?, updated_at = ? WHERE short_code = ?`

	row := s.db.QueryRow(q, url, time.Now().Unix(), shortCode)
	var u ShortenStatResponse
	if err := row.Scan(&u.ID, &u.OriginalURL, &u.ShortCode, &u.CreatedAt, &u.UpdatedAt, &u.AccessCount); err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *SQLStore) Delete(shortCode string) error {
	q := `DELETE FROM urls WHERE short_code = ?`

	if _, err := s.db.Exec(q, shortCode); err != nil {
		return err
	}

	return nil
}

func (s *SQLStore) Exists(shortCode string) bool {
	q := `SELECT id WHERE short_code = ?`

	r, _ := s.db.Query(q, shortCode)
	if !r.Next() {
		return false
	}
	defer r.Close()

	return true
}

func (s *SQLStore) IncrementAccessCount(shortCode string, currentCount int) error {
	q := `UPDATE urls SET access_count = ?, updated_at = ? WHERE short_code = ?`

	_, err := s.db.Exec(q, currentCount+1, shortCode)
	if err != nil {
		return err
	}

	return nil
}
