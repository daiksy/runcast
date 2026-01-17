package running

import (
	"runcast/internal/types"
)

// DistanceCategory is an alias for types.DistanceCategory for backward compatibility
type DistanceCategory = types.DistanceCategory

// GetDistanceCategories returns all available distance categories
func GetDistanceCategories() []types.DistanceCategory {
	return []types.DistanceCategory{
		{
			Key:             "5k",
			DisplayName:     "5キロ",
			Description:     "短距離ランニング - 比較的軽い負荷",
			MinKm:           3.0,
			MaxKm:           7.0,
			TempPenalty:     0,
			HumidityPenalty: 0,
			WindPenalty:     0,
			HeatIndexPenalty: 0,
		},
		{
			Key:             "10k",
			DisplayName:     "10キロ",
			Description:     "中距離ランニング - 中程度の負荷",
			MinKm:           8.0,
			MaxKm:           12.0,
			TempPenalty:     3,
			HumidityPenalty: 2,
			WindPenalty:     1,
			HeatIndexPenalty: 5,
		},
		{
			Key:             "half",
			DisplayName:     "ハーフマラソン",
			Description:     "長距離ランニング - 高い負荷",
			MinKm:           19.0,
			MaxKm:           23.0,
			TempPenalty:     7,
			HumidityPenalty: 5,
			WindPenalty:     3,
			HeatIndexPenalty: 10,
		},
		{
			Key:             "full",
			DisplayName:     "フルマラソン",
			Description:     "超長距離ランニング - 非常に高い負荷",
			MinKm:           40.0,
			MaxKm:           44.0,
			TempPenalty:     15,
			HumidityPenalty: 10,
			WindPenalty:     5,
			HeatIndexPenalty: 20,
		},
	}
}

// GetDistanceCategory returns distance category by key
func GetDistanceCategory(distance string) *types.DistanceCategory {
	categories := GetDistanceCategories()
	for _, category := range categories {
		if category.Key == distance {
			return &category
		}
	}
	return nil
}

// AssessRunningCondition evaluates running conditions
func AssessRunningCondition(temp, apparentTemp, humidity float64, windSpeed, precipitation float64, weatherCode int) types.RunningCondition {
	score := 100
	var warnings []string
	var clothing []string
	
	// Temperature assessment
	if temp < 5 {
		score -= 30
		warnings = append(warnings, "🥶 低温注意: 防寒対策を十分に行ってください")
		clothing = append(clothing, "長袖", "ロングパンツ", "手袋", "帽子")
	} else if temp < 10 {
		score -= 15
		warnings = append(warnings, "🌡️ 寒冷注意: 適切な服装で体温調節してください")
		clothing = append(clothing, "長袖", "ロングパンツ", "軽い手袋")
	} else if temp < 15 {
		score -= 5
		clothing = append(clothing, "長袖", "ロングパンツ")
	} else if temp < 20 {
		clothing = append(clothing, "薄手の長袖", "ショートパンツ")
	} else if temp < 25 {
		clothing = append(clothing, "薄手の半袖", "ショートパンツ")
	} else if temp < 30 {
		clothing = append(clothing, "薄手の半袖", "帽子推奨")
	} else {
		score -= 20
		warnings = append(warnings, "🔥 高温注意: 早朝や夕方の涼しい時間帯を推奨")
		clothing = append(clothing, "薄手の半袖", "帽子必須", "サングラス")
	}
	
	// Apparent temperature (heat index) assessment
	if apparentTemp > 35 {
		score -= 30
		warnings = append(warnings, "⚠️ 熱中症注意: 体感温度が高すぎます")
	} else if apparentTemp > 32 {
		score -= 15
		warnings = append(warnings, "⚠️ 熱中症注意: 体感温度が高すぎます")
	}
	
	// Humidity assessment
	if humidity > 85 {
		score -= 20
		warnings = append(warnings, "💧 高湿度: 汗が乾きにくい状態です")
	} else if humidity > 70 {
		score -= 10
		warnings = append(warnings, "💧 高湿度: 汗が乾きにくい状態です")
	}
	
	// Wind assessment
	if windSpeed > 10 {
		score -= 25
		warnings = append(warnings, "💨 強風注意: 転倒や怪我のリスクがあります")
	} else if windSpeed > 7 {
		score -= 10
		warnings = append(warnings, "💨 風が強め: 注意してランニングしてください")
	}
	
	// Precipitation assessment
	if precipitation > 5 {
		score -= 40
		warnings = append(warnings, "☔ 大雨: ランニングは控えることをお勧めします")
	} else if precipitation > 1 {
		score -= 25
		warnings = append(warnings, "🌧️ 雨: 滑りやすい路面に注意してください")
	} else if precipitation > 0 {
		score -= 10
		warnings = append(warnings, "🌦️ 小雨: 軽い雨具があると良いでしょう")
	}
	
	// Weather code assessment
	if weatherCode >= 95 {
		score -= 50
		warnings = append(warnings, "⚡ 雷雨: 絶対に屋外でのランニングは避けてください")
	} else if weatherCode >= 80 {
		score -= 30
		warnings = append(warnings, "🌧️ にわか雨: 突然の雨に注意してください")
	}
	
	// Ensure score doesn't go below 0
	if score < 0 {
		score = 0
	}
	
	// Determine level and recommendation
	var level, recommendation string
	switch {
	case score >= 80:
		level = "最高"
		recommendation = "ランニングに最適な天候です！"
	case score >= 60:
		level = "良好"
		recommendation = "良好な天候です。ランニングを楽しんでください"
	case score >= 40:
		level = "普通"
		recommendation = "注意事項を確認してからランニングしてください"
	case score >= 20:
		level = "注意"
		recommendation = "警告事項があります。ランニングは控えめに"
	default:
		level = "危険"
		recommendation = "天候が悪いため、ランニングは控えることをお勧めします"
	}
	
	return types.RunningCondition{
		Score:          score,
		Level:          level,
		Recommendation: recommendation,
		Warnings:       warnings,
		Clothing:       clothing,
	}
}

// AssessDistanceBasedRunningCondition evaluates running conditions with distance-specific penalties
func AssessDistanceBasedRunningCondition(temp, apparentTemp, humidity float64, windSpeed, precipitation float64, weatherCode int, distanceCategory *types.DistanceCategory) types.RunningCondition {
	// Start with base assessment
	condition := AssessRunningCondition(temp, apparentTemp, humidity, windSpeed, precipitation, weatherCode)
	
	if distanceCategory == nil {
		return condition
	}
	
	// Apply distance-specific penalties
	condition.Score -= distanceCategory.TempPenalty
	condition.Score -= distanceCategory.HumidityPenalty
	condition.Score -= distanceCategory.WindPenalty
	condition.Score -= distanceCategory.HeatIndexPenalty
	
	// Distance-specific temperature penalties
	if temp > 28 {
		condition.Score -= distanceCategory.TempPenalty
	}
	if temp > 32 {
		condition.Score -= distanceCategory.TempPenalty * 2
	}
	
	// Distance-specific humidity penalties
	if humidity > 80 {
		condition.Score -= distanceCategory.HumidityPenalty
	}
	if humidity > 90 {
		condition.Score -= distanceCategory.HumidityPenalty * 2
	}
	
	// Distance-specific heat index penalties
	if apparentTemp > 30 {
		condition.Score -= distanceCategory.HeatIndexPenalty
	}
	if apparentTemp > 35 {
		condition.Score -= distanceCategory.HeatIndexPenalty * 2
	}
	
	// Add distance-specific warnings
	if distanceCategory.Key == "half" || distanceCategory.Key == "full" {
		if temp > 25 {
			condition.Warnings = append(condition.Warnings, "🏃‍♂️ 長距離警告: 高温下での長時間運動は危険です")
		}
		if humidity > 70 {
			condition.Warnings = append(condition.Warnings, "💦 長距離警告: 高湿度により脱水リスクが高まります")
		}
		if distanceCategory.Key == "full" && temp > 22 {
			condition.Warnings = append(condition.Warnings, "🏃‍♂️ フルマラソン警告: 高温下での長時間運動は危険です")
		}
	}
	
	// Add distance-specific clothing recommendations
	if distanceCategory.Key == "half" || distanceCategory.Key == "full" {
		if temp > 20 {
			condition.Clothing = append(condition.Clothing, "水分補給用品", "エネルギー補給品")
		}
		if temp > 25 {
			condition.Clothing = append(condition.Clothing, "冷却タオル", "塩分補給品")
		}
	}
	
	// Ensure score doesn't go below 0
	if condition.Score < 0 {
		condition.Score = 0
	}
	
	// Update level and recommendation based on new score
	switch {
	case condition.Score >= 80:
		condition.Level = "最高"
		condition.Recommendation = generateDistanceRecommendation(distanceCategory, "最高")
	case condition.Score >= 60:
		condition.Level = "良好"
		condition.Recommendation = generateDistanceRecommendation(distanceCategory, "良好")
	case condition.Score >= 40:
		condition.Level = "普通"
		condition.Recommendation = generateDistanceRecommendation(distanceCategory, "普通")
	case condition.Score >= 20:
		condition.Level = "注意"
		condition.Recommendation = generateDistanceRecommendation(distanceCategory, "注意")
	default:
		condition.Level = "危険"
		condition.Recommendation = generateDistanceRecommendation(distanceCategory, "危険")
	}
	
	return condition
}

// generateDistanceRecommendation generates distance-specific recommendations
func generateDistanceRecommendation(distanceCategory *types.DistanceCategory, level string) string {
	switch level {
	case "最高":
		return distanceCategory.DisplayName + "実行に最適な天候です！"
	case "良好":
		return distanceCategory.DisplayName + "実行に良好な天候です"
	case "普通":
		return distanceCategory.DisplayName + "実行は慎重に、体調と相談して判断してください"
	case "注意":
		return distanceCategory.DisplayName + "実行は控えめに、短縮も検討してください"
	default:
		return distanceCategory.DisplayName + "実行は控えることをお勧めします"
	}
}

// GetDustPenalty calculates dust penalty for running score
func GetDustPenalty(dustLevel *types.DustLevel) int {
	if dustLevel == nil {
		return 0
	}

	switch dustLevel.Level {
	case 1:
		return 5
	case 2:
		return 15
	case 3:
		return 30
	case 4:
		return 50
	default:
		return 0
	}
}

// GetDistanceDustMultiplier returns dust penalty multiplier for distance
func GetDistanceDustMultiplier(distanceCategory *types.DistanceCategory) float64 {
	if distanceCategory == nil {
		return 1.0
	}

	switch distanceCategory.Key {
	case "10k":
		return 1.2
	case "half":
		return 1.5
	case "full":
		return 2.0
	default:
		return 1.0
	}
}

// ApplyDustPenalty applies dust penalty to running condition
func ApplyDustPenalty(condition *types.RunningCondition, dustLevel *types.DustLevel, distanceCategory *types.DistanceCategory) {
	if dustLevel == nil || dustLevel.Level == 0 {
		return
	}

	basePenalty := GetDustPenalty(dustLevel)
	multiplier := GetDistanceDustMultiplier(distanceCategory)
	totalPenalty := int(float64(basePenalty) * multiplier)

	condition.Score -= totalPenalty
	if condition.Score < 0 {
		condition.Score = 0
	}

	// Add dust-related warnings
	if dustLevel.Level >= 2 {
		condition.Warnings = append(condition.Warnings, "🌫️ 黄砂が飛来しています。マスク着用を推奨します")
	}
	if dustLevel.Level >= 3 {
		condition.Warnings = append(condition.Warnings, "🌫️ 呼吸器系に不安がある方は屋内トレーニングを検討してください")
	}
	if dustLevel.Level >= 4 {
		condition.Warnings = append(condition.Warnings, "⚠️ 黄砂が非常に多いため、屋外でのランニングは避けてください")
	}

	// Add dust-related clothing recommendations
	if dustLevel.Level >= 2 {
		condition.Clothing = append(condition.Clothing, "スポーツマスク")
	}
	if dustLevel.Level >= 3 {
		condition.Clothing = append(condition.Clothing, "サングラス（目の保護）")
	}

	// Update level and recommendation based on new score
	switch {
	case condition.Score >= 80:
		condition.Level = "最高"
		condition.Recommendation = "ランニングに最適な天候です！"
	case condition.Score >= 60:
		condition.Level = "良好"
		condition.Recommendation = "良好な天候です。ランニングを楽しんでください"
	case condition.Score >= 40:
		condition.Level = "普通"
		condition.Recommendation = "注意事項を確認してからランニングしてください"
	case condition.Score >= 20:
		condition.Level = "注意"
		condition.Recommendation = "警告事項があります。ランニングは控えめに"
	default:
		condition.Level = "危険"
		condition.Recommendation = "天候が悪いため、ランニングは控えることをお勧めします"
	}
}