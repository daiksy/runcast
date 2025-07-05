package display

import (
	"testing"
)

func TestGetRunningTempIcon(t *testing.T) {
	tests := []struct {
		name     string
		temp     float64
		expected string
	}{
		{
			name:     "temp_icon",
			temp:     32.0,
			expected: "🔥 ",
		},
		{
			name:     "temp_icon",
			temp:     27.0,
			expected: "☀️ ",
		},
		{
			name:     "temp_icon",
			temp:     22.0,
			expected: "🌤️ ",
		},
		{
			name:     "temp_icon",
			temp:     17.0,
			expected: "🌥️ ",
		},
		{
			name:     "temp_icon",
			temp:     12.0,
			expected: "☁️ ",
		},
		{
			name:     "temp_icon",
			temp:     8.0,
			expected: "❄️ ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetRunningTempIcon(tt.temp)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s for temperature %.1f", tt.expected, result, tt.temp)
			}
		})
	}
}