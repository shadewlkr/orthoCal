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

// CalendarData holds all loaded calendar data.
type CalendarData struct {
	FixedFeasts    []models.Feast
	MoveableFeasts []models.Feast
	Saints         []models.Saint
	FastingRules   []models.FastingRule
	Quotes         []models.Quote
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

	return &d, nil
}

// Ensure embed import is used (for go:embed directives).
var _ embed.FS
