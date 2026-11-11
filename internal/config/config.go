package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DBPath    string
	JWTSecret string
	TokenTTL  time.Duration
}

// read settings from the .env file (if it exists) and environment variables
func Load() Config {
	if err := godotenv.Load(); err == nil {
		log.Println("config: loaded settings from .env")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("config: JWT_SECRET is not set (copy .env.example to .env)")
	}
	if secret == "replace-with-your-own-long-random-text" {
		log.Println("config: WARNING JWT_SECRET is still the example value, change it in .env")
	}

	ttl, err := time.ParseDuration(getEnv("TOKEN_TTL", "15m"))
	if err != nil {
		ttl = 15 * time.Minute
	}
	return Config{
		Port:      getEnv("PORT", "8080"),
		DBPath:    getEnv("DB_PATH", "lab04.db"),
		JWTSecret: secret,
		TokenTTL:  ttl,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
