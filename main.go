package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	apiURL = "https://api.open-meteo.com/v1/jma"
)

type WeatherData struct {
	Current struct {
		Temperature        float64 `json:"temperature_2m"`
		ApparentTemp       float64 `json:"apparent_temperature"`
		Humidity           int     `json:"relative_humidity_2m"`
		WindSpeed          float64 `json:"wind_speed_10m"`
		WindDirection      float64 `json:"wind_direction_10m"`
		WeatherCode        int     `json:"weather_code"`
		Precipitation      float64 `json:"precipitation"`
		Dewpoint           float64 `json:"dewpoint_2m"`
	} `json:"current"`
	Daily struct {
		Time           []string  `json:"time"`
		TemperatureMax []float64 `json:"temperature_2m_max"`
		TemperatureMin []float64 `json:"temperature_2m_min"`
		WeatherCode    []int     `json:"weather_code"`
		WindSpeedMax   []float64 `json:"wind_speed_10m_max"`
		PrecipitationSum []float64 `json:"precipitation_sum"`
	} `json:"daily"`
	Hourly struct {
		Time           []string  `json:"time"`
		Temperature    []float64 `json:"temperature_2m"`
		ApparentTemp   []float64 `json:"apparent_temperature"`
		Humidity       []int     `json:"relative_humidity_2m"`
		WindSpeed      []float64 `json:"wind_speed_10m"`
		WindDirection  []float64 `json:"wind_direction_10m"`
		WeatherCode    []int     `json:"weather_code"`
		Precipitation  []float64 `json:"precipitation"`
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
	var runningMode bool
	
	flag.StringVar(&city, "city", "Tokyo", "都市名を指定")
	flag.IntVar(&days, "days", 0, "予報日数を指定（1-7日、0は現在の天気のみ）")
	flag.BoolVar(&runningMode, "running", false, "ランニング向け情報を表示")
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

	if runningMode {
		if days == 0 {
			displayRunningWeather(weather, coord.Name)
		} else {
			displayRunningForecast(weather, coord.Name, days)
		}
	} else {
		if days == 0 {
			displayCurrentWeather(weather, coord.Name)
		} else {
			displayForecastWeather(weather, coord.Name, days)
		}
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
	currentParams := "temperature_2m,apparent_temperature,relative_humidity_2m,wind_speed_10m,wind_direction_10m,weather_code,precipitation,dewpoint_2m"
	dailyParams := "temperature_2m_max,temperature_2m_min,weather_code,wind_speed_10m_max,precipitation_sum"
	hourlyParams := "temperature_2m,apparent_temperature,relative_humidity_2m,wind_speed_10m,wind_direction_10m,weather_code,precipitation"
	
	if days == 0 {
		// 現在の天気のみ
		url = fmt.Sprintf("%s?latitude=%s&longitude=%s&current=%s&timezone=Asia/Tokyo", 
			apiURL, 
			strconv.FormatFloat(lat, 'f', 4, 64), 
			strconv.FormatFloat(lon, 'f', 4, 64),
			currentParams)
	} else {
		// 予報データ
		url = fmt.Sprintf("%s?latitude=%s&longitude=%s&current=%s&daily=%s&hourly=%s&timezone=Asia/Tokyo&forecast_days=%d", 
			apiURL, 
			strconv.FormatFloat(lat, 'f', 4, 64), 
			strconv.FormatFloat(lon, 'f', 4, 64),
			currentParams,
			dailyParams,
			hourlyParams,
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
	if len(dateStr) >= 10 && dateStr[4] == '-' && dateStr[7] == '-' {
		month := dateStr[5:7]
		day := dateStr[8:10]
		return fmt.Sprintf("%s月%s日", month, day)
	}
	return dateStr
}

// Running condition assessment
type RunningCondition struct {
	Score       int    // 0-100 score
	Level       string // Excellent, Good, Fair, Poor, Dangerous
	Recommendation string
	Clothing    []string
	Warnings    []string
}

func assessRunningCondition(temp, apparentTemp, humidity float64, windSpeed, precipitation float64, weatherCode int) RunningCondition {
	condition := RunningCondition{
		Clothing: []string{},
		Warnings: []string{},
	}
	
	score := 100
	
	// Temperature assessment
	if apparentTemp >= 30 {
		score -= 30
		condition.Warnings = append(condition.Warnings, "⚠️  熱中症注意: 体感温度が高すぎます")
		condition.Clothing = append(condition.Clothing, "薄手の半袖", "帽子必須", "サングラス")
	} else if apparentTemp >= 25 {
		score -= 10
		condition.Clothing = append(condition.Clothing, "半袖", "帽子推奨")
	} else if apparentTemp >= 15 {
		// Ideal temperature
		condition.Clothing = append(condition.Clothing, "半袖または薄手の長袖")
	} else if apparentTemp >= 5 {
		score -= 10
		condition.Clothing = append(condition.Clothing, "長袖", "薄手のジャケット")
	} else {
		score -= 20
		condition.Clothing = append(condition.Clothing, "防寒着", "手袋", "ニット帽")
		condition.Warnings = append(condition.Warnings, "❄️  寒さ注意: 十分な防寒対策を")
	}
	
	// Humidity assessment
	if humidity >= 80 {
		score -= 15
		condition.Warnings = append(condition.Warnings, "💧 高湿度: 汗が乾きにくい状態です")
	} else if humidity <= 30 {
		score -= 5
		condition.Warnings = append(condition.Warnings, "🏜️  低湿度: 水分補給をこまめに")
	}
	
	// Wind assessment
	if windSpeed >= 10 {
		score -= 15
		condition.Warnings = append(condition.Warnings, "💨 強風注意")
	} else if windSpeed >= 5 {
		score -= 5
	}
	
	// Precipitation assessment
	if precipitation > 0 {
		score -= 25
		condition.Warnings = append(condition.Warnings, "🌧️  降水中: 滑りやすい路面に注意")
	}
	
	// Weather code assessment
	if weatherCode >= 95 { // Thunderstorm
		score -= 50
		condition.Warnings = append(condition.Warnings, "⛈️  雷雨: ランニング中止を強く推奨")
	} else if weatherCode >= 80 { // Heavy rain
		score -= 30
	} else if weatherCode >= 61 { // Rain
		score -= 20
	}
	
	condition.Score = max(0, score)
	
	// Determine level and recommendation based on score AND warnings
	hasWarnings := len(condition.Warnings) > 0
	hasSevereWarnings := false
	
	// Check for severe warnings
	for _, warning := range condition.Warnings {
		if strings.Contains(warning, "熱中症注意") || 
		   strings.Contains(warning, "雷雨") || 
		   strings.Contains(warning, "強風注意") {
			hasSevereWarnings = true
			break
		}
	}
	
	if hasSevereWarnings {
		// Override recommendation if there are severe warnings
		if condition.Score >= 60 {
			condition.Level = "注意"
			condition.Recommendation = "警告事項があります。十分注意してランニングしてください"
		} else if condition.Score >= 40 {
			condition.Level = "注意"
			condition.Recommendation = "警告事項があります。ランニングは控えめに"
		} else {
			condition.Level = "危険"
			condition.Recommendation = "危険な状況です。ランニング中止を強く推奨します"
		}
	} else if hasWarnings {
		// Moderate warnings present
		if condition.Score >= 80 {
			condition.Level = "良好"
			condition.Recommendation = "注意事項がありますが、ランニング可能です"
		} else if condition.Score >= 60 {
			condition.Level = "良好"
			condition.Recommendation = "注意事項を確認してからランニングしてください"
		} else if condition.Score >= 40 {
			condition.Level = "普通"
			condition.Recommendation = "注意しながらランニング可能です"
		} else if condition.Score >= 20 {
			condition.Level = "注意"
			condition.Recommendation = "ランニングは控えめに、安全第一で"
		} else {
			condition.Level = "危険"
			condition.Recommendation = "ランニング中止を推奨します"
		}
	} else {
		// No warnings - original scoring system
		if condition.Score >= 80 {
			condition.Level = "最高"
			condition.Recommendation = "ランニングに最適な天候です！"
		} else if condition.Score >= 60 {
			condition.Level = "良好"
			condition.Recommendation = "ランニングに適した天候です"
		} else if condition.Score >= 40 {
			condition.Level = "普通"
			condition.Recommendation = "注意しながらランニング可能です"
		} else if condition.Score >= 20 {
			condition.Level = "注意"
			condition.Recommendation = "ランニングは控えめに、安全第一で"
		} else {
			condition.Level = "危険"
			condition.Recommendation = "ランニング中止を推奨します"
		}
	}
	
	return condition
}

func getWindDirection(degrees float64) string {
	directions := []string{"北", "北北東", "北東", "東北東", "東", "東南東", "南東", "南南東", "南", "南南西", "南西", "西南西", "西", "西北西", "北西", "北北西"}
	index := int((degrees + 11.25) / 22.5) % 16
	return directions[index]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func displayRunningWeather(weather *WeatherData, cityName string) {
	fmt.Printf("🏃‍♂️ %s のランニング情報\n", cityName)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	condition := assessRunningCondition(
		weather.Current.Temperature,
		weather.Current.ApparentTemp,
		float64(weather.Current.Humidity),
		weather.Current.WindSpeed,
		weather.Current.Precipitation,
		weather.Current.WeatherCode,
	)
	
	// Running condition display
	fmt.Printf("🏆 ランニング指数: %d/100 (%s)\n", condition.Score, condition.Level)
	fmt.Printf("💡 %s\n", condition.Recommendation)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	// Detailed weather info
	fmt.Printf("🌡️  気温: %.1f°C (体感: %.1f°C)\n", weather.Current.Temperature, weather.Current.ApparentTemp)
	fmt.Printf("💧 湿度: %d%%\n", weather.Current.Humidity)
	fmt.Printf("🌬️  風: %s %.1f m/s\n", getWindDirection(weather.Current.WindDirection), weather.Current.WindSpeed)
	fmt.Printf("☁️  天気: %s\n", getWeatherDescription(weather.Current.WeatherCode))
	if weather.Current.Precipitation > 0 {
		fmt.Printf("🌧️  降水量: %.1f mm/h\n", weather.Current.Precipitation)
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
		fmt.Printf("⚠️  注意事項:\n")
		for _, warning := range condition.Warnings {
			fmt.Printf("   %s\n", warning)
		}
	}
	
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

func displayRunningForecast(weather *WeatherData, cityName string, days int) {
	fmt.Printf("🏃‍♂️ %s の%d日間ランニング予報\n", cityName, days)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	// Current condition
	currentCondition := assessRunningCondition(
		weather.Current.Temperature,
		weather.Current.ApparentTemp,
		float64(weather.Current.Humidity),
		weather.Current.WindSpeed,
		weather.Current.Precipitation,
		weather.Current.WeatherCode,
	)
	fmt.Printf("📅 現在: %.1f°C | %s | 指数: %d/100\n", weather.Current.Temperature, getWeatherDescription(weather.Current.WeatherCode), currentCondition.Score)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	// Daily forecast with running assessment
	for i := 0; i < len(weather.Daily.Time) && i < days; i++ {
		date := weather.Daily.Time[i]
		maxTemp := weather.Daily.TemperatureMax[i]
		minTemp := weather.Daily.TemperatureMin[i]
		weatherCode := weather.Daily.WeatherCode[i]
		maxWind := weather.Daily.WindSpeedMax[i]
		precipitation := weather.Daily.PrecipitationSum[i]
		
		// Estimate daily running condition (using average temperature)
		avgTemp := (maxTemp + minTemp) / 2
		dailyCondition := assessRunningCondition(avgTemp, avgTemp, 60, maxWind, precipitation, weatherCode)
		
		fmt.Printf("📅 %s\n", formatDate(date))
		fmt.Printf("   🌡️  %s%.1f°C〜%.1f°C\n", getRunningTempIcon(avgTemp), minTemp, maxTemp)
		fmt.Printf("   🏆 ランニング指数: %d/100 (%s)\n", dailyCondition.Score, dailyCondition.Level)
		fmt.Printf("   ☁️  %s\n", getWeatherDescription(weatherCode))
		if precipitation > 0 {
			fmt.Printf("   🌧️  降水量: %.1f mm\n", precipitation)
		}
		
		if i < len(weather.Daily.Time)-1 && i < days-1 {
			fmt.Printf("   ────────────────────────────\n")
		}
	}
	
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

func getRunningTempIcon(temp float64) string {
	if temp >= 30 {
		return "🔥 "
	} else if temp >= 25 {
		return "🌡️  "
	} else if temp >= 15 {
		return "👌 "
	} else if temp >= 5 {
		return "🧥 "
	} else {
		return "❄️  "
	}
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