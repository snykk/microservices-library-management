package models

type RegisterRequest struct {
	Email    string
	Username string
	Password string
}

type LoginRequest struct {
	Email    string
	Password string
}

type VerifyEmailRequest struct {
	Email string
	OTP   string
}

type ValidateTokenRequest struct {
	Token string
}
