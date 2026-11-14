package api

import (
	"net/http"
	"strings"
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
