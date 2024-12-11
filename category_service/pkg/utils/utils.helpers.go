package utils

import (
	"category_service/internal/constants"
	"context"
	"fmt"
	"path/filepath"
	"runtime"

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

// GetRequestIDFromContext retrieves the request ID from gRPC metadata.
func GetRequestIDFromContext(ctx context.Context) string {
	requestId, ok := GetMetadataValue(ctx, constants.ContextProtoRequestIDKey)
	if !ok {
		return "unknown"
	}
	return requestId
}
