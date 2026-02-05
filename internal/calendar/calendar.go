package calendar

import (
	"greekOrtho/internal/data"
	"greekOrtho/internal/models"
	"greekOrtho/internal/pascha"
	"time"
)

// Calendar provides methods to look up liturgical info for any date.
type Calendar struct {
	data *data.CalendarData
}

// New creates a new Calendar with the embedded data.
func New(d *data.CalendarData) *Calendar {
	return &Calendar{data: d}
}

// GetDayInfo returns the complete liturgical information for a given date.
func (c *Calendar) GetDayInfo(date time.Time) models.DayInfo {
	p := pascha.Compute(date.Year())

	feasts := c.findFeasts(date, p)
	saints := c.findSaints(date)
	fastingLevel, fastingReason := ResolveFasting(date, p, c.data.FastingRules, feasts)
	readings := ResolveReadings(date, p, c.data, feasts)
	quote := c.selectQuote(date)

	return models.DayInfo{
		Date:          date,
		Feasts:        feasts,
		Saints:        saints,
		FastingLevel:  fastingLevel,
		FastingReason: fastingReason,
		Readings:      readings,
		Quote:         quote,
	}
}

// findFeasts returns all feasts (fixed and moveable) that fall on the given date.
func (c *Calendar) findFeasts(date time.Time, p time.Time) []models.Feast {
	var result []models.Feast

	for _, f := range c.data.FixedFeasts {
		if f.Month != nil && f.Day != nil {
			if int(date.Month()) == *f.Month && date.Day() == *f.Day {
				result = append(result, f)
			}
		}
	}

	daysFromPascha := int(date.Sub(p).Hours() / 24)
	for _, f := range c.data.MoveableFeasts {
		if f.PaschaOffset != nil {
			if daysFromPascha == *f.PaschaOffset {
				result = append(result, f)
			}
		}
	}

	return result
}

// findSaints returns all saints commemorated on the given date.
func (c *Calendar) findSaints(date time.Time) []models.Saint {
	var result []models.Saint
	for _, s := range c.data.Saints {
		if int(date.Month()) == s.Month && date.Day() == s.Day {
			result = append(result, s)
		}
	}
	return result
}

// selectQuote returns a deterministic quote for the given date (day-of-year modulo).
func (c *Calendar) selectQuote(date time.Time) models.Quote {
	if len(c.data.Quotes) == 0 {
		return models.Quote{Text: "Lord, have mercy.", Author: "The Church"}
	}
	doy := date.YearDay()
	idx := doy % len(c.data.Quotes)
	return c.data.Quotes[idx]
}
