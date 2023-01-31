package apiserver

import "os"

type Config struct {
	BindAddr    string
	LogLevel    string
	DatabaseURL string
	SessioKey   string
}

// NewConfig return new Config instance
func NewConfig() *Config {
	return &Config{
		BindAddr:    os.Getenv("BIND_ADDR"),
		LogLevel:    os.Getenv("LOG_LEVEL"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		SessioKey:   os.Getenv("SESSION_KEY"),
	}
}
