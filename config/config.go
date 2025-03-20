package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	ApiUrl      string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		ApiUrl:      os.Getenv("API_URL"),
	}
}
