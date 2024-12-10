package clients

import (
	"api_gateway/internal/constants"
	"api_gateway/internal/datatransfers"
	"api_gateway/pkg/logger"
	protoAuth "api_gateway/proto/auth_service"
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient interface {
	Register(ctx context.Context, dto datatransfers.RegisterRequest) (datatransfers.RegisterResponse, error)
	Login(ctx context.Context, dto datatransfers.LoginRequest) (datatransfers.LoginResponse, error)
	SendOtp(ctx context.Context, dto datatransfers.SendOtpRequest) (datatransfers.SendOtpResponse, error)
	VerifyEmail(ctx context.Context, dto datatransfers.VerifyEmailRequest) (datatransfers.VerifyEmailResponse, error)
	ValidateToken(ctx context.Context, dto datatransfers.ValidateTokenRequest) (datatransfers.ValidateTokenResponse, error)
}

type authClient struct {
	client protoAuth.AuthServiceClient
}

func NewAuthClient() (AuthClient, error) {
	conn, err := grpc.NewClient("auth-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Error("Failed to create AuthClient",
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryConnection),
		)

		return nil, err
	}
	client := protoAuth.NewAuthServiceClient(conn)

	logger.Log.Info("Success to create AuthClient",
		zap.String(constants.LoggerCategory, constants.LoggerCategoryConnection),
	)
	return &authClient{
		client: client,
	}, nil
}

func (authC *authClient) Register(ctx context.Context, dto datatransfers.RegisterRequest) (datatransfers.RegisterResponse, error) {
	reqProto := protoAuth.RegisterRequest{
		Email:    dto.Email,
		Username: dto.Username,
		Password: dto.Password,
	}

	logger.Log.Info("Sending Register request to Auth Service",
		zap.String("email", dto.Email),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := authC.client.Register(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("Register request failed",
			zap.String("email", dto.Email),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
		return datatransfers.RegisterResponse{}, err
	}

	logger.Log.Info("Register request succeeded",
		zap.String("email", resp.User.Email),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	response := datatransfers.RegisterResponse{
		Id:        resp.User.Id,
		Email:     resp.User.Email,
		Username:  resp.User.Username,
		Password:  resp.User.Password,
		Verified:  resp.User.Verified,
		CreatedAt: time.Unix(resp.User.CreatedAt, 0),
		UpdatedAt: time.Unix(resp.User.UpdatedAt, 0),
	}

	return response, nil
}

func (authC *authClient) SendOtp(ctx context.Context, dto datatransfers.SendOtpRequest) (datatransfers.SendOtpResponse, error) {
	reqProto := protoAuth.SendOTPRequest{
		Email: dto.Email,
	}

	logger.Log.Info("Sending SendOTP request to Auth Service",
		zap.String("email", dto.Email),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := authC.client.SendOTP(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("SendOTP request failed",
			zap.String("email", dto.Email),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
		return datatransfers.SendOtpResponse{}, err
	}

	logger.Log.Info("SendOTP request succeeded",
		zap.String("email", dto.Email),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return datatransfers.SendOtpResponse{
		Message: resp.Message,
	}, nil
}

func (authC *authClient) VerifyEmail(ctx context.Context, dto datatransfers.VerifyEmailRequest) (datatransfers.VerifyEmailResponse, error) {
	reqProto := protoAuth.VerifyEmailRequest{
		Email: dto.Email,
		Otp:   dto.OTP,
	}

	logger.Log.Info("Sending VerifyEmail request to Auth Service",
		zap.String("email", dto.Email),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := authC.client.VerifyEmail(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("VerifyEmail request failed",
			zap.String("email", dto.Email),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
		return datatransfers.VerifyEmailResponse{}, err
	}

	logger.Log.Info("VerifyEmail request succeeded",
		zap.String("email", dto.Email),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return datatransfers.VerifyEmailResponse{
		Message: resp.Message,
	}, nil
}

func (authC *authClient) Login(ctx context.Context, dto datatransfers.LoginRequest) (datatransfers.LoginResponse, error) {
	reqProto := protoAuth.LoginRequest{
		Email:    dto.Email,
		Password: dto.Password,
	}

	logger.Log.Info("Sending Login request to Auth Service",
		zap.String("email", dto.Email),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := authC.client.Login(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("Login request failed",
			zap.String("email", dto.Email),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
		return datatransfers.LoginResponse{}, err
	}

	logger.Log.Info("Login request succeeded",
		zap.String("email", dto.Email),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return datatransfers.LoginResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		Message:      resp.Message,
	}, nil
}

func (authC *authClient) ValidateToken(ctx context.Context, dto datatransfers.ValidateTokenRequest) (datatransfers.ValidateTokenResponse, error) {
	reqProto := protoAuth.ValidateTokenRequest{
		Token: dto.Token,
	}

	logger.Log.Info("Sending ValidateToken request to Auth Service",
		zap.String("token", dto.Token),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := authC.client.ValidateToken(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("ValidateToken request failed",
			zap.String("token", dto.Token),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
		return datatransfers.ValidateTokenResponse{}, err
	}

	logger.Log.Info("ValidateToken request succeeded",
		zap.String("token", dto.Token),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return datatransfers.ValidateTokenResponse{
		Valid:  resp.Valid,
		UserID: resp.UserId,
		Role:   resp.Role,
		Email:  resp.Email,
	}, nil
}
