package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config holds all configuration values
type Config struct {
	AppPort string
	DBHost  string
	DBUser  string
	DBPass  string
}

// Global variable to hold the loaded configuration
var AppConfig *Config

// LoadConfig loads environment variables into the AppConfig global variable
func LoadConfig() {
	v := viper.New()
	v.SetConfigFile(".env") // Specify the .env file
	v.AutomaticEnv()        // Allow system environment variable overrides

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Populate the global AppConfig variable
	AppConfig = &Config{
		AppPort: v.GetString("APP_PORT"),
		DBHost:  v.GetString("DB_HOST"),
		DBUser:  v.GetString("DB_USER"),
		DBPass:  v.GetString("DB_PASS"),
	}
}
