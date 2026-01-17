package types

// WeatherData represents weather information from API
type WeatherData struct {
	Current struct {
		Temperature   float64 `json:"temperature_2m"`
		ApparentTemp  float64 `json:"apparent_temperature"`
		Humidity      int     `json:"relative_humidity_2m"`
		WindSpeed     float64 `json:"wind_speed_10m"`
		WindDirection float64 `json:"wind_direction_10m"`
		Precipitation float64 `json:"precipitation"`
		WeatherCode   int     `json:"weather_code"`
	} `json:"current"`
	Hourly struct {
		Time          []string  `json:"time"`
		Temperature   []float64 `json:"temperature_2m"`
		ApparentTemp  []float64 `json:"apparent_temperature"`
		Humidity      []int     `json:"relative_humidity_2m"`
		WindSpeed     []float64 `json:"wind_speed_10m"`
		WindDirection []float64 `json:"wind_direction_10m"`
		Precipitation []float64 `json:"precipitation"`
		WeatherCode   []int     `json:"weather_code"`
	} `json:"hourly"`
	Daily struct {
		Time                []string  `json:"time"`
		TemperatureMax      []float64 `json:"temperature_2m_max"`
		TemperatureMin      []float64 `json:"temperature_2m_min"`
		WindSpeedMax        []float64 `json:"wind_speed_10m_max"`
		WindGustMax         []float64 `json:"wind_gusts_10m_max"`
		PrecipitationSum    []float64 `json:"precipitation_sum"`
		WeatherCode         []int     `json:"weather_code"`
		SunriseTime         []string  `json:"sunrise"`
		SunsetTime          []string  `json:"sunset"`
		DaylightDuration    []float64 `json:"daylight_duration"`
		SunshineDuration    []float64 `json:"sunshine_duration"`
		UvIndexMax          []float64 `json:"uv_index_max"`
		UvIndexClearSkyMax  []float64 `json:"uv_index_clear_sky_max"`
		PrecipitationHours  []float64 `json:"precipitation_hours"`
		PrecipitationProbabilityMax []float64 `json:"precipitation_probability_max"`
	} `json:"daily"`
}

// CityCoordinate represents city name and coordinates
type CityCoordinate struct {
	Name string
	Lat  float64
	Lon  float64
}

// TimeBasedWeather represents weather data for specific time
type TimeBasedWeather struct {
	Time          string
	Temperature   float64
	ApparentTemp  float64
	Humidity      int
	WindSpeed     float64
	WindDirection float64
	Precipitation float64
	WeatherCode   int
}

// TimePeriod represents time period definition
type TimePeriod struct {
	Key         string
	DisplayName string
	StartHour   int
	EndHour     int
}

// DistanceCategory represents running distance category
type DistanceCategory struct {
	Key             string
	DisplayName     string
	Description     string
	MinKm           float64
	MaxKm           float64
	TempPenalty     int
	HumidityPenalty int
	WindPenalty     int
	HeatIndexPenalty int
}

// RunningCondition represents running condition assessment
type RunningCondition struct {
	Score          int
	Level          string
	Recommendation string
	Warnings       []string
	Clothing       []string
}

// AirQualityData represents air quality information from API
type AirQualityData struct {
	Hourly struct {
		Time  []string  `json:"time"`
		Dust  []float64 `json:"dust"`
		PM10  []float64 `json:"pm10"`
		PM2_5 []float64 `json:"pm2_5"`
	} `json:"hourly"`
}

// DustLevel represents dust concentration level
type DustLevel struct {
	Level       int
	DisplayName string
	Description string
	Dust        float64
	PM10        float64
	PM2_5       float64
}