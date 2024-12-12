package grpc_server

import (
	"context"
	"fmt"
	"user_service/internal/constants"
	"user_service/internal/service"
	"user_service/pkg/logger"
	"user_service/pkg/utils"
	protoUser "user_service/proto/user_service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userGRPCServer struct {
	userService service.UserService
	logger      *logger.Logger
	protoUser.UnimplementedUserServiceServer
}

func NewUserGRPCServer(userService service.UserService, logger *logger.Logger) protoUser.UserServiceServer {
	return &userGRPCServer{
		userService: userService,
		logger:      logger,
	}
}

func (s *userGRPCServer) GetUserById(ctx context.Context, req *protoUser.GetUserByIdRequest) (*protoUser.GetUserByIdResponse, error) {
	requestID := utils.GetRequestIDFromMetadataContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received GetUserById request", map[string]interface{}{"user_id": req.UserId}, nil)

	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid GetUserById request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	user, err := s.userService.GetUserById(ctx, req.UserId)
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, fmt.Sprintf("Failed to retrieve user with id '%s'", req.UserId), nil, err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to retrieve user with id '%s'", req.UserId))
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "User retrieved successfully", map[string]interface{}{"user_id": user.Id}, nil)

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
	requestID := utils.GetRequestIDFromMetadataContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received GetUserByEmail request", map[string]interface{}{"email": req.Email}, nil)

	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid GetUserByEmail request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	user, err := s.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, fmt.Sprintf("Failed to retrieve user with email '%s'", req.Email), nil, err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to retrieve user with email '%s'", req.Email))
	}

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "User retrieved successfully", map[string]interface{}{"user_id": user.Id}, nil)

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
	requestID := utils.GetRequestIDFromMetadataContext(ctx)
	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "Received ListUsers request", nil, nil)

	if err := req.Validate(); err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Invalid ListUsers request", nil, err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	users, err := s.userService.ListUsers(ctx)
	if err != nil {
		s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelError, "Failed to retrieve user list", nil, err)
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

	s.logger.LogMessage(utils.GetLocation(), requestID, constants.LogLevelInfo, "User list retrieved successfully", nil, nil)

	return &protoUser.ListUsersResponse{
		Users: protoUsers,
	}, nil
}
