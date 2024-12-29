package exception

import (
	"auth_service/internal/service"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GRPCErrorFormatter(err error) error {
	// Mapping error to gRPC codes
	switch {
	case errors.Is(err, service.ErrEmailAlreadyRegistered):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, service.ErrFailedHashPassword):
		return status.Error(codes.Internal, err.Error())
	case errors.Is(err, service.ErrEmailAlreadyVerified):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, service.ErrGetUserByEmail):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, service.ErrCreateUser):
		return status.Error(codes.Internal, err.Error())
	case errors.Is(err, service.ErrGenerateOTPCode):
		return status.Error(codes.Internal, err.Error())
	case errors.Is(err, service.ErrSendOtpWithMailer):
		return status.Error(codes.Internal, err.Error())
	case errors.Is(err, service.ErrMismatchOTPCode):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrUpdateUserVerification):
		return status.Error(codes.Internal, err.Error())
	case errors.Is(err, service.ErrEmailNotVerified):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, service.ErrInvalidPassword):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, service.ErrGenerateAccessToken):
		return status.Error(codes.Internal, err.Error())
	case errors.Is(err, service.ErrGenerateRefreshToken):
		return status.Error(codes.Internal, err.Error())
	case errors.Is(err, service.ErrPareseToken):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrInvalidRefreshToken):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, service.ErrUpdateRefreshToken):
		return status.Error(codes.Internal, err.Error())
	case errors.Is(err, service.ErrUpdateLastLogin):
		return status.Error(codes.Internal, err.Error())
	default:
		// Fallback for unknown errors
		return status.Error(codes.Unknown, err.Error())
	}
}
