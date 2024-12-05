package clients

import (
	"api_gateway/internal/datatransfers"
	protoUser "api_gateway/proto/user_service"
	"context"
	"time"

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
	conn, err := grpc.Dial("user-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := protoUser.NewUserServiceClient(conn)
	return &userClient{
		client: client,
	}, nil
}

func (u *userClient) GetUserById(ctx context.Context, userId string) (datatransfers.UserResponse, error) {
	reqProto := protoUser.GetUserByIdRequest{
		UserId: userId,
	}

	resp, err := u.client.GetUserById(ctx, &reqProto)
	if err != nil {
		return datatransfers.UserResponse{}, err
	}

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

	resp, err := u.client.GetUserByEmail(ctx, &reqProto)
	if err != nil {
		return datatransfers.UserResponse{}, err
	}

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

	resp, err := u.client.ListUsers(ctx, &reqProto)
	if err != nil {
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

	return users, nil
}
