package main

import (
	"os"
	"strconv"
)

//Made in reference to https://dev.to/craicoverflow/a-no-nonsense-guide-to-environment-variables-in-go-a2f

type DatabaseConfig struct {
	SqliteFolderPath string
	SqliteFileName   string
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
			SqliteFolderPath: getEnv("SQLITE_DB_FOLDER", ""),
			SqliteFileName:   getEnv("SQLITE_DB_NAME", ""),
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
