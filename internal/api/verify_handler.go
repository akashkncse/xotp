package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/akashkncse/xotp/internal/otp"
)

func (s *Server) verifyOTPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request verifyOTPRequest
	if err := decodeRequest(r, &request); err != nil || strings.TrimSpace(request.UserID) == "" {
		http.Error(w, "userId and otp are required", http.StatusBadRequest)
		return
	}

	match := false
	if request.OTP >= 0 && request.OTP <= 999999 {
		storedHash, exists, err := s.store.Find(r.Context(), request.UserID)
		if err != nil {
			http.Error(w, "failed to retrieve OTP", http.StatusInternalServerError)
			return
		}
		match = exists && storedHash == otp.Hash(fmt.Sprintf("%06d", request.OTP))
	}

	writeJSON(w, http.StatusOK, verifyOTPResponse{Match: match})
}
