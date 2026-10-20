package weather

import (
	"runcast/internal/types"
)

// GetTimePeriods returns all available time periods
func GetTimePeriods() map[string]types.TimePeriod {
	return map[string]types.TimePeriod{
		"morning": {"morning", "早朝", 5, 9},
		"noon":    {"noon", "昼", 11, 15},
		"evening": {"evening", "夕方", 17, 19},
		"night":   {"night", "夜", 21, 23},
	}
}

// ExtractTimeBasedWeather extracts weather data for specific time period
func ExtractTimeBasedWeather(weather *types.WeatherData, timeOfDay string, days int) []types.TimeBasedWeather {
	periods := GetTimePeriods()
	period, exists := periods[timeOfDay]
	if !exists {
		return nil
	}
	
	var timeData []types.TimeBasedWeather
	
	for i := 0; i < len(weather.Hourly.Time) && len(timeData) < days*5; i++ {
		hour := ExtractHour(weather.Hourly.Time[i])
		if hour == "" {
			continue
		}
		
		hourInt := 0
		if len(hour) >= 2 {
			if hour[0] == '0' || hour[0] == '1' || hour[0] == '2' {
				if hour[0] == '0' {
					hourInt = int(hour[1] - '0')
				} else if hour[0] == '1' {
					hourInt = 10 + int(hour[1] - '0')
				} else if hour[0] == '2' {
					hourInt = 20 + int(hour[1] - '0')
				}
			}
		}
		
		if hourInt >= period.StartHour && hourInt <= period.EndHour {
			timeData = append(timeData, types.TimeBasedWeather{
				Time:          weather.Hourly.Time[i],
				Temperature:   weather.Hourly.Temperature[i],
				ApparentTemp:  weather.Hourly.ApparentTemp[i],
				Humidity:      weather.Hourly.Humidity[i],
				WindSpeed:     weather.Hourly.WindSpeed[i],
				WindDirection: weather.Hourly.WindDirection[i],
				Precipitation: weather.Hourly.Precipitation[i],
				WeatherCode:   weather.Hourly.WeatherCode[i],
			})
		}
	}
	
	return timeData
}

// ExtractHour extracts hour from ISO time string
func ExtractHour(timeStr string) string {
	// Extract hour from ISO time string (YYYY-MM-DDTHH:MM)
	if len(timeStr) >= 13 {
		return timeStr[11:13]
	}
	return timeStr
}

// ExtractHourInt extracts hour as integer from ISO time string
func ExtractHourInt(timeStr string) int {
	hour := ExtractHour(timeStr)
	if len(hour) >= 2 {
		if hour[0] == '0' {
			return int(hour[1] - '0')
		} else if hour[0] == '1' {
			return 10 + int(hour[1]-'0')
		} else if hour[0] == '2' {
			return 20 + int(hour[1]-'0')
		}
	}
	return 0
}

// ExtractDateBasedWeather extracts weather data for specific date
func ExtractDateBasedWeather(weather *types.WeatherData, dayOffset int) *types.WeatherData {
	if dayOffset == 0 || len(weather.Daily.Time) == 0 {
		return weather
	}
	
	if dayOffset >= len(weather.Daily.Time) {
		return weather
	}
	
	// Check if all required arrays have sufficient length
	if len(weather.Daily.TemperatureMax) <= dayOffset ||
		len(weather.Daily.TemperatureMin) <= dayOffset ||
		len(weather.Daily.WindSpeedMax) <= dayOffset ||
		len(weather.Daily.PrecipitationSum) <= dayOffset ||
		len(weather.Daily.WeatherCode) <= dayOffset {
		return weather
	}
	
	targetDate := weather.Daily.Time[dayOffset]

	// Filter hourly data down to only the hours belonging to the target date,
	// so that time-of-day extraction (e.g. morning/evening) for this date
	// does not accidentally pick up another day's hours.
	var hourlyTime []string
	var hourlyTemp []float64
	var hourlyApparentTemp []float64
	var hourlyHumidity []int
	var hourlyWindSpeed []float64
	var hourlyWindDirection []float64
	var hourlyPrecipitation []float64
	var hourlyWeatherCode []int

	for i, t := range weather.Hourly.Time {
		if len(t) >= 10 && t[:10] == targetDate {
			hourlyTime = append(hourlyTime, t)
			hourlyTemp = append(hourlyTemp, safeFloat64At(weather.Hourly.Temperature, i))
			hourlyApparentTemp = append(hourlyApparentTemp, safeFloat64At(weather.Hourly.ApparentTemp, i))
			hourlyHumidity = append(hourlyHumidity, safeIntAt(weather.Hourly.Humidity, i))
			hourlyWindSpeed = append(hourlyWindSpeed, safeFloat64At(weather.Hourly.WindSpeed, i))
			hourlyWindDirection = append(hourlyWindDirection, safeFloat64At(weather.Hourly.WindDirection, i))
			hourlyPrecipitation = append(hourlyPrecipitation, safeFloat64At(weather.Hourly.Precipitation, i))
			hourlyWeatherCode = append(hourlyWeatherCode, safeIntAt(weather.Hourly.WeatherCode, i))
		}
	}

	// Create new weather data with selected day
	dateSpecificWeather := &types.WeatherData{
		Current: weather.Current,
		Hourly: struct {
			Time          []string  `json:"time"`
			Temperature   []float64 `json:"temperature_2m"`
			ApparentTemp  []float64 `json:"apparent_temperature"`
			Humidity      []int     `json:"relative_humidity_2m"`
			WindSpeed     []float64 `json:"wind_speed_10m"`
			WindDirection []float64 `json:"wind_direction_10m"`
			Precipitation []float64 `json:"precipitation"`
			WeatherCode   []int     `json:"weather_code"`
		}{
			Time:          hourlyTime,
			Temperature:   hourlyTemp,
			ApparentTemp:  hourlyApparentTemp,
			Humidity:      hourlyHumidity,
			WindSpeed:     hourlyWindSpeed,
			WindDirection: hourlyWindDirection,
			Precipitation: hourlyPrecipitation,
			WeatherCode:   hourlyWeatherCode,
		},
		Daily: struct {
			Time                        []string  `json:"time"`
			TemperatureMax              []float64 `json:"temperature_2m_max"`
			TemperatureMin              []float64 `json:"temperature_2m_min"`
			WindSpeedMax                []float64 `json:"wind_speed_10m_max"`
			WindGustMax                 []float64 `json:"wind_gusts_10m_max"`
			PrecipitationSum            []float64 `json:"precipitation_sum"`
			WeatherCode                 []int     `json:"weather_code"`
			SunriseTime                 []string  `json:"sunrise"`
			SunsetTime                  []string  `json:"sunset"`
			DaylightDuration            []float64 `json:"daylight_duration"`
			SunshineDuration            []float64 `json:"sunshine_duration"`
			UvIndexMax                  []float64 `json:"uv_index_max"`
			UvIndexClearSkyMax          []float64 `json:"uv_index_clear_sky_max"`
			PrecipitationHours          []float64 `json:"precipitation_hours"`
			PrecipitationProbabilityMax []float64 `json:"precipitation_probability_max"`
		}{
			Time:                        []string{weather.Daily.Time[dayOffset]},
			TemperatureMax:              []float64{weather.Daily.TemperatureMax[dayOffset]},
			TemperatureMin:              []float64{weather.Daily.TemperatureMin[dayOffset]},
			WindSpeedMax:                []float64{weather.Daily.WindSpeedMax[dayOffset]},
			WindGustMax:                 safeFloat64Slice(weather.Daily.WindGustMax, dayOffset),
			PrecipitationSum:            []float64{weather.Daily.PrecipitationSum[dayOffset]},
			WeatherCode:                 []int{weather.Daily.WeatherCode[dayOffset]},
			SunriseTime:                 safeStringSlice(weather.Daily.SunriseTime, dayOffset),
			SunsetTime:                  safeStringSlice(weather.Daily.SunsetTime, dayOffset),
			DaylightDuration:            safeFloat64Slice(weather.Daily.DaylightDuration, dayOffset),
			SunshineDuration:            safeFloat64Slice(weather.Daily.SunshineDuration, dayOffset),
			UvIndexMax:                  safeFloat64Slice(weather.Daily.UvIndexMax, dayOffset),
			UvIndexClearSkyMax:          safeFloat64Slice(weather.Daily.UvIndexClearSkyMax, dayOffset),
			PrecipitationHours:          safeFloat64Slice(weather.Daily.PrecipitationHours, dayOffset),
			PrecipitationProbabilityMax: safeFloat64Slice(weather.Daily.PrecipitationProbabilityMax, dayOffset),
		},
	}

	return dateSpecificWeather
}

// GetDateDisplayName returns Japanese display name for date specification
func GetDateDisplayName(dateSpec string) string {
	switch dateSpec {
	case "today":
		return "今日の"
	case "tomorrow":
		return "明日の"
	case "day-after-tomorrow":
		return "明後日の"
	default:
		return dateSpec + "の"
	}
}

// GetDateOffset returns day offset for date specification
func GetDateOffset(dateSpec string) int {
	switch dateSpec {
	case "today":
		return 0
	case "tomorrow":
		return 1
	case "day-after-tomorrow":
		return 2
	default:
		return 0
	}
}

// ValidateDateSpec validates if the date specification is valid
func ValidateDateSpec(dateSpec string) bool {
	validDates := []string{"today", "tomorrow", "day-after-tomorrow"}
	for _, valid := range validDates {
		if dateSpec == valid {
			return true
		}
	}
	return false
}

// ValidateTimeSpec validates if the time specification is valid
func ValidateTimeSpec(timeSpec string) bool {
	periods := GetTimePeriods()
	_, exists := periods[timeSpec]
	return exists
}

// safeStringSlice safely accesses string slice and returns slice with single element or empty slice
func safeStringSlice(slice []string, index int) []string {
	if len(slice) > index {
		return []string{slice[index]}
	}
	return []string{}
}

// safeFloat64At safely accesses a float64 slice element, returning 0 if out of range
func safeFloat64At(slice []float64, index int) float64 {
	if index >= 0 && index < len(slice) {
		return slice[index]
	}
	return 0
}

// safeIntAt safely accesses an int slice element, returning 0 if out of range
func safeIntAt(slice []int, index int) int {
	if index >= 0 && index < len(slice) {
		return slice[index]
	}
	return 0
}

// safeFloat64Slice safely accesses float64 slice and returns slice with single element or empty slice
func safeFloat64Slice(slice []float64, index int) []float64 {
	if len(slice) > index {
		return []float64{slice[index]}
	}
	return []float64{}
}