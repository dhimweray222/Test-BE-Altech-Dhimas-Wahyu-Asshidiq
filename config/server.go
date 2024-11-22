package config

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type serverConfig struct {
	URI  string
	Port string
	Host string
}

func NewServerConfig() serverConfig {
	return serverConfig{
		URI:  os.Getenv("SERVER_URI"),
		Port: os.Getenv("SERVER_PORT"),
		Host: fmt.Sprintf("%s:%s", os.Getenv("SERVER_URI"), os.Getenv("SERVER_PORT")),
	}
}
