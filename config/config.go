package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port            string
	JWTSecret       string
	ExpirationHours int
	DatabaseURL     string
	APIBaseURL      string
	TokenExpiration time.Duration
}

var AppConfig *Config

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	AppConfig = &Config{
		Port:            getEnv("PORT", ":8080"),
		JWTSecret:       getEnv("JWT_SECRET", "default-secret-key"),
		ExpirationHours: getEnvAsInt("EXPIRATION_HOURS", 2),
		DatabaseURL:     getEnv("DB", "host=localhost user=postgres password=postgres port=5432 sslmode=disable"),
		APIBaseURL:      getEnv("APIBaseURL", "http://localhost:8080/api/"),
	}

	AppConfig.TokenExpiration = time.Duration(AppConfig.ExpirationHours) * time.Hour
}

func Get() *Config {
	return AppConfig
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	strValue := getEnv(key, "")
	if strValue == "" {
		return defaultValue
	}

	if intValue, err := strconv.Atoi(strValue); err == nil {
		return intValue
	}

	log.Printf("Invalid integer value for %s: %s, using default: %d", key, strValue, defaultValue)
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	strValue := getEnv(key, "")
	if strValue == "" {
		return defaultValue
	}

	if boolValue, err := strconv.ParseBool(strValue); err == nil {
		return boolValue
	}

	log.Printf("Invalid boolean value for %s: %s, using default: %v", key, strValue, defaultValue)
	return defaultValue
}
