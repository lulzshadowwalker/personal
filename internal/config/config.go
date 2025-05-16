package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		slog.Error("failed to load .env file", "error", err)
	}
}

// Get retrieves the value of the environment variable with the given key.
// If the variable is not set, it returns the provided fallback value.
func Get(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}

	return value
}

func MustGet(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		panic(fmt.Errorf("env variable %s not set", key))
	}

	return value
}

func Port() string {
	return Get("PORT", "8080")
}
