package utils

import (
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
