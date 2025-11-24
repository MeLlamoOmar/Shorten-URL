package main

import (
    "encoding/json"
    "errors"
    "net/http"
    "net/http/httptest"
    "net/url"
    "strings"
    "testing"

    "github.com/labstack/echo/v4"
)

type MockServiceForHandlers struct {
    CreateResp *ShortenStatResponse
    GetResp    *ShortenStatResponse
    UpdateResp *ShortenStatResponse
    Err        error
    IncrementCalled bool
}

func (m *MockServiceForHandlers) CreateShortURL(originalURL string) (*ShortenStatResponse, error) {
    return m.CreateResp, m.Err
}

func (m *MockServiceForHandlers) GetShortURL(shortCode string) (*ShortenStatResponse, error) {
    return m.GetResp, m.Err
}

func (m *MockServiceForHandlers) UpdateShortURL(shortCode string, newURL string) (*ShortenStatResponse, error) {
    return m.UpdateResp, m.Err
}

func (m *MockServiceForHandlers) DeleteShortURL(shortCode string) error {
    return m.Err
}

func (m *MockServiceForHandlers) IncrementAccessCount(shortCode string, currentCount int) error {
    if m.Err != nil {
        return m.Err
    }
    m.IncrementCalled = true
    return nil
}

func TestHandlePostSuccess(t *testing.T) {
    svc := &MockServiceForHandlers{
        CreateResp: &ShortenStatResponse{OriginalURL: "http://example.com", ShortCode: "abc123"},
    }

    e := echo.New()
    h := NewHandler(svc)
    reqBody := url.Values{}
    reqBody.Set("url", "http://example.com")

    req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(reqBody.Encode()))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    if err := h.HandlePost(c); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if rec.Code != http.StatusCreated {
        t.Fatalf("expected status 201, got %d", rec.Code)
    }

    var resp ShortenResponse
    if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
        t.Fatalf("invalid json response: %v", err)
    }

    if resp.ShortCode != "abc123" {
        t.Fatalf("expected shortCode abc123, got %s", resp.ShortCode)
    }
}

func TestHandlePostBadRequest(t *testing.T) {
    svc := &MockServiceForHandlers{}
    e := echo.New()
    h := NewHandler(svc)

    req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(""))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    if err := h.HandlePost(c); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if rec.Code != http.StatusBadRequest {
        t.Fatalf("expected status 400, got %d", rec.Code)
    }
}

func TestHandleGetSuccess(t *testing.T) {
    svc := &MockServiceForHandlers{
        GetResp: &ShortenStatResponse{OriginalURL: "http://site", ShortCode: "c1", AccessCount: 2},
    }
    e := echo.New()
    h := NewHandler(svc)

    req := httptest.NewRequest(http.MethodGet, "/shorten/c1", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("shortCode")
    c.SetParamValues("c1")

    if err := h.HandleGet(c); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if rec.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", rec.Code)
    }

    var resp ShortenResponse
    if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
        t.Fatalf("invalid json response: %v", err)
    }

    if resp.ShortCode != "c1" {
        t.Fatalf("expected shortCode c1, got %s", resp.ShortCode)
    }

    if !svc.IncrementCalled {
        t.Fatalf("expected IncrementAccessCount to be called")
    }
}

func TestHandleGetNotFound(t *testing.T) {
    svc := &MockServiceForHandlers{Err: errors.New("not found")}
    e := echo.New()
    h := NewHandler(svc)

    req := httptest.NewRequest(http.MethodGet, "/shorten/x", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("shortCode")
    c.SetParamValues("x")

    if err := h.HandleGet(c); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if rec.Code != http.StatusNotFound {
        t.Fatalf("expected status 404, got %d", rec.Code)
    }
}

func TestHandlePutSuccess(t *testing.T) {
    svc := &MockServiceForHandlers{
        UpdateResp: &ShortenStatResponse{OriginalURL: "new", ShortCode: "u1"},
    }
    e := echo.New()
    h := NewHandler(svc)

    form := url.Values{}
    form.Set("url", "new")
    req := httptest.NewRequest(http.MethodPut, "/shorten/u1", strings.NewReader(form.Encode()))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("shortCode")
    c.SetParamValues("u1")

    if err := h.HandlePut(c); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if rec.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", rec.Code)
    }
}

func TestHandlePutBadRequest(t *testing.T) {
    svc := &MockServiceForHandlers{}
    e := echo.New()
    h := NewHandler(svc)

    req := httptest.NewRequest(http.MethodPut, "/shorten/u1", strings.NewReader(""))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("shortCode")
    c.SetParamValues("u1")

    if err := h.HandlePut(c); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if rec.Code != http.StatusBadRequest {
        t.Fatalf("expected status 400, got %d", rec.Code)
    }
}

func TestHandleDeleteSuccess(t *testing.T) {
    svc := &MockServiceForHandlers{}
    e := echo.New()
    h := NewHandler(svc)

    req := httptest.NewRequest(http.MethodDelete, "/shorten/d1", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("shortCode")
    c.SetParamValues("d1")

    if err := h.HandleDelete(c); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if rec.Code != http.StatusNoContent {
        t.Fatalf("expected status 204, got %d", rec.Code)
    }
}

func TestHandleDeleteNotFound(t *testing.T) {
    svc := &MockServiceForHandlers{Err: errors.New("not found")}
    e := echo.New()
    h := NewHandler(svc)

    req := httptest.NewRequest(http.MethodDelete, "/shorten/d1", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("shortCode")
    c.SetParamValues("d1")

    if err := h.HandleDelete(c); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if rec.Code != http.StatusNotFound {
        t.Fatalf("expected status 404, got %d", rec.Code)
    }
}

func TestHandleStatsSuccess(t *testing.T) {
    svc := &MockServiceForHandlers{GetResp: &ShortenStatResponse{OriginalURL: "s", ShortCode: "st1", AccessCount: 7}}
    e := echo.New()
    h := NewHandler(svc)

    req := httptest.NewRequest(http.MethodGet, "/shorten/st1/stats", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("shortCode")
    c.SetParamValues("st1")

    if err := h.HandleStats(c); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if rec.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", rec.Code)
    }

    var resp ShortenStatResponse
    if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
        t.Fatalf("invalid json response: %v", err)
    }

    if resp.AccessCount != 7 {
        t.Fatalf("expected accessCount 7, got %d", resp.AccessCount)
    }
}
