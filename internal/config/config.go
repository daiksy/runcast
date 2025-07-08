package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"runcast/internal/types"
)

// Config represents the configuration file structure
type Config struct {
	Locations map[string]types.CityCoordinate `toml:"locations"`
}

// LoadConfig loads configuration from available config files
func LoadConfig() (*Config, error) {
	configPaths := getConfigPaths()
	
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			return loadConfigFromFile(path)
		}
	}
	
	// Return empty config if no config file found
	return &Config{
		Locations: make(map[string]types.CityCoordinate),
	}, nil
}

// getConfigPaths returns possible config file paths in order of priority
func getConfigPaths() []string {
	var paths []string
	
	// 1. Current directory
	if cwd, err := os.Getwd(); err == nil {
		paths = append(paths, filepath.Join(cwd, ".runcast.conf"))
	}
	
	// 2. Home directory
	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(home, ".runcast.conf"))
		paths = append(paths, filepath.Join(home, ".config", "runcast", "config.toml"))
	}
	
	return paths
}

// loadConfigFromFile loads configuration from a specific file
func loadConfigFromFile(path string) (*Config, error) {
	var config Config
	
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}
	
	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration in %s: %w", path, err)
	}
	
	return &config, nil
}

// validateConfig validates the configuration structure
func validateConfig(config *Config) error {
	if config.Locations == nil {
		config.Locations = make(map[string]types.CityCoordinate)
	}
	
	for name, location := range config.Locations {
		if name == "" {
			return fmt.Errorf("location name cannot be empty")
		}
		if location.Name == "" {
			return fmt.Errorf("location '%s' must have a name", name)
		}
		if location.Lat < -90 || location.Lat > 90 {
			return fmt.Errorf("location '%s' has invalid latitude: %f", name, location.Lat)
		}
		if location.Lon < -180 || location.Lon > 180 {
			return fmt.Errorf("location '%s' has invalid longitude: %f", name, location.Lon)
		}
	}
	
	return nil
}

// GetCustomLocation returns a custom location by name
func (c *Config) GetCustomLocation(name string) (*types.CityCoordinate, bool) {
	location, exists := c.Locations[name]
	if !exists {
		return nil, false
	}
	return &location, true
}

// GetCustomLocationNames returns all custom location names
func (c *Config) GetCustomLocationNames() []string {
	names := make([]string, 0, len(c.Locations))
	for name := range c.Locations {
		names = append(names, name)
	}
	return names
}