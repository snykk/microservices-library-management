package utils

import (
	"api_gateway/internal/constants"
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

func GetRequestIDFromContext(ctx context.Context) string {
	requestID, ok := ctx.Value(constants.ContextRequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}

	return requestID
}

// GetProtoContext adds a request ID to the gRPC metadata context.
func GetProtoContext(ctx context.Context) context.Context {
	requestID := GetRequestIDFromContext(ctx)

	md := metadata.Pairs(constants.ContextProtoRequestIDKey, requestID) // all upercase md key automatically convert to lower

	return metadata.NewOutgoingContext(ctx, md)
}
