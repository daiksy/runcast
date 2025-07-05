package main

import (
	"testing"
	"weather-cli/internal/running"
	"weather-cli/internal/types"
	"weather-cli/internal/weather"
)

func TestGetTimePeriods(t *testing.T) {
	periods := weather.GetTimePeriods()
	
	expectedPeriods := map[string]struct {
		displayName string
		startHour   int
		endHour     int
	}{
		"morning": {"早朝", 5, 9},
		"noon":    {"昼", 11, 15},
		"evening": {"夕方", 17, 19},
		"night":   {"夜", 21, 23},
	}
	
	for name, expected := range expectedPeriods {
		period, exists := periods[name]
		if !exists {
			t.Errorf("Period %s not found", name)
			continue
		}
		
		if period.DisplayName != expected.displayName {
			t.Errorf("Expected display name %s, got %s", expected.displayName, period.DisplayName)
		}
		
		if period.StartHour != expected.startHour {
			t.Errorf("Expected start hour %d, got %d", expected.startHour, period.StartHour)
		}
		
		if period.EndHour != expected.endHour {
			t.Errorf("Expected end hour %d, got %d", expected.endHour, period.EndHour)
		}
	}
}

func TestExtractHour(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"2025-07-04T05:00", "05"},
		{"2025-07-04T15:30", "15"},
		{"2025-07-04T23:45", "23"},
		{"2025-07-04T09:15:30Z", "09"},
		{"invalid", "invalid"},
		{"", ""},
		{"2025-07-04", "2025-07-04"},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := weather.ExtractHour(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestExtractTimeBasedWeather(t *testing.T) {
	// Create mock weather data
	weatherData := &types.WeatherData{}
	weatherData.Hourly.Time = make([]string, 24)
	weatherData.Hourly.Temperature = make([]float64, 24)
	weatherData.Hourly.ApparentTemp = make([]float64, 24)
	weatherData.Hourly.Humidity = make([]int, 24)
	weatherData.Hourly.WindSpeed = make([]float64, 24)
	weatherData.Hourly.WindDirection = make([]float64, 24)
	weatherData.Hourly.WeatherCode = make([]int, 24)
	weatherData.Hourly.Precipitation = make([]float64, 24)
	
	// Fill with mock data for 24 hours
	for i := 0; i < 24; i++ {
		weatherData.Hourly.Time[i] = "2025-07-04T" + formatHour(i) + ":00"
		weatherData.Hourly.Temperature[i] = 20.0 + float64(i)
		weatherData.Hourly.ApparentTemp[i] = 22.0 + float64(i)
		weatherData.Hourly.Humidity[i] = 50 + i
		weatherData.Hourly.WindSpeed[i] = 2.0 + float64(i)*0.1
		weatherData.Hourly.WindDirection[i] = float64(i * 15)
		weatherData.Hourly.WeatherCode[i] = 1
		weatherData.Hourly.Precipitation[i] = 0.0
	}
	
	// Test morning period (5-9)
	morningData := weather.ExtractTimeBasedWeather(weatherData, "morning", 1)
	expectedMorningHours := 5
	if len(morningData) != expectedMorningHours {
		t.Errorf("Expected %d morning hours, got %d", expectedMorningHours, len(morningData))
	}
	
	// Test first morning hour
	if len(morningData) > 0 {
		firstMorning := morningData[0]
		if firstMorning.Temperature != 25.0 { // 20 + 5 (hour 5)
			t.Errorf("Expected temperature 25.0, got %.1f", firstMorning.Temperature)
		}
	}
	
	// Test evening period (17-19)
	eveningData := weather.ExtractTimeBasedWeather(weatherData, "evening", 1)
	expectedEveningHours := 3
	if len(eveningData) != expectedEveningHours {
		t.Errorf("Expected %d evening hours, got %d", expectedEveningHours, len(eveningData))
	}
	
	// Test invalid time period
	invalidData := weather.ExtractTimeBasedWeather(weatherData, "invalid", 1)
	if invalidData != nil {
		t.Error("Expected nil for invalid time period")
	}
}

func TestTimeBasedWeatherValidation(t *testing.T) {
	validTimes := []string{"morning", "noon", "evening", "night"}
	invalidTimes := []string{"afternoon", "midnight", "dawn", "invalid", ""}
	
	for _, validTime := range validTimes {
		periods := weather.GetTimePeriods()
		if _, exists := periods[validTime]; !exists {
			t.Errorf("Valid time %s should exist in periods", validTime)
		}
	}
	
	for _, invalidTime := range invalidTimes {
		periods := weather.GetTimePeriods()
		if _, exists := periods[invalidTime]; exists {
			t.Errorf("Invalid time %s should not exist in periods", invalidTime)
		}
	}
}

func TestTimeBasedRunningRecommendation(t *testing.T) {
	// Create weather data with varying conditions for different hours
	testWeather := &types.WeatherData{}
	testWeather.Hourly.Time = make([]string, 24)
	testWeather.Hourly.Temperature = make([]float64, 24)
	testWeather.Hourly.ApparentTemp = make([]float64, 24)
	testWeather.Hourly.Humidity = make([]int, 24)
	testWeather.Hourly.WindSpeed = make([]float64, 24)
	testWeather.Hourly.WindDirection = make([]float64, 24)
	testWeather.Hourly.WeatherCode = make([]int, 24)
	testWeather.Hourly.Precipitation = make([]float64, 24)
	
	// Set up data where hour 6 has better conditions than hour 8
	for i := 0; i < 24; i++ {
		testWeather.Hourly.Time[i] = "2025-07-04T" + formatHour(i) + ":00"
		testWeather.Hourly.Temperature[i] = 20.0
		testWeather.Hourly.ApparentTemp[i] = 22.0
		testWeather.Hourly.Humidity[i] = 60
		testWeather.Hourly.WindSpeed[i] = 2.0
		testWeather.Hourly.WindDirection[i] = 180.0
		testWeather.Hourly.WeatherCode[i] = 1
		testWeather.Hourly.Precipitation[i] = 0.0
	}
	
	// Make hour 8 have worse conditions (higher humidity)
	testWeather.Hourly.Humidity[8] = 90 // High humidity should lower score
	
	morningData := weather.ExtractTimeBasedWeather(testWeather, "morning", 1)
	if len(morningData) == 0 {
		t.Fatal("No morning data extracted")
	}
	
	// Find best and worst conditions
	bestScore := -1
	worstScore := 101
	
	for _, data := range morningData {
		condition := running.AssessRunningCondition(
			data.Temperature,
			data.ApparentTemp,
			float64(data.Humidity),
			data.WindSpeed,
			data.Precipitation,
			data.WeatherCode,
		)
		
		if condition.Score > bestScore {
			bestScore = condition.Score
		}
		if condition.Score < worstScore {
			worstScore = condition.Score
		}
	}
	
	// There should be variation in scores due to different humidity levels
	if bestScore == worstScore {
		t.Error("Expected variation in running condition scores across different hours")
	}
}

// Helper function to format hour as two digits
func formatHour(hour int) string {
	if hour < 10 {
		return "0" + string(rune('0'+hour))
	}
	return string(rune('0'+hour/10)) + string(rune('0'+hour%10))
}