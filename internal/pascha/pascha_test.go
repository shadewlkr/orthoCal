package pascha

import (
	"testing"
	"time"
)

func TestCompute(t *testing.T) {
	// Known Orthodox Pascha dates (Gregorian)
	tests := []struct {
		year int
		want time.Time
	}{
		{2024, time.Date(2024, 5, 5, 0, 0, 0, 0, time.UTC)},
		{2025, time.Date(2025, 4, 20, 0, 0, 0, 0, time.UTC)},
		{2026, time.Date(2026, 4, 12, 0, 0, 0, 0, time.UTC)},
		{2027, time.Date(2027, 5, 2, 0, 0, 0, 0, time.UTC)},
		{2028, time.Date(2028, 4, 16, 0, 0, 0, 0, time.UTC)},
		{2029, time.Date(2029, 4, 8, 0, 0, 0, 0, time.UTC)},
		{2030, time.Date(2030, 4, 28, 0, 0, 0, 0, time.UTC)},
	}

	for _, tt := range tests {
		got := Compute(tt.year)
		if !got.Equal(tt.want) {
			t.Errorf("Compute(%d) = %s, want %s", tt.year, got.Format("2006-01-02"), tt.want.Format("2006-01-02"))
		}
	}
}
