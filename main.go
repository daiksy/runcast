package main

import (
	"flag"
	"fmt"
	"log"

	"weather-cli/internal/display"
	"weather-cli/internal/running"
	"weather-cli/internal/types"
	"weather-cli/internal/weather"
)

func main() {
	city := flag.String("city", "tokyo", "都市名を指定")
	days := flag.Int("days", 0, "予報日数を指定 (0は現在の天気のみ)")
	runningMode := flag.Bool("running", false, "ランニング向け情報を表示")
	timeOfDay := flag.String("time", "", "時間帯を指定 (morning, noon, evening, night)")
	dateSpec := flag.String("date", "", "日付を指定 (today, tomorrow, day-after-tomorrow)")
	distanceFlag := flag.String("distance", "", "目標距離を指定 (5k, 10k, half, full)")
	flag.Parse()

	// Distance category processing
	var distanceCategory *types.DistanceCategory
	if *distanceFlag != "" {
		distanceCategory = running.GetDistanceCategory(*distanceFlag)
		if distanceCategory == nil {
			fmt.Printf("無効な距離です: %s\n", *distanceFlag)
			fmt.Println("有効な距離: 5k, 10k, half, full")
			return
		}
		*runningMode = true // Auto-enable running mode
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
	requiredDays := *days
	if *dateSpec != "" {
		dayOffset := weather.GetDateOffset(*dateSpec)
		// Ensure we have enough data for the requested date
		if requiredDays <= dayOffset {
			requiredDays = dayOffset + 1
		}
	} else if *timeOfDay != "" && requiredDays == 0 {
		// Time-specific queries need at least 1 day of forecast data for hourly data
		requiredDays = 1
	}
	
	// Get weather data
	weatherData, err := weather.GetWeather(coord.Lat, coord.Lon, requiredDays)
	if err != nil {
		log.Fatal(err)
	}

	// Display logic
	if *dateSpec != "" {
		dayOffset := weather.GetDateOffset(*dateSpec)
		
		if *timeOfDay != "" {
			// Date + time specific weather
			if *runningMode {
				display.DisplayDateTimeBasedRunningWeatherWithDistance(weatherData, coord.Name, *dateSpec, *timeOfDay, dayOffset, distanceCategory)
			} else {
				display.DisplayDateTimeBasedWeather(weatherData, coord.Name, *dateSpec, *timeOfDay, dayOffset)
			}
		} else {
			// Date specific weather (full day)
			if *runningMode {
				display.DisplayDateBasedRunningWeatherWithDistance(weatherData, coord.Name, *dateSpec, dayOffset, distanceCategory)
			} else {
				display.DisplayDateBasedWeather(weatherData, coord.Name, *dateSpec, dayOffset)
			}
		}
	} else if *timeOfDay != "" {
		// Time-specific weather
		if *runningMode {
			display.DisplayTimeBasedRunningWeatherWithDistance(weatherData, coord.Name, *timeOfDay, requiredDays, distanceCategory)
		} else {
			display.DisplayTimeBasedWeather(weatherData, coord.Name, *timeOfDay, requiredDays)
		}
	} else if *runningMode {
		if *days == 0 {
			display.DisplayRunningWeatherWithDistance(weatherData, coord.Name, distanceCategory)
		} else {
			display.DisplayRunningForecastWithDistance(weatherData, coord.Name, *days, distanceCategory)
		}
	} else {
		if *days == 0 {
			display.DisplayCurrentWeather(weatherData, coord.Name)
		} else {
			display.DisplayForecastWeather(weatherData, coord.Name, *days)
		}
	}
}