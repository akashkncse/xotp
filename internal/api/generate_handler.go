package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/akashkncse/xotp/internal/otp"
)

func (s *Server) generateOTPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request generateOTPRequest
	if err := decodeRequest(r, &request); err != nil || strings.TrimSpace(request.UserID) == "" {
		http.Error(w, "userId is required", http.StatusBadRequest)
		return
	}

	code, err := otp.Generate()
	if err != nil {
		http.Error(w, "failed to generate OTP", http.StatusInternalServerError)
		return
	}

	codeNumber, err := strconv.Atoi(code)
	if err != nil {
		http.Error(w, "failed to encode OTP", http.StatusInternalServerError)
		return
	}

	s.store.Save(request.UserID, otp.Hash(code))
	writeJSON(w, http.StatusOK, generateOTPResponse{OTP: codeNumber})
}
