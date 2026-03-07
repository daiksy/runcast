package weather

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"runcast/internal/types"
)

const pollenAPIURL = "https://wxtech.weathernews.com/opendata/v1/pollen"

// cityCodes maps city keys to Weathernews citycode
var cityCodes = map[string]string{
	"tokyo":     "13101",
	"osaka":     "27128",
	"kyoto":     "26104",
	"yokohama":  "14103",
	"nagoya":    "23105",
	"sapporo":   "01101",
	"fukuoka":   "40131",
	"sendai":    "04101",
	"hiroshima": "34101",
	"naha":      "47201",
	"kobe":      "28110",
	"shiga":     "25201",
}

// GetCityCode returns Weathernews citycode for a given city key
func GetCityCode(cityKey string) (string, bool) {
	code, ok := cityCodes[cityKey]
	return code, ok
}

// GetPollen fetches pollen data from Weathernews API
func GetPollen(cityCode string, days int) ([]types.PollenData, error) {
	base := time.Now().AddDate(0, 0, days)
	dateStr := base.Format("20060102")

	url := fmt.Sprintf("%s?citycode=%s&start=%s&end=%s", pollenAPIURL, cityCode, dateStr, dateStr)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("pollen API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pollen API request failed with status: %d", resp.StatusCode)
	}

	r := csv.NewReader(resp.Body)
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse pollen CSV: %w", err)
	}

	var result []types.PollenData
	for _, record := range records[1:] { // skip header
		if len(record) < 3 {
			continue
		}
		pollen, err := strconv.Atoi(strings.TrimSpace(record[2]))
		if err != nil || pollen == -9999 {
			continue
		}
		result = append(result, types.PollenData{
			Time:   record[1],
			Pollen: pollen,
		})
	}

	return result, nil
}

// GetCurrentPollenLevel returns current pollen level from pollen data
func GetCurrentPollenLevel(pollenData []types.PollenData) *types.PollenLevel {
	if len(pollenData) == 0 {
		return nil
	}

	now := time.Now()
	currentHour := now.Format("2006-01-02T15")

	// Find the entry matching the current hour
	for _, d := range pollenData {
		if strings.HasPrefix(d.Time, currentHour) {
			return createPollenLevel(d.Pollen)
		}
	}

	// Fall back to latest valid entry
	return createPollenLevel(pollenData[len(pollenData)-1].Pollen)
}

// GetHourlyPollenLevel returns pollen level for a specific hour
func GetHourlyPollenLevel(pollenData []types.PollenData, hour int) *types.PollenLevel {
	if len(pollenData) == 0 {
		return nil
	}

	target := fmt.Sprintf("T%02d:", hour)
	for _, d := range pollenData {
		if strings.Contains(d.Time, target) {
			return createPollenLevel(d.Pollen)
		}
	}

	return nil
}

// createPollenLevel creates a PollenLevel from a raw pollen count
func createPollenLevel(pollen int) *types.PollenLevel {
	level := 0
	displayName := "なし"

	switch {
	case pollen >= 200:
		level = 4
		displayName = "非常に多い"
	case pollen >= 100:
		level = 3
		displayName = "多い"
	case pollen >= 30:
		level = 2
		displayName = "やや多い"
	case pollen >= 1:
		level = 1
		displayName = "少ない"
	}

	return &types.PollenLevel{
		Level:       level,
		DisplayName: displayName,
		Pollen:      pollen,
	}
}
