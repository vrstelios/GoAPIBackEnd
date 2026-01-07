package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"sync"
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
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

type APIConfig struct {
	BaseURL string
}

var (
	appConfig *Config
	configMu  sync.RWMutex
)

func init() {
	loadConfig()
}

// GetConfig loads config from YAML based on APP_ENV
func GetConfig() *Config {
	configMu.RLock()
	defer configMu.RUnlock()
	return appConfig
}

func loadConfig() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	filePath := getConfigFile(env)

	v := viper.New()
	v.SetConfigFile(filePath)
	v.SetConfigType("yml")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file %s: %v", filePath, err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to parse config file: %v", err)
	}

	cfg.TokenExpiration = time.Duration(time.Now().Add(7 * 24 * time.Hour).Unix())

	configMu.Lock()
	appConfig = &cfg
	configMu.Unlock()
}

func getConfigFile(env string) string {
	switch env {
	case "production":
		return "config/config-production.yml"
	/*case "docker":
	return "config/config-docker.yml"*/
	default:
		return "config/config-development.yml"
	}
}
