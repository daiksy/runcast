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
	fmt.Println("例:")
	fmt.Println("  runcast -city=osaka")
	fmt.Println("  runcast -city=tokyo -time=morning")
	fmt.Println("  runcast -city=kyoto -date=tomorrow -distance=10k")
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

	// Display logic - always in running mode
	if *dateSpec != "" {
		dayOffset := weather.GetDateOffset(*dateSpec)
		
		if *timeOfDay != "" {
			// Date + time specific running weather
			display.DisplayDateTimeBasedRunningWeatherWithDistance(weatherData, coord.Name, *dateSpec, *timeOfDay, dayOffset, distanceCategory)
		} else {
			// Date specific running weather (full day)
			display.DisplayDateBasedRunningWeatherWithDistance(weatherData, coord.Name, *dateSpec, dayOffset, distanceCategory)
		}
	} else if *timeOfDay != "" {
		// Time-specific running weather
		display.DisplayTimeBasedRunningWeatherWithDistance(weatherData, coord.Name, *timeOfDay, requiredDays, distanceCategory)
	} else {
		// Current running weather
		display.DisplayRunningWeatherWithDistance(weatherData, coord.Name, distanceCategory)
	}
}