package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandle_ValidInvoice(t *testing.T) {
	body := `{"amount": 100, "currency": "USD", "due_date": "2025-08-01"}`

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	Handle(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", rr.Code)
	}

	var resp map[string]bool
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	if valid, ok := resp["valid"]; !ok || !valid {
		t.Errorf("expected valid: true, got: %v", resp)
	}
}

func TestHandle_InvalidInvoice(t *testing.T) {
	body := `{"amount": -5, "curreny": "usd", "due_date":"Aug 1"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	Handle(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422 Unprocessable Entity, got %d", rr.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	if valid, ok := resp["valid"]; ok && valid.(bool) {
		t.Errorf("expected valid: false, got: true")
	}
}
