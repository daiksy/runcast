package weather

import (
	"testing"
	"weather-cli/internal/types"
)

func TestGetTimePeriods(t *testing.T) {
	periods := GetTimePeriods()
	
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
		name     string
		timeStr  string
		expected string
	}{
		{
			name:     "2025-07-04T05:00",
			timeStr:  "2025-07-04T05:00",
			expected: "05",
		},
		{
			name:     "2025-07-04T15:30",
			timeStr:  "2025-07-04T15:30",
			expected: "15",
		},
		{
			name:     "2025-07-04T23:45",
			timeStr:  "2025-07-04T23:45",
			expected: "23",
		},
		{
			name:     "2025-07-04T09:15:30Z",
			timeStr:  "2025-07-04T09:15:30Z",
			expected: "09",
		},
		{
			name:     "invalid",
			timeStr:  "invalid",
			expected: "invalid",
		},
		{
			name:     "",
			timeStr:  "",
			expected: "",
		},
		{
			name:     "2025-07-04",
			timeStr:  "2025-07-04",
			expected: "2025-07-04",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractHour(tt.timeStr)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestExtractTimeBasedWeather(t *testing.T) {
	// Create mock weather data
	weather := &types.WeatherData{
		Hourly: struct {
			Time          []string  `json:"time"`
			Temperature   []float64 `json:"temperature_2m"`
			ApparentTemp  []float64 `json:"apparent_temperature"`
			Humidity      []int     `json:"relative_humidity_2m"`
			WindSpeed     []float64 `json:"wind_speed_10m"`
			WindDirection []float64 `json:"wind_direction_10m"`
			Precipitation []float64 `json:"precipitation"`
			WeatherCode   []int     `json:"weather_code"`
		}{
			Time:          []string{"2025-07-05T05:00", "2025-07-05T06:00", "2025-07-05T07:00", "2025-07-05T08:00", "2025-07-05T09:00", "2025-07-05T10:00", "2025-07-05T11:00"},
			Temperature:   []float64{20.0, 21.0, 22.0, 23.0, 24.0, 25.0, 26.0},
			ApparentTemp:  []float64{19.0, 20.0, 21.0, 22.0, 23.0, 24.0, 25.0},
			Humidity:      []int{80, 75, 70, 65, 60, 55, 50},
			WindSpeed:     []float64{2.0, 2.5, 3.0, 3.5, 4.0, 4.5, 5.0},
			WindDirection: []float64{180, 185, 190, 195, 200, 205, 210},
			Precipitation: []float64{0, 0, 0, 0, 0, 0, 0},
			WeatherCode:   []int{0, 0, 1, 1, 2, 2, 3},
		},
	}

	// Test morning period extraction
	timeData := ExtractTimeBasedWeather(weather, "morning", 1)
	
	// Morning is 5-9, so should get 5 entries
	expectedCount := 5
	if len(timeData) != expectedCount {
		t.Errorf("Expected %d time entries for morning, got %d", expectedCount, len(timeData))
	}
	
	// Check first entry
	if len(timeData) > 0 {
		first := timeData[0]
		if first.Temperature != 20.0 {
			t.Errorf("Expected first temperature 20.0, got %f", first.Temperature)
		}
	}
}

func TestGetDateDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		dateSpec string
		expected string
	}{
		{
			name:     "today",
			dateSpec: "today",
			expected: "今日の",
		},
		{
			name:     "tomorrow",
			dateSpec: "tomorrow",
			expected: "明日の",
		},
		{
			name:     "day-after-tomorrow",
			dateSpec: "day-after-tomorrow",
			expected: "明後日の",
		},
		{
			name:     "invalid",
			dateSpec: "invalid",
			expected: "invalidの",
		},
		{
			name:     "",
			dateSpec: "",
			expected: "の",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetDateDisplayName(tt.dateSpec)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetDateOffset(t *testing.T) {
	tests := []struct {
		name     string
		dateSpec string
		expected int
	}{
		{
			name:     "today",
			dateSpec: "today",
			expected: 0,
		},
		{
			name:     "tomorrow",
			dateSpec: "tomorrow",
			expected: 1,
		},
		{
			name:     "day-after-tomorrow",
			dateSpec: "day-after-tomorrow",
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetDateOffset(tt.dateSpec)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestExtractDateBasedWeather(t *testing.T) {
	// Create mock weather data with daily data
	weather := &types.WeatherData{
		Current: struct {
			Temperature   float64 `json:"temperature_2m"`
			ApparentTemp  float64 `json:"apparent_temperature"`
			Humidity      int     `json:"relative_humidity_2m"`
			WindSpeed     float64 `json:"wind_speed_10m"`
			WindDirection float64 `json:"wind_direction_10m"`
			Precipitation float64 `json:"precipitation"`
			WeatherCode   int     `json:"weather_code"`
		}{
			Temperature: 25.0,
		},
		Daily: struct {
			Time                        []string  `json:"time"`
			TemperatureMax              []float64 `json:"temperature_2m_max"`
			TemperatureMin              []float64 `json:"temperature_2m_min"`
			WindSpeedMax                []float64 `json:"wind_speed_10m_max"`
			WindGustMax                 []float64 `json:"wind_gusts_10m_max"`
			PrecipitationSum            []float64 `json:"precipitation_sum"`
			WeatherCode                 []int     `json:"weather_code"`
			SunriseTime                 []string  `json:"sunrise"`
			SunsetTime                  []string  `json:"sunset"`
			DaylightDuration            []float64 `json:"daylight_duration"`
			SunshineDuration            []float64 `json:"sunshine_duration"`
			UvIndexMax                  []float64 `json:"uv_index_max"`
			UvIndexClearSkyMax          []float64 `json:"uv_index_clear_sky_max"`
			PrecipitationHours          []float64 `json:"precipitation_hours"`
			PrecipitationProbabilityMax []float64 `json:"precipitation_probability_max"`
		}{
			Time:           []string{"2025-07-05", "2025-07-06", "2025-07-07"},
			TemperatureMax: []float64{30.0, 28.0, 26.0},
			TemperatureMin: []float64{20.0, 18.0, 16.0},
			WindSpeedMax:   []float64{5.0, 4.0, 3.0},
			PrecipitationSum: []float64{0.0, 1.0, 2.0},
			WeatherCode:    []int{0, 1, 2},
		},
		Hourly: struct {
			Time          []string  `json:"time"`
			Temperature   []float64 `json:"temperature_2m"`
			ApparentTemp  []float64 `json:"apparent_temperature"`
			Humidity      []int     `json:"relative_humidity_2m"`
			WindSpeed     []float64 `json:"wind_speed_10m"`
			WindDirection []float64 `json:"wind_direction_10m"`
			Precipitation []float64 `json:"precipitation"`
			WeatherCode   []int     `json:"weather_code"`
		}{
			Time: []string{"2025-07-05T00:00", "2025-07-05T01:00", "2025-07-06T00:00"},
		},
	}

	// Test dayOffset 0 (today)
	result := ExtractDateBasedWeather(weather, 0)
	if result != weather {
		t.Error("Expected same weather data for dayOffset 0")
	}

	// Test dayOffset 1 (tomorrow)
	result = ExtractDateBasedWeather(weather, 1)
	if len(result.Daily.Time) != 1 {
		t.Errorf("Expected 1 daily entry for tomorrow, got %d", len(result.Daily.Time))
	}
	if len(result.Daily.Time) > 0 && result.Daily.Time[0] != "2025-07-06" {
		t.Errorf("Expected date 2025-07-06, got %s", result.Daily.Time[0])
	}
}