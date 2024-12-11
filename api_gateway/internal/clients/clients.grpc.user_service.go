package clients

import (
	"api_gateway/internal/constants"
	"api_gateway/internal/datatransfers"
	protoUser "api_gateway/proto/user_service"
	"context"
	"log"
	"time"

	"api_gateway/pkg/logger"
	"api_gateway/pkg/utils"

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
	logger *logger.Logger
}

func NewUserClient(logger *logger.Logger) (UserClient, error) {
	conn, err := grpc.NewClient("user-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Failed to create UserClient:", err)
		return nil, err
	}
	client := protoUser.NewUserServiceClient(conn)

	log.Println("Successfully created UserClient")

	return &userClient{
		client: client,
		logger: logger,
	}, nil
}

func (u *userClient) GetUserById(ctx context.Context, userId string) (datatransfers.UserResponse, error) {
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoUser.GetUserByIdRequest{
		UserId: userId,
	}

	extra := map[string]interface{}{
		"user_id": userId,
	}

	u.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending GetUserById request to User Service", extra, nil)

	resp, err := u.client.GetUserById(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		u.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "GetUserById request failed", extra, err)
		return datatransfers.UserResponse{}, err
	}

	u.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "GetUserById request succeeded", extra, nil)

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
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoUser.GetUserByEmailRequest{
		Email: email,
	}

	extra := map[string]interface{}{
		"email": email,
	}

	u.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending GetUserByEmail request to User Service", extra, nil)

	resp, err := u.client.GetUserByEmail(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		u.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "GetUserByEmail request failed", extra, err)
		return datatransfers.UserResponse{}, err
	}

	u.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "GetUserByEmail request succeeded", extra, nil)

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
	requestID := utils.GetRequestIDFromContext(ctx)

	reqProto := protoUser.ListUsersRequest{}

	u.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Sending ListUsers request to User Service", nil, nil)

	resp, err := u.client.ListUsers(utils.GetProtoContext(ctx), &reqProto)
	if err != nil {
		u.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "ListUsers request failed", nil, err)
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

	u.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "ListUsers request succeeded", nil, nil)

	return users, nil
}
