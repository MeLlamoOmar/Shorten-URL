package main

import (
	"database/sql"
	"log"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// inMemoryStore := NewInMemoryStore()
	db, err := sql.Open("sqlite3", "./urlShortener.db")
	if err != nil {
		log.Fatal("Cannot open DB")
	}
	defer db.Close()
	
	store := NewSQLStore(db)
	urlService := NewURLService(store)
	handler := NewHandler(urlService)
	e := echo.New()

	shorten := e.Group("/shorten")
	shorten.POST("", handler.HandlePost)
	shorten.GET("/:shortCode", handler.HandleGet)
	shorten.PUT("/:shortCode", handler.HandlePut)
	shorten.DELETE("/:shortCode", handler.HandleDelete)
	shorten.GET("/:shortCode/stats", handler.HandleStats)

	e.Logger.Fatal(e.Start(":8000"))
}
