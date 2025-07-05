package display

import (
	"fmt"
	"weather-cli/internal/running"
	"weather-cli/internal/types"
	"weather-cli/internal/weather"
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

// DisplayForecastWeather displays forecast weather information
func DisplayForecastWeather(weatherData *types.WeatherData, cityName string, days int) {
	fmt.Printf("🌤️ %s の%d日間天気予報\n", cityName, days)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	// Current weather
	fmt.Printf("📅 現在: %.1f°C | %s\n", weatherData.Current.Temperature, weather.GetWeatherDescription(weatherData.Current.WeatherCode))
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	// Daily forecast
	for i := 0; i < len(weatherData.Daily.Time) && i < days; i++ {
		date := weatherData.Daily.Time[i]
		maxTemp := weatherData.Daily.TemperatureMax[i]
		minTemp := weatherData.Daily.TemperatureMin[i]
		weatherCode := weatherData.Daily.WeatherCode[i]
		precipitation := weatherData.Daily.PrecipitationSum[i]
		
		fmt.Printf("📅 %s\n", weather.FormatDate(date))
		fmt.Printf("   🌡️ %.1f°C〜%.1f°C\n", minTemp, maxTemp)
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

// DisplayRunningWeatherWithDistance displays running weather with distance consideration
func DisplayRunningWeatherWithDistance(weatherData *types.WeatherData, cityName string, distanceCategory *types.DistanceCategory) {
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