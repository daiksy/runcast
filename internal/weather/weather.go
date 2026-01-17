package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"
	"runcast/internal/config"
	"runcast/internal/types"
)

const apiURL = "https://api.open-meteo.com/v1/jma"
const airQualityAPIURL = "https://air-quality-api.open-meteo.com/v1/air-quality"

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
	// Check built-in cities first
	if coord, exists := Cities[city]; exists {
		return &coord, nil
	}
	
	// Check custom locations from config
	cfg, err := config.LoadConfig()
	if err != nil {
		// If config loading fails, continue with built-in cities only
		fmt.Printf("警告: 設定ファイルの読み込みに失敗しました: %v\n", err)
	} else {
		if coord, exists := cfg.GetCustomLocation(city); exists {
			return coord, nil
		}
	}
	
	// Generate error message with all available locations
	supportedCities := GetSupportedCities()
	allLocations := make([]string, len(supportedCities))
	copy(allLocations, supportedCities)
	
	if cfg != nil {
		customLocations := cfg.GetCustomLocationNames()
		if len(customLocations) > 0 {
			sort.Strings(customLocations)
			allLocations = append(allLocations, customLocations...)
			sort.Strings(allLocations)
		}
	}
	
	return nil, fmt.Errorf("都市が見つかりません: %s\n対応都市: %v", city, allLocations)
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

// GetAirQuality fetches air quality data from API
func GetAirQuality(lat, lon float64) (*types.AirQualityData, error) {
	url := fmt.Sprintf("%s?latitude=%s&longitude=%s&hourly=dust,pm10,pm2_5&timezone=Asia/Tokyo&forecast_days=1",
		airQualityAPIURL,
		strconv.FormatFloat(lat, 'f', 4, 64),
		strconv.FormatFloat(lon, 'f', 4, 64))

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Air Quality API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Air Quality API request failed with status: %d", resp.StatusCode)
	}

	var airQuality types.AirQualityData
	if err := json.NewDecoder(resp.Body).Decode(&airQuality); err != nil {
		return nil, fmt.Errorf("failed to decode air quality response: %w", err)
	}

	return &airQuality, nil
}

// GetCurrentDustLevel returns current dust level based on air quality data
func GetCurrentDustLevel(airQuality *types.AirQualityData) *types.DustLevel {
	if airQuality == nil || len(airQuality.Hourly.Time) == 0 {
		return nil
	}

	// Find current hour data
	now := time.Now()
	currentHour := now.Format("2006-01-02T15:00")

	for i, t := range airQuality.Hourly.Time {
		if t == currentHour {
			dust := 0.0
			pm10 := 0.0
			pm2_5 := 0.0

			if i < len(airQuality.Hourly.Dust) {
				dust = airQuality.Hourly.Dust[i]
			}
			if i < len(airQuality.Hourly.PM10) {
				pm10 = airQuality.Hourly.PM10[i]
			}
			if i < len(airQuality.Hourly.PM2_5) {
				pm2_5 = airQuality.Hourly.PM2_5[i]
			}

			return createDustLevel(dust, pm10, pm2_5)
		}
	}

	// If current hour not found, use first available data
	dust := 0.0
	pm10 := 0.0
	pm2_5 := 0.0

	if len(airQuality.Hourly.Dust) > 0 {
		dust = airQuality.Hourly.Dust[0]
	}
	if len(airQuality.Hourly.PM10) > 0 {
		pm10 = airQuality.Hourly.PM10[0]
	}
	if len(airQuality.Hourly.PM2_5) > 0 {
		pm2_5 = airQuality.Hourly.PM2_5[0]
	}

	return createDustLevel(dust, pm10, pm2_5)
}

// GetHourlyDustLevel returns dust level for a specific hour
func GetHourlyDustLevel(airQuality *types.AirQualityData, hour int, days int) *types.DustLevel {
	if airQuality == nil || len(airQuality.Hourly.Time) == 0 {
		return nil
	}

	targetDate := time.Now().AddDate(0, 0, days)
	targetTime := fmt.Sprintf("%sT%02d:00", targetDate.Format("2006-01-02"), hour)

	for i, t := range airQuality.Hourly.Time {
		if t == targetTime {
			dust := 0.0
			pm10 := 0.0
			pm2_5 := 0.0

			if i < len(airQuality.Hourly.Dust) {
				dust = airQuality.Hourly.Dust[i]
			}
			if i < len(airQuality.Hourly.PM10) {
				pm10 = airQuality.Hourly.PM10[i]
			}
			if i < len(airQuality.Hourly.PM2_5) {
				pm2_5 = airQuality.Hourly.PM2_5[i]
			}

			return createDustLevel(dust, pm10, pm2_5)
		}
	}

	return nil
}

// createDustLevel creates DustLevel from raw values
func createDustLevel(dust, pm10, pm2_5 float64) *types.DustLevel {
	level := 0
	displayName := "なし"
	description := "黄砂の影響なし"

	if dust > 500 {
		level = 4
		displayName = "非常に多い"
		description = "屋外活動は控えるべき"
	} else if dust > 200 {
		level = 3
		displayName = "多い"
		description = "外出時に注意が必要"
	} else if dust > 100 {
		level = 2
		displayName = "やや多い"
		description = "視程に影響の可能性"
	} else if dust > 50 {
		level = 1
		displayName = "少ない"
		description = "わずかに飛来"
	}

	return &types.DustLevel{
		Level:       level,
		DisplayName: displayName,
		Description: description,
		Dust:        dust,
		PM10:        pm10,
		PM2_5:       pm2_5,
	}
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