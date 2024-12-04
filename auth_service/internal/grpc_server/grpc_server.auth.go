package grpc_server

import (
	"auth_service/internal/models"
	"auth_service/internal/service"
	"auth_service/pkg/redis"
	protoAuth "auth_service/proto"
	"context"
	"fmt"
)

type authServer struct {
	authService service.AuthService
	redisCache  redis.RedisCache
	protoAuth.UnimplementedAuthServiceServer
}

func NewAuthServer(authService service.AuthService, redisCache redis.RedisCache) protoAuth.AuthServiceServer {
	return &authServer{
		authService: authService,
		redisCache:  redisCache,
	}
}

func (s *authServer) Register(ctx context.Context, req *protoAuth.RegisterRequest) (*protoAuth.RegisterResponse, error) {
	result, err := s.authService.Register(ctx, models.RegisterRequest{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	user := result.User
	return &protoAuth.RegisterResponse{
		User: &protoAuth.User{
			Id:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			Role:      user.Role,
			Verified:  user.Verified,
			CreatedAt: user.CreatedAt.Unix(),
			UpdatedAt: user.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *authServer) SendOTP(ctx context.Context, req *protoAuth.SendOTPRequest) (*protoAuth.SendOTPResponse, error) {
	otpCode, err := s.authService.SendOTP(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to send OTP email: %v", err)
	}

	otpKey := fmt.Sprintf("user_otp:%s", req.Email)
	err = s.redisCache.Set(otpKey, otpCode)
	if err != nil {
		return nil, err
	}

	return &protoAuth.SendOTPResponse{
		Message: "OTP sent successfully",
	}, nil
}

func (s *authServer) VerifyEmail(ctx context.Context, req *protoAuth.VerifyEmailRequest) (*protoAuth.VerifyEmailResponse, error) {
	verifyEmailRequest := models.VerifyEmailRequest{
		Email: req.Email,
		OTP:   req.Otp,
	}

	fmt.Println("verify req", verifyEmailRequest)

	otpKey := fmt.Sprintf("user_otp:%s", verifyEmailRequest.Email)
	redisOtp, err := s.redisCache.Get(otpKey)
	if err != nil {
		return nil, err
	}

	fmt.Println("otp key", otpKey)

	result, err := s.authService.VerifyEmail(ctx, verifyEmailRequest, redisOtp)
	if err != nil {
		return nil, err
	}

	err = s.redisCache.Del(req.Email)
	if err != nil {
		return nil, err
	}

	return &protoAuth.VerifyEmailResponse{
		Message: result.Message,
	}, nil
}

func (s *authServer) Login(ctx context.Context, req *protoAuth.LoginRequest) (*protoAuth.LoginResponse, error) {
	result, err := s.authService.Login(ctx, models.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &protoAuth.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		Message:      result.Message,
	}, nil
}

func (s *authServer) ValidateToken(ctx context.Context, req *protoAuth.ValidateTokenRequest) (*protoAuth.ValidateTokenResponse, error) {
	result, err := s.authService.ValidateToken(ctx, models.ValidateTokenRequest{Token: req.Token})
	if err != nil {
		return nil, err
	}

	return &protoAuth.ValidateTokenResponse{
		Valid:  result.Valid,
		UserId: result.UserID,
	}, nil
}
