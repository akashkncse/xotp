package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/akashkncse/xotp/internal/otp"
)

func TestAliveHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	NewServer().aliveHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var body aliveResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !body.Alive {
		t.Fatal("expected alive to be true")
	}
}

func TestGenerateAndVerifyOTPHandlers(t *testing.T) {
	server := NewServer()

	generateRequest := httptest.NewRequest(http.MethodGet, "/generateOtp", bytes.NewBufferString(`{"userId":"user-123"}`))
	generateResponse := httptest.NewRecorder()
	server.generateOTPHandler(generateResponse, generateRequest)

	if generateResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, generateResponse.Code)
	}

	var generated generateOTPResponse
	if err := json.NewDecoder(generateResponse.Body).Decode(&generated); err != nil {
		t.Fatalf("decode generate response: %v", err)
	}

	verifyRequest := httptest.NewRequest(
		http.MethodPost,
		"/verifyOtp",
		bytes.NewBufferString(`{"userId":"user-123","otp":`+strconv.Itoa(generated.OTP)+`}`),
	)
	verifyResponse := httptest.NewRecorder()
	server.verifyOTPHandler(verifyResponse, verifyRequest)

	var verified verifyOTPResponse
	if err := json.NewDecoder(verifyResponse.Body).Decode(&verified); err != nil {
		t.Fatalf("decode verify response: %v", err)
	}
	if !verified.Match {
		t.Fatal("expected generated OTP to match")
	}
}

func TestExpiryHandlerDeletesOTP(t *testing.T) {
	server := NewServer()
	server.store.Save("user-123", otp.Hash("123456"))

	request := httptest.NewRequest(http.MethodDelete, "/expiry", bytes.NewBufferString(`{"userId":"user-123"}`))
	response := httptest.NewRecorder()
	server.expiryHandler(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, response.Code)
	}
	if _, exists := server.store.Find("user-123"); exists {
		t.Fatal("expected OTP to be deleted")
	}
}
