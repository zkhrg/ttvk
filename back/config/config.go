package config

import "os"

type Config struct {
	DatabaseURL string
	Port        string
}

func LoadConfig() *Config {
	return &Config{
		DatabaseURL: os.Getenv("DB_URL"),
		Port:        os.Getenv("PORT"),
	}
}
