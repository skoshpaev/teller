package server

import (
	"Teller/internal/config"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestPublishWithoutJWT(t *testing.T) {
	cfg := &config.Config{
		JwtSecret: "H1g$eCr3t!2S#cUr3T@256-bSecr3tIt",
		Port:      "8080",
	}
	s := New(cfg)

	payload := map[string]interface{}{
		"channel": "test-channel",
		"message": map[string]interface{}{
			"key": "value",
		},
	}

	jsonPayload, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/publish", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handlePublish)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Expected status Unauthorized, got %v", status)
	}
}

func TestPublishWithInvalidJWT(t *testing.T) {
	cfg := &config.Config{
		JwtSecret: "H1g$eCr3t!2S#cUr3T@256-bSecr3tIt",
		Port:      "8080",
	}
	s := New(cfg)

	payload := map[string]interface{}{
		"channel": "test-channel",
		"message": map[string]interface{}{
			"key": "value",
		},
	}

	jsonPayload, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/publish", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer invalid-token")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handlePublish)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Expected status Unauthorized, got %v", status)
	}
}

func TestPublishWithValidJWT(t *testing.T) {
	cfg := &config.Config{
		JwtSecret: "H1g$eCr3t!2S#cUr3T@256-bSecr3tIt",
		Port:      "8080",
	}
	s := New(cfg)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "1234567890",
		"name": "John Doe",
		"iat":  time.Now().Unix(),
	})

	validToken, err := token.SignedString([]byte(cfg.JwtSecret))
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	payload := map[string]interface{}{
		"channel": "test-channel",
		"message": map[string]interface{}{
			"key": "value",
		},
	}

	jsonPayload, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/publish", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+validToken)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handlePublish)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status OK, got %v", status)
	}
}

func TestPublishWithFormEncodedJWT(t *testing.T) {
	cfg := &config.Config{
		JwtSecret: "H1g$eCr3t!2S#cUr3T@256-bSecr3tIt",
		Port:      "8080",
	}
	s := New(cfg)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "1234567890",
		"name": "John Doe",
		"iat":  time.Now().Unix(),
	})

	validToken, err := token.SignedString([]byte(cfg.JwtSecret))
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	formData := "topic=test-channel&data=Your+message+content"
	req, err := http.NewRequest("POST", "/publish", strings.NewReader(formData))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+validToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handlePublish)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status OK, got %v", status)
	}
}

func TestSubscribeWithInvalidJWT(t *testing.T) {
	cfg := &config.Config{
		JwtSecret: "H1g$eCr3t!2S#cUr3T@256-bSecr3tIt",
		Port:      "8080",
	}
	s := New(cfg)

	req, err := http.NewRequest("GET", "/subscribe?channel=test-channel", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer invalid-token")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handleSubscribe)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Expected status Unauthorized, got %v", status)
	}
}
