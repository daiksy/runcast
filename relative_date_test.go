package main

import (
	"testing"
	"weather-cli/internal/types"
	"weather-cli/internal/weather"
)

func TestGetDateDisplayName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"today", "今日の"},
		{"tomorrow", "明日の"},
		{"day-after-tomorrow", "明後日の"},
		{"invalid", "invalidの"},
		{"", "の"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := weather.GetDateDisplayName(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestExtractDateBasedWeather(t *testing.T) {
	// Create mock weather data
	weatherData := &types.WeatherData{}
	weatherData.Current.Temperature = 25.0
	weatherData.Current.ApparentTemp = 27.0
	weatherData.Current.Humidity = 60
	weatherData.Current.WindSpeed = 3.0
	weatherData.Current.WindDirection = 180.0
	weatherData.Current.WeatherCode = 1
	weatherData.Current.Precipitation = 0.0
	
	weatherData.Daily.Time = []string{"2025-07-04", "2025-07-05", "2025-07-06"}
	weatherData.Daily.TemperatureMax = []float64{30.0, 32.0, 28.0}
	weatherData.Daily.TemperatureMin = []float64{20.0, 22.0, 18.0}
	weatherData.Daily.WeatherCode = []int{1, 2, 3}
	weatherData.Daily.WindSpeedMax = []float64{5.0, 6.0, 4.0}
	weatherData.Daily.WindGustMax = []float64{7.0, 8.0, 6.0}
	weatherData.Daily.PrecipitationSum = []float64{0.0, 1.0, 2.0}
	weatherData.Daily.SunriseTime = []string{"05:30", "05:31", "05:32"}
	weatherData.Daily.SunsetTime = []string{"19:00", "19:01", "19:02"}
	weatherData.Daily.DaylightDuration = []float64{13.5, 13.5, 13.5}
	weatherData.Daily.SunshineDuration = []float64{8.0, 6.0, 10.0}
	weatherData.Daily.UvIndexMax = []float64{7.0, 6.0, 8.0}
	weatherData.Daily.UvIndexClearSkyMax = []float64{9.0, 8.0, 10.0}
	weatherData.Daily.PrecipitationHours = []float64{0.0, 2.0, 1.0}
	weatherData.Daily.PrecipitationProbabilityMax = []float64{10.0, 80.0, 30.0}
	
	weatherData.Hourly.Time = make([]string, 72)  // 3 days * 24 hours
	weatherData.Hourly.Temperature = make([]float64, 72)
	weatherData.Hourly.ApparentTemp = make([]float64, 72)
	weatherData.Hourly.Humidity = make([]int, 72)
	weatherData.Hourly.WindSpeed = make([]float64, 72)
	weatherData.Hourly.WindDirection = make([]float64, 72)
	weatherData.Hourly.WeatherCode = make([]int, 72)
	weatherData.Hourly.Precipitation = make([]float64, 72)

	// Fill hourly data for 3 days
	for i := 0; i < 72; i++ {
		day := i / 24
		hour := i % 24
		weatherData.Hourly.Time[i] = "2025-07-0" + string(rune('4'+day)) + "T" + formatHour(hour) + ":00"
		weatherData.Hourly.Temperature[i] = 20.0 + float64(day*2) + float64(hour)*0.5
		weatherData.Hourly.ApparentTemp[i] = 22.0 + float64(day*2) + float64(hour)*0.5
		weatherData.Hourly.Humidity[i] = 50 + day*5 + hour
		weatherData.Hourly.WindSpeed[i] = 2.0 + float64(day)*0.5
		weatherData.Hourly.WindDirection[i] = float64(hour * 15)
		weatherData.Hourly.WeatherCode[i] = 1 + day
		weatherData.Hourly.Precipitation[i] = float64(day)
	}

	// Test dayOffset = 0 (today)
	todayWeather := weather.ExtractDateBasedWeather(weatherData, 0)
	if todayWeather != weatherData {
		t.Error("dayOffset 0 should return original weather data")
	}

	// Test dayOffset = 1 (tomorrow)
	tomorrowWeather := weather.ExtractDateBasedWeather(weatherData, 1)
	if len(tomorrowWeather.Daily.Time) != 1 {
		t.Errorf("Expected 1 daily entry for tomorrow, got %d", len(tomorrowWeather.Daily.Time))
	}
	if tomorrowWeather.Daily.Time[0] != "2025-07-05" {
		t.Errorf("Expected tomorrow date 2025-07-05, got %s", tomorrowWeather.Daily.Time[0])
	}
	if len(tomorrowWeather.Hourly.Time) != 24 {
		t.Errorf("Expected 24 hourly entries for tomorrow, got %d", len(tomorrowWeather.Hourly.Time))
	}

	// Test dayOffset = 2 (day-after-tomorrow)
	dayAfterWeather := weather.ExtractDateBasedWeather(weatherData, 2)
	if len(dayAfterWeather.Daily.Time) != 1 {
		t.Errorf("Expected 1 daily entry for day-after-tomorrow, got %d", len(dayAfterWeather.Daily.Time))
	}
	if dayAfterWeather.Daily.Time[0] != "2025-07-06" {
		t.Errorf("Expected day-after-tomorrow date 2025-07-06, got %s", dayAfterWeather.Daily.Time[0])
	}

	// Verify hourly data extraction for specific day
	firstHourTomorrow := tomorrowWeather.Hourly.Temperature[0]
	expectedTemp := 20.0 + float64(1*2) + float64(0)*0.5 // day=1, hour=0
	if firstHourTomorrow != expectedTemp {
		t.Errorf("Expected first hour temperature %.1f, got %.1f", expectedTemp, firstHourTomorrow)
	}
}

func TestRelativeDateValidation(t *testing.T) {
	validDates := []string{"today", "tomorrow", "day-after-tomorrow"}
	invalidDates := []string{"yesterday", "next-week", "invalid", ""}

	for _, validDate := range validDates {
		displayName := weather.GetDateDisplayName(validDate)
		if displayName == "" {
			t.Errorf("Valid date %s should have a display name", validDate)
		}
	}

	for _, invalidDate := range invalidDates {
		displayName := weather.GetDateDisplayName(invalidDate)
		if displayName != "" && invalidDate != "" {
			t.Errorf("Invalid date %s should not have a display name, got %s", invalidDate, displayName)
		}
	}
}

func TestDateOffsetCalculation(t *testing.T) {
	tests := []struct {
		dateSpec string
		expected int
	}{
		{"today", 0},
		{"tomorrow", 1},
		{"day-after-tomorrow", 2},
	}

	for _, tt := range tests {
		t.Run(tt.dateSpec, func(t *testing.T) {
			var dayOffset int
			switch tt.dateSpec {
			case "today":
				dayOffset = 0
			case "tomorrow":
				dayOffset = 1
			case "day-after-tomorrow":
				dayOffset = 2
			}

			if dayOffset != tt.expected {
				t.Errorf("Expected offset %d for %s, got %d", tt.expected, tt.dateSpec, dayOffset)
			}
		})
	}
}

func TestExtractDateBasedWeatherEdgeCases(t *testing.T) {
	// Test with empty weather data
	emptyWeather := &types.WeatherData{}
	result := weather.ExtractDateBasedWeather(emptyWeather, 1)
	if len(result.Daily.Time) != 0 {
		t.Error("Empty weather data should result in empty daily data")
	}

	// Test with dayOffset beyond available data
	testWeather := &types.WeatherData{}
	testWeather.Daily.Time = []string{"2025-07-04"}
	testWeather.Daily.TemperatureMax = []float64{30.0}
	testWeather.Daily.TemperatureMin = []float64{20.0}
	testWeather.Daily.WeatherCode = []int{1}
	testWeather.Daily.WindSpeedMax = []float64{5.0}
	testWeather.Daily.WindGustMax = []float64{7.0}
	testWeather.Daily.PrecipitationSum = []float64{0.0}
	testWeather.Daily.SunriseTime = []string{"05:30"}
	testWeather.Daily.SunsetTime = []string{"19:00"}
	testWeather.Daily.DaylightDuration = []float64{13.5}
	testWeather.Daily.SunshineDuration = []float64{8.0}
	testWeather.Daily.UvIndexMax = []float64{7.0}
	testWeather.Daily.UvIndexClearSkyMax = []float64{9.0}
	testWeather.Daily.PrecipitationHours = []float64{0.0}
	testWeather.Daily.PrecipitationProbabilityMax = []float64{10.0}
	testWeather.Hourly.Time = make([]string, 24)

	result = weather.ExtractDateBasedWeather(testWeather, 5) // Beyond available data
	if len(result.Daily.Time) != 1 {
		t.Error("dayOffset beyond available data should return original weather data")
	}
}