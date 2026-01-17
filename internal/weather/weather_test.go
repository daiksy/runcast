package weather

import (
	"testing"
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
			coord, err := GetCityCoordinate(tt.city)
			
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
				t.Errorf("Expected lat %f, got %f", tt.expectedLat, coord.Lat)
			}
			
			if coord.Lon != tt.expectedLon {
				t.Errorf("Expected lon %f, got %f", tt.expectedLon, coord.Lon)
			}
		})
	}
}

func TestGetWeatherDescription(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		expected string
	}{
		{
			name:     "快晴",
			code:     0,
			expected: "快晴",
		},
		{
			name:     "晴れ",
			code:     1,
			expected: "晴れ",
		},
		{
			name:     "一部曇り",
			code:     2,
			expected: "一部曇り",
		},
		{
			name:     "曇り",
			code:     3,
			expected: "曇り",
		},
		{
			name:     "霧",
			code:     45,
			expected: "霧",
		},
		{
			name:     "弱い霧雨",
			code:     51,
			expected: "弱い霧雨",
		},
		{
			name:     "弱い雨",
			code:     61,
			expected: "弱い雨",
		},
		{
			name:     "弱い雪",
			code:     71,
			expected: "弱い雪",
		},
		{
			name:     "雷雨",
			code:     95,
			expected: "雷雨",
		},
		{
			name:     "不明",
			code:     999,
			expected: "不明",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetWeatherDescription(tt.code)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
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
			input:    "2025-07-05",
			expected: "07月05日",
		},
		{
			name:     "Another valid date",
			input:    "2025-12-25",
			expected: "12月25日",
		},
		{
			name:     "Date with time",
			input:    "2025-07-05T10:30:00",
			expected: "07月05日",
		},
		{
			name:     "Short string",
			input:    "2025-07",
			expected: "2025-07",
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
			result := FormatDate(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetWindDirection(t *testing.T) {
	tests := []struct {
		name      string
		direction float64
		expected  string
	}{
		{
			name:      "北",
			direction: 0,
			expected:  "北",
		},
		{
			name:      "北東",
			direction: 45,
			expected:  "北東",
		},
		{
			name:      "東",
			direction: 90,
			expected:  "東",
		},
		{
			name:      "南東",
			direction: 135,
			expected:  "南東",
		},
		{
			name:      "南",
			direction: 180,
			expected:  "南",
		},
		{
			name:      "南西",
			direction: 225,
			expected:  "南西",
		},
		{
			name:      "西",
			direction: 270,
			expected:  "西",
		},
		{
			name:      "北西",
			direction: 315,
			expected:  "北西",
		},
		{
			name:      "北",
			direction: 360,
			expected:  "北",
		},
		{
			name:      "北北東",
			direction: 22.5,
			expected:  "北北東",
		},
		{
			name:      "東北東",
			direction: 67.5,
			expected:  "東北東",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetWindDirection(tt.direction)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestCreateDustLevel(t *testing.T) {
	tests := []struct {
		name                string
		dust                float64
		pm10                float64
		pm2_5               float64
		expectedLevel       int
		expectedDisplayName string
	}{
		{
			name:                "No dust",
			dust:                30,
			pm10:                40,
			pm2_5:               20,
			expectedLevel:       0,
			expectedDisplayName: "なし",
		},
		{
			name:                "Low dust",
			dust:                75,
			pm10:                60,
			pm2_5:               30,
			expectedLevel:       1,
			expectedDisplayName: "少ない",
		},
		{
			name:                "Moderate dust",
			dust:                150,
			pm10:                100,
			pm2_5:               50,
			expectedLevel:       2,
			expectedDisplayName: "やや多い",
		},
		{
			name:                "High dust",
			dust:                300,
			pm10:                200,
			pm2_5:               80,
			expectedLevel:       3,
			expectedDisplayName: "多い",
		},
		{
			name:                "Very high dust",
			dust:                600,
			pm10:                400,
			pm2_5:               150,
			expectedLevel:       4,
			expectedDisplayName: "非常に多い",
		},
		{
			name:                "Boundary 50",
			dust:                50,
			pm10:                40,
			pm2_5:               20,
			expectedLevel:       0,
			expectedDisplayName: "なし",
		},
		{
			name:                "Boundary 51",
			dust:                51,
			pm10:                40,
			pm2_5:               20,
			expectedLevel:       1,
			expectedDisplayName: "少ない",
		},
		{
			name:                "Boundary 100",
			dust:                100,
			pm10:                80,
			pm2_5:               40,
			expectedLevel:       1,
			expectedDisplayName: "少ない",
		},
		{
			name:                "Boundary 101",
			dust:                101,
			pm10:                80,
			pm2_5:               40,
			expectedLevel:       2,
			expectedDisplayName: "やや多い",
		},
		{
			name:                "Boundary 200",
			dust:                200,
			pm10:                150,
			pm2_5:               60,
			expectedLevel:       2,
			expectedDisplayName: "やや多い",
		},
		{
			name:                "Boundary 201",
			dust:                201,
			pm10:                150,
			pm2_5:               60,
			expectedLevel:       3,
			expectedDisplayName: "多い",
		},
		{
			name:                "Boundary 500",
			dust:                500,
			pm10:                350,
			pm2_5:               120,
			expectedLevel:       3,
			expectedDisplayName: "多い",
		},
		{
			name:                "Boundary 501",
			dust:                501,
			pm10:                350,
			pm2_5:               120,
			expectedLevel:       4,
			expectedDisplayName: "非常に多い",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := createDustLevel(tt.dust, tt.pm10, tt.pm2_5)

			if result.Level != tt.expectedLevel {
				t.Errorf("Expected level %d, got %d", tt.expectedLevel, result.Level)
			}

			if result.DisplayName != tt.expectedDisplayName {
				t.Errorf("Expected display name %s, got %s", tt.expectedDisplayName, result.DisplayName)
			}

			if result.Dust != tt.dust {
				t.Errorf("Expected dust %f, got %f", tt.dust, result.Dust)
			}

			if result.PM10 != tt.pm10 {
				t.Errorf("Expected PM10 %f, got %f", tt.pm10, result.PM10)
			}

			if result.PM2_5 != tt.pm2_5 {
				t.Errorf("Expected PM2.5 %f, got %f", tt.pm2_5, result.PM2_5)
			}
		})
	}
}