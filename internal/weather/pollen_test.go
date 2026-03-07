package weather

import (
	"testing"

	"runcast/internal/types"
)

func TestGetCityCode(t *testing.T) {
	tests := []struct {
		city     string
		wantCode string
		wantOk   bool
	}{
		{"tokyo", "13101", true},
		{"osaka", "27128", true},
		{"sapporo", "01101", true},
		{"naha", "47201", true},
		{"unknown", "", false},
		{"home", "", false}, // custom location
	}

	for _, tt := range tests {
		t.Run(tt.city, func(t *testing.T) {
			code, ok := GetCityCode(tt.city)
			if ok != tt.wantOk {
				t.Errorf("GetCityCode(%q) ok = %v, want %v", tt.city, ok, tt.wantOk)
			}
			if ok && code != tt.wantCode {
				t.Errorf("GetCityCode(%q) code = %v, want %v", tt.city, code, tt.wantCode)
			}
		})
	}
}

func TestCreatePollenLevel(t *testing.T) {
	tests := []struct {
		pollen      int
		wantLevel   int
		wantDisplay string
	}{
		{0, 0, "なし"},
		{1, 1, "少ない"},
		{29, 1, "少ない"},
		{30, 2, "やや多い"},
		{99, 2, "やや多い"},
		{100, 3, "多い"},
		{199, 3, "多い"},
		{200, 4, "非常に多い"},
		{500, 4, "非常に多い"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			pl := createPollenLevel(tt.pollen)
			if pl.Level != tt.wantLevel {
				t.Errorf("pollen=%d: level = %d, want %d", tt.pollen, pl.Level, tt.wantLevel)
			}
			if pl.DisplayName != tt.wantDisplay {
				t.Errorf("pollen=%d: display = %q, want %q", tt.pollen, pl.DisplayName, tt.wantDisplay)
			}
			if pl.Pollen != tt.pollen {
				t.Errorf("pollen=%d: stored pollen = %d", tt.pollen, pl.Pollen)
			}
		})
	}
}

func TestGetCurrentPollenLevel(t *testing.T) {
	t.Run("empty data returns nil", func(t *testing.T) {
		pl := GetCurrentPollenLevel(nil)
		if pl != nil {
			t.Errorf("expected nil, got %v", pl)
		}
	})

	t.Run("falls back to last entry", func(t *testing.T) {
		data := []types.PollenData{
			{Time: "2000-01-01T00:00:00+09:00", Pollen: 5},
			{Time: "2000-01-01T01:00:00+09:00", Pollen: 10},
		}
		pl := GetCurrentPollenLevel(data)
		if pl == nil {
			t.Fatal("expected non-nil PollenLevel")
		}
		if pl.Pollen != 10 {
			t.Errorf("expected pollen=10, got %d", pl.Pollen)
		}
	})
}

func TestGetHourlyPollenLevel(t *testing.T) {
	data := []types.PollenData{
		{Time: "2026-03-07T06:00:00+09:00", Pollen: 5},
		{Time: "2026-03-07T07:00:00+09:00", Pollen: 50},
		{Time: "2026-03-07T08:00:00+09:00", Pollen: 120},
	}

	t.Run("matches correct hour", func(t *testing.T) {
		pl := GetHourlyPollenLevel(data, 7)
		if pl == nil {
			t.Fatal("expected non-nil")
		}
		if pl.Pollen != 50 {
			t.Errorf("expected 50, got %d", pl.Pollen)
		}
		if pl.Level != 2 {
			t.Errorf("expected level 2 (やや多い), got %d", pl.Level)
		}
	})

	t.Run("no match returns nil", func(t *testing.T) {
		pl := GetHourlyPollenLevel(data, 12)
		if pl != nil {
			t.Errorf("expected nil, got %v", pl)
		}
	})

	t.Run("nil data returns nil", func(t *testing.T) {
		pl := GetHourlyPollenLevel(nil, 7)
		if pl != nil {
			t.Errorf("expected nil, got %v", pl)
		}
	})
}
