package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	inMemoryStore := NewInMemoryStore()
	urlService := NewURLService(inMemoryStore)
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