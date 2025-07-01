package config

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DriverAPIUrl   string
	DriverAPIToken string
	Ctx            context.Context
	SecretKey      string
}

func NewConfig() *Config {
	if err := loadEnv(); err != nil {
		log.Printf("Error loading env: %v", err)
		return nil
	}

	ctx := context.Background()
	cfg := &Config{
		Port:           getEnvWithDefault("PORT", "10002"),
		Ctx:            ctx,
		DriverAPIUrl:   getEnvWithDefault("DRIVER_LOCATION_BASE_URL", "http://127.0.0.1:10001"),
		DriverAPIToken: getEnvWithDefault("DRIVER_LOCATION_API_TOKEN", ""),
		SecretKey:      getEnvWithDefault("SECRET_KEY", "driver-rider-matching-api-secret-key"),
	}

	return cfg
}

func getEnvWithDefault(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func loadEnv() error {
	return godotenv.Load()
}
