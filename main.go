package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	apiURL = "https://api.openweathermap.org/data/2.5/weather"
	apiKey = "demo" // デモ用、実際は環境変数から取得
)

type WeatherData struct {
	Name string `json:"name"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
}

func main() {
	var city string
	var units string
	
	flag.StringVar(&city, "city", "Tokyo", "都市名を指定")
	flag.StringVar(&units, "units", "metric", "単位 (metric, imperial, kelvin)")
	flag.Parse()

	if city == "" {
		fmt.Println("都市名を指定してください: -city <都市名>")
		os.Exit(1)
	}

	weather, err := getWeather(city, units)
	if err != nil {
		fmt.Printf("天気情報の取得に失敗しました: %v\n", err)
		os.Exit(1)
	}

	displayWeather(weather, units)
}

func getWeather(city, units string) (*WeatherData, error) {
	url := fmt.Sprintf("%s?q=%s&appid=%s&units=%s", apiURL, city, apiKey, units)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weather WeatherData
	if err := json.Unmarshal(body, &weather); err != nil {
		return nil, err
	}

	return &weather, nil
}

func displayWeather(weather *WeatherData, units string) {
	fmt.Printf("🌤️  %s の天気情報\n", weather.Name)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	// 温度単位の表示
	tempUnit := "°C"
	if units == "imperial" {
		tempUnit = "°F"
	} else if units == "kelvin" {
		tempUnit = "K"
	}
	
	fmt.Printf("🌡️  気温: %.1f%s\n", weather.Main.Temp, tempUnit)
	fmt.Printf("💧 湿度: %d%%\n", weather.Main.Humidity)
	fmt.Printf("🌬️  風速: %.1f m/s\n", weather.Wind.Speed)
	
	if len(weather.Weather) > 0 {
		fmt.Printf("☁️  天気: %s (%s)\n", weather.Weather[0].Main, weather.Weather[0].Description)
	}
	
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}