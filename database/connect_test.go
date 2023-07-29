package database

import (
	"github.com/spf13/viper"
	"testing"
)

func TestConnect(t *testing.T) {
	// Load configuration from file
	viper.SetConfigFile("../config/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		t.Errorf("Failed to read configuration file: %v", err)
	}

	// Connect to the database
	db, err := Connect()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}

	// Check that the database connection is valid
	err = db.Exec("SELECT 1").Error
	if err != nil {
		t.Errorf("Failed to run query on database: %v", err)
	}

	// Close the database connection
	err = Close()
	if err != nil {
		t.Errorf("Failed to close database connection: %v", err)
	}
}
