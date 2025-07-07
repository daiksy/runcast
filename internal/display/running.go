package display

import (
	"fmt"
	"rancast/internal/running"
	"rancast/internal/types"
	"rancast/internal/weather"
)

// DisplayRunningForecastWithDistance displays running forecast with distance consideration
func DisplayRunningForecastWithDistance(weatherData *types.WeatherData, cityName string, days int, distanceCategory *types.DistanceCategory) {
	var titleSuffix string
	if distanceCategory != nil {
		titleSuffix = fmt.Sprintf("(%s)", distanceCategory.DisplayName)
	} else {
		titleSuffix = ""
	}
	
	fmt.Printf("🏃‍♂️ %s の%d日間ランニング予報%s\n", cityName, days, titleSuffix)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	// Distance category info
	if distanceCategory != nil {
		fmt.Printf("📏 目標距離: %s (%.1f-%.1fkm)\n", distanceCategory.DisplayName, distanceCategory.MinKm, distanceCategory.MaxKm)
		fmt.Printf("💭 %s\n", distanceCategory.Description)
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	}
	
	// Current condition
	var currentCondition types.RunningCondition
	if distanceCategory != nil {
		currentCondition = running.AssessDistanceBasedRunningCondition(
			weatherData.Current.Temperature,
			weatherData.Current.ApparentTemp,
			float64(weatherData.Current.Humidity),
			weatherData.Current.WindSpeed,
			weatherData.Current.Precipitation,
			weatherData.Current.WeatherCode,
			distanceCategory,
		)
	} else {
		currentCondition = running.AssessRunningCondition(
			weatherData.Current.Temperature,
			weatherData.Current.ApparentTemp,
			float64(weatherData.Current.Humidity),
			weatherData.Current.WindSpeed,
			weatherData.Current.Precipitation,
			weatherData.Current.WeatherCode,
		)
	}
	
	fmt.Printf("📅 現在: %.1f°C | %s | 指数: %d/100\n", weatherData.Current.Temperature, weather.GetWeatherDescription(weatherData.Current.WeatherCode), currentCondition.Score)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	// Daily forecast with running assessment
	for i := 0; i < len(weatherData.Daily.Time) && i < days; i++ {
		date := weatherData.Daily.Time[i]
		maxTemp := weatherData.Daily.TemperatureMax[i]
		minTemp := weatherData.Daily.TemperatureMin[i]
		weatherCode := weatherData.Daily.WeatherCode[i]
		maxWind := weatherData.Daily.WindSpeedMax[i]
		precipitation := weatherData.Daily.PrecipitationSum[i]
		
		// Estimate daily running condition (using average temperature)
		avgTemp := (maxTemp + minTemp) / 2
		var dailyCondition types.RunningCondition
		if distanceCategory != nil {
			dailyCondition = running.AssessDistanceBasedRunningCondition(avgTemp, avgTemp, 60, maxWind, precipitation, weatherCode, distanceCategory)
		} else {
			dailyCondition = running.AssessRunningCondition(avgTemp, avgTemp, 60, maxWind, precipitation, weatherCode)
		}
		
		fmt.Printf("📅 %s\n", weather.FormatDate(date))
		fmt.Printf("   🌡️ %s%.1f°C〜%.1f°C\n", GetRunningTempIcon(avgTemp), minTemp, maxTemp)
		fmt.Printf("   🏆 ランニング指数: %d/100 (%s)\n", dailyCondition.Score, dailyCondition.Level)
		fmt.Printf("   ☁️ %s\n", weather.GetWeatherDescription(weatherCode))
		if precipitation > 0 {
			fmt.Printf("   🌧️ 降水量: %.1f mm\n", precipitation)
		}
		
		if i < len(weatherData.Daily.Time)-1 && i < days-1 {
			fmt.Printf("   ────────────────────────────\n")
		}
	}
	
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

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