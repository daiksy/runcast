package running

import (
	"testing"
)

func TestGetDistanceCategories(t *testing.T) {
	categories := GetDistanceCategories()
	
	expectedCategories := []string{"5k", "10k", "half", "full"}
	
	categoriesMap := make(map[string]bool)
	for _, cat := range categories {
		categoriesMap[cat.Key] = true
	}
	
	for _, expected := range expectedCategories {
		if !categoriesMap[expected] {
			t.Errorf("Expected distance category %s not found", expected)
		}
	}
	
	// Test 5k category (should have no penalties)
	category5k := GetDistanceCategory("5k")
	if category5k != nil {
		if category5k.TempPenalty != 0 || category5k.HumidityPenalty != 0 || category5k.WindPenalty != 0 {
			t.Errorf("5k category should have no penalties, got temp=%d, humidity=%d, wind=%d", 
				category5k.TempPenalty, category5k.HumidityPenalty, category5k.WindPenalty)
		}
		if category5k.DisplayName != "5キロ" {
			t.Errorf("Expected display name '5キロ', got '%s'", category5k.DisplayName)
		}
	}
	
	// Test full marathon category (should have highest penalties)
	categoryFull := GetDistanceCategory("full")
	if categoryFull != nil {
		if categoryFull.TempPenalty != 15 || categoryFull.HumidityPenalty != 10 || categoryFull.WindPenalty != 5 {
			t.Errorf("Full category penalties incorrect, got temp=%d, humidity=%d, wind=%d", 
				categoryFull.TempPenalty, categoryFull.HumidityPenalty, categoryFull.WindPenalty)
		}
		if categoryFull.DisplayName != "フルマラソン" {
			t.Errorf("Expected display name 'フルマラソン', got '%s'", categoryFull.DisplayName)
		}
	}
}

func TestGetDistanceCategory(t *testing.T) {
	tests := []struct {
		name        string
		distance    string
		expectFound bool
		displayName string
	}{
		{
			name:        "5k",
			distance:    "5k",
			expectFound: true,
			displayName: "5キロ",
		},
		{
			name:        "10k",
			distance:    "10k",
			expectFound: true,
			displayName: "10キロ",
		},
		{
			name:        "half",
			distance:    "half",
			expectFound: true,
			displayName: "ハーフマラソン",
		},
		{
			name:        "full",
			distance:    "full",
			expectFound: true,
			displayName: "フルマラソン",
		},
		{
			name:        "invalid",
			distance:    "invalid",
			expectFound: false,
		},
		{
			name:        "",
			distance:    "",
			expectFound: false,
		},
		{
			name:        "3k",
			distance:    "3k",
			expectFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category := GetDistanceCategory(tt.distance)
			
			if tt.expectFound {
				if category == nil {
					t.Errorf("Expected to find category for %s, but got nil", tt.distance)
					return
				}
				if category.DisplayName != tt.displayName {
					t.Errorf("Expected display name %s, got %s", tt.displayName, category.DisplayName)
				}
			} else {
				if category != nil {
					t.Errorf("Expected nil for %s, but got category", tt.distance)
				}
			}
		})
	}
}

func TestAssessRunningCondition(t *testing.T) {
	tests := []struct {
		name          string
		temp          float64
		apparentTemp  float64
		humidity      float64
		windSpeed     float64
		precipitation float64
		weatherCode   int
		expectedLevel string
		minScore      int
		maxScore      int
	}{
		{
			name:          "Perfect conditions",
			temp:          22.0,
			apparentTemp:  24.0,
			humidity:      50.0,
			windSpeed:     2.0,
			precipitation: 0.0,
			weatherCode:   0,
			expectedLevel: "最高",
			minScore:      80,
			maxScore:      100,
		},
		{
			name:          "Hot weather",
			temp:          35.0,
			apparentTemp:  40.0,
			humidity:      60.0,
			windSpeed:     1.0,
			precipitation: 0.0,
			weatherCode:   0,
			expectedLevel: "普通",
			minScore:      40,
			maxScore:      59,
		},
		{
			name:          "Cold weather",
			temp:          2.0,
			apparentTemp:  -1.0,
			humidity:      40.0,
			windSpeed:     1.0,
			precipitation: 0.0,
			weatherCode:   0,
			expectedLevel: "良好",
			minScore:      60,
			maxScore:      79,
		},
		{
			name:          "Rainy conditions",
			temp:          20.0,
			apparentTemp:  22.0,
			humidity:      90.0,
			windSpeed:     3.0,
			precipitation: 2.0,
			weatherCode:   61,
			expectedLevel: "普通",
			minScore:      40,
			maxScore:      59,
		},
		{
			name:          "Thunderstorm",
			temp:          25.0,
			apparentTemp:  28.0,
			humidity:      80.0,
			windSpeed:     5.0,
			precipitation: 10.0,
			weatherCode:   95,
			expectedLevel: "危険",
			minScore:      0,
			maxScore:      19,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition := AssessRunningCondition(tt.temp, tt.apparentTemp, tt.humidity, tt.windSpeed, tt.precipitation, tt.weatherCode)
			
			if condition.Level != tt.expectedLevel {
				t.Errorf("Expected level %s, got %s", tt.expectedLevel, condition.Level)
			}
			
			if condition.Score < tt.minScore || condition.Score > tt.maxScore {
				t.Errorf("Expected score between %d and %d, got %d", tt.minScore, tt.maxScore, condition.Score)
			}
		})
	}
}

func TestAssessDistanceBasedRunningCondition(t *testing.T) {
	category10k := GetDistanceCategory("10k")
	categoryFull := GetDistanceCategory("full")
	
	// Test same conditions with different distances
	temp := 28.0
	apparentTemp := 32.0
	humidity := 75.0
	windSpeed := 3.0
	precipitation := 0.0
	weatherCode := 1
	
	// Base condition (no distance)
	baseCondition := AssessRunningCondition(temp, apparentTemp, humidity, windSpeed, precipitation, weatherCode)
	
	// 10k condition
	condition10k := AssessDistanceBasedRunningCondition(temp, apparentTemp, humidity, windSpeed, precipitation, weatherCode, category10k)
	
	// Full marathon condition
	conditionFull := AssessDistanceBasedRunningCondition(temp, apparentTemp, humidity, windSpeed, precipitation, weatherCode, categoryFull)
	
	// Scores should decrease with distance (more penalties)
	if condition10k.Score >= baseCondition.Score {
		t.Errorf("10k score should be lower than base score due to penalties")
	}
	
	if conditionFull.Score >= condition10k.Score {
		t.Errorf("Full marathon score should be lower than 10k score due to higher penalties")
	}
	
	// Full marathon should have additional warnings
	if len(conditionFull.Warnings) <= len(condition10k.Warnings) {
		t.Errorf("Full marathon should have more warnings than 10k")
	}
}