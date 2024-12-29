package service

import "errors"

var (
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrFailedHashPassword     = errors.New("failed to hash password")
	ErrEmailAlreadyVerified   = errors.New("email already verified")
	ErrGetUserByEmail         = errors.New("error get user by email")
	ErrCreateUser             = errors.New("error create user")
	ErrGenerateOTPCode        = errors.New("error generate otp code")
	ErrSendOtpWithMailer      = errors.New("error send otp with mailer service")
	ErrMismatchOTPCode        = errors.New("mismatch otp code")
	ErrUpdateUserVerification = errors.New("error update user verification")
	ErrEmailNotVerified       = errors.New("user email not verified")
	ErrInvalidPassword        = errors.New("password is not valid")
	ErrGenerateAccessToken    = errors.New("error generating access token")
	ErrGenerateRefreshToken   = errors.New("error generating refresh token")
	ErrParseToken             = errors.New("error parse token")
	ErrMarshalOTPMessage      = errors.New("error marshal otp message")
	ErrPublishToQueue         = errors.New("error publish message to queue")
	ErrInvalidRefreshToken    = errors.New("invalid refresh token")
	ErrUpdateRefreshToken     = errors.New("error update refresh token")
	ErrUpdateLastLogin        = errors.New("error update last login")
	ErrLogoutFailed           = errors.New("error logout")
)
