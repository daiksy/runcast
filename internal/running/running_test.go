package running

import (
	"runcast/internal/types"
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

func TestGetDustPenalty(t *testing.T) {
	tests := []struct {
		name            string
		dustLevel       *types.DustLevel
		expectedPenalty int
	}{
		{
			name:            "nil dust level",
			dustLevel:       nil,
			expectedPenalty: 0,
		},
		{
			name:            "level 0 (none)",
			dustLevel:       &types.DustLevel{Level: 0},
			expectedPenalty: 0,
		},
		{
			name:            "level 1 (low)",
			dustLevel:       &types.DustLevel{Level: 1},
			expectedPenalty: 5,
		},
		{
			name:            "level 2 (moderate)",
			dustLevel:       &types.DustLevel{Level: 2},
			expectedPenalty: 15,
		},
		{
			name:            "level 3 (high)",
			dustLevel:       &types.DustLevel{Level: 3},
			expectedPenalty: 30,
		},
		{
			name:            "level 4 (very high)",
			dustLevel:       &types.DustLevel{Level: 4},
			expectedPenalty: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			penalty := GetDustPenalty(tt.dustLevel)
			if penalty != tt.expectedPenalty {
				t.Errorf("Expected penalty %d, got %d", tt.expectedPenalty, penalty)
			}
		})
	}
}

func TestGetDistanceDustMultiplier(t *testing.T) {
	tests := []struct {
		name               string
		distanceCategory   *types.DistanceCategory
		expectedMultiplier float64
	}{
		{
			name:               "nil category",
			distanceCategory:   nil,
			expectedMultiplier: 1.0,
		},
		{
			name:               "5k",
			distanceCategory:   GetDistanceCategory("5k"),
			expectedMultiplier: 1.0,
		},
		{
			name:               "10k",
			distanceCategory:   GetDistanceCategory("10k"),
			expectedMultiplier: 1.2,
		},
		{
			name:               "half",
			distanceCategory:   GetDistanceCategory("half"),
			expectedMultiplier: 1.5,
		},
		{
			name:               "full",
			distanceCategory:   GetDistanceCategory("full"),
			expectedMultiplier: 2.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			multiplier := GetDistanceDustMultiplier(tt.distanceCategory)
			if multiplier != tt.expectedMultiplier {
				t.Errorf("Expected multiplier %f, got %f", tt.expectedMultiplier, multiplier)
			}
		})
	}
}

func TestApplyDustPenalty(t *testing.T) {
	// Test with moderate dust level
	condition := types.RunningCondition{
		Score:          80,
		Level:          "最高",
		Recommendation: "ランニングに最適な天候です！",
		Warnings:       []string{},
		Clothing:       []string{"薄手の半袖"},
	}

	dustLevel := &types.DustLevel{
		Level:       2,
		DisplayName: "やや多い",
		Description: "視程に影響の可能性",
		Dust:        150,
		PM10:        80,
		PM2_5:       35,
	}

	ApplyDustPenalty(&condition, dustLevel, nil)

	// Score should be reduced by 15 (level 2 penalty)
	if condition.Score != 65 {
		t.Errorf("Expected score 65, got %d", condition.Score)
	}

	// Should have dust warning
	hasWarning := false
	for _, warning := range condition.Warnings {
		if warning == "🌫️ 黄砂が飛来しています。マスク着用を推奨します" {
			hasWarning = true
			break
		}
	}
	if !hasWarning {
		t.Errorf("Expected dust warning not found")
	}

	// Should have sports mask in clothing
	hasMask := false
	for _, item := range condition.Clothing {
		if item == "スポーツマスク" {
			hasMask = true
			break
		}
	}
	if !hasMask {
		t.Errorf("Expected sports mask in clothing recommendations")
	}

	// Level should be updated to 良好 (score 65)
	if condition.Level != "良好" {
		t.Errorf("Expected level '良好', got '%s'", condition.Level)
	}
}

func TestGetPM25Penalty(t *testing.T) {
	tests := []struct {
		name            string
		pm25            float64
		expectedPenalty int
	}{
		{
			name:            "Good (below 35)",
			pm25:            30,
			expectedPenalty: 0,
		},
		{
			name:            "Boundary 35",
			pm25:            35,
			expectedPenalty: 0,
		},
		{
			name:            "Slightly elevated (36-50)",
			pm25:            45,
			expectedPenalty: 5,
		},
		{
			name:            "Boundary 50",
			pm25:            50,
			expectedPenalty: 5,
		},
		{
			name:            "High (51-70)",
			pm25:            60,
			expectedPenalty: 15,
		},
		{
			name:            "Boundary 70",
			pm25:            70,
			expectedPenalty: 15,
		},
		{
			name:            "Very high (71+)",
			pm25:            85,
			expectedPenalty: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			penalty := GetPM25Penalty(tt.pm25)
			if penalty != tt.expectedPenalty {
				t.Errorf("Expected penalty %d, got %d", tt.expectedPenalty, penalty)
			}
		})
	}
}

func TestApplyDustPenaltyWithPM25Warning(t *testing.T) {
	// Test with high PM2.5 level
	condition := types.RunningCondition{
		Score:          100,
		Level:          "最高",
		Recommendation: "ランニングに最適な天候です！",
		Warnings:       []string{},
		Clothing:       []string{},
	}

	dustLevel := &types.DustLevel{
		Level:       0, // No dust
		DisplayName: "なし",
		Description: "黄砂の影響なし",
		Dust:        10,
		PM10:        60,
		PM2_5:       55, // Above 50, should trigger warning
	}

	ApplyDustPenalty(&condition, dustLevel, nil)

	// Score should be reduced by PM2.5 penalty (15)
	if condition.Score != 85 {
		t.Errorf("Expected score 85, got %d", condition.Score)
	}

	// Should have PM2.5 warning
	hasWarning := false
	for _, warning := range condition.Warnings {
		if warning == "😷 PM2.5が高め(50μg/m³超)です。長時間の屋外運動に注意してください" {
			hasWarning = true
			break
		}
	}
	if !hasWarning {
		t.Errorf("Expected PM2.5 warning not found, warnings: %v", condition.Warnings)
	}

	// Should have sports mask recommendation
	hasMask := false
	for _, item := range condition.Clothing {
		if item == "スポーツマスク" {
			hasMask = true
			break
		}
	}
	if !hasMask {
		t.Errorf("Expected sports mask in clothing recommendations")
	}
}

func TestApplyDustPenaltyWithAlertLevelPM25(t *testing.T) {
	// Test with alert level PM2.5
	condition := types.RunningCondition{
		Score:          100,
		Level:          "最高",
		Recommendation: "ランニングに最適な天候です！",
		Warnings:       []string{},
		Clothing:       []string{},
	}

	dustLevel := &types.DustLevel{
		Level:       0,
		DisplayName: "なし",
		Description: "黄砂の影響なし",
		Dust:        5,
		PM10:        100,
		PM2_5:       75, // Above 70, alert level
	}

	ApplyDustPenalty(&condition, dustLevel, nil)

	// Score should be reduced by PM2.5 penalty (30)
	if condition.Score != 70 {
		t.Errorf("Expected score 70, got %d", condition.Score)
	}

	// Should have alert level warning
	hasWarning := false
	for _, warning := range condition.Warnings {
		if warning == "⚠️ PM2.5が注意喚起レベル(70μg/m³超)です。屋外での激しい運動は避けてください" {
			hasWarning = true
			break
		}
	}
	if !hasWarning {
		t.Errorf("Expected PM2.5 alert warning not found, warnings: %v", condition.Warnings)
	}
}

func TestApplyDustPenaltyWithDistance(t *testing.T) {
	// Test with high dust level and full marathon
	condition := types.RunningCondition{
		Score:          100,
		Level:          "最高",
		Recommendation: "ランニングに最適な天候です！",
		Warnings:       []string{},
		Clothing:       []string{},
	}

	dustLevel := &types.DustLevel{
		Level:       3,
		DisplayName: "多い",
		Description: "外出時に注意が必要",
		Dust:        300,
		PM10:        150,
		PM2_5:       30, // Below 35, no PM2.5 penalty
	}

	categoryFull := GetDistanceCategory("full")
	ApplyDustPenalty(&condition, dustLevel, categoryFull)

	// Score should be reduced by dust penalty only: 30 * 2.0 = 60
	if condition.Score != 40 {
		t.Errorf("Expected score 40, got %d", condition.Score)
	}

	// Should have both dust warnings
	warningCount := 0
	for _, warning := range condition.Warnings {
		if warning == "🌫️ 黄砂が飛来しています。マスク着用を推奨します" ||
			warning == "🌫️ 呼吸器系に不安がある方は屋内トレーニングを検討してください" {
			warningCount++
		}
	}
	if warningCount != 2 {
		t.Errorf("Expected 2 dust warnings, got %d", warningCount)
	}

	// Should have both mask and sunglasses
	hasMask := false
	hasSunglasses := false
	for _, item := range condition.Clothing {
		if item == "スポーツマスク" {
			hasMask = true
		}
		if item == "サングラス（目の保護）" {
			hasSunglasses = true
		}
	}
	if !hasMask || !hasSunglasses {
		t.Errorf("Expected sports mask and sunglasses in clothing")
	}
}