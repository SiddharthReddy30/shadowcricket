package config

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DataDir     string
	TokenSecret string
}

func Load(envFile string) (Config, error) {
	godotenv.Load(envFile)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "data"
	}

	secret := os.Getenv("TOKEN_SECRET")
	if len(secret) != 32 {
		return Config{}, fmt.Errorf("TOKEN_SECRET must be exactly 32 bytes, got %d", len(secret))
	}

	return Config{
		Port:        port,
		DataDir:     dataDir,
		TokenSecret: secret,
	}, nil
}
