package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestLoad(t *testing.T) {
	// Reset viper for clean test
	viper.Reset()

	// Test loading with defaults
	config, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test default values
	if config.Server.Port != "8080" {
		t.Errorf("Expected default port 8080, got %s", config.Server.Port)
	}

	if config.Database.Host != "localhost" {
		t.Errorf("Expected default database host localhost, got %s", config.Database.Host)
	}

	if config.JWT.ExpireTime != 168 {
		t.Errorf("Expected default JWT expire time 168, got %d", config.JWT.ExpireTime)
	}
}

func TestLoadWithEnvVars(t *testing.T) {
	// Reset viper for clean test
	viper.Reset()

	// Set environment variables
	os.Setenv("SERVER_PORT", "9000")
	os.Setenv("DATABASE_HOST", "testhost")
	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("DATABASE_HOST")
	}()

	config, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test environment variable override
	if config.Server.Port != "9000" {
		t.Errorf("Expected port 9000 from env var, got %s", config.Server.Port)
	}

	if config.Database.Host != "testhost" {
		t.Errorf("Expected database host testhost from env var, got %s", config.Database.Host)
	}
}

func TestGetDatabaseURL(t *testing.T) {
	config := &Config{
		Database: DatabaseConfig{
			Username: "testuser",
			Password: "testpass",
			Host:     "testhost",
			Port:     3306,
			Database: "testdb",
		},
	}

	expected := "testuser:testpass@tcp(testhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	actual := config.GetDatabaseURL()

	if actual != expected {
		t.Errorf("Expected database URL %s, got %s", expected, actual)
	}
}

func TestGetServerAddress(t *testing.T) {
	config := &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: "8080",
		},
	}

	expected := "0.0.0.0:8080"
	actual := config.GetServerAddress()

	if actual != expected {
		t.Errorf("Expected server address %s, got %s", expected, actual)
	}
}