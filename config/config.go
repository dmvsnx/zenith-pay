package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	AppName string
	AppEnv string
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	CloudinaryURL string

	JwtSecretKey string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		AppName:    loadEnv("APP_NAME"),
		AppEnv:     loadEnv("APP_ENV"),
		AppPort:    loadEnv("APP_PORT"),

		DBHost:     loadEnv("DB_HOST"),
		DBPort:     loadEnv("DB_PORT"),
		DBUser:     loadEnv("DB_USER"),
		DBPassword: loadEnv("DB_PASSWORD"),
		DBName:     loadEnv("DB_NAME"),

		CloudinaryURL: loadEnv("CLOUDINARY_URL"),
		
		JwtSecretKey: loadEnv("JWT_SECRET"),
	}
}

// loadEnv retrieves the value of the environment variable named by the key.
// It panics if the variable is not set.
func loadEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic("Environment variable " + key + " not set")
	}

	return value
}
