package config

import (
	"fmt"
	"wallet-service/internal/logger"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	AppPort    string `mapstructure:"APP_PORT"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.Log.Errorf("Error reading config file: %v", err)
		return nil, fmt.Errorf("error reading config file: %v", err)
	}
	var config Config

	logger.Log.Debug("Unmarshalling config data into struct...")

	if err := viper.Unmarshal(&config); err != nil {
		logger.Log.Errorf("Error unmarshalling config data: %v", err)
		return nil, fmt.Errorf("error unmarshalling config data: %v", err)
	}

	logger.Log.Info("Configuration loaded successfully")
	return &config, nil
}
