package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerHost           string
	ServerPort           string
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
}

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

func (cfg *Config) loadFlags() {
	const (
		defaultAddress     = "localhost:8080"
		defaultPostgresURI = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
		defaultAccrualURI  = "http://localhost:8081"
	)

	// Load environment variables
	cfg.RunAddress = getEnv("RUN_ADDRESS", defaultAddress)
	cfg.DatabaseURI = getEnv("DATABASE_URI", defaultPostgresURI)
	cfg.AccrualSystemAddress = getEnv("ACCRUAL_SYSTEM_ADDRESS", defaultAccrualURI)

	// Define command-line flags
	flag.StringVar(&cfg.RunAddress, "a", cfg.RunAddress, "run address (default: localhost:8080)")
	flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "database URI (default: postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable)")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", cfg.AccrualSystemAddress, "accrual system address (default: localhost:8081)")
	flag.Parse()
}
