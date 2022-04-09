package env

import (
	"fmt"
	"os"
	"strings"
)

// IsEmpty returns true if either the environment variable does not exist, or
// it does exist, but is empty.
func IsEmpty(key string) bool {
	return os.Getenv(key) == ""
}

// GetOrPanic wraps os.LookupEnv. It will panic if an environment variable is
// not present. It will not panic if the environment variable exists, but is
// empty.
func GetOrPanic(key string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		panic(fmt.Errorf("'%s' does not exist in environment variables", key))
	}
	return val
}

// GetDefault will return the default value if the given environment variable
// is empty. Otherwise it returns the environment variable.
func GetDefault(key, deflt string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return deflt
}

// GetBoolean returns true if the environment variable is set to "true", "yes", or "1"
// regardless of case.
func GetBoolean(key string) bool {
	val := strings.ToLower(os.Getenv(key))
	return val == "true" || val == "yes" || val == "1"
}
