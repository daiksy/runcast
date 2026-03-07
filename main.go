package main

import (
	"flag"
	"fmt"
	"log"

	"runcast/internal/display"
	"runcast/internal/running"
	"runcast/internal/types"
	"runcast/internal/weather"
)

func showHelp() {
	fmt.Println("🏃‍♂️ runcast - ランニング天気予報")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("使用方法:")
	fmt.Println("  runcast [オプション]")
	fmt.Println()
	fmt.Println("オプション:")
	fmt.Println("  -city string")
	fmt.Println("      都市名を指定 (デフォルト: tokyo)")
	fmt.Println("  -time string")
	fmt.Println("      時間帯を指定 (morning, noon, evening, night)")
	fmt.Println("  -date string")
	fmt.Println("      日付を指定 (today, tomorrow, day-after-tomorrow)")
	fmt.Println("  -distance string")
	fmt.Println("      目標距離を指定 (5k, 10k, half, full)")
	fmt.Println("  -help")
	fmt.Println("      このヘルプを表示")
	fmt.Println()
	fmt.Println("対応都市:")
	supportedCities := weather.GetSupportedCities()
	for i, city := range supportedCities {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(city)
	}
	fmt.Println()
	fmt.Println()
	fmt.Println("カスタム位置設定:")
	fmt.Println("  .runcast.conf ファイルを作成することで任意の位置を追加できます")
	fmt.Println("  設定ファイルの配置場所（優先順）:")
	fmt.Println("    1. カレントディレクトリ: .runcast.conf")
	fmt.Println("    2. ホームディレクトリ: ~/.runcast.conf")
	fmt.Println("    3. 設定ディレクトリ: ~/.config/runcast/config.toml")
	fmt.Println()
	fmt.Println("  設定ファイルの例:")
	fmt.Println("    [locations]")
	fmt.Println("    home = { name = \"自宅\", lat = 35.6762, lon = 139.6503 }")
	fmt.Println("    office = { name = \"会社\", lat = 35.6584, lon = 139.7016 }")
	fmt.Println()
	fmt.Println("例:")
	fmt.Println("  runcast -city=osaka")
	fmt.Println("  runcast -city=tokyo -time=morning")
	fmt.Println("  runcast -city=kyoto -date=tomorrow -distance=10k")
	fmt.Println("  runcast -city=home    # カスタム位置を使用")
}

func main() {
	city := flag.String("city", "tokyo", "都市名を指定")
	timeOfDay := flag.String("time", "", "時間帯を指定 (morning, noon, evening, night)")
	dateSpec := flag.String("date", "", "日付を指定 (today, tomorrow, day-after-tomorrow)")
	distanceFlag := flag.String("distance", "", "目標距離を指定 (5k, 10k, half, full)")
	help := flag.Bool("help", false, "ヘルプを表示")
	flag.Parse()

	// Show help if requested
	if *help {
		showHelp()
		return
	}

	// Distance category processing
	var distanceCategory *types.DistanceCategory
	if *distanceFlag != "" {
		distanceCategory = running.GetDistanceCategory(*distanceFlag)
		if distanceCategory == nil {
			fmt.Printf("無効な距離です: %s\n", *distanceFlag)
			fmt.Println("有効な距離: 5k, 10k, half, full")
			return
		}
	}

	// Get city coordinates
	coord, err := weather.GetCityCoordinate(*city)
	if err != nil {
		log.Fatal(err)
	}

	// Validate date specification if provided
	if *dateSpec != "" && !weather.ValidateDateSpec(*dateSpec) {
		fmt.Printf("無効な日付指定です: %s\n", *dateSpec)
		fmt.Println("有効な日付: today, tomorrow, day-after-tomorrow")
		return
	}

	// Validate time specification if provided
	if *timeOfDay != "" && !weather.ValidateTimeSpec(*timeOfDay) {
		fmt.Printf("無効な時間指定です: %s\n", *timeOfDay)
		fmt.Println("有効な時間: morning, noon, evening, night")
		return
	}

	// Determine required forecast days
	requiredDays := 1 // Default to 1 day for running forecasts
	if *dateSpec != "" {
		dayOffset := weather.GetDateOffset(*dateSpec)
		// Ensure we have enough data for the requested date
		if requiredDays <= dayOffset {
			requiredDays = dayOffset + 1
		}
	}
	
	// Get weather data
	weatherData, err := weather.GetWeather(coord.Lat, coord.Lon)
	if err != nil {
		log.Fatal(err)
	}

	// Get air quality data
	airQuality, err := weather.GetAirQuality(coord.Lat, coord.Lon)
	if err != nil {
		// Air quality data is optional, continue without it
		fmt.Printf("警告: 大気質データの取得に失敗しました: %v\n", err)
		airQuality = nil
	}

	// Get pollen data (only for supported cities)
	var pollenData []types.PollenData
	if cityCode, ok := weather.GetCityCode(*city); ok {
		dayOffset := 0
		if *dateSpec != "" {
			dayOffset = weather.GetDateOffset(*dateSpec)
		}
		pollenData, err = weather.GetPollen(cityCode, dayOffset)
		if err != nil {
			// Pollen data is optional, continue without it
			fmt.Printf("警告: 花粉データの取得に失敗しました: %v\n", err)
			pollenData = nil
		}
	}

	// Display logic - always in running mode
	if *dateSpec != "" {
		dayOffset := weather.GetDateOffset(*dateSpec)

		if *timeOfDay != "" {
			// Date + time specific running weather
			display.DisplayDateTimeBasedRunningWeatherFull(weatherData, coord.Name, *dateSpec, *timeOfDay, dayOffset, distanceCategory, airQuality, pollenData)
		} else {
			// Date specific running weather (full day)
			dustLevel := weather.GetCurrentDustLevel(airQuality)
			pollenLevel := weather.GetCurrentPollenLevel(pollenData)
			display.DisplayDateBasedRunningWeatherFull(weatherData, coord.Name, *dateSpec, dayOffset, distanceCategory, dustLevel, pollenLevel)
		}
	} else if *timeOfDay != "" {
		// Time-specific running weather
		display.DisplayTimeBasedRunningWeatherFull(weatherData, coord.Name, *timeOfDay, requiredDays, distanceCategory, airQuality, pollenData)
	} else {
		// Current running weather
		dustLevel := weather.GetCurrentDustLevel(airQuality)
		pollenLevel := weather.GetCurrentPollenLevel(pollenData)
		display.DisplayRunningWeatherFull(weatherData, coord.Name, distanceCategory, dustLevel, pollenLevel)
	}
}