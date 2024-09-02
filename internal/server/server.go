package server

import (
	"Teller/internal/config"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Server struct {
	Config      *config.Config
	Subscribers map[string][]chan json.RawMessage
	SubsMutex   sync.RWMutex
}

func New(cfg *config.Config) *Server {
	return &Server{
		Config:      cfg,
		Subscribers: make(map[string][]chan json.RawMessage),
	}
}

func (s *Server) Start() {
	http.HandleFunc("/publish", s.handlePublish)
	http.HandleFunc("/subscribe", s.handleSubscribe)
	server := &http.Server{
		Addr:    ":" + s.Config.Port,
		Handler: http.DefaultServeMux,
	}

	fmt.Printf("Starting server on port %s\n", s.Config.Port)
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Server error:", err)
		os.Exit(1) // Exit the application with a non-zero status code on error
	}
}

func (s *Server) handlePublish(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	tokenString := s.getTokenFromHeader(w, r)
	if tokenString == "" {
		return
	}

	token, err := s.validateToken(tokenString)
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var payload struct {
		Channel string          `json:"channel"`
		Message json.RawMessage `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if payload.Channel == "" || !json.Valid(payload.Message) {
		http.Error(w, "Invalid channel or message", http.StatusBadRequest)
		return
	}

	s.SubsMutex.RLock()
	defer s.SubsMutex.RUnlock()

	if chans, found := s.Subscribers[payload.Channel]; found {
		for _, ch := range chans {
			ch <- payload.Message
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message received successfully"))
}

func (s *Server) handleSubscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	tokenString := s.getTokenFromHeader(w, r)
	if tokenString == "" {
		return
	}

	token, err := s.validateToken(tokenString)
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	channel := r.URL.Query().Get("channel")
	if channel == "" {
		http.Error(w, "Missing channel parameter", http.StatusBadRequest)
		return
	}

	if r.URL.Query().Get("test") != "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // First, set the response status

		fmt.Fprintf(w, `{"test": true}`)

		flusher, ok := w.(http.Flusher)
		if ok {
			flusher.Flush() // Send the data to the client
		}

		time.Sleep(3 * time.Second)
		return
	}

	messageChan := make(chan json.RawMessage)
	s.SubsMutex.Lock()
	s.Subscribers[channel] = append(s.Subscribers[channel], messageChan)
	s.SubsMutex.Unlock()

	defer func() {
		s.SubsMutex.Lock()
		subscribers := s.Subscribers[channel]
		for i, ch := range subscribers {
			if ch == messageChan {
				s.Subscribers[channel] = append(subscribers[:i], subscribers[i+1:]...)
				break
			}
		}
		s.SubsMutex.Unlock()
		close(messageChan)
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	for msg := range messageChan {
		fmt.Fprintf(w, "data: %s\n\n", msg)
		flusher.Flush()
	}
}

func (s *Server) getTokenFromHeader(w http.ResponseWriter, r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return ""
	}
	return authHeader[7:]
}

func (s *Server) validateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.Config.JwtSecret), nil
	})
}
