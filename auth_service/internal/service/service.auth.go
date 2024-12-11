package service

import (
	"auth_service/internal/constants"
	"auth_service/internal/models"
	"auth_service/internal/repository"
	"auth_service/pkg/jwt"
	"auth_service/pkg/mailer"
	"auth_service/pkg/rabbitmq"
	"auth_service/pkg/utils"
	"context"
	"log"
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
		log.Printf("[%s] User with email %s already exists\n", utils.GetLocation(), req.Email)
		return nil, ErrEmailAlreadyRegistered
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("[%s] Failed to hash password: %v\n", utils.GetLocation(), err)
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
		log.Printf("[%s] Failed to create user: %v\n", utils.GetLocation(), err)
		return nil, ErrCreateUser
	}

	log.Printf("[%s] User %s created successfully\n", utils.GetLocation(), req.Email)
	return &models.RegisterResponse{
		User: *createdUser,
	}, nil
}

func (s *authService) SendOTP(ctx context.Context, email string) (*string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("[%s] Failed when get email %s: %v\n", utils.GetLocation(), email, err)
		return nil, ErrGetUserByEmail
	}

	if user.Verified {
		log.Printf("[%s] Email already verified\n", utils.GetLocation())
		return nil, ErrEmailAlreadyVerified
	}

	otp, err := utils.GenerateOTPCode(6)
	if err != nil {
		log.Printf("[%s] Failed when generate OTP code: %v\n", utils.GetLocation(), err)
		return nil, ErrGenerateOTPCode
	}

	// Publish to RabbitMQ
	err = s.publisher.Publish(constants.EmailExchange, constants.OTPQueue, map[string]string{
		"email": email,
		"otp":   otp,
	})
	if err != nil {
		log.Printf("[%s] Failed to publish to RabbitMQ: %v\n", utils.GetLocation(), err)
		return nil, ErrPublishToQueue
	}

	return &otp, nil
}

func (s *authService) VerifyEmail(ctx context.Context, req *models.VerifyEmailRequest, redisOtp string) (*models.VerifyEmailResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("[%s] Failed to get user by email %s: %v\n", utils.GetLocation(), req.Email, err)
		return nil, ErrGetUserByEmail
	}

	if user.Verified {
		log.Printf("[%s] Email %s is already verified\n", utils.GetLocation(), req.Email)
		return nil, ErrEmailAlreadyVerified
	}

	if req.OTP != redisOtp {
		log.Printf("[%s] OTP mismatch for email %s\n", utils.GetLocation(), req.Email)
		return nil, ErrMismatchOTPCode
	}

	verified := true
	err = s.repo.UpdateUserVerification(ctx, req.Email, verified)
	if err != nil {
		log.Printf("[%s] Failed to update user verification for email %s: %v\n", utils.GetLocation(), req.Email, err)
		return nil, ErrUpdateUserVerification
	}

	log.Printf("[%s] Email %s successfully verified\n", utils.GetLocation(), req.Email)
	return &models.VerifyEmailResponse{Message: "Email verified successfully"}, nil
}

func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("[%s] Failed to get user by email %s: %v\n", utils.GetLocation(), req.Email, err)
		return nil, ErrGetUserByEmail
	}

	if !user.Verified {
		log.Printf("[%s] Email %s is not verified\n", utils.GetLocation(), req.Email)
		return nil, ErrEmailNotVerified
	}

	// Check password
	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		log.Printf("[%s] Invalid password for email %s\n", utils.GetLocation(), req.Email)
		return nil, ErrInvalidPassword
	}

	// Generate Access Token
	accessToken, err := s.jwtService.GenerateToken(user.ID, user.Role, user.Email)
	if err != nil {
		log.Printf("[%s] Failed to generate access token for email %s: %v\n", utils.GetLocation(), req.Email, err)
		return nil, ErrGenerateAccessToken
	}

	// Generate Refresh Token
	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, user.Role, user.Email)
	if err != nil {
		log.Printf("[%s] Failed to generate refresh token for email %s: %v\n", utils.GetLocation(), req.Email, err)
		return nil, ErrGenerateRefreshToken
	}

	log.Printf("[%s] Login successful for email %s\n", utils.GetLocation(), req.Email)
	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "Login successful",
	}, nil
}

func (s *authService) ValidateToken(ctx context.Context, req *models.ValidateTokenRequest) (*models.ValidateTokenResponse, error) {
	claims, err := s.jwtService.ParseToken(req.Token)
	if err != nil {
		log.Printf("[%s] Failed to parse token: %v\n", utils.GetLocation(), err)
		return nil, ErrPareseToken
	}

	log.Printf("[%s] Token validation successful for user ID %s\n", utils.GetLocation(), claims.UserID)
	return &models.ValidateTokenResponse{
		Valid:  true,
		UserID: claims.UserID,
		Role:   claims.Role,
		Email:  claims.Email,
	}, nil
}
