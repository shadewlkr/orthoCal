package models

import "time"

// FastingLevel represents the severity of a fast day.
type FastingLevel string

const (
	FastingStrict    FastingLevel = "strict"     // No meat, dairy, fish, oil, or wine
	FastingOilWine   FastingLevel = "oil_wine"   // Oil and wine permitted
	FastingFish      FastingLevel = "fish"       // Fish, oil, and wine permitted
	FastingDairyFish FastingLevel = "dairy_fish" // Dairy and fish permitted (Cheesefare)
	FastingNone      FastingLevel = "none"       // No fasting
)

// FastingLevelSeverity returns a numeric severity (lower = more strict).
func FastingLevelSeverity(level FastingLevel) int {
	switch level {
	case FastingStrict:
		return 0
	case FastingOilWine:
		return 1
	case FastingFish:
		return 2
	case FastingDairyFish:
		return 3
	case FastingNone:
		return 4
	default:
		return -1
	}
}

// WeekdayOverride allows different fasting levels on specific weekdays within a period.
type WeekdayOverride struct {
	Weekday time.Weekday `json:"weekday"`
	Level   FastingLevel `json:"level"`
}

// FastingRule defines a fasting period with priority-based resolution.
type FastingRule struct {
	Name              string            `json:"name"`
	Level             FastingLevel      `json:"level"`
	Priority          int               `json:"priority"`
	PaschaOffsetStart *int              `json:"pascha_offset_start,omitempty"` // Days from Pascha
	PaschaOffsetEnd   *int              `json:"pascha_offset_end,omitempty"`
	FixedStartMonth   *int              `json:"fixed_start_month,omitempty"`
	FixedStartDay     *int              `json:"fixed_start_day,omitempty"`
	FixedEndMonth     *int              `json:"fixed_end_month,omitempty"`
	FixedEndDay       *int              `json:"fixed_end_day,omitempty"`
	WeekdayOnly       *int              `json:"weekday_only,omitempty"` // 0=Sunday .. 6=Saturday
	WeekdayOverrides  []WeekdayOverride `json:"weekday_overrides,omitempty"`
	Description       string            `json:"description"`
}

// FeastRank indicates the importance of a feast day.
type FeastRank string

const (
	RankGreat FeastRank = "great"
	RankMajor FeastRank = "major"
	RankMinor FeastRank = "minor"
)

// Feast represents a fixed or moveable feast day.
type Feast struct {
	Name            string        `json:"name"`
	GreekName       string        `json:"greek_name,omitempty"`
	Description     string        `json:"description,omitempty"`
	Rank            FeastRank     `json:"rank"`
	Month           *int          `json:"month,omitempty"`         // For fixed feasts
	Day             *int          `json:"day,omitempty"`           // For fixed feasts
	PaschaOffset    *int          `json:"pascha_offset,omitempty"` // For moveable feasts
	FastingOverride *FastingLevel `json:"fasting_override,omitempty"`
}

// Saint represents a commemorated saint on a given date.
type Saint struct {
	Name        string `json:"name"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Month       int    `json:"month"`
	Day         int    `json:"day"`
}

// Quote represents a Church Father or saint's quote.
type Quote struct {
	Text   string `json:"text"`
	Author string `json:"author"`
	Source string `json:"source,omitempty"`
}

// ScriptureReading represents a single scripture citation.
type ScriptureReading struct {
	Book    string `json:"book"`    // e.g., "Hebrews"
	Passage string `json:"passage"` // e.g., "13:17-21"
}

// DayReadings represents a pair of epistle and gospel readings with their source.
type DayReadings struct {
	Epistle *ScriptureReading `json:"epistle,omitempty"`
	Gospel  *ScriptureReading `json:"gospel,omitempty"`
	Source  string            `json:"source,omitempty"` // "Cycle" or "Feast: Name"
}

// DayInfo is the composite result returned by GetDayInfo for display.
type DayInfo struct {
	Date          time.Time
	Feasts        []Feast
	Saints        []Saint
	FastingLevel  FastingLevel
	FastingReason string
	Readings      []DayReadings
	Quote         Quote
}
