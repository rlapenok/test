package config

import (
	"os"

	"github.com/gookit/slog"
)

type Config struct {
	Token string
	Port  string
}

func New() *Config {
	//ENV для токена тг бота
	token := os.Getenv("TOKEN")
	if token == "" {
		slog.Fatal("ENV TOKEN not set")
		os.Exit(1)
	}
	//ENV для gRPC сервера
	port := os.Getenv("PORT")
	if token == "" {
		slog.Fatal("ENV PORT  not set")
		os.Exit(1)
	}
	config := Config{
		Token: token,
		Port:  port,
	}
	return &config
}
