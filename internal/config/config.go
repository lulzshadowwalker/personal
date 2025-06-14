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

func Host() string {
	return Get("HOST", "localhost")
}

func DBHost() string {
	return Get("DB_HOST", "localhost")
}

func DBPort() string {
	return Get("DB_PORT", "5432")
}

func DBUsername() string {
	return Get("DB_USERNAME", "postgres")
}

func DBPassword() string {
	return Get("DB_PASSWORD", "postgres")
}

func DBName() string {
	return Get("DB_NAME", "personal")
}

func DBSSLMode() string {
	return Get("DB_SSLMODE", "disable")
}

func DBConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		DBHost(), DBPort(), DBUsername(), DBPassword(), DBName(), DBSSLMode())
}

func Development() bool {
	return os.Getenv("GO_ENV") == "production"
}

func Production() bool {
	return os.Getenv("GO_ENV") == "production"
}
