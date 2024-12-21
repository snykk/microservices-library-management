package grpc_server

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type healthGRPCServer struct {
	grpc_health_v1.UnimplementedHealthServer
}

// NewHealthGRPCServer initializes the gRPC Health server.
func NewHealthGRPCServer() grpc_health_v1.HealthServer {
	return &healthGRPCServer{}
}

// Check implements the standard gRPC Health Check.
func (s *healthGRPCServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	// Ensure that the requested service is auth_service
	if req.Service == "" {
		// If the service is empty, return status UNKNOWN
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_UNKNOWN,
		}, status.Errorf(codes.InvalidArgument, "Service name is required")
	}

	if req.Service == "auth_service" {
		// If the service is auth_service, return status SERVING
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_SERVING,
		}, nil
	}

	// If the service is not recognized, return status UNKNOWN
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_UNKNOWN,
	}, nil
}

// Watch streams the health status periodically.
func (s *healthGRPCServer) Watch(req *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	// Ensure that the service name is included in the request
	if req.Service == "" {
		return status.Errorf(codes.InvalidArgument, "Service name is required")
	}

	// Simulate monitoring the health check status periodically
	for {
		if req.Service == "auth_service" {
			// Only auth_service will send the SERVING status
			if err := stream.Send(&grpc_health_v1.HealthCheckResponse{
				Status: grpc_health_v1.HealthCheckResponse_SERVING,
			}); err != nil {
				return status.Errorf(codes.Internal, "Failed to send health check response: %v", err)
			}
		} else {
			// For services other than auth_service, send the UNKNOWN status
			if err := stream.Send(&grpc_health_v1.HealthCheckResponse{
				Status: grpc_health_v1.HealthCheckResponse_UNKNOWN,
			}); err != nil {
				return status.Errorf(codes.Internal, "Failed to send health check response: %v", err)
			}
		}

		// Wait for 5 seconds before sending the status again
		time.Sleep(5 * time.Second)
	}
}
