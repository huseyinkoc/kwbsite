package models

// TokenResponse represents the response for successful authentication
type TokenResponse struct {
	Token     string `json:"token" example:"eyJhbGciOiJ..."`
	ExpiresIn int    `json:"expires_in" example:"900"`
	Message   string `json:"message" example:"Login successful"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid credentials"`
}

// LockedResponse represents account locked response
type LockedResponse struct {
	Error       string `json:"error" example:"Account locked"`
	LockedUntil string `json:"locked_until" example:"2023-12-31T23:59:59Z"`
}
