package display

import (
	"fmt"
	"runcast/internal/running"
	"runcast/internal/types"
	"runcast/internal/weather"
)

// DisplayTimeBasedWeather displays time-based weather information
func DisplayTimeBasedWeather(weatherData *types.WeatherData, cityName, timeOfDay string, days int) {
	periods := weather.GetTimePeriods()
	period := periods[timeOfDay]
	
	timeData := weather.ExtractTimeBasedWeather(weatherData, timeOfDay, days)
	if len(timeData) == 0 {
		fmt.Println("指定された時間帯のデータが見つかりません")
		return
	}
	
	fmt.Printf("🌤️ %s の%s時間帯天気情報\n", cityName, period.DisplayName)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	for i, data := range timeData {
		hour := weather.ExtractHour(data.Time)
		temp := data.Temperature
		weatherDesc := weather.GetWeatherDescription(data.WeatherCode)
		
		fmt.Printf("📅 %s時: %.1f°C | %s", hour, temp, weatherDesc)
		if data.Precipitation > 0 {
			fmt.Printf(" | 🌧️ %.1fmm", data.Precipitation)
		}
		fmt.Printf("\n")
		
		if (i+1)%3 == 0 && i < len(timeData)-1 {
			fmt.Printf("   ────────────────────────────\n")
		}
	}
	
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

// DisplayDateBasedWeather displays date-based weather information
func DisplayDateBasedWeather(weatherData *types.WeatherData, cityName, dateSpec string, dayOffset int) {
	dateSpecificWeather := weather.ExtractDateBasedWeather(weatherData, dayOffset)
	
	dateDisplayName := weather.GetDateDisplayName(dateSpec)
	fmt.Printf("🌤️ %s の%s天気情報\n", cityName, dateDisplayName)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	if len(dateSpecificWeather.Daily.Time) == 0 {
		fmt.Println("指定された日付のデータが見つかりません")
		return
	}
	
	// Daily summary
	date := dateSpecificWeather.Daily.Time[0]
	maxTemp := dateSpecificWeather.Daily.TemperatureMax[0]
	minTemp := dateSpecificWeather.Daily.TemperatureMin[0]
	weatherCode := dateSpecificWeather.Daily.WeatherCode[0]
	precipitation := dateSpecificWeather.Daily.PrecipitationSum[0]
	
	fmt.Printf("📅 %s (%s)\n", weather.FormatDate(date), dateDisplayName)
	fmt.Printf("🌡️ %.1f°C〜%.1f°C\n", minTemp, maxTemp)
	fmt.Printf("☁️ %s\n", weather.GetWeatherDescription(weatherCode))
	if precipitation > 0 {
		fmt.Printf("🌧️ 降水量: %.1f mm\n", precipitation)
	}
	
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

// DisplayDateBasedRunningWeatherWithDistance displays date-based running weather with distance consideration
func DisplayDateBasedRunningWeatherWithDistance(weatherData *types.WeatherData, cityName, dateSpec string, dayOffset int, distanceCategory *types.DistanceCategory) {
	DisplayDateBasedRunningWeatherWithDistanceAndDust(weatherData, cityName, dateSpec, dayOffset, distanceCategory, nil)
}

// DisplayDateBasedRunningWeatherWithDistanceAndDust displays date-based running weather with distance and dust consideration
func DisplayDateBasedRunningWeatherWithDistanceAndDust(weatherData *types.WeatherData, cityName, dateSpec string, dayOffset int, distanceCategory *types.DistanceCategory, dustLevel *types.DustLevel) {
	dateSpecificWeather := weather.ExtractDateBasedWeather(weatherData, dayOffset)
	
	dateDisplayName := weather.GetDateDisplayName(dateSpec)
	var titleSuffix string
	if distanceCategory != nil {
		titleSuffix = fmt.Sprintf("(%s)", distanceCategory.DisplayName)
	} else {
		titleSuffix = ""
	}
	
	fmt.Printf("🏃‍♂️ %s の%sランニング情報%s\n", cityName, dateDisplayName, titleSuffix)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	if len(dateSpecificWeather.Daily.Time) == 0 {
		fmt.Println("指定された日付のデータが見つかりません")
		return
	}
	
	// Distance category info
	if distanceCategory != nil {
		fmt.Printf("📏 目標距離: %s (%.1f-%.1fkm)\n", distanceCategory.DisplayName, distanceCategory.MinKm, distanceCategory.MaxKm)
		fmt.Printf("💭 %s\n", distanceCategory.Description)
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	}
	
	// Daily summary
	date := dateSpecificWeather.Daily.Time[0]
	maxTemp := dateSpecificWeather.Daily.TemperatureMax[0]
	minTemp := dateSpecificWeather.Daily.TemperatureMin[0]
	weatherCode := dateSpecificWeather.Daily.WeatherCode[0]
	maxWind := dateSpecificWeather.Daily.WindSpeedMax[0]
	precipitation := dateSpecificWeather.Daily.PrecipitationSum[0]
	
	// Estimate daily running condition (using average temperature)
	avgTemp := (maxTemp + minTemp) / 2
	var dailyCondition types.RunningCondition
	if distanceCategory != nil {
		dailyCondition = running.AssessDistanceBasedRunningCondition(avgTemp, avgTemp, 60, maxWind, precipitation, weatherCode, distanceCategory)
	} else {
		dailyCondition = running.AssessRunningCondition(avgTemp, avgTemp, 60, maxWind, precipitation, weatherCode)
	}

	// Apply dust penalty
	running.ApplyDustPenalty(&dailyCondition, dustLevel, distanceCategory)

	fmt.Printf("📅 %s (%s)\n", weather.FormatDate(date), dateDisplayName)
	fmt.Printf("🏆 ランニング指数: %d/100 (%s)\n", dailyCondition.Score, dailyCondition.Level)
	fmt.Printf("💡 %s\n", dailyCondition.Recommendation)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	fmt.Printf("🌡️ %s%.1f°C〜%.1f°C\n", GetRunningTempIcon(avgTemp), minTemp, maxTemp)
	fmt.Printf("☁️ %s\n", weather.GetWeatherDescription(weatherCode))
	if precipitation > 0 {
		fmt.Printf("🌧️ 降水量: %.1f mm\n", precipitation)
	}

	// Dust information
	if dustLevel != nil {
		fmt.Printf("🌫️ 黄砂: %s (%.0f μg/m³)\n", dustLevel.DisplayName, dustLevel.Dust)
		fmt.Printf("   PM2.5: %.0f μg/m³ / PM10: %.0f μg/m³\n", dustLevel.PM2_5, dustLevel.PM10)
	}

	// Clothing recommendations
	if len(dailyCondition.Clothing) > 0 {
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("👕 推奨ウェア:\n")
		for _, item := range dailyCondition.Clothing {
			fmt.Printf("   • %s\n", item)
		}
	}

	// Warnings
	if len(dailyCondition.Warnings) > 0 {
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("⚠️ 注意事項:\n")
		for _, warning := range dailyCondition.Warnings {
			fmt.Printf("   %s\n", warning)
		}
	}

	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

// DisplayDateTimeBasedWeather displays date and time based weather information
func DisplayDateTimeBasedWeather(weatherData *types.WeatherData, cityName, dateSpec, timeOfDay string, dayOffset int) {
	dateSpecificWeather := weather.ExtractDateBasedWeather(weatherData, dayOffset)
	
	periods := weather.GetTimePeriods()
	period := periods[timeOfDay]
	dateDisplayName := weather.GetDateDisplayName(dateSpec)
	
	timeData := weather.ExtractTimeBasedWeather(dateSpecificWeather, timeOfDay, 1)
	if len(timeData) == 0 {
		fmt.Println("指定された日付・時間帯のデータが見つかりません")
		return
	}
	
	fmt.Printf("🌤️ %s の%s%s時間帯天気情報\n", cityName, dateDisplayName, period.DisplayName)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	for i, data := range timeData {
		hour := weather.ExtractHour(data.Time)
		temp := data.Temperature
		weatherDesc := weather.GetWeatherDescription(data.WeatherCode)
		
		fmt.Printf("📅 %s時: %.1f°C | %s", hour, temp, weatherDesc)
		if data.Precipitation > 0 {
			fmt.Printf(" | 🌧️ %.1fmm", data.Precipitation)
		}
		fmt.Printf("\n")
		
		if (i+1)%3 == 0 && i < len(timeData)-1 {
			fmt.Printf("   ────────────────────────────\n")
		}
	}
	
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

// DisplayDateTimeBasedRunningWeatherWithDistance displays date and time based running weather with distance consideration
func DisplayDateTimeBasedRunningWeatherWithDistance(weatherData *types.WeatherData, cityName, dateSpec, timeOfDay string, dayOffset int, distanceCategory *types.DistanceCategory) {
	DisplayDateTimeBasedRunningWeatherWithDistanceAndDust(weatherData, cityName, dateSpec, timeOfDay, dayOffset, distanceCategory, nil)
}

// DisplayDateTimeBasedRunningWeatherWithDistanceAndDust displays date and time based running weather with distance and dust consideration
func DisplayDateTimeBasedRunningWeatherWithDistanceAndDust(weatherData *types.WeatherData, cityName, dateSpec, timeOfDay string, dayOffset int, distanceCategory *types.DistanceCategory, airQuality *types.AirQualityData) {
	dateSpecificWeather := weather.ExtractDateBasedWeather(weatherData, dayOffset)
	
	periods := weather.GetTimePeriods()
	period := periods[timeOfDay]
	dateDisplayName := weather.GetDateDisplayName(dateSpec)
	
	timeData := weather.ExtractTimeBasedWeather(dateSpecificWeather, timeOfDay, 1)
	if len(timeData) == 0 {
		fmt.Println("指定された日付・時間帯のデータが見つかりません")
		return
	}
	
	var titleSuffix string
	if distanceCategory != nil {
		titleSuffix = fmt.Sprintf("(%s)", distanceCategory.DisplayName)
	} else {
		titleSuffix = ""
	}
	
	fmt.Printf("🏃‍♂️ %s の%s%s時間帯ランニング情報%s\n", cityName, dateDisplayName, period.DisplayName, titleSuffix)
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
	
	fmt.Printf("⏰ %s%s時間帯詳細 (%d:00-%d:00)\n", dateDisplayName, period.DisplayName, period.StartHour, period.EndHour)
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

		// Get dust level for this hour
		hour := weather.ExtractHour(data.Time)
		hourInt := weather.ExtractHourInt(data.Time)
		dustLevel := weather.GetHourlyDustLevel(airQuality, hourInt, dayOffset)
		running.ApplyDustPenalty(&condition, dustLevel, distanceCategory)

		fmt.Printf("🕐 %s時: %d/100 (%s)\n", hour, condition.Score, condition.Level)
		fmt.Printf("   🌡️ %.1f°C (体感: %.1f°C) | 💧 %d%%\n",
			data.Temperature, data.ApparentTemp, data.Humidity)
		fmt.Printf("   ☁️ %s", weather.GetWeatherDescription(data.WeatherCode))
		if data.Precipitation > 0 {
			fmt.Printf(" | 🌧️ %.1fmm", data.Precipitation)
		}
		if dustLevel != nil {
			fmt.Printf(" | 🌫️ %s", dustLevel.DisplayName)
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