package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func (s *Server) expiryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request expiryRequest
	if err := decodeRequest(r, &request); err != nil || strings.TrimSpace(request.UserID) == "" {
		http.Error(w, "userId is required", http.StatusBadRequest)
		return
	}

	if _, err := s.store.Delete(r.Context(), request.UserID); err != nil {
		http.Error(w, "failed to delete OTP", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) setTTLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		TTL int64 `json:"ttl"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.TTL <= 0 {
		http.Error(w, "ttl must be a positive number of seconds", http.StatusBadRequest)
		return
	}

	s.store.SetDefaultTTL(time.Duration(request.TTL) * time.Second)
	w.WriteHeader(http.StatusNoContent)
}
