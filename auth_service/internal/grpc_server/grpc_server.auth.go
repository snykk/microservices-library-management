package grpc_server

import (
	"auth_service/internal/constants"
	"auth_service/internal/exception"
	"auth_service/internal/models"
	"auth_service/internal/service"
	"auth_service/pkg/logger"
	"auth_service/pkg/redis"
	"auth_service/pkg/utils"
	protoAuth "auth_service/proto/auth_service"
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authServer struct {
	authService service.AuthService
	redisCache  redis.RedisCache
	logger      *logger.Logger
	protoAuth.UnimplementedAuthServiceServer
}

func NewAuthServer(authService service.AuthService, redisCache redis.RedisCache, logger *logger.Logger) protoAuth.AuthServiceServer {
	return &authServer{
		authService: authService,
		redisCache:  redisCache,
		logger:      logger,
	}
}

func (s *authServer) Register(ctx context.Context, req *protoAuth.RegisterRequest) (*protoAuth.RegisterResponse, error) {
	requestID := utils.GetRequestIDFromMetadataContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received Register request", map[string]interface{}{"email": req.Email, "username": req.Username}, nil)

	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid Register request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	result, err := s.authService.Register(ctx, &models.RegisterRequest{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Register service failed", nil, err)
		return nil, exception.GRPCErrorFormatter(err)
	}
	user := result.User

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Register service succeeded", map[string]interface{}{"created_user": user}, nil)
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
	requestID := utils.GetRequestIDFromMetadataContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received SendOTP request", map[string]interface{}{"email": req.Email}, nil)

	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid SendOTP request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	otpCode, err := s.authService.SendOTP(context.WithValue(ctx, constants.ContextRequestIDKey, requestID), req.Email)
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to send OTP", nil, err)
		return nil, exception.GRPCErrorFormatter(err)
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "OTP generated successfully", nil, nil)
	otpKey := fmt.Sprintf("user_otp:%s", req.Email)
	err = s.redisCache.Set(otpKey, otpCode)
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to cache OTP", nil, err)
		return nil, exception.GRPCErrorFormatter(err)
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "OTP cached successfully", nil, nil)
	return &protoAuth.SendOTPResponse{
		Message: "OTP sent successfully",
	}, nil
}

func (s *authServer) VerifyEmail(ctx context.Context, req *protoAuth.VerifyEmailRequest) (*protoAuth.VerifyEmailResponse, error) {
	requestID := utils.GetRequestIDFromMetadataContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received VerifyEmail request", map[string]interface{}{"email": req.Email}, nil)

	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid VerifyEmail request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	verifyEmailRequest := models.VerifyEmailRequest{
		Email: req.Email,
		OTP:   req.Otp,
	}

	otpKey := fmt.Sprintf("user_otp:%s", verifyEmailRequest.Email)
	redisOtp, err := s.redisCache.Get(otpKey)
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to get OTP from cache", nil, err)
		return nil, exception.GRPCErrorFormatter(err)
	}

	result, err := s.authService.VerifyEmail(ctx, &verifyEmailRequest, redisOtp)
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "VerifyEmail service failed", nil, err)
		return nil, exception.GRPCErrorFormatter(err)
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Email verification succeeded", nil, nil)
	err = s.redisCache.Del(otpKey)
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to delete OTP from cache", nil, err)
		return nil, exception.GRPCErrorFormatter(err)
	}

	return &protoAuth.VerifyEmailResponse{
		Message: result.Message,
	}, nil
}

func (s *authServer) Login(ctx context.Context, req *protoAuth.LoginRequest) (*protoAuth.LoginResponse, error) {
	requestID := utils.GetRequestIDFromMetadataContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received Login request", map[string]interface{}{"email": req.Email}, nil)

	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid Login request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	result, err := s.authService.Login(ctx, &models.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Login service failed", nil, err)
		return nil, exception.GRPCErrorFormatter(err)
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Login success, access token & refresh token generated", nil, nil)
	return &protoAuth.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		Message:      result.Message,
	}, nil
}

func (s *authServer) ValidateToken(ctx context.Context, req *protoAuth.ValidateTokenRequest) (*protoAuth.ValidateTokenResponse, error) {
	requestID := utils.GetRequestIDFromMetadataContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received ValidateToken request", map[string]interface{}{"token": req.Token}, nil)

	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid ValidateToken request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	result, err := s.authService.ValidateToken(ctx, &models.ValidateTokenRequest{Token: req.Token})
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "ValidateToken service failed", nil, err)
		return nil, exception.GRPCErrorFormatter(err)
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Token validation succeeded", nil, nil)
	return &protoAuth.ValidateTokenResponse{
		Valid:  result.Valid,
		UserId: result.UserID,
		Role:   result.Role,
		Email:  result.Email,
	}, nil
}
