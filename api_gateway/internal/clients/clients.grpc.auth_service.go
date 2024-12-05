package clients

import (
	"api_gateway/internal/datatransfers"
	protoAuth "api_gateway/proto/auth_service"
	"context"
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
}

func NewAuthClient() (AuthClient, error) {
	conn, err := grpc.NewClient("auth-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := protoAuth.NewAuthServiceClient(conn)
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

	resp, err := authC.client.Register(ctx, &reqProto)
	if err != nil {
		return datatransfers.RegisterResponse{}, err
	}

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

	resp, err := authC.client.SendOTP(ctx, &reqProto)
	if err != nil {
		return datatransfers.SendOtpResponse{}, err
	}

	return datatransfers.SendOtpResponse{
		Message: resp.Message,
	}, nil
}

func (authC *authClient) VerifyEmail(ctx context.Context, dto datatransfers.VerifyEmailRequest) (datatransfers.VerifyEmailResponse, error) {
	reqProto := protoAuth.VerifyEmailRequest{
		Email: dto.Email,
		Otp:   dto.OTP,
	}

	resp, err := authC.client.VerifyEmail(ctx, &reqProto)
	if err != nil {
		return datatransfers.VerifyEmailResponse{}, err
	}

	return datatransfers.VerifyEmailResponse{
		Message: resp.Message,
	}, nil
}

func (authC *authClient) Login(ctx context.Context, dto datatransfers.LoginRequest) (datatransfers.LoginResponse, error) {
	reqProto := protoAuth.LoginRequest{
		Email:    dto.Email,
		Password: dto.Password,
	}

	resp, err := authC.client.Login(ctx, &reqProto)
	if err != nil {
		return datatransfers.LoginResponse{}, err
	}

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

	resp, err := authC.client.ValidateToken(ctx, &reqProto)
	if err != nil {
		return datatransfers.ValidateTokenResponse{}, err
	}

	return datatransfers.ValidateTokenResponse{
		Valid:  resp.Valid,
		UserID: resp.UserId,
		Role:   resp.Role,
	}, nil
}
