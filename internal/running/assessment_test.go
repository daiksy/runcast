package running

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
		expectWarnings       bool
	}{
		{
			name:                 "Hot weather with heat stroke warning",
			temp:                 38.0,
			apparentTemp:         42.0,
			humidity:             60.0,
			windSpeed:            2.0,
			precipitation:        0.0,
			weatherCode:          0,
			expectedLevel:        "普通",
			expectedRecommendation: "注意事項を確認してからランニングしてください",
			expectWarnings:       true,
		},
		{
			name:                 "Thunderstorm conditions",
			temp:                 25.0,
			apparentTemp:         28.0,
			humidity:             80.0,
			windSpeed:            8.0,
			precipitation:        15.0,
			weatherCode:          95,
			expectedLevel:        "危険",
			expectedRecommendation: "天候が悪いため、ランニングは控えることをお勧めします",
			expectWarnings:       true,
		},
		{
			name:                 "Perfect conditions without warnings",
			temp:                 20.0,
			apparentTemp:         22.0,
			humidity:             50.0,
			windSpeed:            2.0,
			precipitation:        0.0,
			weatherCode:          1,
			expectedLevel:        "最高",
			expectedRecommendation: "ランニングに最適な天候です！",
			expectWarnings:       false,
		},
		{
			name:                 "High humidity with moderate warnings",
			temp:                 25.0,
			apparentTemp:         28.0,
			humidity:             90.0,
			windSpeed:            1.0,
			precipitation:        0.0,
			weatherCode:          2,
			expectedLevel:        "最高",
			expectedRecommendation: "ランニングに最適な天候です！",
			expectWarnings:       true,
		},
		{
			name:                 "Strong wind conditions",
			temp:                 18.0,
			apparentTemp:         16.0,
			humidity:             60.0,
			windSpeed:            12.0,
			precipitation:        0.0,
			weatherCode:          1,
			expectedLevel:        "良好",
			expectedRecommendation: "良好な天候です。ランニングを楽しんでください",
			expectWarnings:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition := AssessRunningCondition(tt.temp, tt.apparentTemp, tt.humidity, tt.windSpeed, tt.precipitation, tt.weatherCode)
			
			if condition.Level != tt.expectedLevel {
				t.Errorf("Expected level %s, got %s", tt.expectedLevel, condition.Level)
			}
			
			if condition.Recommendation != tt.expectedRecommendation {
				t.Errorf("Expected recommendation '%s', got '%s'", tt.expectedRecommendation, condition.Recommendation)
			}
			
			hasWarnings := len(condition.Warnings) > 0
			if hasWarnings != tt.expectWarnings {
				t.Errorf("Expected warnings: %v, got warnings: %v", tt.expectWarnings, hasWarnings)
			}
		})
	}
}

func TestSevereWarningDetection(t *testing.T) {
	// Test thunderstorm warning
	condition := AssessRunningCondition(25.0, 28.0, 80.0, 5.0, 10.0, 95)
	
	// Should contain thunderstorm warning
	foundThunderstormWarning := false
	for _, warning := range condition.Warnings {
		if warning == "⚡ 雷雨: 絶対に屋外でのランニングは避けてください" {
			foundThunderstormWarning = true
			break
		}
	}
	
	if !foundThunderstormWarning {
		t.Error("Expected thunderstorm warning not found")
	}
}

func TestWarningCategorization(t *testing.T) {
	tests := []struct {
		name            string
		temp            float64
		apparentTemp    float64
		humidity        float64
		windSpeed       float64
		precipitation   float64
		weatherCode     int
		expectedWarnings []string
	}{
		{
			name:          "Cold weather warnings",
			temp:          2.0,
			apparentTemp:  -1.0,
			humidity:      40.0,
			windSpeed:     1.0,
			precipitation: 0.0,
			weatherCode:   0,
			expectedWarnings: []string{"🥶 低温注意: 防寒対策を十分に行ってください"},
		},
		{
			name:          "High humidity warnings",
			temp:          25.0,
			apparentTemp:  28.0,
			humidity:      88.0,
			windSpeed:     1.0,
			precipitation: 0.0,
			weatherCode:   0,
			expectedWarnings: []string{"💧 高湿度: 汗が乾きにくい状態です"},
		},
		{
			name:          "Strong wind warnings",
			temp:          20.0,
			apparentTemp:  22.0,
			humidity:      50.0,
			windSpeed:     11.0,
			precipitation: 0.0,
			weatherCode:   0,
			expectedWarnings: []string{"💨 強風注意: 転倒や怪我のリスクがあります"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition := AssessRunningCondition(tt.temp, tt.apparentTemp, tt.humidity, tt.windSpeed, tt.precipitation, tt.weatherCode)
			
			for _, expectedWarning := range tt.expectedWarnings {
				found := false
				for _, warning := range condition.Warnings {
					if warning == expectedWarning {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected warning '%s' not found", expectedWarning)
				}
			}
		})
	}
}

func TestDistanceSpecificWarnings(t *testing.T) {
	category := GetDistanceCategory("full")
	
	// Test high temperature with full marathon
	condition := AssessDistanceBasedRunningCondition(26.0, 30.0, 60.0, 2.0, 0.0, 0, category)
	
	// Should contain full marathon specific warning
	foundFullWarning := false
	for _, warning := range condition.Warnings {
		if warning == "🏃‍♂️ フルマラソン警告: 高温下での長時間運動は危険です" {
			foundFullWarning = true
			break
		}
	}
	
	if !foundFullWarning {
		t.Error("Expected full marathon specific warning not found")
	}
}

func TestDistanceSpecificClothing(t *testing.T) {
	categoryHalf := GetDistanceCategory("half")
	
	// Test warm weather with half marathon
	condition := AssessDistanceBasedRunningCondition(25.0, 28.0, 60.0, 2.0, 0.0, 0, categoryHalf)
	
	// Should contain distance-specific clothing recommendations
	foundHydrationGear := false
	for _, item := range condition.Clothing {
		if item == "水分補給用品" {
			foundHydrationGear = true
			break
		}
	}
	
	if !foundHydrationGear {
		t.Error("Expected hydration gear recommendation for long distance running")
	}
}

func TestRunningConditionClothing(t *testing.T) {
	// Test cold weather clothing
	condition := AssessRunningCondition(8.0, 6.0, 60.0, 2.0, 0.0, 0)
	
	// Should recommend warm clothing
	foundWarmClothing := false
	for _, item := range condition.Clothing {
		if item == "長袖" || item == "ロングパンツ" {
			foundWarmClothing = true
			break
		}
	}
	
	if !foundWarmClothing {
		t.Error("Expected warm clothing recommendation for cold weather")
	}
}

func TestDistancePenalties(t *testing.T) {
	// Test that distance penalties are applied correctly
	baseTemp := 30.0
	baseHumidity := 80.0
	baseCondition := AssessRunningCondition(baseTemp, baseTemp, baseHumidity, 2.0, 0.0, 0)
	
	category10k := GetDistanceCategory("10k")
	categoryFull := GetDistanceCategory("full")
	
	condition10k := AssessDistanceBasedRunningCondition(baseTemp, baseTemp, baseHumidity, 2.0, 0.0, 0, category10k)
	conditionFull := AssessDistanceBasedRunningCondition(baseTemp, baseTemp, baseHumidity, 2.0, 0.0, 0, categoryFull)
	
	// Scores should decrease with distance penalties
	if condition10k.Score >= baseCondition.Score {
		t.Error("10k condition should have lower score than base due to penalties")
	}
	
	if conditionFull.Score >= condition10k.Score {
		t.Error("Full marathon condition should have lower score than 10k due to higher penalties")
	}
}

func TestNilDistanceCategory(t *testing.T) {
	// Test that nil distance category works correctly
	condition := AssessDistanceBasedRunningCondition(25.0, 28.0, 60.0, 2.0, 0.0, 0, nil)
	baseCondition := AssessRunningCondition(25.0, 28.0, 60.0, 2.0, 0.0, 0)
	
	// Should be identical to base condition
	if condition.Score != baseCondition.Score {
		t.Errorf("Expected same score for nil distance category, got %d vs %d", condition.Score, baseCondition.Score)
	}
	
	if condition.Level != baseCondition.Level {
		t.Errorf("Expected same level for nil distance category, got %s vs %s", condition.Level, baseCondition.Level)
	}
}