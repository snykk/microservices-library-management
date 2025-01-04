package configs

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	GrpcPort                   string
	DSN                        string
	RedisHost                  string
	RedisPassword              string
	RedisDefaultExp            time.Duration // minute unit
	RedisDB                    int
	RedisPort                  string
	JwtIssuer                  string
	JwtSecret                  string
	JwtExpAccessToken          time.Duration // minute unit
	JwtExpRefreshToken         time.Duration // minute unit
	EmailSenderContainerFile   string
	EmailPasswordContainerFile string
	RabbitMQURL                string
	LoggerWorkerType           string
	LoggerWorkerNum            int
	LoggerWorkerBufferSize     int
}

var AppConfig Config

func InitializeAppConfig() error {
	viper.AutomaticEnv()

	// Helper functions for integer and duration values
	getIntEnv := func(key string) (int, error) {
		val := viper.GetString(key)
		if val == "" {
			return 0, fmt.Errorf("%s is required", key)
		}
		num, err := strconv.Atoi(val)
		if err != nil {
			return 0, fmt.Errorf("%s must be a valid integer", key)
		}
		return num, nil
	}

	getDurationEnv := func(key string, unit time.Duration) (time.Duration, error) {
		val, err := getIntEnv(key)
		if err != nil {
			return 0, err
		}
		return time.Duration(val) * unit, nil
	}

	// Map for required string values
	requiredStringKeys := map[string]*string{
		"GRPC_PORT":                     &AppConfig.GrpcPort,
		"DSN":                           &AppConfig.DSN,
		"REDIS_HOST":                    &AppConfig.RedisHost,
		"REDIS_PASSWORD":                &AppConfig.RedisPassword,
		"REDIS_PORT":                    &AppConfig.RedisPort,
		"JWT_ISSUER":                    &AppConfig.JwtIssuer,
		"JWT_SECRET":                    &AppConfig.JwtSecret,
		"EMAIL_SENDER_CONTAINER_FILE":   &AppConfig.EmailSenderContainerFile,
		"EMAIL_PASSWORD_CONTAINER_FILE": &AppConfig.EmailPasswordContainerFile,
		"RABBITMQ_URL":                  &AppConfig.RabbitMQURL,
		"LOGGER_WORKER_TYPE":            &AppConfig.LoggerWorkerType,
	}

	// Assign string values
	for key, ref := range requiredStringKeys {
		*ref = viper.GetString(key)
		if *ref == "" {
			return fmt.Errorf("%s is required", key)
		}
	}

	// Assign integer values
	var err error
	AppConfig.RedisDB, err = getIntEnv("REDIS_DB")
	if err != nil {
		return err
	}

	// Assign duration values
	AppConfig.RedisDefaultExp, err = getDurationEnv("REDIS_DEFAULT_EXP", time.Minute)
	if err != nil {
		return err
	}

	AppConfig.JwtExpAccessToken, err = getDurationEnv("JWT_EXP_ACCESS_TOKEN", time.Minute)
	if err != nil {
		return err
	}

	AppConfig.JwtExpRefreshToken, err = getDurationEnv("JWT_EXP_REFRESH_TOKEN", time.Minute)
	if err != nil {
		return err
	}

	AppConfig.LoggerWorkerNum, err = getIntEnv("LOGGER_WORKER_NUM")
	if err != nil {
		return err
	}

	AppConfig.LoggerWorkerBufferSize, err = getIntEnv("LOGGER_WORKER_BUFFER_SIZE")
	if err != nil {
		return err
	}

	return nil
}
