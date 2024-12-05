package grpc_server

import (
	"context"
	"fmt"
	"user_service/internal/service"
	protoUser "user_service/proto/user_service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userGRPCServer struct {
	userService service.UserService
	protoUser.UnimplementedUserServiceServer
}

func NewUserGRPCServer(userService service.UserService) protoUser.UserServiceServer {
	return &userGRPCServer{
		userService: userService,
	}
}

func (s *userGRPCServer) GetUserById(ctx context.Context, req *protoUser.GetUserByIdRequest) (*protoUser.GetUserByIdResponse, error) {
	user, err := s.userService.GetUserById(ctx, &req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to retrieve user with id '%s'", req.UserId))
	}

	return &protoUser.GetUserByIdResponse{
		User: &protoUser.User{
			Id:        user.Id,
			Email:     user.Email,
			Username:  user.Username,
			Password:  user.Password,
			Verified:  user.Verified,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Unix(),
			UpdatedAt: user.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *userGRPCServer) GetUserByEmail(ctx context.Context, req *protoUser.GetUserByEmailRequest) (*protoUser.GetUserByEmailResponse, error) {
	user, err := s.userService.GetUserByEmail(ctx, &req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to retrieve user with email '%s'", req.Email))
	}

	return &protoUser.GetUserByEmailResponse{
		User: &protoUser.User{
			Id:        user.Id,
			Email:     user.Email,
			Username:  user.Username,
			Password:  user.Password,
			Verified:  user.Verified,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Unix(),
			UpdatedAt: user.UpdatedAt.Unix(),
		},
	}, nil
}

func (s *userGRPCServer) ListUsers(ctx context.Context, req *protoUser.ListUsersRequest) (*protoUser.ListUsersResponse, error) {
	users, err := s.userService.ListUsers(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to retrieve user list")
	}

	var protoUsers []*protoUser.User
	for _, user := range users {
		protoUsers = append(protoUsers, &protoUser.User{
			Id:        user.Id,
			Email:     user.Email,
			Username:  user.Username,
			Password:  user.Password,
			Verified:  user.Verified,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Unix(),
			UpdatedAt: user.UpdatedAt.Unix(),
		})
	}

	return &protoUser.ListUsersResponse{
		Users: protoUsers,
	}, nil
}
