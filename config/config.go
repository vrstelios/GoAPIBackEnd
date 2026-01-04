package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server          ServerConfig
	Database        DatabaseConfig
	JWT             JWTConfig
	API             APIConfig
	TokenExpiration time.Duration
}

type ServerConfig struct {
	Port string
	Host string
	Env  string
}

type DatabaseConfig struct {
	//URL          string // Κρατάς την παλιά μορφή για συμβατότητα
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	JWTSecret       string
	ExpirationHours int
}

type APIConfig struct {
	BaseURL string
}

var AppConfig *Config

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	AppConfig = &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", ":8080"),
			Host: getEnv("HOST", "localhost"),
		},
		Database: parseDatabaseConfig(),
		JWT: JWTConfig{
			JWTSecret:       getEnv("JWT_SECRET", "default-secret-key-change-in-production"),
			ExpirationHours: getEnvAsInt("EXPIRATION_HOURS", 2),
		},
		API: APIConfig{
			BaseURL: getEnv("API_BASE_URL", "http://localhost:8080/api/"),
		},
	}

	log.Printf("Configuration loaded for environment: %s", AppConfig.Server.Env)
}

func parseDatabaseConfig() DatabaseConfig {
	dbURL := getEnv("DB", "host=localhost user=postgres password=postgres port=5432 sslmode=disable")

	params := parseConnectionString(dbURL)

	return DatabaseConfig{
		Host:     getParam(params, "host", "localhost"),
		Port:     getParam(params, "port", "5432"),
		User:     getParam(params, "user", "postgres"),
		Password: getParam(params, "password", "postgres"),
		Name:     getParam(params, "dbname", ""),
		SSLMode:  getParam(params, "sslmode", "disable"),
	}
}

func parseConnectionString(connStr string) map[string]string {
	params := make(map[string]string)

	pairs := strings.Fields(connStr)
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			params[kv[0]] = kv[1]
		}
	}

	return params
}

func getParam(params map[string]string, key, defaultValue string) string {
	if value, ok := params[key]; ok {
		return value
	}
	return defaultValue
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
