package service

import (
	"auth_service/internal/models"
	"auth_service/internal/repository"
	"auth_service/pkg/jwt"
	"auth_service/pkg/mailer"
	"auth_service/pkg/utils"
	"context"
	"errors"
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
		return models.RegisterResponse{}, errors.New("email already registered")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return models.RegisterResponse{}, errors.New("failed to hash password")
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
		return models.RegisterResponse{}, err
	}

	return models.RegisterResponse{
		User: createdUser,
	}, nil
}

func (s *authService) SendOTP(ctx context.Context, email string) (string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	if user.Verified {
		return "", errors.New("email already verified")
	}

	otp, err := utils.GenerateOTPCode(6)
	if err != nil {
		return "", err
	}

	if err = s.mailer.SendOTP(otp, email); err != nil {
		return "", err
	}

	return otp, nil
}

func (s *authService) VerifyEmail(ctx context.Context, req models.VerifyEmailRequest, redisOtp string) (models.VerifyEmailResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return models.VerifyEmailResponse{}, err
	}

	if user.Verified {
		return models.VerifyEmailResponse{}, errors.New("email already verified")
	}

	if req.OTP != redisOtp {
		return models.VerifyEmailResponse{Message: "Invalid OTP"}, nil
	}

	err = s.repo.UpdateUserVerification(req.Email, true)
	if err != nil {
		return models.VerifyEmailResponse{}, err
	}

	return models.VerifyEmailResponse{Message: "Email verified successfully"}, nil
}

func (s *authService) Login(ctx context.Context, req models.LoginRequest) (models.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return models.LoginResponse{}, errors.New("user not found")
	}

	if !user.Verified {
		return models.LoginResponse{}, errors.New("user email not verified")
	}

	// Check password
	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		return models.LoginResponse{}, errors.New("invalid password")
	}

	// Generate Access Token
	accessToken, err := s.jwtService.GenerateToken(user.ID, user.Role, user.Email)
	if err != nil {
		return models.LoginResponse{}, errors.New("error generating access token")
	}

	// Generate Refresh Token
	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, user.Role, user.Email)
	if err != nil {
		return models.LoginResponse{}, errors.New("error generating refresh token")
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
		return models.ValidateTokenResponse{Valid: false}, err
	}

	return models.ValidateTokenResponse{
		Valid:  true,
		UserID: claims.UserID,
	}, nil
}
