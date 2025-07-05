package main

import (
	"testing"
)

func TestGetDistanceCategories(t *testing.T) {
	categories := getDistanceCategories()
	
	expectedCategories := []string{"5k", "10k", "half", "full"}
	
	for _, expected := range expectedCategories {
		if _, exists := categories[expected]; !exists {
			t.Errorf("Expected distance category %s not found", expected)
		}
	}
	
	// Test 5k category (should have no penalties)
	if category, exists := categories["5k"]; exists {
		if category.TempPenalty != 0 || category.HumidityPenalty != 0 || category.WindPenalty != 0 {
			t.Errorf("5k category should have no penalties, got temp=%d, humidity=%d, wind=%d", 
				category.TempPenalty, category.HumidityPenalty, category.WindPenalty)
		}
		if category.DisplayName != "5キロ" {
			t.Errorf("Expected 5k display name '5キロ', got '%s'", category.DisplayName)
		}
	}
	
	// Test full marathon category (should have highest penalties)
	if category, exists := categories["full"]; exists {
		if category.TempPenalty <= 0 || category.HumidityPenalty <= 0 || category.WindPenalty <= 0 {
			t.Errorf("Full marathon category should have penalties, got temp=%d, humidity=%d, wind=%d", 
				category.TempPenalty, category.HumidityPenalty, category.WindPenalty)
		}
		if category.DisplayName != "フルマラソン" {
			t.Errorf("Expected full marathon display name 'フルマラソン', got '%s'", category.DisplayName)
		}
	}
}

func TestGetDistanceCategory(t *testing.T) {
	tests := []struct {
		distance string
		expected bool
	}{
		{"5k", true},
		{"10k", true},
		{"half", true},
		{"full", true},
		{"invalid", false},
		{"", false},
		{"3k", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.distance, func(t *testing.T) {
			result := getDistanceCategory(tt.distance)
			exists := result != nil
			if exists != tt.expected {
				t.Errorf("Expected exists=%v for distance '%s', got %v", tt.expected, tt.distance, exists)
			}
		})
	}
}

func TestAssessDistanceBasedRunningCondition(t *testing.T) {
	// Test conditions: moderate temperature, high humidity
	temp := 25.0
	apparentTemp := 28.0
	humidity := 85.0
	windSpeed := 3.0
	precipitation := 0.0
	weatherCode := 1
	
	// Test 5k vs full marathon scoring difference
	category5k := getDistanceCategory("5k")
	categoryFull := getDistanceCategory("full")
	
	condition5k := assessDistanceBasedRunningCondition(temp, apparentTemp, humidity, windSpeed, precipitation, weatherCode, category5k)
	conditionFull := assessDistanceBasedRunningCondition(temp, apparentTemp, humidity, windSpeed, precipitation, weatherCode, categoryFull)
	
	// Full marathon should have lower score due to penalties
	if conditionFull.Score >= condition5k.Score {
		t.Errorf("Full marathon score (%d) should be lower than 5k score (%d) under same conditions", 
			conditionFull.Score, condition5k.Score)
	}
	
	// Check recommendations are distance-specific
	if !contains(condition5k.Recommendation, "5キロ") {
		t.Errorf("5k recommendation should mention '5キロ', got: %s", condition5k.Recommendation)
	}
	if !contains(conditionFull.Recommendation, "フルマラソン") {
		t.Errorf("Full marathon recommendation should mention 'フルマラソン', got: %s", conditionFull.Recommendation)
	}
}

func TestDistanceSpecificWarnings(t *testing.T) {
	// Test high temperature conditions for full marathon
	temp := 30.0
	apparentTemp := 33.0
	humidity := 90.0
	windSpeed := 2.0
	precipitation := 0.0
	weatherCode := 1
	
	categoryFull := getDistanceCategory("full")
	condition := assessDistanceBasedRunningCondition(temp, apparentTemp, humidity, windSpeed, precipitation, weatherCode, categoryFull)
	
	// Should have full marathon specific warning
	hasFullMarathonWarning := false
	hasLongDistanceWarning := false
	
	for _, warning := range condition.Warnings {
		if contains(warning, "フルマラソン警告") {
			hasFullMarathonWarning = true
		}
		if contains(warning, "長距離警告") {
			hasLongDistanceWarning = true
		}
	}
	
	if !hasFullMarathonWarning {
		t.Error("Full marathon under high temperature should have specific warning")
	}
	if !hasLongDistanceWarning {
		t.Error("Full marathon under high humidity should have long distance warning")
	}
}

func TestDistanceSpecificClothing(t *testing.T) {
	// Test moderate temperature for long distances
	temp := 22.0
	apparentTemp := 24.0
	humidity := 60.0
	windSpeed := 2.0
	precipitation := 0.0
	weatherCode := 1
	
	category5k := getDistanceCategory("5k")
	categoryHalf := getDistanceCategory("half")
	
	condition5k := assessDistanceBasedRunningCondition(temp, apparentTemp, humidity, windSpeed, precipitation, weatherCode, category5k)
	conditionHalf := assessDistanceBasedRunningCondition(temp, apparentTemp, humidity, windSpeed, precipitation, weatherCode, categoryHalf)
	
	// Half marathon should have additional gear recommendations
	has5kHydration := false
	hasHalfHydration := false
	
	for _, item := range condition5k.Clothing {
		if contains(item, "水分補給用品") {
			has5kHydration = true
		}
	}
	
	for _, item := range conditionHalf.Clothing {
		if contains(item, "水分補給用品") {
			hasHalfHydration = true
		}
	}
	
	if has5kHydration {
		t.Error("5k should not require hydration gear under moderate conditions")
	}
	if !hasHalfHydration {
		t.Error("Half marathon should require hydration gear")
	}
}

func TestDistancePenalties(t *testing.T) {
	// Test that penalties are applied correctly
	temp := 27.0
	apparentTemp := 30.0
	humidity := 85.0
	windSpeed := 7.0
	precipitation := 0.0
	weatherCode := 1
	
	// Test each distance category
	distances := []string{"5k", "10k", "half", "full"}
	var scores []int
	
	for _, distance := range distances {
		category := getDistanceCategory(distance)
		condition := assessDistanceBasedRunningCondition(temp, apparentTemp, humidity, windSpeed, precipitation, weatherCode, category)
		scores = append(scores, condition.Score)
	}
	
	// Scores should generally decrease as distance increases (more penalties)
	for i := 1; i < len(scores); i++ {
		if scores[i] > scores[i-1] {
			t.Errorf("Score for %s (%d) should not be higher than %s (%d) under challenging conditions", 
				distances[i], scores[i], distances[i-1], scores[i-1])
		}
	}
}

func TestNilDistanceCategory(t *testing.T) {
	// Test that function handles nil distance category gracefully
	temp := 25.0
	apparentTemp := 27.0
	humidity := 60.0
	windSpeed := 3.0
	precipitation := 0.0
	weatherCode := 1
	
	condition := assessDistanceBasedRunningCondition(temp, apparentTemp, humidity, windSpeed, precipitation, weatherCode, nil)
	baseCondition := assessRunningCondition(temp, apparentTemp, humidity, windSpeed, precipitation, weatherCode)
	
	// Should behave the same as base assessment when distance category is nil
	if condition.Score != baseCondition.Score {
		t.Errorf("Nil distance category should produce same score as base assessment, got %d vs %d", 
			condition.Score, baseCondition.Score)
	}
}

