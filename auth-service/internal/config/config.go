package config

import (
	"auth-service/internal/logger"
	"os"
)

type Config struct {
	Port        string `mapstructure:"PORT"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
}

func LoadConfig() *Config {
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DATABASE_URL")
	jwt := os.Getenv("JWT_SECRET")

	if port == "" || dbURL == "" || jwt == "" {
		logger.Log.Fatal("Missing required environment variables")
	}

	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
		JWTSecret:   jwt,
	}
}
