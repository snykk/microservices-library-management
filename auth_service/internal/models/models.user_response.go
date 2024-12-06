package models

type RegisterResponse struct {
	User UserRecord
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
	Message      string
}

type VerifyEmailResponse struct {
	Message string
}

type ValidateTokenResponse struct {
	Valid  bool
	UserID string
	Role   string
	Email  string
}
