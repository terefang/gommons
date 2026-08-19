package xenv

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Get returns the value of the environment variable named by key.
// It returns an empty string when the variable is unset.
func Get(key string) string {
	return os.Getenv(key)
}

// Exists reports whether the environment variable named by key is set.
// A variable set to an empty string still exists.
func Exists(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}

// GetDefault returns the value of the environment variable named by key.
// It returns fallback when the variable is unset.
func GetDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Require returns the value of the environment variable named by key.
// It returns an error when the variable is unset.
func Require(key string) (string, error) {
	if value, ok := os.LookupEnv(key); ok {
		return value, nil
	}
	return "", fmt.Errorf("xenv: required environment variable %q is not set", key)
}

// GetInt returns the parsed int value of the environment variable named by key.
// It returns fallback when the variable is unset or cannot be parsed as an int.
func GetInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return fallback
	}
	return parsed
}

// GetBool returns the parsed bool value of the environment variable named by key.
// It returns fallback when the variable is unset or cannot be parsed as a bool.
func GetBool(key string, fallback bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	parsed, err := strconv.ParseBool(strings.TrimSpace(value))
	if err != nil {
		return fallback
	}
	return parsed
}

// GetDuration returns the parsed duration value of the environment variable named by key.
// It returns fallback when the variable is unset or cannot be parsed by time.ParseDuration.
func GetDuration(key string, fallback time.Duration) time.Duration {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	parsed, err := time.ParseDuration(strings.TrimSpace(value))
	if err != nil {
		return fallback
	}
	return parsed
}

// GetSlice splits the environment variable named by key by sep.
// It returns fallback when the variable is unset.
// Empty items are removed after trimming surrounding spaces.
func GetSlice(key, sep string, fallback []string) []string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	if sep == "" {
		sep = ","
	}

	parts := strings.Split(value, sep)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}
