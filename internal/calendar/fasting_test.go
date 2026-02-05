package calendar

import (
	"greekOrtho/internal/data"
	"greekOrtho/internal/models"
	"greekOrtho/internal/pascha"
	"testing"
	"time"
)

func mustLoad(t *testing.T) *data.CalendarData {
	t.Helper()
	d, err := data.Load()
	if err != nil {
		t.Fatalf("failed to load data: %v", err)
	}
	return d
}

func TestFasting_GreatLentWeekday(t *testing.T) {
	d := mustLoad(t)
	// 2026 Pascha = April 12. Clean Monday = Feb 23. A Monday in Lent: March 2, 2026
	p := pascha.Compute(2026)
	date := time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC)
	level, _ := ResolveFasting(date, p, d.FastingRules, nil)
	if level != models.FastingStrict {
		t.Errorf("Great Lent weekday: got %s, want strict", level)
	}
}

func TestFasting_GreatLentSaturday(t *testing.T) {
	d := mustLoad(t)
	p := pascha.Compute(2026)
	// Saturday March 7, 2026 — should be oil_wine
	date := time.Date(2026, 3, 7, 0, 0, 0, 0, time.UTC)
	level, _ := ResolveFasting(date, p, d.FastingRules, nil)
	if level != models.FastingOilWine {
		t.Errorf("Great Lent Saturday: got %s, want oil_wine", level)
	}
}

func TestFasting_HolyWeek(t *testing.T) {
	d := mustLoad(t)
	p := pascha.Compute(2026)
	// Holy Wednesday = April 8, 2026
	date := time.Date(2026, 4, 8, 0, 0, 0, 0, time.UTC)
	level, _ := ResolveFasting(date, p, d.FastingRules, nil)
	if level != models.FastingStrict {
		t.Errorf("Holy Week: got %s, want strict", level)
	}
}

func TestFasting_BrightWeek(t *testing.T) {
	d := mustLoad(t)
	p := pascha.Compute(2026)
	// Bright Wednesday = April 15, 2026
	date := time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC)
	level, _ := ResolveFasting(date, p, d.FastingRules, nil)
	if level != models.FastingNone {
		t.Errorf("Bright Week: got %s, want none", level)
	}
}

func TestFasting_WednesdayRegular(t *testing.T) {
	d := mustLoad(t)
	p := pascha.Compute(2026)
	// A regular Wednesday outside any fast period: July 1, 2026
	date := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	level, _ := ResolveFasting(date, p, d.FastingRules, nil)
	if level != models.FastingOilWine {
		t.Errorf("Regular Wednesday: got %s, want oil_wine", level)
	}
}

func TestFasting_RegularTuesday(t *testing.T) {
	d := mustLoad(t)
	p := pascha.Compute(2026)
	// A regular Tuesday outside any fast period: June 30, 2026
	date := time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)
	level, _ := ResolveFasting(date, p, d.FastingRules, nil)
	if level != models.FastingNone {
		t.Errorf("Regular Tuesday: got %s, want none", level)
	}
}

func TestFasting_AnnunciationDuringLent(t *testing.T) {
	d := mustLoad(t)
	p := pascha.Compute(2026)
	date := time.Date(2026, 3, 25, 0, 0, 0, 0, time.UTC)
	// Pass Annunciation feast with fish override
	fishLevel := models.FastingFish
	feasts := []models.Feast{
		{Name: "Annunciation of the Theotokos", FastingOverride: &fishLevel},
	}
	level, _ := ResolveFasting(date, p, d.FastingRules, feasts)
	if level != models.FastingFish {
		t.Errorf("Annunciation during Lent: got %s, want fish", level)
	}
}

func TestFasting_BeheadingStrictOverride(t *testing.T) {
	d := mustLoad(t)
	p := pascha.Compute(2026)
	date := time.Date(2026, 8, 29, 0, 0, 0, 0, time.UTC) // Saturday
	strictLevel := models.FastingStrict
	feasts := []models.Feast{
		{Name: "Beheading of John the Baptist", FastingOverride: &strictLevel},
	}
	level, _ := ResolveFasting(date, p, d.FastingRules, feasts)
	if level != models.FastingStrict {
		t.Errorf("Beheading of John Baptist: got %s, want strict", level)
	}
}

func TestFasting_DormitionFast(t *testing.T) {
	d := mustLoad(t)
	p := pascha.Compute(2026)
	// Aug 5, 2026 is a Wednesday — weekday in Dormition fast
	date := time.Date(2026, 8, 5, 0, 0, 0, 0, time.UTC)
	level, _ := ResolveFasting(date, p, d.FastingRules, nil)
	if level != models.FastingStrict {
		t.Errorf("Dormition Fast weekday: got %s, want strict", level)
	}
}

func TestFasting_NativityFast(t *testing.T) {
	d := mustLoad(t)
	p := pascha.Compute(2026)
	// Dec 1, 2026 is a Tuesday — should be fish
	date := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)
	level, _ := ResolveFasting(date, p, d.FastingRules, nil)
	if level != models.FastingFish {
		t.Errorf("Nativity Fast Tuesday: got %s, want fish", level)
	}
}

func TestFasting_CheesfareWeek(t *testing.T) {
	d := mustLoad(t)
	p := pascha.Compute(2026)
	// Cheesefare week 2026: Feb 16-22 (Pascha offset -55 to -49)
	date := time.Date(2026, 2, 18, 0, 0, 0, 0, time.UTC) // Wednesday
	level, _ := ResolveFasting(date, p, d.FastingRules, nil)
	if level != models.FastingDairyFish {
		t.Errorf("Cheesefare Week: got %s, want dairy_fish", level)
	}
}

func TestFasting_Christmas(t *testing.T) {
	d := mustLoad(t)
	p := pascha.Compute(2026)
	date := time.Date(2026, 12, 25, 0, 0, 0, 0, time.UTC)
	level, _ := ResolveFasting(date, p, d.FastingRules, nil)
	if level != models.FastingNone {
		t.Errorf("Christmas: got %s, want none", level)
	}
}
