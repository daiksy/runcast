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
	Daily struct {
		Time           []string  `json:"time"`
		TemperatureMax []float64 `json:"temperature_2m_max"`
		TemperatureMin []float64 `json:"temperature_2m_min"`
		WeatherCode    []int     `json:"weather_code"`
		WindSpeedMax   []float64 `json:"wind_speed_10m_max"`
	} `json:"daily"`
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
	var days int
	
	flag.StringVar(&city, "city", "Tokyo", "都市名を指定")
	flag.IntVar(&days, "days", 0, "予報日数を指定（1-7日、0は現在の天気のみ）")
	flag.Parse()

	if city == "" {
		fmt.Println("都市名を指定してください: -city <都市名>")
		os.Exit(1)
	}
	
	if days < 0 || days > 7 {
		fmt.Println("予報日数は0-7日の範囲で指定してください")
		os.Exit(1)
	}

	coord, err := getCityCoordinate(city)
	if err != nil {
		fmt.Printf("都市の座標取得に失敗しました: %v\n", err)
		os.Exit(1)
	}

	weather, err := getWeather(coord.Lat, coord.Lon, days)
	if err != nil {
		fmt.Printf("天気情報の取得に失敗しました: %v\n", err)
		os.Exit(1)
	}

	if days == 0 {
		displayCurrentWeather(weather, coord.Name)
	} else {
		displayForecastWeather(weather, coord.Name, days)
	}
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

func getWeather(lat, lon float64, days int) (*WeatherData, error) {
	forecastDays := 1
	if days > 0 {
		forecastDays = days
	}
	
	var url string
	if days == 0 {
		// 現在の天気のみ
		url = fmt.Sprintf("%s?latitude=%s&longitude=%s&current=temperature_2m,relative_humidity_2m,wind_speed_10m,weather_code&timezone=Asia/Tokyo", 
			apiURL, 
			strconv.FormatFloat(lat, 'f', 4, 64), 
			strconv.FormatFloat(lon, 'f', 4, 64))
	} else {
		// 予報データ
		url = fmt.Sprintf("%s?latitude=%s&longitude=%s&current=temperature_2m,relative_humidity_2m,wind_speed_10m,weather_code&daily=temperature_2m_max,temperature_2m_min,weather_code,wind_speed_10m_max&timezone=Asia/Tokyo&forecast_days=%d", 
			apiURL, 
			strconv.FormatFloat(lat, 'f', 4, 64), 
			strconv.FormatFloat(lon, 'f', 4, 64),
			forecastDays)
	}
	
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

func displayCurrentWeather(weather *WeatherData, cityName string) {
	fmt.Printf("🌤️  %s の現在の天気\n", cityName)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	fmt.Printf("🌡️  気温: %.1f°C\n", weather.Current.Temperature)
	fmt.Printf("💧 湿度: %d%%\n", weather.Current.Humidity)
	fmt.Printf("🌬️  風速: %.1f m/s\n", weather.Current.WindSpeed)
	fmt.Printf("☁️  天気: %s\n", getWeatherDescription(weather.Current.WeatherCode))
	
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

func displayForecastWeather(weather *WeatherData, cityName string, days int) {
	fmt.Printf("🌤️  %s の%d日間天気予報\n", cityName, days)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	// 現在の天気を表示
	fmt.Printf("📅 現在: %.1f°C | %s\n", weather.Current.Temperature, getWeatherDescription(weather.Current.WeatherCode))
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	// 予報データを表示
	for i := 0; i < len(weather.Daily.Time) && i < days; i++ {
		date := weather.Daily.Time[i]
		maxTemp := weather.Daily.TemperatureMax[i]
		minTemp := weather.Daily.TemperatureMin[i]
		weatherCode := weather.Daily.WeatherCode[i]
		maxWind := weather.Daily.WindSpeedMax[i]
		
		fmt.Printf("📅 %s\n", formatDate(date))
		fmt.Printf("   🌡️  最高: %.1f°C / 最低: %.1f°C\n", maxTemp, minTemp)
		fmt.Printf("   🌬️  最大風速: %.1f m/s\n", maxWind)
		fmt.Printf("   ☁️  天気: %s\n", getWeatherDescription(weatherCode))
		
		if i < len(weather.Daily.Time)-1 && i < days-1 {
			fmt.Printf("   ────────────────────────────\n")
		}
	}
	
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

func formatDate(dateStr string) string {
	// 簡単な日付パース（YYYY-MM-DD形式想定）
	if len(dateStr) >= 10 {
		month := dateStr[5:7]
		day := dateStr[8:10]
		return fmt.Sprintf("%s月%s日", month, day)
	}
	return dateStr
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