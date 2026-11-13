package api

import (
	"encoding/json"
	"net/http"

	"github.com/akashkncse/xotp/internal/otp"
)

type Server struct {
	store *otp.Store
}

func NewServer() *Server {
	return &Server{store: otp.NewStore()}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.aliveHandler)
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
