package main

import (
	"testing"
)

// Integration tests for API calls
// These tests make actual API calls, so they require internet connection
// and may be slower than unit tests

func TestGetWeatherIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tests := []struct {
		name string
		lat  float64
		lon  float64
		days int
	}{
		{
			name: "Tokyo current weather",
			lat:  35.6762,
			lon:  139.6503,
			days: 0,
		},
		{
			name: "Tokyo 3-day forecast",
			lat:  35.6762,
			lon:  139.6503,
			days: 3,
		},
		{
			name: "Osaka 7-day forecast",
			lat:  34.6937,
			lon:  135.5023,
			days: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weather, err := getWeather(tt.lat, tt.lon, tt.days)
			if err != nil {
				t.Fatalf("API call failed: %v", err)
			}

			// Verify current weather data is present
			if weather.Current.Temperature == 0 && weather.Current.Humidity == 0 {
				t.Error("Current weather data appears to be empty")
			}

			if tt.days > 0 {
				// Verify daily forecast data is present
				if len(weather.Daily.Time) == 0 {
					t.Error("Daily forecast data is missing")
				}

				if len(weather.Daily.Time) != len(weather.Daily.TemperatureMax) ||
					len(weather.Daily.Time) != len(weather.Daily.TemperatureMin) ||
					len(weather.Daily.Time) != len(weather.Daily.WeatherCode) {
					t.Error("Daily forecast data arrays have inconsistent lengths")
				}

				// Verify we got at least the requested number of days (or fewer if API limits)
				expectedDays := tt.days
				if len(weather.Daily.Time) < expectedDays {
					t.Logf("Warning: Expected %d days, got %d days", expectedDays, len(weather.Daily.Time))
				}
			}

			// Verify weather codes are valid
			if weather.Current.WeatherCode < 0 {
				t.Errorf("Invalid weather code: %d", weather.Current.WeatherCode)
			}

			// Verify temperature seems reasonable (between -50 and 60 Celsius)
			if weather.Current.Temperature < -50 || weather.Current.Temperature > 60 {
				t.Errorf("Temperature seems unreasonable: %.1f°C", weather.Current.Temperature)
			}

			// Verify humidity is in valid range (0-100%)
			if weather.Current.Humidity < 0 || weather.Current.Humidity > 100 {
				t.Errorf("Humidity out of range: %d%%", weather.Current.Humidity)
			}
		})
	}
}

func TestGetWeatherWithInvalidCoordinates(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test with coordinates that are way out of range
	_, err := getWeather(999.0, 999.0, 0)
	
	// The API might still return data or give an error
	// We mainly want to ensure our code doesn't crash
	if err != nil {
		t.Logf("Expected behavior: API returned error for invalid coordinates: %v", err)
	}
}

// Helper function to test the full workflow
func TestFullWorkflowIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test the complete workflow: city -> coordinates -> weather
	city := "tokyo"
	
	coord, err := getCityCoordinate(city)
	if err != nil {
		t.Fatalf("Failed to get coordinates for %s: %v", city, err)
	}

	weather, err := getWeather(coord.Lat, coord.Lon, 1)
	if err != nil {
		t.Fatalf("Failed to get weather for %s: %v", city, err)
	}

	// Verify we got reasonable data
	if weather.Current.Temperature == 0 && weather.Current.Humidity == 0 {
		t.Error("Weather data appears to be empty")
	}

	// Test weather description conversion
	description := getWeatherDescription(weather.Current.WeatherCode)
	if description == "" {
		t.Error("Weather description is empty")
	}

	t.Logf("Successfully retrieved weather for %s: %.1f°C, %s", 
		coord.Name, weather.Current.Temperature, description)
}