package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/akashkncse/xotp/internal/otp"
)

type otpStore interface {
	Save(ctx context.Context, userID, codeHash string) error
	Find(ctx context.Context, userID string) (string, bool, error)
	Delete(ctx context.Context, userID string) (bool, error)
	SetDefaultTTL(ttl time.Duration)
}

type Server struct {
	store otpStore
}

func NewServer() *Server {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	return NewServerWithStore(otp.NewStore(addr))
}

func NewServerWithStore(store otpStore) *Server {
	return &Server{store: store}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.aliveHandler)
	mux.HandleFunc("/generateOtp", s.generateOTPHandler)
	mux.HandleFunc("/verifyOtp", s.verifyOTPHandler)
	mux.HandleFunc("/expiry", s.expiryHandler)
	mux.HandleFunc("/setTtl", s.setTTLHandler)
	return mux
}

func decodeRequest(r *http.Request, target any) error {
	return json.NewDecoder(r.Body).Decode(target)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
