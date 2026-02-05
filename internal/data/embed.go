package data

import (
	"embed"
	"encoding/json"
	"fmt"
	"greekOrtho/internal/models"
)

//go:embed fixed_feasts.json
var fixedFeastsJSON []byte

//go:embed moveable_feasts.json
var moveableFeastsJSON []byte

//go:embed saints.json
var saintsJSON []byte

//go:embed fasting_rules.json
var fastingRulesJSON []byte

//go:embed quotes.json
var quotesJSON []byte

//go:embed epistle_cycle.json
var epistleCycleJSON []byte

//go:embed gospel_cycle.json
var gospelCycleJSON []byte

//go:embed feast_readings.json
var feastReadingsJSON []byte

// EpistleCycle maps week-of-Pentecost (string) → weekday (string "0"-"6") → reading.
type EpistleCycle map[string]map[string]models.ScriptureReading

// GospelCycle contains the four gospel series.
type GospelCycle struct {
	John    map[string]map[string]models.ScriptureReading `json:"john"`
	Matthew map[string]map[string]models.ScriptureReading `json:"matthew"`
	Luke    map[string]map[string]models.ScriptureReading `json:"luke"`
	Lenten  map[string]map[string]models.ScriptureReading `json:"lenten"`
}

// FeastReadingEntry holds a single feast's readings.
type FeastReadingEntry struct {
	Rank    string                   `json:"rank"`
	Epistle *models.ScriptureReading `json:"epistle,omitempty"`
	Gospel  *models.ScriptureReading `json:"gospel,omitempty"`
}

// FeastReadings holds fixed (by month/day) and moveable (by pascha-offset) feast readings.
type FeastReadings struct {
	Fixed    map[string]FeastReadingEntry `json:"fixed"`
	Moveable map[string]FeastReadingEntry `json:"moveable"`
}

// CalendarData holds all loaded calendar data.
type CalendarData struct {
	FixedFeasts    []models.Feast
	MoveableFeasts []models.Feast
	Saints         []models.Saint
	FastingRules   []models.FastingRule
	Quotes         []models.Quote
	EpistleCycle   EpistleCycle
	GospelCycle    GospelCycle
	FeastReadings  FeastReadings
}

// Load parses all embedded JSON data and returns a CalendarData struct.
func Load() (*CalendarData, error) {
	var d CalendarData

	if err := json.Unmarshal(fixedFeastsJSON, &d.FixedFeasts); err != nil {
		return nil, fmt.Errorf("parsing fixed_feasts.json: %w", err)
	}
	if err := json.Unmarshal(moveableFeastsJSON, &d.MoveableFeasts); err != nil {
		return nil, fmt.Errorf("parsing moveable_feasts.json: %w", err)
	}
	if err := json.Unmarshal(saintsJSON, &d.Saints); err != nil {
		return nil, fmt.Errorf("parsing saints.json: %w", err)
	}
	if err := json.Unmarshal(fastingRulesJSON, &d.FastingRules); err != nil {
		return nil, fmt.Errorf("parsing fasting_rules.json: %w", err)
	}
	if err := json.Unmarshal(quotesJSON, &d.Quotes); err != nil {
		return nil, fmt.Errorf("parsing quotes.json: %w", err)
	}
	if err := json.Unmarshal(epistleCycleJSON, &d.EpistleCycle); err != nil {
		return nil, fmt.Errorf("parsing epistle_cycle.json: %w", err)
	}
	if err := json.Unmarshal(gospelCycleJSON, &d.GospelCycle); err != nil {
		return nil, fmt.Errorf("parsing gospel_cycle.json: %w", err)
	}
	if err := json.Unmarshal(feastReadingsJSON, &d.FeastReadings); err != nil {
		return nil, fmt.Errorf("parsing feast_readings.json: %w", err)
	}

	return &d, nil
}

// Ensure embed import is used (for go:embed directives).
var _ embed.FS
