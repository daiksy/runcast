package main

import (
	"testing"
)

func TestAssessRunningCondition(t *testing.T) {
	tests := []struct {
		name           string
		temp           float64
		apparentTemp   float64
		humidity       float64
		windSpeed      float64
		precipitation  float64
		weatherCode    int
		expectedLevel  string
		minScore       int
		maxScore       int
	}{
		{
			name:          "Perfect conditions",
			temp:          20,
			apparentTemp:  20,
			humidity:      50,
			windSpeed:     2,
			precipitation: 0,
			weatherCode:   1,
			expectedLevel: "最高",
			minScore:      80,
			maxScore:      100,
		},
		{
			name:          "Hot weather",
			temp:          35,
			apparentTemp:  40,
			humidity:      80,
			windSpeed:     1,
			precipitation: 0,
			weatherCode:   1,
			expectedLevel: "普通",
			minScore:      40,
			maxScore:      60,
		},
		{
			name:          "Cold weather",
			temp:          0,
			apparentTemp:  -5,
			humidity:      60,
			windSpeed:     3,
			precipitation: 0,
			weatherCode:   1,
			expectedLevel: "最高",
			minScore:      80,
			maxScore:      100,
		},
		{
			name:          "Rainy conditions",
			temp:          20,
			apparentTemp:  20,
			humidity:      90,
			windSpeed:     8,
			precipitation: 5,
			weatherCode:   63,
			expectedLevel: "注意",
			minScore:      20,
			maxScore:      50,
		},
		{
			name:          "Thunderstorm",
			temp:          25,
			apparentTemp:  28,
			humidity:      85,
			windSpeed:     12,
			precipitation: 10,
			weatherCode:   95,
			expectedLevel: "危険",
			minScore:      0,
			maxScore:      20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition := assessRunningCondition(
				tt.temp,
				tt.apparentTemp,
				tt.humidity,
				tt.windSpeed,
				tt.precipitation,
				tt.weatherCode,
			)

			if condition.Level != tt.expectedLevel {
				t.Errorf("Expected level %s, got %s", tt.expectedLevel, condition.Level)
			}

			if condition.Score < tt.minScore || condition.Score > tt.maxScore {
				t.Errorf("Expected score between %d and %d, got %d", tt.minScore, tt.maxScore, condition.Score)
			}

			// Verify that warnings are present for dangerous conditions
			if tt.weatherCode == 95 {
				found := false
				for _, warning := range condition.Warnings {
					if len(warning) > 0 {
						found = true
						break
					}
				}
				if !found {
					t.Error("Expected warnings for thunderstorm conditions")
				}
			}

			// Verify clothing recommendations are present
			if len(condition.Clothing) == 0 {
				t.Error("Expected clothing recommendations")
			}
		})
	}
}

func TestGetWindDirection(t *testing.T) {
	tests := []struct {
		degrees  float64
		expected string
	}{
		{0, "北"},
		{45, "北東"},
		{90, "東"},
		{135, "南東"},
		{180, "南"},
		{225, "南西"},
		{270, "西"},
		{315, "北西"},
		{360, "北"},
		{22.5, "北北東"},
		{67.5, "東北東"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := getWindDirection(tt.degrees)
			if result != tt.expected {
				t.Errorf("Expected %s for %f degrees, got %s", tt.expected, tt.degrees, result)
			}
		})
	}
}

func TestGetRunningTempIcon(t *testing.T) {
	tests := []struct {
		temp     float64
		expected string
	}{
		{35, "🔥 "},
		{28, "🌡️  "},
		{20, "👌 "},
		{10, "🧥 "},
		{0, "❄️  "},
	}

	for _, tt := range tests {
		t.Run("temp_icon", func(t *testing.T) {
			result := getRunningTempIcon(tt.temp)
			if result != tt.expected {
				t.Errorf("Expected %s for %.1f°C, got %s", tt.expected, tt.temp, result)
			}
		})
	}
}

func TestRunningConditionWarnings(t *testing.T) {
	// Test high temperature warnings
	condition := assessRunningCondition(35, 40, 50, 2, 0, 1)
	if len(condition.Warnings) == 0 {
		t.Error("Expected warnings for high temperature")
	}

	// Test high humidity warnings
	condition = assessRunningCondition(25, 25, 90, 2, 0, 1)
	found := false
	for _, warning := range condition.Warnings {
		if len(warning) > 0 {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected warnings for high humidity")
	}

	// Test precipitation warnings
	condition = assessRunningCondition(20, 20, 50, 2, 5, 63)
	found = false
	for _, warning := range condition.Warnings {
		if len(warning) > 0 {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected warnings for precipitation")
	}
}

func TestRunningConditionClothing(t *testing.T) {
	// Test hot weather clothing
	condition := assessRunningCondition(35, 40, 50, 2, 0, 1)
	found := false
	for _, item := range condition.Clothing {
		if item == "薄手の半袖" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected light clothing recommendation for hot weather")
	}

	// Test cold weather clothing
	condition = assessRunningCondition(0, -5, 50, 2, 0, 1)
	found = false
	for _, item := range condition.Clothing {
		if item == "防寒着" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected warm clothing recommendation for cold weather")
	}
}

// Benchmark tests for running functions
func BenchmarkAssessRunningCondition(b *testing.B) {
	for i := 0; i < b.N; i++ {
		assessRunningCondition(20, 22, 60, 3, 0, 1)
	}
}

func BenchmarkGetWindDirection(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getWindDirection(180)
	}
}