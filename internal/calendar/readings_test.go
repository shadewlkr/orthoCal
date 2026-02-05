package calendar

import (
	"greekOrtho/internal/data"
	"greekOrtho/internal/models"
	"testing"
	"time"
)

func TestResolveReadings_Pascha(t *testing.T) {
	d, err := data.Load()
	if err != nil {
		t.Fatalf("failed to load data: %v", err)
	}

	// Pascha 2026 is April 12
	pascha := time.Date(2026, 4, 12, 0, 0, 0, 0, time.UTC)
	date := pascha

	feasts := []models.Feast{{Name: "Pascha", Rank: models.RankGreat}}
	readings := ResolveReadings(date, pascha, d, feasts)

	if len(readings) == 0 {
		t.Fatal("expected readings for Pascha, got none")
	}

	r := readings[0]
	if r.Gospel == nil {
		t.Fatal("expected Gospel reading for Pascha")
	}
	if r.Gospel.Book != "John" || r.Gospel.Passage != "1:1-17" {
		t.Errorf("expected John 1:1-17, got %s %s", r.Gospel.Book, r.Gospel.Passage)
	}
}

func TestResolveReadings_Pentecost(t *testing.T) {
	d, err := data.Load()
	if err != nil {
		t.Fatalf("failed to load data: %v", err)
	}

	// Pascha 2026 is April 12, Pentecost is 49 days later = May 31
	pascha := time.Date(2026, 4, 12, 0, 0, 0, 0, time.UTC)
	date := time.Date(2026, 5, 31, 0, 0, 0, 0, time.UTC)

	feasts := []models.Feast{{Name: "Pentecost", Rank: models.RankGreat}}
	readings := ResolveReadings(date, pascha, d, feasts)

	if len(readings) == 0 {
		t.Fatal("expected readings for Pentecost, got none")
	}

	r := readings[0]
	if r.Gospel == nil {
		t.Fatal("expected Gospel reading for Pentecost")
	}
	if r.Gospel.Book != "John" {
		t.Errorf("expected John gospel for Pentecost, got %s", r.Gospel.Book)
	}
}

func TestResolveReadings_FixedFeast(t *testing.T) {
	d, err := data.Load()
	if err != nil {
		t.Fatalf("failed to load data: %v", err)
	}

	// Theophany - January 6
	pascha := time.Date(2026, 4, 12, 0, 0, 0, 0, time.UTC)
	date := time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC)

	feasts := []models.Feast{{Name: "Theophany", Rank: models.RankGreat}}
	readings := ResolveReadings(date, pascha, d, feasts)

	if len(readings) == 0 {
		t.Fatal("expected readings for Theophany, got none")
	}

	r := readings[0]
	if r.Gospel == nil {
		t.Fatal("expected Gospel reading for Theophany")
	}
	if r.Gospel.Book != "Matthew" || r.Gospel.Passage != "3:13-17" {
		t.Errorf("expected Matthew 3:13-17, got %s %s", r.Gospel.Book, r.Gospel.Passage)
	}
}

func TestResolveReadings_ElevationOfCross(t *testing.T) {
	d, err := data.Load()
	if err != nil {
		t.Fatalf("failed to load data: %v", err)
	}

	// Elevation of the Cross - September 14
	pascha := time.Date(2026, 4, 12, 0, 0, 0, 0, time.UTC)
	date := time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC)

	feasts := []models.Feast{{Name: "Elevation of the Holy Cross", Rank: models.RankGreat}}
	readings := ResolveReadings(date, pascha, d, feasts)

	if len(readings) == 0 {
		t.Fatal("expected readings for Elevation, got none")
	}

	r := readings[0]
	if r.Epistle == nil {
		t.Fatal("expected Epistle reading for Elevation")
	}
	if r.Epistle.Book != "1 Corinthians" {
		t.Errorf("expected 1 Corinthians epistle, got %s", r.Epistle.Book)
	}
}

func TestResolveReadings_RegularSunday(t *testing.T) {
	d, err := data.Load()
	if err != nil {
		t.Fatalf("failed to load data: %v", err)
	}

	// A regular Sunday after Pentecost (June 7, 2026 = 1 week after Pentecost)
	pascha := time.Date(2026, 4, 12, 0, 0, 0, 0, time.UTC)
	date := time.Date(2026, 6, 7, 0, 0, 0, 0, time.UTC)

	readings := ResolveReadings(date, pascha, d, nil)

	if len(readings) == 0 {
		t.Fatal("expected cycle readings for regular Sunday, got none")
	}

	r := readings[0]
	if r.Source != "Cycle" {
		t.Errorf("expected Cycle source, got %s", r.Source)
	}
	if r.Epistle == nil {
		t.Error("expected Epistle reading")
	}
	if r.Gospel == nil {
		t.Error("expected Gospel reading")
	}
}

func TestResolveReadings_Annunciation(t *testing.T) {
	d, err := data.Load()
	if err != nil {
		t.Fatalf("failed to load data: %v", err)
	}

	// Annunciation - March 25
	pascha := time.Date(2026, 4, 12, 0, 0, 0, 0, time.UTC)
	date := time.Date(2026, 3, 25, 0, 0, 0, 0, time.UTC)

	feasts := []models.Feast{{Name: "Annunciation of the Theotokos", Rank: models.RankGreat}}
	readings := ResolveReadings(date, pascha, d, feasts)

	if len(readings) == 0 {
		t.Fatal("expected readings for Annunciation, got none")
	}

	r := readings[0]
	if r.Gospel == nil {
		t.Fatal("expected Gospel reading for Annunciation")
	}
	if r.Gospel.Book != "Luke" {
		t.Errorf("expected Luke gospel for Annunciation, got %s", r.Gospel.Book)
	}
}
