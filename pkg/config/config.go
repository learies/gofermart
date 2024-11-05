package config

// Config struct holds the configuration values for the application
type Config struct {
	ServerPort string
}

// LoadConfig loads the application configuration
func LoadConfig() *Config {
	return &Config{
		ServerPort: ":8080",
	}
}
