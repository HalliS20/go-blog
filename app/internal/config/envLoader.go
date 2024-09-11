package config

import (
	"github.com/joho/godotenv"
	"os"
)

func LoadConfig() error {
	// Try to load .env file
	if err := godotenv.Load(); err != nil {
		// If .env file doesn't exist, it's not an error
		// We'll use the environment variables set by the system
		if !os.IsNotExist(err) {
			return err
		}
	}

	// At this point, environment variables are either loaded from .env or already set by the system
	return nil
}

// GetEnv retrieves an environment variable or returns a default value
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
