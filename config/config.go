package config

import (
	"github.com/goccy/go-json"
	"log"
	"os"
)

var App *Config

type Config struct {
	MasterAPIDomain string    `json:"MasterAPIDomain"`
	JWT             JWTConfig `json:"jwt"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type JWTConfig struct {
	ExpirationHours int16  `json:"expiration_hours"`
	Secret          string `json:"secret"`
}

func Load() {
	configData, err := os.ReadFile("config/config.json")
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = json.Unmarshal(configData, &config)
	if err != nil {
		log.Fatal(err)
	}

	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		App.JWT.Secret = secret
	}

	App = &config
}
