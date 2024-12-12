package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	EmailSenderContainerFile   string
	EmailPasswordContainerFile string
	RabbitMQURL                string
}

var AppConfig Config

func InitializeAppConfig() error {
	viper.AutomaticEnv()

	// Map for required string values
	requiredStringKeys := map[string]*string{
		"EMAIL_SENDER_CONTAINER_FILE":   &AppConfig.EmailSenderContainerFile,
		"EMAIL_PASSWORD_CONTAINER_FILE": &AppConfig.EmailPasswordContainerFile,
		"RABBITMQ_URL":                  &AppConfig.RabbitMQURL,
	}

	// Assign string values
	for key, ref := range requiredStringKeys {
		*ref = viper.GetString(key)
		if *ref == "" {
			return fmt.Errorf("%s is required", key)
		}
	}
	return nil
}
