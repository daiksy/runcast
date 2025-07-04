package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

const (
	apiURL = "https://api.open-meteo.com/v1/jma"
)

type WeatherData struct {
	Current struct {
		Temperature float64 `json:"temperature_2m"`
		Humidity    int     `json:"relative_humidity_2m"`
		WindSpeed   float64 `json:"wind_speed_10m"`
		WeatherCode int     `json:"weather_code"`
	} `json:"current"`
	Hourly struct {
		Time        []string  `json:"time"`
		Temperature []float64 `json:"temperature_2m"`
		Humidity    []int     `json:"relative_humidity_2m"`
		WindSpeed   []float64 `json:"wind_speed_10m"`
		WeatherCode []int     `json:"weather_code"`
	} `json:"hourly"`
}

type CityCoordinate struct {
	Name string
	Lat  float64
	Lon  float64
}

func main() {
	var city string
	
	flag.StringVar(&city, "city", "Tokyo", "都市名を指定")
	flag.Parse()

	if city == "" {
		fmt.Println("都市名を指定してください: -city <都市名>")
		os.Exit(1)
	}

	coord, err := getCityCoordinate(city)
	if err != nil {
		fmt.Printf("都市の座標取得に失敗しました: %v\n", err)
		os.Exit(1)
	}

	weather, err := getWeather(coord.Lat, coord.Lon)
	if err != nil {
		fmt.Printf("天気情報の取得に失敗しました: %v\n", err)
		os.Exit(1)
	}

	displayWeather(weather, coord.Name)
}

func getCityCoordinate(city string) (*CityCoordinate, error) {
	cities := map[string]CityCoordinate{
		"tokyo":    {"東京", 35.6762, 139.6503},
		"osaka":    {"大阪", 34.6937, 135.5023},
		"kyoto":    {"京都", 35.0116, 135.7681},
		"yokohama": {"横浜", 35.4437, 139.6380},
		"nagoya":   {"名古屋", 35.1815, 136.9066},
		"sapporo":  {"札幌", 43.0642, 141.3469},
		"fukuoka":  {"福岡", 33.5904, 130.4017},
		"sendai":   {"仙台", 38.2682, 140.8694},
		"hiroshima":{"広島", 34.3853, 132.4553},
		"naha":     {"那覇", 26.2124, 127.6792},
	}
	
	if coord, exists := cities[city]; exists {
		return &coord, nil
	}
	
	return nil, fmt.Errorf("都市が見つかりません: %s", city)
}

func getWeather(lat, lon float64) (*WeatherData, error) {
	url := fmt.Sprintf("%s?latitude=%s&longitude=%s&current=temperature_2m,relative_humidity_2m,wind_speed_10m,weather_code&hourly=temperature_2m,relative_humidity_2m,wind_speed_10m,weather_code&timezone=Asia/Tokyo&forecast_days=1", 
		apiURL, 
		strconv.FormatFloat(lat, 'f', 4, 64), 
		strconv.FormatFloat(lon, 'f', 4, 64))
	
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

func displayWeather(weather *WeatherData, cityName string) {
	fmt.Printf("🌤️  %s の天気情報\n", cityName)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	fmt.Printf("🌡️  気温: %.1f°C\n", weather.Current.Temperature)
	fmt.Printf("💧 湿度: %d%%\n", weather.Current.Humidity)
	fmt.Printf("🌬️  風速: %.1f m/s\n", weather.Current.WindSpeed)
	fmt.Printf("☁️  天気: %s\n", getWeatherDescription(weather.Current.WeatherCode))
	
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

func getWeatherDescription(code int) string {
	weatherCodes := map[int]string{
		0:  "快晴",
		1:  "晴れ",
		2:  "一部曇り",
		3:  "曇り",
		45: "霧",
		48: "着氷性の霧",
		51: "弱い霧雨",
		53: "中程度の霧雨",
		55: "強い霧雨",
		61: "弱い雨",
		63: "中程度の雨",
		65: "強い雨",
		71: "弱い雪",
		73: "中程度の雪",
		75: "強い雪",
		80: "弱いにわか雨",
		81: "中程度のにわか雨",
		82: "強いにわか雨",
		95: "雷雨",
		96: "雹を伴う雷雨",
		99: "大粒の雹を伴う雷雨",
	}
	
	if desc, exists := weatherCodes[code]; exists {
		return desc
	}
	return "不明"
}