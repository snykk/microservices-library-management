package configs

import (
	"fmt"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
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

	// Map for required string values
	requiredStringKeys := map[string]*string{
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

	// Assign integer values with validation
	var err error
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
