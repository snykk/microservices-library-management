package service

import (
	"auth_service/internal/models"
	"auth_service/internal/repository"
	"auth_service/pkg/jwt"
	"auth_service/pkg/mailer"
	"auth_service/pkg/utils"
	"context"
	"time"
)

type AuthService interface {
	Register(ctx context.Context, req models.RegisterRequest) (models.RegisterResponse, error)
	SendOTP(ctx context.Context, email string) (string, error)
	VerifyEmail(ctx context.Context, req models.VerifyEmailRequest, redisOtp string) (models.VerifyEmailResponse, error)
	Login(ctx context.Context, req models.LoginRequest) (models.LoginResponse, error)
	ValidateToken(ctx context.Context, req models.ValidateTokenRequest) (models.ValidateTokenResponse, error)
}

type authService struct {
	repo       repository.AuthRepository
	jwtService jwt.JWTService
	mailer     mailer.OTPMailer
}

func NewAuthService(repo repository.AuthRepository, jwtService jwt.JWTService, mailer mailer.OTPMailer) AuthService {
	return &authService{
		repo:       repo,
		jwtService: jwtService,
		mailer:     mailer,
	}
}

func (s *authService) Register(ctx context.Context, req models.RegisterRequest) (models.RegisterResponse, error) {
	_, err := s.repo.GetUserByEmail(req.Email)
	if err == nil {
		return models.RegisterResponse{}, ErrEmailAlreadyRegistered
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return models.RegisterResponse{}, ErrFailedHashPassword
	}

	user := models.UserRecord{
		Email:     req.Email,
		Username:  req.Username,
		Password:  hashedPassword,
		Verified:  false,
		Role:      "user", // Default role
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// Create user
	createdUser, err := s.repo.CreateUser(user)
	if err != nil {
		return models.RegisterResponse{}, ErrCreateUser
	}

	return models.RegisterResponse{
		User: createdUser,
	}, nil
}

func (s *authService) SendOTP(ctx context.Context, email string) (string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", ErrGetUserByEmail
	}

	if user.Verified {
		return "", ErrEmailAlreadyVerified
	}

	otp, err := utils.GenerateOTPCode(6)
	if err != nil {
		return "", ErrGenerateOTPCode
	}

	if err = s.mailer.SendOTP(otp, email); err != nil {
		return "", ErrSendOtpWithMailer
	}

	return otp, nil
}

func (s *authService) VerifyEmail(ctx context.Context, req models.VerifyEmailRequest, redisOtp string) (models.VerifyEmailResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return models.VerifyEmailResponse{}, ErrGetUserByEmail
	}

	if user.Verified {
		return models.VerifyEmailResponse{}, ErrEmailAlreadyVerified
	}

	if req.OTP != redisOtp {
		return models.VerifyEmailResponse{}, ErrMismatchOTPCode
	}

	err = s.repo.UpdateUserVerification(req.Email, true)
	if err != nil {
		return models.VerifyEmailResponse{}, ErrUpdateUserVerification
	}

	return models.VerifyEmailResponse{Message: "Email verified successfully"}, nil
}

func (s *authService) Login(ctx context.Context, req models.LoginRequest) (models.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return models.LoginResponse{}, ErrGetUserByEmail
	}

	if !user.Verified {
		return models.LoginResponse{}, ErrEmailNotVerified
	}

	// Check password
	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		return models.LoginResponse{}, ErrInvalidPassword
	}

	// Generate Access Token
	accessToken, err := s.jwtService.GenerateToken(user.ID, user.Role, user.Email)
	if err != nil {
		return models.LoginResponse{}, ErrGenerateAccessToken
	}

	// Generate Refresh Token
	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, user.Role, user.Email)
	if err != nil {
		return models.LoginResponse{}, ErrGenerateRefreshToken
	}

	return models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "Login successful",
	}, nil
}

func (s *authService) ValidateToken(ctx context.Context, req models.ValidateTokenRequest) (models.ValidateTokenResponse, error) {
	claims, err := s.jwtService.ParseToken(req.Token)
	if err != nil {
		return models.ValidateTokenResponse{Valid: false}, ErrPareseToken
	}

	return models.ValidateTokenResponse{
		Valid:  true,
		UserID: claims.UserID,
	}, nil
}
