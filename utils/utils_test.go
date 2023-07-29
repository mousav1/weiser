package utils

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestInitLogger(t *testing.T) {
	// Set up any necessary test configurations
	viper.Set("logging.path", "test.log")

	// Call the function being tested
	err := InitLogger()

	// Verify the results
	if err != nil {
		t.Errorf("InitLogger returned an error: %s", err.Error())
	}

	// Clean up any resources used by the test
	os.Remove("test.log")
}

func TestLoggingFunctions(t *testing.T) {
	// Set up any necessary test configurations
	viper.Set("logging.path", "test.log")

	// Call the functions being tested
	Info("This is an info message")
	Warn("This is a warning message")
	Error("This is an error message")

	// Clean up any resources used by the test
	os.Remove("test.log")
}

func TestGenerateAndCompareHash(t *testing.T) {
	// Define test password
	password := "testPassword"

	// Generate hash
	hash, err := GenerateHash(password)
	if err != nil {
		t.Errorf("Failed to generate hash: %v", err)
	}

	// Compare hash with the original password
	err = CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		t.Errorf("Failed to compare hash and password: %v", err)
	}
}
