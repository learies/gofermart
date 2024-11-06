package config

import (
	"flag"
	"os"
)

// Config struct holds the configuration values for the application
type Config struct {
	ServerHost  string
	ServerPort  string
	RunAddress  string
	DatabaseURI string
}

// NewConfig creates a new Config instance and loads the configuration
func NewConfig() *Config {
	cfg := &Config{}
	cfg.loadFlags()
	return cfg
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// loadFlags parses command-line flags and populates the Config struct
func (cfg *Config) loadFlags() {
	const (
		defaultAddress     = "localhost:8080"
		defaultPostgresURI = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	)

	// Load environment variables
	cfg.RunAddress = getEnv("RUN_ADDRESS", defaultAddress)
	cfg.DatabaseURI = getEnv("DATABASE_URI", defaultPostgresURI)

	// Define command-line flags
	flag.StringVar(&cfg.RunAddress, "a", cfg.RunAddress, "run address (default: localhost:8080)")
	flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "database URI (default: postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable)")
	flag.Parse()
}
