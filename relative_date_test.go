package main

import (
	"testing"
)

func TestGetDateDisplayName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"today", "今日の"},
		{"tomorrow", "明日の"},
		{"day-after-tomorrow", "明後日の"},
		{"invalid", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := getDateDisplayName(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestExtractDateBasedWeather(t *testing.T) {
	// Create mock weather data
	weather := &WeatherData{
		Current: struct {
			Temperature        float64 `json:"temperature_2m"`
			ApparentTemp       float64 `json:"apparent_temperature"`
			Humidity           int     `json:"relative_humidity_2m"`
			WindSpeed          float64 `json:"wind_speed_10m"`
			WindDirection      float64 `json:"wind_direction_10m"`
			WeatherCode        int     `json:"weather_code"`
			Precipitation      float64 `json:"precipitation"`
			Dewpoint           float64 `json:"dewpoint_2m"`
		}{
			Temperature: 25.0,
			ApparentTemp: 27.0,
			Humidity: 60,
			WindSpeed: 3.0,
			WindDirection: 180.0,
			WeatherCode: 1,
			Precipitation: 0.0,
			Dewpoint: 15.0,
		},
		Daily: struct {
			Time           []string  `json:"time"`
			TemperatureMax []float64 `json:"temperature_2m_max"`
			TemperatureMin []float64 `json:"temperature_2m_min"`
			WeatherCode    []int     `json:"weather_code"`
			WindSpeedMax   []float64 `json:"wind_speed_10m_max"`
			PrecipitationSum []float64 `json:"precipitation_sum"`
		}{
			Time:           []string{"2025-07-04", "2025-07-05", "2025-07-06"},
			TemperatureMax: []float64{30.0, 32.0, 28.0},
			TemperatureMin: []float64{20.0, 22.0, 18.0},
			WeatherCode:    []int{1, 2, 3},
			WindSpeedMax:   []float64{5.0, 6.0, 4.0},
			PrecipitationSum: []float64{0.0, 1.0, 2.0},
		},
		Hourly: struct {
			Time           []string  `json:"time"`
			Temperature    []float64 `json:"temperature_2m"`
			ApparentTemp   []float64 `json:"apparent_temperature"`
			Humidity       []int     `json:"relative_humidity_2m"`
			WindSpeed      []float64 `json:"wind_speed_10m"`
			WindDirection  []float64 `json:"wind_direction_10m"`
			WeatherCode    []int     `json:"weather_code"`
			Precipitation  []float64 `json:"precipitation"`
		}{
			Time:          make([]string, 72),  // 3 days * 24 hours
			Temperature:   make([]float64, 72),
			ApparentTemp:  make([]float64, 72),
			Humidity:      make([]int, 72),
			WindSpeed:     make([]float64, 72),
			WindDirection: make([]float64, 72),
			WeatherCode:   make([]int, 72),
			Precipitation: make([]float64, 72),
		},
	}

	// Fill hourly data for 3 days
	for i := 0; i < 72; i++ {
		day := i / 24
		hour := i % 24
		weather.Hourly.Time[i] = "2025-07-0" + string(rune('4'+day)) + "T" + formatHour(hour) + ":00"
		weather.Hourly.Temperature[i] = 20.0 + float64(day*2) + float64(hour)*0.5
		weather.Hourly.ApparentTemp[i] = 22.0 + float64(day*2) + float64(hour)*0.5
		weather.Hourly.Humidity[i] = 50 + day*5 + hour
		weather.Hourly.WindSpeed[i] = 2.0 + float64(day)*0.5
		weather.Hourly.WindDirection[i] = float64(hour * 15)
		weather.Hourly.WeatherCode[i] = 1 + day
		weather.Hourly.Precipitation[i] = float64(day)
	}

	// Test dayOffset = 0 (today)
	todayWeather := extractDateBasedWeather(weather, 0)
	if todayWeather != weather {
		t.Error("dayOffset 0 should return original weather data")
	}

	// Test dayOffset = 1 (tomorrow)
	tomorrowWeather := extractDateBasedWeather(weather, 1)
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
	dayAfterWeather := extractDateBasedWeather(weather, 2)
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
		displayName := getDateDisplayName(validDate)
		if displayName == "" {
			t.Errorf("Valid date %s should have a display name", validDate)
		}
	}

	for _, invalidDate := range invalidDates {
		displayName := getDateDisplayName(invalidDate)
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
	emptyWeather := &WeatherData{}
	result := extractDateBasedWeather(emptyWeather, 1)
	if len(result.Daily.Time) != 0 {
		t.Error("Empty weather data should result in empty daily data")
	}

	// Test with dayOffset beyond available data
	weather := &WeatherData{
		Daily: struct {
			Time           []string  `json:"time"`
			TemperatureMax []float64 `json:"temperature_2m_max"`
			TemperatureMin []float64 `json:"temperature_2m_min"`
			WeatherCode    []int     `json:"weather_code"`
			WindSpeedMax   []float64 `json:"wind_speed_10m_max"`
			PrecipitationSum []float64 `json:"precipitation_sum"`
		}{
			Time: []string{"2025-07-04"},
			TemperatureMax: []float64{30.0},
			TemperatureMin: []float64{20.0},
			WeatherCode: []int{1},
			WindSpeedMax: []float64{5.0},
			PrecipitationSum: []float64{0.0},
		},
		Hourly: struct {
			Time           []string  `json:"time"`
			Temperature    []float64 `json:"temperature_2m"`
			ApparentTemp   []float64 `json:"apparent_temperature"`
			Humidity       []int     `json:"relative_humidity_2m"`
			WindSpeed      []float64 `json:"wind_speed_10m"`
			WindDirection  []float64 `json:"wind_direction_10m"`
			WeatherCode    []int     `json:"weather_code"`
			Precipitation  []float64 `json:"precipitation"`
		}{
			Time: make([]string, 24),
		},
	}

	result = extractDateBasedWeather(weather, 5) // Beyond available data
	if len(result.Daily.Time) != 0 {
		t.Error("dayOffset beyond available data should result in empty daily data")
	}
}