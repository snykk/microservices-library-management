package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.Logger

// LoggerConfig holds configuration for the logger.
type LoggerConfig struct {
	OutputPaths []string // Output paths for logs
	MaxSize     int      // Max size of log file in MB
	MaxBackups  int      // Max number of backup files
	MaxAge      int      // Max age of a log file in days
	Compress    bool     // Compress backup log files
	IsDev       bool     // Flag for development environment
	ServiceName string   // Name of the service
}

// DefaultLoggerConfig provides default configuration for the logger.
var DefaultLoggerConfig = LoggerConfig{
	OutputPaths: []string{"stdout", "../../logs/service.log"},
	MaxSize:     10, // 10 MB
	MaxBackups:  5,
	MaxAge:      30, // 30 days
	Compress:    true,
	IsDev:       false,
	ServiceName: "default-service",
}

// // getRelativePath returns the absolute path relative to the location of this file.
// func getRelativePath(relativePath string) string { // running in local (no container)
// 	_, b, _, _ := runtime.Caller(0) // Get the path to this file (logger.go)
// 	basePath := filepath.Dir(b)
// 	return filepath.Join(basePath, relativePath)
// }

// Initialize initializes the logger with custom configuration.
func Initialize(config LoggerConfig) error {
	// Determine log level
	atomicLevel := zap.NewAtomicLevel()

	// // Adjust paths to be relative to logger package
	// for i, path := range config.OutputPaths {
	// 	if path != "stdout" && path != "stderr" {
	// 		config.OutputPaths[i] = getRelativePath(path)
	// 	}
	// }

	// Set up encoder based on environment
	var encoder zapcore.Encoder
	if config.IsDev {
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	} else {
		encoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		})
	}

	// Create log writers
	writers := make([]zapcore.WriteSyncer, 0)
	for _, path := range config.OutputPaths {
		if path == "stdout" {
			writers = append(writers, zapcore.AddSync(os.Stdout))
		} else {
			writers = append(writers, zapcore.AddSync(&lumberjack.Logger{
				Filename:   path,
				MaxSize:    config.MaxSize,
				MaxBackups: config.MaxBackups,
				MaxAge:     config.MaxAge,
				Compress:   config.Compress,
			}))
		}
	}
	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writers...), atomicLevel)

	// Build logger with default "service" field
	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.Fields(zap.String("service", config.ServiceName)), // Add default service field
	)
	Log = logger

	// Confirm logger initialization
	Log.Info("Logger initialized",
		zap.Strings("outputPaths", config.OutputPaths),
		zap.String("service", config.ServiceName),
		zap.Bool("isDev", config.IsDev),
	)
	return nil
}

// Sync flushes any buffered log entries.
func Sync() {
	if err := Log.Sync(); err != nil {
		Log.Error("Failed to sync logger", zap.Error(err))
	}
}

// // Info logs an info message with optional fields.
// func Info(msg string, fields ...zap.Field) {
// 	log.Info(msg, fields...)
// }

// // Debug logs a debug message with optional fields.
// func Debug(msg string, fields ...zap.Field) {
// 	log.Debug(msg, fields...)
// }

// // Warn logs a warning message with optional fields.
// func Warn(msg string, fields ...zap.Field) {
// 	log.Warn(msg, fields...)
// }

// // Error logs an error message with optional fields.
// func Error(msg string, fields ...zap.Field) {
// 	log.Error(msg, fields...)
// }

// // Panic logs a panic message with optional fields and then panics.
// func Panic(msg string, fields ...zap.Field) {
// 	log.Panic(msg, fields...)
// }

// // Fatal logs a fatal message with optional fields and then exits the application.
// func Fatal(msg string, fields ...zap.Field) {
// 	log.Fatal(msg, fields...)
// }
