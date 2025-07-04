package main

import (
	"testing"
)

func TestRunningRecommendationWithWarnings(t *testing.T) {
	tests := []struct {
		name                 string
		temp                 float64
		apparentTemp         float64
		humidity             float64
		windSpeed            float64
		precipitation        float64
		weatherCode          int
		expectedLevel        string
		expectedRecommendation string
		shouldHaveWarnings   bool
	}{
		{
			name:                 "Hot weather with heat stroke warning",
			temp:                 35,
			apparentTemp:         40,
			humidity:             80,
			windSpeed:            2,
			precipitation:        0,
			weatherCode:          1,
			expectedLevel:        "注意",
			expectedRecommendation: "警告事項があります。ランニングは控えめに",
			shouldHaveWarnings:   true,
		},
		{
			name:                 "Thunderstorm conditions",
			temp:                 25,
			apparentTemp:         28,
			humidity:             85,
			windSpeed:            12,
			precipitation:        10,
			weatherCode:          95,
			expectedLevel:        "危険",
			expectedRecommendation: "危険な状況です。ランニング中止を強く推奨します",
			shouldHaveWarnings:   true,
		},
		{
			name:                 "Perfect conditions without warnings",
			temp:                 20,
			apparentTemp:         20,
			humidity:             50,
			windSpeed:            2,
			precipitation:        0,
			weatherCode:          1,
			expectedLevel:        "最高",
			expectedRecommendation: "ランニングに最適な天候です！",
			shouldHaveWarnings:   false,
		},
		{
			name:                 "High humidity with moderate warnings",
			temp:                 25,
			apparentTemp:         25,
			humidity:             85,
			windSpeed:            3,
			precipitation:        0,
			weatherCode:          1,
			expectedLevel:        "良好",
			expectedRecommendation: "注意事項を確認してからランニングしてください",
			shouldHaveWarnings:   true,
		},
		{
			name:                 "Strong wind conditions",
			temp:                 20,
			apparentTemp:         18,
			humidity:             60,
			windSpeed:            11,
			precipitation:        0,
			weatherCode:          1,
			expectedLevel:        "注意",
			expectedRecommendation: "警告事項があります。十分注意してランニングしてください",
			shouldHaveWarnings:   true,
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

			if condition.Recommendation != tt.expectedRecommendation {
				t.Errorf("Expected recommendation '%s', got '%s'", 
					tt.expectedRecommendation, condition.Recommendation)
			}

			hasWarnings := len(condition.Warnings) > 0
			if hasWarnings != tt.shouldHaveWarnings {
				t.Errorf("Expected warnings: %v, got warnings: %v (count: %d)", 
					tt.shouldHaveWarnings, hasWarnings, len(condition.Warnings))
			}

			// Verify that severe warnings override positive recommendations
			if tt.shouldHaveWarnings && tt.temp >= 30 {
				// Should not recommend running with heat stroke warning
				if condition.Recommendation == "ランニングに最適な天候です！" ||
				   condition.Recommendation == "ランニングに適した天候です" {
					t.Errorf("Should not recommend running with heat stroke warning. Got: %s", 
						condition.Recommendation)
				}
			}
		})
	}
}

func TestSevereWarningDetection(t *testing.T) {
	// Test heat stroke warning
	condition := assessRunningCondition(35, 40, 50, 2, 0, 1)
	
	foundHeatWarning := false
	for _, warning := range condition.Warnings {
		if contains(warning, "熱中症注意") {
			foundHeatWarning = true
			break
		}
	}
	
	if !foundHeatWarning {
		t.Error("Expected heat stroke warning for high temperature")
	}
	
	// Should not recommend "最適" or "適した" with severe warnings
	if condition.Recommendation == "ランニングに最適な天候です！" ||
	   condition.Recommendation == "ランニングに適した天候です" {
		t.Errorf("Should not give positive recommendation with severe warnings. Got: %s", 
			condition.Recommendation)
	}
}

func TestWarningCategorization(t *testing.T) {
	// Test that different warning types are properly categorized
	
	// Severe warning (heat stroke)
	condition1 := assessRunningCondition(35, 40, 50, 2, 0, 1)
	if condition1.Level == "最高" || condition1.Level == "良好" {
		t.Error("Severe warning should prevent highest ratings")
	}
	
	// Moderate warning (high humidity)
	condition2 := assessRunningCondition(22, 22, 85, 3, 0, 1)
	if len(condition2.Warnings) == 0 {
		t.Error("High humidity should generate warnings")
	}
	
	// No warnings
	condition3 := assessRunningCondition(20, 20, 50, 2, 0, 1)
	if len(condition3.Warnings) > 0 {
		t.Error("Perfect conditions should not generate warnings")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || 
		   (len(s) > len(substr) && contains(s[1:], substr))
}