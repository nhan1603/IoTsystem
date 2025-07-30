package env

import "os"

// GetwithDefault gets an environment variable or returns a default value
func GetwithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
