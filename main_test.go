package main

import (
	"testing"
	"weather-cli/internal/weather"
)

func TestGetCityCoordinate(t *testing.T) {
	tests := []struct {
		name          string
		city          string
		expectedName  string
		expectedLat   float64
		expectedLon   float64
		expectError   bool
	}{
		{
			name:         "Tokyo",
			city:         "tokyo",
			expectedName: "東京",
			expectedLat:  35.6762,
			expectedLon:  139.6503,
			expectError:  false,
		},
		{
			name:         "Osaka",
			city:         "osaka",
			expectedName: "大阪",
			expectedLat:  34.6937,
			expectedLon:  135.5023,
			expectError:  false,
		},
		{
			name:         "Kyoto",
			city:         "kyoto",
			expectedName: "京都",
			expectedLat:  35.0116,
			expectedLon:  135.7681,
			expectError:  false,
		},
		{
			name:        "Invalid city",
			city:        "invalid",
			expectError: true,
		},
		{
			name:        "Empty city",
			city:        "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coord, err := weather.GetCityCoordinate(tt.city)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for city %s, but got none", tt.city)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for city %s: %v", tt.city, err)
				return
			}

			if coord.Name != tt.expectedName {
				t.Errorf("Expected name %s, got %s", tt.expectedName, coord.Name)
			}

			if coord.Lat != tt.expectedLat {
				t.Errorf("Expected latitude %f, got %f", tt.expectedLat, coord.Lat)
			}

			if coord.Lon != tt.expectedLon {
				t.Errorf("Expected longitude %f, got %f", tt.expectedLon, coord.Lon)
			}
		})
	}
}

func TestGetWeatherDescription(t *testing.T) {
	tests := []struct {
		code        int
		expected    string
	}{
		{0, "快晴"},
		{1, "晴れ"},
		{2, "一部曇り"},
		{3, "曇り"},
		{45, "霧"},
		{51, "弱い霧雨"},
		{61, "弱い雨"},
		{71, "弱い雪"},
		{95, "雷雨"},
		{999, "不明"}, // Unknown code
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := weather.GetWeatherDescription(tt.code)
			if result != tt.expected {
				t.Errorf("Expected %s for code %d, got %s", tt.expected, tt.code, result)
			}
		})
	}
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid date format",
			input:    "2025-07-04",
			expected: "07月04日",
		},
		{
			name:     "Another valid date",
			input:    "2025-12-25",
			expected: "12月25日",
		},
		{
			name:     "Date with time",
			input:    "2025-01-01T00:00:00Z",
			expected: "01月01日",
		},
		{
			name:     "Short string",
			input:    "2025",
			expected: "2025",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Invalid format",
			input:    "invalid-date",
			expected: "invalid-date",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := weather.FormatDate(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// Benchmark tests
func BenchmarkGetCityCoordinate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		weather.GetCityCoordinate("tokyo")
	}
}

func BenchmarkGetWeatherDescription(b *testing.B) {
	for i := 0; i < b.N; i++ {
		weather.GetWeatherDescription(1)
	}
}

func BenchmarkFormatDate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		weather.FormatDate("2025-07-04")
	}
}