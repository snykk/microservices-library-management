package exception

import (
	"book_service/internal/service"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GRPCErrorFormatter(err error) error {
	switch {
	case errors.Is(err, service.ErrGetBook):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, service.ErrGetListBook):
		return status.Error(codes.Internal, err.Error())
	case errors.Is(err, service.ErrUpdateBook):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrDeleteBook):
		return status.Error(codes.Internal, err.Error())
	default:
		return status.Error(codes.Unknown, err.Error())
	}
}
