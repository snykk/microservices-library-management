package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	LogPath         string
	RabbitMQURL     string
	MongoURL        string
	MongoDB         string
	MongoCollection string
}

var AppConfig Config

func InitializeAppConfig() error {
	viper.AutomaticEnv()

	// Direct assignments for string values
	requiredStringKeys := map[string]*string{
		"LOG_PATH":         &AppConfig.LogPath,
		"RABBITMQ_URL":     &AppConfig.RabbitMQURL,
		"MONGO_URL":        &AppConfig.MongoURL,
		"MONGO_DB":         &AppConfig.MongoDB,
		"MONGO_COLLECTION": &AppConfig.MongoCollection,
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
