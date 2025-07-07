package display

import (
	"fmt"
	"runcast/internal/running"
	"runcast/internal/types"
	"runcast/internal/weather"
)

// DisplayTimeBasedRunningWeatherWithDistance displays time-based running weather with distance consideration
func DisplayTimeBasedRunningWeatherWithDistance(weatherData *types.WeatherData, cityName, timeOfDay string, days int, distanceCategory *types.DistanceCategory) {
	periods := weather.GetTimePeriods()
	period := periods[timeOfDay]
	
	timeData := weather.ExtractTimeBasedWeather(weatherData, timeOfDay, days)
	if len(timeData) == 0 {
		fmt.Println("指定された時間帯のデータが見つかりません")
		return
	}
	
	var titleSuffix string
	if distanceCategory != nil {
		titleSuffix = fmt.Sprintf("(%s)", distanceCategory.DisplayName)
	} else {
		titleSuffix = ""
	}
	
	fmt.Printf("🏃‍♂️ %s の%s時間帯ランニング情報%s\n", cityName, period.DisplayName, titleSuffix)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	// Distance category info
	if distanceCategory != nil {
		fmt.Printf("📏 目標距離: %s (%.1f-%.1fkm)\n", distanceCategory.DisplayName, distanceCategory.MinKm, distanceCategory.MaxKm)
		fmt.Printf("💭 %s\n", distanceCategory.Description)
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	}
	
	bestCondition := types.TimeBasedWeather{}
	bestScore := -1
	bestTime := ""
	
	fmt.Printf("⏰ %s時間帯詳細 (%d:00-%d:00)\n", period.DisplayName, period.StartHour, period.EndHour)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	for _, data := range timeData {
		var condition types.RunningCondition
		if distanceCategory != nil {
			condition = running.AssessDistanceBasedRunningCondition(
				data.Temperature,
				data.ApparentTemp,
				float64(data.Humidity),
				data.WindSpeed,
				data.Precipitation,
				data.WeatherCode,
				distanceCategory,
			)
		} else {
			condition = running.AssessRunningCondition(
				data.Temperature,
				data.ApparentTemp,
				float64(data.Humidity),
				data.WindSpeed,
				data.Precipitation,
				data.WeatherCode,
			)
		}
		
		hour := weather.ExtractHour(data.Time)
		fmt.Printf("🕐 %s時: %d/100 (%s)\n", hour, condition.Score, condition.Level)
		fmt.Printf("   🌡️ %.1f°C (体感: %.1f°C) | 💧 %d%%\n", 
			data.Temperature, data.ApparentTemp, data.Humidity)
		fmt.Printf("   ☁️ %s", weather.GetWeatherDescription(data.WeatherCode))
		if data.Precipitation > 0 {
			fmt.Printf(" | 🌧️ %.1fmm", data.Precipitation)
		}
		fmt.Printf("\n")
		
		if condition.Score > bestScore {
			bestScore = condition.Score
			bestCondition = data
			bestTime = hour
		}
		
		fmt.Printf("   ────────────────────────────\n")
	}
	
	// Best time recommendation
	if bestScore >= 0 {
		var bestRunningCondition types.RunningCondition
		if distanceCategory != nil {
			bestRunningCondition = running.AssessDistanceBasedRunningCondition(
				bestCondition.Temperature,
				bestCondition.ApparentTemp,
				float64(bestCondition.Humidity),
				bestCondition.WindSpeed,
				bestCondition.Precipitation,
				bestCondition.WeatherCode,
				distanceCategory,
			)
		} else {
			bestRunningCondition = running.AssessRunningCondition(
				bestCondition.Temperature,
				bestCondition.ApparentTemp,
				float64(bestCondition.Humidity),
				bestCondition.WindSpeed,
				bestCondition.Precipitation,
				bestCondition.WeatherCode,
			)
		}
		
		fmt.Printf("🏆 最適時間: %s時 (スコア: %d/100)\n", bestTime, bestScore)
		fmt.Printf("💡 %s\n", bestRunningCondition.Recommendation)
		
		if len(bestRunningCondition.Warnings) > 0 {
			fmt.Printf("⚠️ 注意事項:\n")
			for _, warning := range bestRunningCondition.Warnings {
				fmt.Printf("   %s\n", warning)
			}
		}
	}
	
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}