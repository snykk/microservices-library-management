package configs

import (
	"fmt"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort                string
	RabbitMQURL            string
	ReadTimeout            int
	WriteTimeout           int
	AuthServiceURL         string
	AuthorServiceURL       string
	BookServiceURL         string
	CategoryServiceURL     string
	LoanServiceURL         string
	UserServiceURL         string
	LoggerWorkerType       string
	LoggerWorkerNum        int
	LoggerWorkerBufferSize int
	MaxRequestPerMinute    int
} // mapstrucuture issue: should assign manually

var AppConfig Config

func InitializeAppConfig() error {
	viper.AutomaticEnv()

	// Define a helper function for assigning integer values
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

	// Direct assignments for string values
	requiredStringKeys := map[string]*string{
		"APP_PORT":             &AppConfig.AppPort,
		"RABBITMQ_URL":         &AppConfig.RabbitMQURL,
		"AUTH_SERVICE_URL":     &AppConfig.AuthServiceURL,
		"AUTHOR_SERVICE_URL":   &AppConfig.AuthorServiceURL,
		"BOOK_SERVICE_URL":     &AppConfig.BookServiceURL,
		"CATEGORY_SERVICE_URL": &AppConfig.CategoryServiceURL,
		"LOAN_SERVICE_URL":     &AppConfig.LoanServiceURL,
		"USER_SERVICE_URL":     &AppConfig.UserServiceURL,
		"LOGGER_WORKER_TYPE":   &AppConfig.LoggerWorkerType,
	}

	for key, ref := range requiredStringKeys {
		*ref = viper.GetString(key)
		if *ref == "" {
			return fmt.Errorf("%s is required", key)
		}
	}

	// Assign integer values with validation
	var err error
	AppConfig.ReadTimeout, err = getIntEnv("READ_TIMEOUT")
	if err != nil || AppConfig.ReadTimeout == 0 {
		return err
	}

	AppConfig.WriteTimeout, err = getIntEnv("WRITE_TIMEOUT")
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

	AppConfig.MaxRequestPerMinute, err = getIntEnv("MAX_REQUEST_PER_MINUTE")
	if err != nil {
		return err
	}

	return nil
}
