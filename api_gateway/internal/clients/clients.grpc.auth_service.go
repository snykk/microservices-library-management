package clients

import (
	"api_gateway/configs"
	"api_gateway/internal/constants"
	"api_gateway/internal/datatransfers"
	"api_gateway/pkg/logger"
	"api_gateway/pkg/utils"
	protoAuth "api_gateway/proto/auth_service"
	"context"
	"log"
	"time"

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
	logger *logger.Logger
}

func NewAuthClient(logger *logger.Logger) (AuthClient, error) {
	conn, err := grpc.NewClient(configs.AppConfig.AuthServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Failed to create AuthClient:", err)
		return nil, err
	}
	client := protoAuth.NewAuthServiceClient(conn)

	log.Println("Successfully created AuthClient")

	return &authClient{
		client: client,
		logger: logger,
	}, nil
}

func (authC *authClient) Register(ctx context.Context, dto datatransfers.RegisterRequest) (datatransfers.RegisterResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoAuth.RegisterRequest{
		Email:    dto.Email,
		Username: dto.Username,
		Password: dto.Password,
	}

	extra := map[string]interface{}{
		"email":    dto.Email,
		"username": dto.Username,
	}

	authC.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending Register request to Auth Service", extra, nil)

	resp, err := authC.client.Register(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		authC.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Register request failed", extra, err)
		return datatransfers.RegisterResponse{}, err
	}

	authC.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Register request succeeded", extra, nil)

	return datatransfers.RegisterResponse{
		Id:        resp.User.Id,
		Email:     resp.User.Email,
		Username:  resp.User.Username,
		Password:  dto.Password, // Password as sent, not returned by the service
		Verified:  resp.User.Verified,
		CreatedAt: time.Unix(resp.User.CreatedAt, 0),
		UpdatedAt: time.Unix(resp.User.UpdatedAt, 0),
	}, nil
}

func (authC *authClient) Login(ctx context.Context, dto datatransfers.LoginRequest) (datatransfers.LoginResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoAuth.LoginRequest{
		Email:    dto.Email,
		Password: dto.Password,
	}

	extra := map[string]interface{}{
		"email": dto.Email,
	}

	authC.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending Login request to Auth Service", extra, nil)

	resp, err := authC.client.Login(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		authC.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Login request failed", extra, err)
		return datatransfers.LoginResponse{}, err
	}

	authC.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Login request succeeded", extra, nil)

	return datatransfers.LoginResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		Message:      resp.Message,
	}, nil
}

func (authC *authClient) SendOtp(ctx context.Context, dto datatransfers.SendOtpRequest) (datatransfers.SendOtpResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoAuth.SendOTPRequest{
		Email: dto.Email,
	}

	extra := map[string]interface{}{
		"email": dto.Email,
	}

	authC.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending SendOtp request to Auth Service", extra, nil)

	resp, err := authC.client.SendOTP(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		authC.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "SendOtp request failed", extra, err)
		return datatransfers.SendOtpResponse{}, err
	}

	authC.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "SendOtp request succeeded", extra, nil)

	return datatransfers.SendOtpResponse{
		Message: resp.Message,
	}, nil
}

func (authC *authClient) VerifyEmail(ctx context.Context, dto datatransfers.VerifyEmailRequest) (datatransfers.VerifyEmailResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoAuth.VerifyEmailRequest{
		Email: dto.Email,
		Otp:   dto.OTP,
	}

	extra := map[string]interface{}{
		"email": dto.Email,
		"otp":   dto.OTP,
	}

	authC.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending VerifyEmail request to Auth Service", extra, nil)

	resp, err := authC.client.VerifyEmail(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		authC.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "VerifyEmail request failed", extra, err)
		return datatransfers.VerifyEmailResponse{}, err
	}

	authC.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "VerifyEmail request succeeded", extra, nil)

	return datatransfers.VerifyEmailResponse{
		Message: resp.Message,
	}, nil
}

func (authC *authClient) ValidateToken(ctx context.Context, dto datatransfers.ValidateTokenRequest) (datatransfers.ValidateTokenResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoAuth.ValidateTokenRequest{
		Token: dto.Token,
	}

	extra := map[string]interface{}{
		"token": dto.Token,
	}

	authC.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending ValidateToken request to Auth Service", extra, nil)

	resp, err := authC.client.ValidateToken(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		authC.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "ValidateToken request failed", extra, err)
		return datatransfers.ValidateTokenResponse{}, err
	}

	authC.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "ValidateToken request succeeded", extra, nil)

	return datatransfers.ValidateTokenResponse{
		Valid:  resp.Valid,
		UserID: resp.UserId,
		Role:   resp.Role,
		Email:  resp.Email,
	}, nil
}
