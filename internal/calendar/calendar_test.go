package calendar

import (
	"greekOrtho/internal/data"
	"greekOrtho/internal/models"
	"testing"
	"time"
)

func newCalendar(t *testing.T) *Calendar {
	t.Helper()
	d, err := data.Load()
	if err != nil {
		t.Fatalf("failed to load data: %v", err)
	}
	return New(d)
}

func TestGetDayInfo_Pascha2026(t *testing.T) {
	cal := newCalendar(t)
	date := time.Date(2026, 4, 12, 0, 0, 0, 0, time.UTC)
	info := cal.GetDayInfo(date)

	foundPascha := false
	for _, f := range info.Feasts {
		if f.Name == "Pascha (Resurrection of Christ)" {
			foundPascha = true
			break
		}
	}
	if !foundPascha {
		t.Error("expected Pascha feast on April 12, 2026")
	}
	if info.FastingLevel != models.FastingNone {
		t.Errorf("Pascha fasting: got %s, want none", info.FastingLevel)
	}
}

func TestGetDayInfo_Annunciation2026(t *testing.T) {
	cal := newCalendar(t)
	date := time.Date(2026, 3, 25, 0, 0, 0, 0, time.UTC)
	info := cal.GetDayInfo(date)

	foundFeast := false
	for _, f := range info.Feasts {
		if f.Name == "Annunciation of the Theotokos" {
			foundFeast = true
			break
		}
	}
	if !foundFeast {
		t.Error("expected Annunciation feast on March 25, 2026")
	}
	// During Great Lent but fish override
	if info.FastingLevel != models.FastingFish {
		t.Errorf("Annunciation fasting: got %s, want fish", info.FastingLevel)
	}
}

func TestGetDayInfo_Christmas(t *testing.T) {
	cal := newCalendar(t)
	date := time.Date(2026, 12, 25, 0, 0, 0, 0, time.UTC)
	info := cal.GetDayInfo(date)

	foundChristmas := false
	for _, f := range info.Feasts {
		if f.Name == "Nativity of Christ (Christmas)" {
			foundChristmas = true
			break
		}
	}
	if !foundChristmas {
		t.Error("expected Christmas feast on Dec 25")
	}
	if info.FastingLevel != models.FastingNone {
		t.Errorf("Christmas fasting: got %s, want none", info.FastingLevel)
	}
}

func TestGetDayInfo_BeheadingAug29(t *testing.T) {
	cal := newCalendar(t)
	date := time.Date(2026, 8, 29, 0, 0, 0, 0, time.UTC)
	info := cal.GetDayInfo(date)

	if info.FastingLevel != models.FastingStrict {
		t.Errorf("Beheading fasting: got %s, want strict", info.FastingLevel)
	}

	foundSaint := false
	for _, s := range info.Saints {
		if s.Month == 8 && s.Day == 29 {
			foundSaint = true
			break
		}
	}
	if !foundSaint {
		t.Error("expected saint entry for Aug 29")
	}
}

func TestGetDayInfo_QuoteDeterministic(t *testing.T) {
	cal := newCalendar(t)
	date := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
	q1 := cal.GetDayInfo(date).Quote
	q2 := cal.GetDayInfo(date).Quote
	if q1.Text != q2.Text {
		t.Error("quote should be deterministic for same date")
	}
}

func TestGetDayInfo_PalmSunday2026(t *testing.T) {
	cal := newCalendar(t)
	// Palm Sunday 2026 = April 5
	date := time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)
	info := cal.GetDayInfo(date)

	foundPalm := false
	for _, f := range info.Feasts {
		if f.Name == "Palm Sunday (Entry into Jerusalem)" {
			foundPalm = true
			break
		}
	}
	if !foundPalm {
		t.Error("expected Palm Sunday feast on April 5, 2026")
	}
	if info.FastingLevel != models.FastingFish {
		t.Errorf("Palm Sunday fasting: got %s, want fish", info.FastingLevel)
	}
}
