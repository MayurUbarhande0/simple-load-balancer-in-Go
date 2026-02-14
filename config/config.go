package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config represents the load balancer configuration
type Config struct {
	Port                string        `json:"port"`
	HealthCheckInterval time.Duration `json:"health_check_interval"`
	ReadTimeout         time.Duration `json:"read_timeout"`
	WriteTimeout        time.Duration `json:"write_timeout"`
	IdleTimeout         time.Duration `json:"idle_timeout"`
	Backends            []string      `json:"backends"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Port:                ":8080",
		HealthCheckInterval: 10 * time.Second,
		ReadTimeout:         15 * time.Second,
		WriteTimeout:        15 * time.Second,
		IdleTimeout:         60 * time.Second,
		Backends: []string{
			"http://localhost:8081",
			"http://localhost:8082",
		},
	}
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filename string) (*Config, error) {
	config := DefaultConfig()

	// If file doesn't exist, return default config
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return config, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return config, nil
}

// SaveConfig saves configuration to a JSON file
func SaveConfig(filename string, config *Config) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}
