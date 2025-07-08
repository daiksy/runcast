package config

import (
	"os"
	"path/filepath"
	"testing"

	"runcast/internal/types"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".runcast.conf")
	
	configContent := `[locations]
home = { name = "自宅", lat = 35.6762, lon = 139.6503 }
office = { name = "会社", lat = 35.6584, lon = 139.7016 }`
	
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	
	// Change working directory to temp dir for testing
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)
	
	// Load config
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	
	// Test home location
	homeCoord, exists := config.GetCustomLocation("home")
	if !exists {
		t.Error("Expected 'home' location to exist")
	}
	if homeCoord.Name != "自宅" {
		t.Errorf("Expected name '自宅', got '%s'", homeCoord.Name)
	}
	if homeCoord.Lat != 35.6762 {
		t.Errorf("Expected lat 35.6762, got %f", homeCoord.Lat)
	}
	if homeCoord.Lon != 139.6503 {
		t.Errorf("Expected lon 139.6503, got %f", homeCoord.Lon)
	}
	
	// Test office location
	officeCoord, exists := config.GetCustomLocation("office")
	if !exists {
		t.Error("Expected 'office' location to exist")
	}
	if officeCoord.Name != "会社" {
		t.Errorf("Expected name '会社', got '%s'", officeCoord.Name)
	}
	
	// Test non-existent location
	_, exists = config.GetCustomLocation("nonexistent")
	if exists {
		t.Error("Expected 'nonexistent' location to not exist")
	}
	
	// Test custom location names
	names := config.GetCustomLocationNames()
	if len(names) != 2 {
		t.Errorf("Expected 2 custom locations, got %d", len(names))
	}
}

func TestLoadConfigNoFile(t *testing.T) {
	// Test loading config when no file exists
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)
	
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig should not fail when no config file exists: %v", err)
	}
	
	if len(config.Locations) != 0 {
		t.Errorf("Expected empty locations, got %d", len(config.Locations))
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError bool
	}{
		{
			name: "valid config",
			config: Config{
				Locations: map[string]types.CityCoordinate{
					"home": {"自宅", 35.6762, 139.6503},
				},
			},
			expectError: false,
		},
		{
			name: "invalid latitude",
			config: Config{
				Locations: map[string]types.CityCoordinate{
					"invalid": {"無効", 91.0, 139.6503}, // latitude > 90
				},
			},
			expectError: true,
		},
		{
			name: "invalid longitude",
			config: Config{
				Locations: map[string]types.CityCoordinate{
					"invalid": {"無効", 35.6762, 181.0}, // longitude > 180
				},
			},
			expectError: true,
		},
		{
			name: "empty location name",
			config: Config{
				Locations: map[string]types.CityCoordinate{
					"test": {"", 35.6762, 139.6503}, // empty name
				},
			},
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(&tt.config)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}