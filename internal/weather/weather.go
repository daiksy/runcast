package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"
	"runcast/internal/types"
)

const apiURL = "https://api.open-meteo.com/v1/jma"

// Cities holds all supported cities
var Cities = map[string]types.CityCoordinate{
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
	"kobe":     {"神戸", 34.6901, 135.1956},
	"shiga":    {"滋賀", 35.0044, 135.8686},
}

// GetSupportedCities returns a list of all supported city names
func GetSupportedCities() []string {
	cities := make([]string, 0, len(Cities))
	for key := range Cities {
		cities = append(cities, key)
	}
	sort.Strings(cities)
	return cities
}

// GetCityCoordinate returns city coordinates by city name
func GetCityCoordinate(city string) (*types.CityCoordinate, error) {
	if coord, exists := Cities[city]; exists {
		return &coord, nil
	}
	
	supportedCities := GetSupportedCities()
	return nil, fmt.Errorf("都市が見つかりません: %s\n対応都市: %v", city, supportedCities)
}

// GetWeather fetches weather data from API
func GetWeather(lat, lon float64) (*types.WeatherData, error) {
	forecastDays := 1
	
	var url string
	currentParams := "temperature_2m,apparent_temperature,relative_humidity_2m,wind_speed_10m,wind_direction_10m,weather_code,precipitation,dewpoint_2m"
	dailyParams := "temperature_2m_max,temperature_2m_min,weather_code,wind_speed_10m_max,precipitation_sum"
	hourlyParams := "temperature_2m,apparent_temperature,relative_humidity_2m,wind_speed_10m,wind_direction_10m,weather_code,precipitation"
	

	// 予報データ
	url = fmt.Sprintf("%s?latitude=%s&longitude=%s&current=%s&daily=%s&hourly=%s&timezone=Asia/Tokyo&forecast_days=%d", 
		apiURL, 
		strconv.FormatFloat(lat, 'f', 4, 64), 
		strconv.FormatFloat(lon, 'f', 4, 64),
		currentParams,
		dailyParams,
		hourlyParams,
		forecastDays)
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}
	
	var weather types.WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return &weather, nil
}

// GetWeatherDescription returns Japanese weather description
func GetWeatherDescription(code int) string {
	weatherCodes := map[int]string{
		0:  "快晴",
		1:  "晴れ",
		2:  "一部曇り",
		3:  "曇り",
		45: "霧",
		48: "着氷霧",
		51: "弱い霧雨",
		53: "霧雨",
		55: "強い霧雨",
		56: "軽い着氷霧雨",
		57: "着氷霧雨",
		61: "弱い雨",
		63: "雨",
		65: "強い雨",
		66: "軽い着氷雨",
		67: "着氷雨",
		71: "弱い雪",
		73: "雪",
		75: "強い雪",
		77: "雪つぶ",
		80: "弱いにわか雨",
		81: "にわか雨",
		82: "強いにわか雨",
		85: "弱いにわか雪",
		86: "にわか雪",
		95: "雷雨",
		96: "雹を伴う雷雨",
		99: "強い雹を伴う雷雨",
	}
	
	if desc, exists := weatherCodes[code]; exists {
		return desc
	}
	return "不明"
}

// FormatDate formats date string to Japanese format
func FormatDate(dateStr string) string {
	if len(dateStr) < 10 {
		return dateStr
	}
	
	datePart := dateStr[:10]
	t, err := time.Parse("2006-01-02", datePart)
	if err != nil {
		return dateStr
	}
	
	return t.Format("01月02日")
}

// GetWindDirection converts wind direction to Japanese
func GetWindDirection(direction float64) string {
	directions := []string{
		"北", "北北東", "北東", "東北東", "東", "東南東", "南東", "南南東",
		"南", "南南西", "南西", "西南西", "西", "西北西", "北西", "北北西",
	}
	
	index := int((direction + 11.25) / 22.5) % 16
	return directions[index]
}