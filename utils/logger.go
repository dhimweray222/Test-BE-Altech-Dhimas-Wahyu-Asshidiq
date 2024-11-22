package utils

import (
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

func NewLogger() *zap.SugaredLogger {
	// Determine if the environment is production or not
	isProduction, err := strconv.ParseBool(os.Getenv("IS_PRODUCTION"))
	if err != nil {
		// Default to development mode if the environment variable is not set or invalid
		isProduction = false
	}

	var logger *zap.Logger
	if isProduction {
		logger, err = zap.NewProduction() // Production logger with structured logs
	} else {
		logger, err = zap.NewDevelopment() // Development logger with human-readable logs
	}

	if err != nil {
		// Fall back to a no-op logger in case of initialization failure
		panic("Failed to initialize logger: " + err.Error())
	}

	// Return a sugared logger for easier usage
	return logger.Sugar()
}
