package service

import (
	"auth_service/internal/models"
	"auth_service/internal/repository"
	"auth_service/pkg/jwt"
	"auth_service/pkg/mailer"
	"auth_service/pkg/rabbitmq"
	"auth_service/pkg/utils"
	"context"
	"encoding/json"
)

type AuthService interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.RegisterResponse, error)
	SendOTP(ctx context.Context, email string) (*string, error)
	VerifyEmail(ctx context.Context, req *models.VerifyEmailRequest, redisOtp string) (*models.VerifyEmailResponse, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error)
	ValidateToken(ctx context.Context, req *models.ValidateTokenRequest) (*models.ValidateTokenResponse, error)
}

type authService struct {
	repo       repository.AuthRepository
	jwtService jwt.JWTService
	mailer     mailer.OTPMailer
	publisher  *rabbitmq.Publisher
}

func NewAuthService(repo repository.AuthRepository, jwtService jwt.JWTService, mailer mailer.OTPMailer, publisher *rabbitmq.Publisher) AuthService {
	return &authService{
		repo:       repo,
		jwtService: jwtService,
		mailer:     mailer,
		publisher:  publisher,
	}
}

func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) (*models.RegisterResponse, error) {
	userFromDB, _ := s.repo.GetUserByEmail(ctx, req.Email)
	if userFromDB != nil {
		return nil, ErrEmailAlreadyRegistered
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, ErrFailedHashPassword
	}

	user := models.UserRecord{
		Email:    req.Email,
		Username: req.Username,
		Password: hashedPassword,
		Verified: false,
		Role:     "user", // Default role
	}
	// Create user
	createdUser, err := s.repo.CreateUser(ctx, &user)
	if err != nil {
		return nil, ErrCreateUser
	}

	return &models.RegisterResponse{
		User: *createdUser,
	}, nil
}

func (s *authService) SendOTP(ctx context.Context, email string) (*string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, ErrGetUserByEmail
	}

	if user.Verified {
		return nil, ErrEmailAlreadyVerified
	}

	otp, err := utils.GenerateOTPCode(6)
	if err != nil {
		return nil, ErrGenerateOTPCode
	}

	// if err = s.mailer.SendOTP(otp, email); err != nil { // todo: use rabbitmq to enhance response time
	// 	return nil, ErrSendOtpWithMailer
	// }

	// Prepare message
	message := map[string]string{
		"email": email,
		"otp":   otp,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return nil, ErrMarshalOTPMessage
	}

	// Publish to RabbitMQ
	err = s.publisher.Publish("email_exchange", "otp_code", messageBytes)
	if err != nil {
		return nil, ErrPublishToQueue
	}

	return &otp, nil
}

func (s *authService) VerifyEmail(ctx context.Context, req *models.VerifyEmailRequest, redisOtp string) (*models.VerifyEmailResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrGetUserByEmail
	}

	if user.Verified {
		return nil, ErrEmailAlreadyVerified
	}

	if req.OTP != redisOtp {
		return nil, ErrMismatchOTPCode
	}

	verified := true
	err = s.repo.UpdateUserVerification(ctx, req.Email, verified)
	if err != nil {
		return nil, ErrUpdateUserVerification
	}

	return &models.VerifyEmailResponse{Message: "Email verified successfully"}, nil
}

func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrGetUserByEmail
	}

	if !user.Verified {
		return nil, ErrEmailNotVerified
	}

	// Check password
	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		return nil, ErrInvalidPassword
	}

	// Generate Access Token
	accessToken, err := s.jwtService.GenerateToken(user.ID, user.Role, user.Email)
	if err != nil {
		return nil, ErrGenerateAccessToken
	}

	// Generate Refresh Token
	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, user.Role, user.Email)
	if err != nil {
		return nil, ErrGenerateRefreshToken
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "Login successful",
	}, nil
}

func (s *authService) ValidateToken(ctx context.Context, req *models.ValidateTokenRequest) (*models.ValidateTokenResponse, error) {
	claims, err := s.jwtService.ParseToken(req.Token)
	if err != nil {
		return nil, ErrPareseToken
	}

	return &models.ValidateTokenResponse{
		Valid:  true,
		UserID: claims.UserID,
		Role:   claims.Role,
		Email:  claims.Email,
	}, nil
}
