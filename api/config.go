package main

import (
	"os"
	"strconv"
)

//Made in reference to https://dev.to/craicoverflow/a-no-nonsense-guide-to-environment-variables-in-go-a2f

type DatabaseConfig struct {
	DBname   string
	User     string
	Password string
	Host     string
	Port     int
}

type FrontendCorsConfig struct {
	Host string
	Port int
}

type Config struct {
	Database DatabaseConfig
	Frontend FrontendCorsConfig
	ApiPort  int
}

func New() *Config {
	return &Config{
		Database: DatabaseConfig{
			DBname:   getEnv("DATABASE_NAME", ""),
			User:     getEnv("DATABASE_USER", ""),
			Password: getEnv("DATABASE_PASSWORD", ""),
			Host:     getEnv("DATABASE_HOST", "localhost"),
			Port:     getEnvAsInt("DATABASE_PORT", 5432),
		},
		Frontend: FrontendCorsConfig{
			Host: getEnv("FRONTEND_HOST", "http://localhost"),
			Port: getEnvAsInt("FRONTEND_PORT", 8000),
		},
		ApiPort: getEnvAsInt("API_PORT", 8080),
	}
}

func getEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueString := getEnv(key, "")
	if value, err := strconv.Atoi(valueString); err == nil {
		return value
	}
	return defaultValue
}
