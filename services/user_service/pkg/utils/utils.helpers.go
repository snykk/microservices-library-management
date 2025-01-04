package utils

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"user_service/internal/constants"

	"google.golang.org/grpc/metadata"
)

func GetLocation() string {
	_, file, line, _ := runtime.Caller(1)

	dir := filepath.Base(filepath.Dir(file))
	base := filepath.Base(file)

	return fmt.Sprintf("%s/%s:%d", dir, base, line)
}

// GetMetadataValue retrieves a value from gRPC metadata by key.
func GetMetadataValue(ctx context.Context, key string) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false // Metadata not available
	}

	values := md[key]
	if len(values) > 0 {
		return values[0], true // Return the first value found for the key
	}

	return "", false // Key not found in metadata
}

// GetRequestIDFromMetadataContext retrieves the request ID from gRPC metadata.
func GetRequestIDFromMetadataContext(ctx context.Context) string {
	requestId, ok := GetMetadataValue(ctx, constants.ContextProtoRequestIDKey)
	if !ok {
		return "unknown"
	}
	return requestId
}

// CalculateTotalPages calculates the total number of pages based on total items and page size.
func CalculateTotalPages(totalItems, pageSize int) int {
	if totalItems == 0 || pageSize == 0 {
		return 0
	}

	totalPages := totalItems / pageSize
	if totalItems%pageSize > 0 {
		totalPages++
	}

	return totalPages
}
