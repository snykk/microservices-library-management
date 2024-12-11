package utils

import (
	"api_gateway/internal/constants"
	"context"
	"fmt"
	"path/filepath"
	"runtime"
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
