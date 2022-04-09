package env

import (
	"os"
	"testing"
)

func withEnvSet(key, value string, f func()) {
	oldValue, exists := os.LookupEnv(key)
	os.Setenv(key, value)
	f()
	if exists {
		os.Setenv(key, oldValue)
	} else {
		os.Unsetenv(key)
	}
}

func TestIsEmpty(t *testing.T) {
	// Create a new channel to read from
	withEnvSet("MY_VAL", "hello", func() {
		if IsEmpty("MY_VAL") {
			t.Errorf("Expected 'MY_VAL' to not be empty, but it was")
		}
	})

	if !IsEmpty("MY_VAL") {
		t.Errorf("Expcted 'MY_OTHER_VAL' to be empty, but it wasn't")
	}
}

func TestGetDefault(t *testing.T) {
	withEnvSet("MY_VAL", "value", func() {
		if GetDefault("MY_VAL", "default") != "value" {
			t.Errorf("Expected GetDefault to return 'value' if it is present.")
		}
	})

	if GetDefault("MY_VAL", "default") != "default" {
		t.Errorf("Expected GetDefault to return 'default' if it is not present.")
	}

	withEnvSet("MY_VAL", "", func() {
		if GetDefault("MY_VAL", "default") != "default" {
			t.Errorf("Expected GetDefault to return 'default' if it is empty.")
		}
	})
}

func TestGetOrPanicActuallyPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected GetOrPanic to panic if env var is not present.")
		}
	}()

	GetOrPanic("NOTHING")
}

func TestGetOrPanicWhenExists(t *testing.T) {
	withEnvSet("MY_VAL", "", func() {
		GetOrPanic("MY_VAL")
	})
}

func TestGetBoolean(t *testing.T) {
	tests := map[string]bool{
		"true":           true,
		"TRUE":           true,
		"yes":            true,
		"YES":            true,
		"1":              true,
		"false":          false,
		"FALSE":          false,
		"no":             false,
		"NO":             false,
		"0":              false,
		"something else": false,
		"":               false,
	}

	for value, expected := range tests {
		withEnvSet("MY_VAL", value, func() {
			if GetBoolean("MY_VAL") != expected {
				t.Errorf("Expected GetBoolean to return %v, but it returned %v", expected, GetBoolean("MY_VAL"))
			}
		})
	}
}
