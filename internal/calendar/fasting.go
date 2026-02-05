package calendar

import (
	"greekOrtho/internal/models"
	"time"
)

// ResolveFasting determines the fasting level and reason for a given date.
// It evaluates all rules, picks the highest-priority matching rule, applies
// weekday overrides, and then applies feast-day fasting overrides (only if more lenient).
func ResolveFasting(date time.Time, pascha time.Time, rules []models.FastingRule, feasts []models.Feast) (models.FastingLevel, string) {
	daysFromPascha := int(date.Sub(pascha).Hours() / 24)

	var bestRule *models.FastingRule
	var bestPriority int = -1

	for i := range rules {
		r := &rules[i]
		if !ruleMatches(date, daysFromPascha, pascha, r) {
			continue
		}
		if r.Priority > bestPriority {
			bestPriority = r.Priority
			bestRule = r
		}
	}

	level := models.FastingNone
	reason := "No fasting today"

	if bestRule != nil {
		level = bestRule.Level

		// Apply weekday overrides within the period
		for _, ov := range bestRule.WeekdayOverrides {
			if date.Weekday() == ov.Weekday {
				level = ov.Level
				break
			}
		}
		reason = bestRule.Description
	}

	// Apply feast-day fasting overrides — only if more lenient
	for _, f := range feasts {
		if f.FastingOverride != nil {
			override := *f.FastingOverride
			if override == models.FastingStrict {
				// Strict overrides (e.g., Beheading, Elevation of Cross) always apply
				if models.FastingLevelSeverity(override) < models.FastingLevelSeverity(level) {
					level = override
					reason = f.Name + " — strict fast day"
				}
			} else if models.FastingLevelSeverity(override) > models.FastingLevelSeverity(level) {
				// More lenient overrides (e.g., Annunciation = fish during Lent)
				level = override
				reason = f.Name + " — fasting relaxed for the feast"
			}
		}
	}

	return level, reason
}

// ruleMatches checks if a fasting rule applies to the given date.
func ruleMatches(date time.Time, daysFromPascha int, pascha time.Time, r *models.FastingRule) bool {
	// Weekday-only rules (Wed/Fri)
	if r.WeekdayOnly != nil {
		return int(date.Weekday()) == *r.WeekdayOnly
	}

	// Pure Pascha-offset range
	if r.PaschaOffsetStart != nil && r.PaschaOffsetEnd != nil &&
		r.FixedStartMonth == nil && r.FixedEndMonth == nil {
		return daysFromPascha >= *r.PaschaOffsetStart && daysFromPascha <= *r.PaschaOffsetEnd
	}

	// Pure fixed date range
	if r.FixedStartMonth != nil && r.FixedStartDay != nil &&
		r.FixedEndMonth != nil && r.FixedEndDay != nil &&
		r.PaschaOffsetStart == nil {
		return inFixedRange(date, *r.FixedStartMonth, *r.FixedStartDay, *r.FixedEndMonth, *r.FixedEndDay)
	}

	// Hybrid: Pascha-offset start, fixed date end (Apostles' Fast)
	if r.PaschaOffsetStart != nil && r.FixedEndMonth != nil && r.FixedEndDay != nil {
		start := pascha.AddDate(0, 0, *r.PaschaOffsetStart)
		end := time.Date(date.Year(), time.Month(*r.FixedEndMonth), *r.FixedEndDay, 0, 0, 0, 0, time.UTC)
		// The fast only exists if start is before end
		if !start.Before(end) {
			return false
		}
		return !date.Before(start) && !date.After(end)
	}

	return false
}

// inFixedRange checks if date falls within a fixed month/day range, handling year boundaries.
func inFixedRange(date time.Time, startMonth, startDay, endMonth, endDay int) bool {
	m := int(date.Month())
	d := date.Day()

	startVal := startMonth*100 + startDay
	endVal := endMonth*100 + endDay
	dateVal := m*100 + d

	if startVal <= endVal {
		// Same-year range (e.g., Aug 1 - Aug 14)
		return dateVal >= startVal && dateVal <= endVal
	}
	// Wraps around year boundary (e.g., Dec 25 - Jan 4)
	return dateVal >= startVal || dateVal <= endVal
}
