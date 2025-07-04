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

type TimeBasedWeather struct {
	Time        string
	Temperature float64
	ApparentTemp float64
	Humidity    int
	WindSpeed   float64
	WindDirection float64
	WeatherCode int
	Precipitation float64
}

type TimePeriod struct {
	Name        string
	DisplayName string
	StartHour   int
	EndHour     int
}

func main() {
	var city string
	var days int
	var runningMode bool
	var timeOfDay string
	var dateSpec string
	
	flag.StringVar(&city, "city", "Tokyo", "都市名を指定")
	flag.IntVar(&days, "days", 0, "予報日数を指定（1-7日、0は現在の天気のみ）")
	flag.BoolVar(&runningMode, "running", false, "ランニング向け情報を表示")
	flag.StringVar(&timeOfDay, "time", "", "時間帯を指定（morning=早朝, noon=昼, evening=夕方, night=夜）")
	flag.StringVar(&dateSpec, "date", "", "日付を指定（today=今日, tomorrow=明日, day-after-tomorrow=明後日）")
	flag.Parse()

	if city == "" {
		fmt.Println("都市名を指定してください: -city <都市名>")
		os.Exit(1)
	}
	
	if days < 0 || days > 7 {
		fmt.Println("予報日数は0-7日の範囲で指定してください")
		os.Exit(1)
	}
	
	// Validate date specification
	var dayOffset int = 0
	if dateSpec != "" {
		validDates := []string{"today", "tomorrow", "day-after-tomorrow"}
		valid := false
		for _, validDate := range validDates {
			if dateSpec == validDate {
				valid = true
				break
			}
		}
		if !valid {
			fmt.Println("日付は today, tomorrow, day-after-tomorrow のいずれかを指定してください")
			os.Exit(1)
		}
		
		// Calculate day offset
		switch dateSpec {
		case "today":
			dayOffset = 0
		case "tomorrow":
			dayOffset = 1
		case "day-after-tomorrow":
			dayOffset = 2
		}
		
		// Date-based queries require forecast data
		if days == 0 {
			days = max(3, dayOffset + 1) // Ensure we have enough forecast days
		}
	}
	
	// Validate time of day
	if timeOfDay != "" {
		validTimes := []string{"morning", "noon", "evening", "night"}
		valid := false
		for _, validTime := range validTimes {
			if timeOfDay == validTime {
				valid = true
				break
			}
		}
		if !valid {
			fmt.Println("時間帯は morning, noon, evening, night のいずれかを指定してください")
			os.Exit(1)
		}
		
		// Time-based queries require forecast data
		if days == 0 {
			days = 1
		}
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

	if dateSpec != "" {
		// Date-specific weather
		if timeOfDay != "" {
			// Date + time specific weather
			if runningMode {
				displayDateTimeBasedRunningWeather(weather, coord.Name, dateSpec, timeOfDay, dayOffset)
			} else {
				displayDateTimeBasedWeather(weather, coord.Name, dateSpec, timeOfDay, dayOffset)
			}
		} else {
			// Date specific weather (full day)
			if runningMode {
				displayDateBasedRunningWeather(weather, coord.Name, dateSpec, dayOffset)
			} else {
				displayDateBasedWeather(weather, coord.Name, dateSpec, dayOffset)
			}
		}
	} else if timeOfDay != "" {
		// Time-specific weather
		if runningMode {
			displayTimeBasedRunningWeather(weather, coord.Name, timeOfDay, days)
		} else {
			displayTimeBasedWeather(weather, coord.Name, timeOfDay, days)
		}
	} else if runningMode {
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

func getTimePeriods() map[string]TimePeriod {
	return map[string]TimePeriod{
		"morning": {"morning", "早朝", 5, 9},
		"noon":    {"noon", "昼", 11, 15},
		"evening": {"evening", "夕方", 17, 19},
		"night":   {"night", "夜", 21, 23},
	}
}

func extractTimeBasedWeather(weather *WeatherData, timeOfDay string, days int) []TimeBasedWeather {
	periods := getTimePeriods()
	period, exists := periods[timeOfDay]
	if !exists {
		return nil
	}

	var results []TimeBasedWeather
	
	for day := 0; day < days && day < len(weather.Hourly.Time)/24; day++ {
		for hour := period.StartHour; hour <= period.EndHour; hour++ {
			index := day*24 + hour
			if index < len(weather.Hourly.Time) {
				results = append(results, TimeBasedWeather{
					Time:          weather.Hourly.Time[index],
					Temperature:   weather.Hourly.Temperature[index],
					ApparentTemp:  weather.Hourly.ApparentTemp[index],
					Humidity:      weather.Hourly.Humidity[index],
					WindSpeed:     weather.Hourly.WindSpeed[index],
					WindDirection: weather.Hourly.WindDirection[index],
					WeatherCode:   weather.Hourly.WeatherCode[index],
					Precipitation: weather.Hourly.Precipitation[index],
				})
			}
		}
	}
	
	return results
}

func displayTimeBasedWeather(weather *WeatherData, cityName, timeOfDay string, days int) {
	periods := getTimePeriods()
	period := periods[timeOfDay]
	
	timeData := extractTimeBasedWeather(weather, timeOfDay, days)
	if len(timeData) == 0 {
		fmt.Println("指定された時間帯のデータが見つかりません")
		return
	}
	
	fmt.Printf("🌤️  %s の%s時間帯天気情報\n", cityName, period.DisplayName)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	for i, data := range timeData {
		hour := extractHour(data.Time)
		temp := data.Temperature
		weather := getWeatherDescription(data.WeatherCode)
		
		fmt.Printf("📅 %s時: %.1f°C | %s", hour, temp, weather)
		if data.Precipitation > 0 {
			fmt.Printf(" | 🌧️ %.1fmm", data.Precipitation)
		}
		fmt.Printf("\n")
		
		if (i+1)%3 == 0 && i < len(timeData)-1 {
			fmt.Printf("   ────────────────────────────\n")
		}
	}
	
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

func displayTimeBasedRunningWeather(weather *WeatherData, cityName, timeOfDay string, days int) {
	periods := getTimePeriods()
	period := periods[timeOfDay]
	
	timeData := extractTimeBasedWeather(weather, timeOfDay, days)
	if len(timeData) == 0 {
		fmt.Println("指定された時間帯のデータが見つかりません")
		return
	}
	
	fmt.Printf("🏃‍♂️ %s の%s時間帯ランニング情報\n", cityName, period.DisplayName)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	bestCondition := TimeBasedWeather{}
	bestScore := -1
	bestTime := ""
	
	fmt.Printf("⏰ %s時間帯詳細 (%d:00-%d:00)\n", period.DisplayName, period.StartHour, period.EndHour)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	for _, data := range timeData {
		condition := assessRunningCondition(
			data.Temperature,
			data.ApparentTemp,
			float64(data.Humidity),
			data.WindSpeed,
			data.Precipitation,
			data.WeatherCode,
		)
		
		hour := extractHour(data.Time)
		fmt.Printf("🕐 %s時: %d/100 (%s)\n", hour, condition.Score, condition.Level)
		fmt.Printf("   🌡️ %.1f°C (体感: %.1f°C) | 💧 %d%%\n", 
			data.Temperature, data.ApparentTemp, data.Humidity)
		fmt.Printf("   ☁️ %s", getWeatherDescription(data.WeatherCode))
		if data.Precipitation > 0 {
			fmt.Printf(" | 🌧️ %.1fmm", data.Precipitation)
		}
		fmt.Printf("\n")
		
		if condition.Score > bestScore {
			bestScore = condition.Score
			bestCondition = data
			bestTime = hour
		}
		
		fmt.Printf("   ────────────────────────────\n")
	}
	
	// Best time recommendation
	if bestScore >= 0 {
		bestRunningCondition := assessRunningCondition(
			bestCondition.Temperature,
			bestCondition.ApparentTemp,
			float64(bestCondition.Humidity),
			bestCondition.WindSpeed,
			bestCondition.Precipitation,
			bestCondition.WeatherCode,
		)
		
		fmt.Printf("🏆 最適時間: %s時 (スコア: %d/100)\n", bestTime, bestScore)
		fmt.Printf("💡 %s\n", bestRunningCondition.Recommendation)
		
		if len(bestRunningCondition.Warnings) > 0 {
			fmt.Printf("⚠️  注意事項:\n")
			for _, warning := range bestRunningCondition.Warnings {
				fmt.Printf("   %s\n", warning)
			}
		}
	}
	
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

func extractHour(timeStr string) string {
	// Extract hour from ISO time string (YYYY-MM-DDTHH:MM)
	if len(timeStr) >= 13 {
		return timeStr[11:13]
	}
	return timeStr
}

func extractDateBasedWeather(weather *WeatherData, dayOffset int) *WeatherData {
	if dayOffset == 0 || len(weather.Daily.Time) == 0 {
		return weather
	}
	
	// Extract specific day data from forecast
	dateSpecificWeather := &WeatherData{
		Current: weather.Current, // Keep current for reference
		Daily: struct {
			Time           []string  `json:"time"`
			TemperatureMax []float64 `json:"temperature_2m_max"`
			TemperatureMin []float64 `json:"temperature_2m_min"`
			WeatherCode    []int     `json:"weather_code"`
			WindSpeedMax   []float64 `json:"wind_speed_10m_max"`
			PrecipitationSum []float64 `json:"precipitation_sum"`
		}{},
		Hourly: struct {
			Time           []string  `json:"time"`
			Temperature    []float64 `json:"temperature_2m"`
			ApparentTemp   []float64 `json:"apparent_temperature"`
			Humidity       []int     `json:"relative_humidity_2m"`
			WindSpeed      []float64 `json:"wind_speed_10m"`
			WindDirection  []float64 `json:"wind_direction_10m"`
			WeatherCode    []int     `json:"weather_code"`
			Precipitation  []float64 `json:"precipitation"`
		}{},
	}
	
	// Extract specific day from daily data
	if dayOffset < len(weather.Daily.Time) {
		dateSpecificWeather.Daily.Time = []string{weather.Daily.Time[dayOffset]}
		dateSpecificWeather.Daily.TemperatureMax = []float64{weather.Daily.TemperatureMax[dayOffset]}
		dateSpecificWeather.Daily.TemperatureMin = []float64{weather.Daily.TemperatureMin[dayOffset]}
		dateSpecificWeather.Daily.WeatherCode = []int{weather.Daily.WeatherCode[dayOffset]}
		dateSpecificWeather.Daily.WindSpeedMax = []float64{weather.Daily.WindSpeedMax[dayOffset]}
		dateSpecificWeather.Daily.PrecipitationSum = []float64{weather.Daily.PrecipitationSum[dayOffset]}
	}
	
	// Extract 24 hours of hourly data for the specific day
	startIndex := dayOffset * 24
	endIndex := startIndex + 24
	if endIndex <= len(weather.Hourly.Time) {
		dateSpecificWeather.Hourly.Time = weather.Hourly.Time[startIndex:endIndex]
		dateSpecificWeather.Hourly.Temperature = weather.Hourly.Temperature[startIndex:endIndex]
		dateSpecificWeather.Hourly.ApparentTemp = weather.Hourly.ApparentTemp[startIndex:endIndex]
		dateSpecificWeather.Hourly.Humidity = weather.Hourly.Humidity[startIndex:endIndex]
		dateSpecificWeather.Hourly.WindSpeed = weather.Hourly.WindSpeed[startIndex:endIndex]
		dateSpecificWeather.Hourly.WindDirection = weather.Hourly.WindDirection[startIndex:endIndex]
		dateSpecificWeather.Hourly.WeatherCode = weather.Hourly.WeatherCode[startIndex:endIndex]
		dateSpecificWeather.Hourly.Precipitation = weather.Hourly.Precipitation[startIndex:endIndex]
	}
	
	return dateSpecificWeather
}

func displayDateBasedWeather(weather *WeatherData, cityName, dateSpec string, dayOffset int) {
	dateSpecificWeather := extractDateBasedWeather(weather, dayOffset)
	
	dateDisplayName := getDateDisplayName(dateSpec)
	fmt.Printf("🌤️  %s の%s天気予報\n", cityName, dateDisplayName)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	if len(dateSpecificWeather.Daily.Time) == 0 {
		fmt.Println("指定された日付のデータが見つかりません")
		return
	}
	
	// Display daily summary
	date := dateSpecificWeather.Daily.Time[0]
	maxTemp := dateSpecificWeather.Daily.TemperatureMax[0]
	minTemp := dateSpecificWeather.Daily.TemperatureMin[0]
	weatherCode := dateSpecificWeather.Daily.WeatherCode[0]
	maxWind := dateSpecificWeather.Daily.WindSpeedMax[0]
	precipitation := dateSpecificWeather.Daily.PrecipitationSum[0]
	
	fmt.Printf("📅 %s (%s)\n", formatDate(date), dateDisplayName)
	fmt.Printf("🌡️  最高: %.1f°C / 最低: %.1f°C\n", maxTemp, minTemp)
	fmt.Printf("🌬️  最大風速: %.1f m/s\n", maxWind)
	fmt.Printf("☁️  天気: %s\n", getWeatherDescription(weatherCode))
	if precipitation > 0 {
		fmt.Printf("🌧️  降水量: %.1f mm\n", precipitation)
	}
	
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

func displayDateBasedRunningWeather(weather *WeatherData, cityName, dateSpec string, dayOffset int) {
	dateSpecificWeather := extractDateBasedWeather(weather, dayOffset)
	
	dateDisplayName := getDateDisplayName(dateSpec)
	fmt.Printf("🏃‍♂️ %s の%sランニング情報\n", cityName, dateDisplayName)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	if len(dateSpecificWeather.Daily.Time) == 0 {
		fmt.Println("指定された日付のデータが見つかりません")
		return
	}
	
	// Daily summary
	date := dateSpecificWeather.Daily.Time[0]
	maxTemp := dateSpecificWeather.Daily.TemperatureMax[0]
	minTemp := dateSpecificWeather.Daily.TemperatureMin[0]
	weatherCode := dateSpecificWeather.Daily.WeatherCode[0]
	maxWind := dateSpecificWeather.Daily.WindSpeedMax[0]
	precipitation := dateSpecificWeather.Daily.PrecipitationSum[0]
	
	// Estimate daily running condition (using average temperature)
	avgTemp := (maxTemp + minTemp) / 2
	dailyCondition := assessRunningCondition(avgTemp, avgTemp, 60, maxWind, precipitation, weatherCode)
	
	fmt.Printf("📅 %s (%s)\n", formatDate(date), dateDisplayName)
	fmt.Printf("🏆 ランニング指数: %d/100 (%s)\n", dailyCondition.Score, dailyCondition.Level)
	fmt.Printf("💡 %s\n", dailyCondition.Recommendation)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	fmt.Printf("🌡️  %s%.1f°C〜%.1f°C\n", getRunningTempIcon(avgTemp), minTemp, maxTemp)
	fmt.Printf("☁️  %s\n", getWeatherDescription(weatherCode))
	if precipitation > 0 {
		fmt.Printf("🌧️  降水量: %.1f mm\n", precipitation)
	}
	
	// Clothing recommendations
	if len(dailyCondition.Clothing) > 0 {
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("👕 推奨ウェア:\n")
		for _, item := range dailyCondition.Clothing {
			fmt.Printf("   • %s\n", item)
		}
	}
	
	// Warnings
	if len(dailyCondition.Warnings) > 0 {
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("⚠️  注意事項:\n")
		for _, warning := range dailyCondition.Warnings {
			fmt.Printf("   %s\n", warning)
		}
	}
	
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

func displayDateTimeBasedWeather(weather *WeatherData, cityName, dateSpec, timeOfDay string, dayOffset int) {
	dateSpecificWeather := extractDateBasedWeather(weather, dayOffset)
	
	periods := getTimePeriods()
	period := periods[timeOfDay]
	dateDisplayName := getDateDisplayName(dateSpec)
	
	timeData := extractTimeBasedWeather(dateSpecificWeather, timeOfDay, 1)
	if len(timeData) == 0 {
		fmt.Println("指定された日付・時間帯のデータが見つかりません")
		return
	}
	
	fmt.Printf("🌤️  %s の%s%s時間帯天気情報\n", cityName, dateDisplayName, period.DisplayName)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	for i, data := range timeData {
		hour := extractHour(data.Time)
		temp := data.Temperature
		weather := getWeatherDescription(data.WeatherCode)
		
		fmt.Printf("📅 %s時: %.1f°C | %s", hour, temp, weather)
		if data.Precipitation > 0 {
			fmt.Printf(" | 🌧️ %.1fmm", data.Precipitation)
		}
		fmt.Printf("\n")
		
		if (i+1)%3 == 0 && i < len(timeData)-1 {
			fmt.Printf("   ────────────────────────────\n")
		}
	}
	
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

func displayDateTimeBasedRunningWeather(weather *WeatherData, cityName, dateSpec, timeOfDay string, dayOffset int) {
	dateSpecificWeather := extractDateBasedWeather(weather, dayOffset)
	
	periods := getTimePeriods()
	period := periods[timeOfDay]
	dateDisplayName := getDateDisplayName(dateSpec)
	
	timeData := extractTimeBasedWeather(dateSpecificWeather, timeOfDay, 1)
	if len(timeData) == 0 {
		fmt.Println("指定された日付・時間帯のデータが見つかりません")
		return
	}
	
	fmt.Printf("🏃‍♂️ %s の%s%s時間帯ランニング情報\n", cityName, dateDisplayName, period.DisplayName)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	bestCondition := TimeBasedWeather{}
	bestScore := -1
	bestTime := ""
	
	fmt.Printf("⏰ %s%s時間帯詳細 (%d:00-%d:00)\n", dateDisplayName, period.DisplayName, period.StartHour, period.EndHour)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	for _, data := range timeData {
		condition := assessRunningCondition(
			data.Temperature,
			data.ApparentTemp,
			float64(data.Humidity),
			data.WindSpeed,
			data.Precipitation,
			data.WeatherCode,
		)
		
		hour := extractHour(data.Time)
		fmt.Printf("🕐 %s時: %d/100 (%s)\n", hour, condition.Score, condition.Level)
		fmt.Printf("   🌡️ %.1f°C (体感: %.1f°C) | 💧 %d%%\n", 
			data.Temperature, data.ApparentTemp, data.Humidity)
		fmt.Printf("   ☁️ %s", getWeatherDescription(data.WeatherCode))
		if data.Precipitation > 0 {
			fmt.Printf(" | 🌧️ %.1fmm", data.Precipitation)
		}
		fmt.Printf("\n")
		
		if condition.Score > bestScore {
			bestScore = condition.Score
			bestCondition = data
			bestTime = hour
		}
		
		fmt.Printf("   ────────────────────────────\n")
	}
	
	// Best time recommendation
	if bestScore >= 0 {
		bestRunningCondition := assessRunningCondition(
			bestCondition.Temperature,
			bestCondition.ApparentTemp,
			float64(bestCondition.Humidity),
			bestCondition.WindSpeed,
			bestCondition.Precipitation,
			bestCondition.WeatherCode,
		)
		
		fmt.Printf("🏆 最適時間: %s時 (スコア: %d/100)\n", bestTime, bestScore)
		fmt.Printf("💡 %s\n", bestRunningCondition.Recommendation)
		
		if len(bestRunningCondition.Warnings) > 0 {
			fmt.Printf("⚠️  注意事項:\n")
			for _, warning := range bestRunningCondition.Warnings {
				fmt.Printf("   %s\n", warning)
			}
		}
	}
	
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
}

func getDateDisplayName(dateSpec string) string {
	switch dateSpec {
	case "today":
		return "今日の"
	case "tomorrow":
		return "明日の"
	case "day-after-tomorrow":
		return "明後日の"
	default:
		return ""
	}
}