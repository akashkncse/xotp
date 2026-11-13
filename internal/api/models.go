package api

type aliveResponse struct {
	Alive bool `json:"alive"`
}

type generateOTPRequest struct {
	UserID string `json:"userId"`
}

type generateOTPResponse struct {
	OTP int `json:"otp"`
}

type verifyOTPRequest struct {
	UserID string `json:"userId"`
	OTP    int    `json:"otp"`
}

type verifyOTPResponse struct {
	Match bool `json:"match"`
}
