package clients

import (
	"api_gateway/internal/constants"
	"api_gateway/internal/datatransfers"
	protoUser "api_gateway/proto/user_service"
	"context"
	"time"

	"api_gateway/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient interface {
	GetUserById(ctx context.Context, userId string) (datatransfers.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (datatransfers.UserResponse, error)
	ListUsers(ctx context.Context) ([]datatransfers.UserResponse, error)
}

type userClient struct {
	client protoUser.UserServiceClient
}

func NewUserClient() (UserClient, error) {
	conn, err := grpc.NewClient("user-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Error("Failed to create UserClient",
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryConnection),
		)
		return nil, err
	}

	client := protoUser.NewUserServiceClient(conn)

	logger.Log.Info("Successfully created UserClient",
		zap.String(constants.LoggerCategory, constants.LoggerCategoryConnection),
	)

	return &userClient{
		client: client,
	}, nil
}

func (u *userClient) GetUserById(ctx context.Context, userId string) (datatransfers.UserResponse, error) {
	reqProto := protoUser.GetUserByIdRequest{
		UserId: userId,
	}

	logger.Log.Info("Sending GetUserById request to User Service",
		zap.String("user_id", userId),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := u.client.GetUserById(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("GetUserById request failed",
			zap.String("user_id", userId),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
		return datatransfers.UserResponse{}, err
	}

	logger.Log.Info("GetUserById request succeeded",
		zap.String("user_id", resp.User.Id),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return datatransfers.UserResponse{
		Id:        resp.User.Id,
		Username:  resp.User.Username,
		Email:     resp.User.Email,
		Verified:  resp.User.Verified,
		Role:      resp.User.Role,
		CreatedAt: time.Unix(resp.User.CreatedAt, 0),
		UpdatedAt: time.Unix(resp.User.UpdatedAt, 0),
	}, nil
}

func (u *userClient) GetUserByEmail(ctx context.Context, email string) (datatransfers.UserResponse, error) {
	reqProto := protoUser.GetUserByEmailRequest{
		Email: email,
	}

	logger.Log.Info("Sending GetUserByEmail request to User Service",
		zap.String("email", email),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := u.client.GetUserByEmail(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("GetUserByEmail request failed",
			zap.String("email", email),
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
		return datatransfers.UserResponse{}, err
	}

	logger.Log.Info("GetUserByEmail request succeeded",
		zap.String("email", resp.User.Email),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return datatransfers.UserResponse{
		Id:        resp.User.Id,
		Username:  resp.User.Username,
		Email:     resp.User.Email,
		Verified:  resp.User.Verified,
		Role:      resp.User.Role,
		CreatedAt: time.Unix(resp.User.CreatedAt, 0),
		UpdatedAt: time.Unix(resp.User.UpdatedAt, 0),
	}, nil
}

func (u *userClient) ListUsers(ctx context.Context) ([]datatransfers.UserResponse, error) {
	reqProto := protoUser.ListUsersRequest{}

	logger.Log.Info("Sending ListUsers request to User Service",
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	resp, err := u.client.ListUsers(ctx, &reqProto)
	if err != nil {
		logger.Log.Error("ListUsers request failed",
			zap.Error(err),
			zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
		)
		return nil, err
	}

	var users []datatransfers.UserResponse
	for _, user := range resp.Users {
		users = append(users, datatransfers.UserResponse{
			Id:        user.Id,
			Username:  user.Username,
			Email:     user.Email,
			Verified:  user.Verified,
			Role:      user.Role,
			CreatedAt: time.Unix(user.CreatedAt, 0),
			UpdatedAt: time.Unix(user.UpdatedAt, 0),
		})
	}

	logger.Log.Info("ListUsers request succeeded",
		zap.Int("users_count", len(users)),
		zap.String(constants.LoggerCategory, constants.LoggerCategoryGrpcClient),
	)

	return users, nil
}
