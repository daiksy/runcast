package display

import (
	"fmt"
	"runcast/internal/running"
	"runcast/internal/types"
	"runcast/internal/weather"
)

// GetRunningTempIcon returns temperature icon for running
func GetRunningTempIcon(temp float64) string {
	if temp >= 30 {
		return "🔥 "
	} else if temp >= 25 {
		return "☀️ "
	} else if temp >= 20 {
		return "🌤️ "
	} else if temp >= 15 {
		return "🌥️ "
	} else if temp >= 10 {
		return "☁️ "
	} else {
		return "❄️ "
	}
}

// DisplayCurrentWeather displays current weather information
func DisplayCurrentWeather(weatherData *types.WeatherData, cityName string) {
	fmt.Printf("🌤️ %s の現在の天気\n", cityName)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("🌡️ 気温: %.1f°C (体感: %.1f°C)\n", weatherData.Current.Temperature, weatherData.Current.ApparentTemp)
	fmt.Printf("💧 湿度: %d%%\n", weatherData.Current.Humidity)
	fmt.Printf("🌬️ 風: %s %.1f m/s\n", weather.GetWindDirection(weatherData.Current.WindDirection), weatherData.Current.WindSpeed)
	fmt.Printf("☁️ 天気: %s\n", weather.GetWeatherDescription(weatherData.Current.WeatherCode))
	if weatherData.Current.Precipitation > 0 {
		fmt.Printf("🌧️ 降水量: %.1f mm\n", weatherData.Current.Precipitation)
	}
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

// DisplayRunningWeatherWithDistance displays running weather with distance consideration
func DisplayRunningWeatherWithDistance(weatherData *types.WeatherData, cityName string, distanceCategory *types.DistanceCategory) {
	DisplayRunningWeatherWithDistanceAndDust(weatherData, cityName, distanceCategory, nil)
}

// DisplayRunningWeatherWithDistanceAndDust displays running weather with distance and dust consideration
func DisplayRunningWeatherWithDistanceAndDust(weatherData *types.WeatherData, cityName string, distanceCategory *types.DistanceCategory, dustLevel *types.DustLevel) {
	var condition types.RunningCondition
	var titleSuffix string

	if distanceCategory != nil {
		titleSuffix = fmt.Sprintf("(%s)", distanceCategory.DisplayName)
		condition = running.AssessDistanceBasedRunningCondition(
			weatherData.Current.Temperature,
			weatherData.Current.ApparentTemp,
			float64(weatherData.Current.Humidity),
			weatherData.Current.WindSpeed,
			weatherData.Current.Precipitation,
			weatherData.Current.WeatherCode,
			distanceCategory,
		)
	} else {
		titleSuffix = ""
		condition = running.AssessRunningCondition(
			weatherData.Current.Temperature,
			weatherData.Current.ApparentTemp,
			float64(weatherData.Current.Humidity),
			weatherData.Current.WindSpeed,
			weatherData.Current.Precipitation,
			weatherData.Current.WeatherCode,
		)
	}

	// Apply dust penalty
	running.ApplyDustPenalty(&condition, dustLevel, distanceCategory)

	fmt.Printf("🏃‍♂️ %s のランニング情報%s\n", cityName, titleSuffix)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	// Distance category info
	if distanceCategory != nil {
		fmt.Printf("📏 目標距離: %s (%.1f-%.1fkm)\n", distanceCategory.DisplayName, distanceCategory.MinKm, distanceCategory.MaxKm)
		fmt.Printf("💭 %s\n", distanceCategory.Description)
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	}

	fmt.Printf("🏆 ランニング指数: %d/100 (%s)\n", condition.Score, condition.Level)
	fmt.Printf("💡 %s\n", condition.Recommendation)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	fmt.Printf("🌡️ 気温: %.1f°C (体感: %.1f°C)\n", weatherData.Current.Temperature, weatherData.Current.ApparentTemp)
	fmt.Printf("💧 湿度: %d%%\n", weatherData.Current.Humidity)
	fmt.Printf("🌬️ 風: %s %.1f m/s\n", weather.GetWindDirection(weatherData.Current.WindDirection), weatherData.Current.WindSpeed)
	fmt.Printf("☁️ 天気: %s\n", weather.GetWeatherDescription(weatherData.Current.WeatherCode))
	if weatherData.Current.Precipitation > 0 {
		fmt.Printf("🌧️ 降水量: %.1f mm\n", weatherData.Current.Precipitation)
	}

	// Dust information
	if dustLevel != nil && dustLevel.Level > 0 {
		fmt.Printf("🌫️ 黄砂: %s (%.0f μg/m³)\n", dustLevel.DisplayName, dustLevel.Dust)
		fmt.Printf("   PM2.5: %.0f μg/m³ / PM10: %.0f μg/m³\n", dustLevel.PM2_5, dustLevel.PM10)
	}

	// Clothing recommendations
	if len(condition.Clothing) > 0 {
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("👕 推奨ウェア:\n")
		for _, item := range condition.Clothing {
			fmt.Printf("   • %s\n", item)
		}
	}

	// Warnings
	if len(condition.Warnings) > 0 {
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("⚠️ 注意事項:\n")
		for _, warning := range condition.Warnings {
			fmt.Printf("   %s\n", warning)
		}
	}

	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}