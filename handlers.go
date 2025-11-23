package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HandlePost(c echo.Context) error {
	originalURL := c.FormValue("url")
	if originalURL == "" {
		return c.String(http.StatusBadRequest, "url is required")
	}

	shorten, err := h.service.CreateShortURL(originalURL)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, TransformToShortenResponse(shorten))
}

func (h *Handler) HandleGet(c echo.Context) error {
	shortCode := c.Param("shortCode")
	shorten, err := h.service.GetShortURL(shortCode)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	if err := h.service.IncrementAccessCount(shortCode, shorten.AccessCount); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, TransformToShortenResponse(shorten))
}

func (h *Handler) HandlePut(c echo.Context) error {
	shortCode := c.Param("shortCode")
	newURL := c.FormValue("url")
	if newURL == "" {
		return c.String(http.StatusBadRequest, "url is required")
	}

	updatedShorten, err := h.service.UpdateShortURL(shortCode, newURL)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, TransformToShortenResponse(updatedShorten))
}

func (h *Handler) HandleDelete(c echo.Context) error {
	shortCode := c.Param("shortCode")
	if err := h.service.DeleteShortURL(shortCode); err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) HandleStats(c echo.Context) error {
	shortCode := c.Param("shortCode")

	shorten, err := h.service.GetShortURL(shortCode)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, shorten)
}
