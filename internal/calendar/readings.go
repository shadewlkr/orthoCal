package calendar

import (
	"fmt"
	"greekOrtho/internal/data"
	"greekOrtho/internal/models"
	"time"
)

// ResolveReadings determines the scripture readings for a given date.
// It checks feast readings first (fixed and moveable), then falls back to the
// lectionary cycle (epistle cycle + gospel series with Lukan Jump computation).
func ResolveReadings(date time.Time, pascha time.Time, d *data.CalendarData, feasts []models.Feast) []models.DayReadings {
	daysFromPascha := int(date.Sub(pascha).Hours() / 24)

	// 1. Check feast readings (both fixed and moveable)
	feastReadings := resolveFeastReadings(date, daysFromPascha, d)

	// 2. Resolve cycle readings
	cycleReadings := resolveCycleReadings(date, pascha, daysFromPascha, d)

	// 3. Combine: great feasts replace, minor/major supplement
	return combineReadings(cycleReadings, feastReadings, feasts)
}

// resolveFeastReadings checks fixed (month/day) and moveable (pascha-offset) feast readings.
func resolveFeastReadings(date time.Time, daysFromPascha int, d *data.CalendarData) *models.DayReadings {
	// Check fixed feasts
	key := fmt.Sprintf("%d/%d", date.Month(), date.Day())
	if entry, ok := d.FeastReadings.Fixed[key]; ok {
		reading := models.DayReadings{
			Epistle: entry.Epistle,
			Gospel:  entry.Gospel,
			Source:  "Feast",
		}
		return &reading
	}

	// Check moveable feasts
	offsetKey := fmt.Sprintf("%d", daysFromPascha)
	if entry, ok := d.FeastReadings.Moveable[offsetKey]; ok {
		reading := models.DayReadings{
			Epistle: entry.Epistle,
			Gospel:  entry.Gospel,
			Source:  "Feast",
		}
		return &reading
	}

	return nil
}

// resolveCycleReadings looks up the epistle and gospel from the lectionary cycle tables.
func resolveCycleReadings(date time.Time, pascha time.Time, daysFromPascha int, d *data.CalendarData) *models.DayReadings {
	weekday := fmt.Sprintf("%d", date.Weekday())

	epistle := resolveEpistle(daysFromPascha, weekday, d)
	gospel := resolveGospel(date, pascha, daysFromPascha, weekday, d)

	if epistle == nil && gospel == nil {
		return nil
	}

	return &models.DayReadings{
		Epistle: epistle,
		Gospel:  gospel,
		Source:  "Cycle",
	}
}

// resolveEpistle looks up the epistle reading from the cycle table.
func resolveEpistle(daysFromPascha int, weekday string, d *data.CalendarData) *models.ScriptureReading {
	if daysFromPascha < 0 {
		lentStart := -48   // Clean Monday
		triodionStart := -70 // Sunday of Publican and Pharisee

		if daysFromPascha >= lentStart && daysFromPascha < -7 {
			// During Great Lent (Clean Monday to Saturday before Palm Sunday)
			lentDay := daysFromPascha - lentStart
			lentWeek := (lentDay / 7) + 33 // Map to tail of epistle cycle
			weekKey := fmt.Sprintf("%d", lentWeek)
			if weekData, ok := d.EpistleCycle[weekKey]; ok {
				if reading, ok := weekData[weekday]; ok {
					return &reading
				}
			}
		} else if daysFromPascha >= triodionStart && daysFromPascha < lentStart {
			// Pre-Lenten weeks (Triodion before Clean Monday)
			// Continue from weeks 30-32 of the epistle cycle
			preLentDay := daysFromPascha - triodionStart
			preLentWeek := (preLentDay / 7) + 30
			weekKey := fmt.Sprintf("%d", preLentWeek)
			if weekData, ok := d.EpistleCycle[weekKey]; ok {
				if reading, ok := weekData[weekday]; ok {
					return &reading
				}
			}
		}
		return nil
	}

	weekOfPentecost := (daysFromPascha / 7) + 1
	weekKey := fmt.Sprintf("%d", weekOfPentecost)

	if weekData, ok := d.EpistleCycle[weekKey]; ok {
		if reading, ok := weekData[weekday]; ok {
			return &reading
		}
	}
	return nil
}

// resolveGospel determines which gospel series applies and looks up the reading.
func resolveGospel(date time.Time, pascha time.Time, daysFromPascha int, weekday string, d *data.CalendarData) *models.ScriptureReading {
	// Before Pascha: Lenten period
	if daysFromPascha < 0 {
		return resolveLentenGospel(daysFromPascha, weekday, d)
	}

	pentecost := 49 // Day 49 from Pascha

	// John series: Pascha to Pentecost (days 0-49)
	if daysFromPascha <= pentecost {
		weekOfPentecost := (daysFromPascha / 7) + 1
		weekKey := fmt.Sprintf("%d", weekOfPentecost)
		if weekData, ok := d.GospelCycle.John[weekKey]; ok {
			if reading, ok := weekData[weekday]; ok {
				return &reading
			}
		}
		return nil
	}

	// After Pentecost: Matthew → Luke series with Lukan Jump
	elevation := time.Date(date.Year(), 9, 14, 0, 0, 0, 0, time.UTC)
	pentecostDate := pascha.AddDate(0, 0, pentecost)

	// Weeks from Pentecost to Elevation of the Cross
	daysToElevation := int(elevation.Sub(pentecostDate).Hours() / 24)
	matthewWeeks := daysToElevation / 7
	if matthewWeeks < 1 {
		matthewWeeks = 1
	}

	// Days since Pentecost (day after Pentecost = day 1 of Matthew series)
	daysSincePentecost := daysFromPascha - pentecost
	matthewWeek := (daysSincePentecost / 7) + 1

	// Lukan Jump: if Matthew series is shorter than 16 weeks, Luke starts earlier
	lukanJump := 0
	if matthewWeeks < 17 {
		lukanJump = 17 - matthewWeeks
	}

	// Are we still in the Matthew series?
	if matthewWeek <= matthewWeeks || matthewWeek <= 16 {
		weekKey := fmt.Sprintf("%d", matthewWeek)
		if weekData, ok := d.GospelCycle.Matthew[weekKey]; ok {
			if reading, ok := weekData[weekday]; ok {
				return &reading
			}
		}
		return nil
	}

	// Luke series
	lukeWeek := matthewWeek - matthewWeeks + lukanJump
	if lukeWeek < 1 {
		lukeWeek = 1
	}

	weekKey := fmt.Sprintf("%d", lukeWeek)
	if weekData, ok := d.GospelCycle.Luke[weekKey]; ok {
		if reading, ok := weekData[weekday]; ok {
			return &reading
		}
	}

	return nil
}

// resolveLentenGospel handles the Lenten and pre-Lenten period gospel readings.
func resolveLentenGospel(daysFromPascha int, weekday string, d *data.CalendarData) *models.ScriptureReading {
	lentStart := -48      // Clean Monday
	triodionStart := -70  // Sunday of Publican and Pharisee
	holyWeekStart := -7

	wd := int(0)
	fmt.Sscanf(weekday, "%d", &wd)

	// Pre-Lenten period (Triodion before Clean Monday)
	if daysFromPascha >= triodionStart && daysFromPascha < lentStart {
		// Continue Luke series from weeks 14-17
		preLentDay := daysFromPascha - triodionStart
		preLentWeek := (preLentDay / 7) + 14
		weekKey := fmt.Sprintf("%d", preLentWeek)
		if weekData, ok := d.GospelCycle.Luke[weekKey]; ok {
			if reading, ok := weekData[weekday]; ok {
				return &reading
			}
		}
		return nil
	}

	// Holy Week — handled by feast readings
	if daysFromPascha >= holyWeekStart {
		return nil
	}

	// Great Lent (Clean Monday to Saturday before Palm Sunday)
	if daysFromPascha >= lentStart {
		lentDay := daysFromPascha - lentStart
		lentWeek := (lentDay / 7) + 1

		// Saturdays and Sundays during Lent continue the Luke series
		if wd == 0 || wd == 6 {
			weekKey := fmt.Sprintf("%d", lentWeek+11)
			if weekData, ok := d.GospelCycle.Luke[weekKey]; ok {
				if reading, ok := weekData[weekday]; ok {
					return &reading
				}
			}
			return nil
		}

		// Weekdays use the Lenten/Mark series
		weekKey := fmt.Sprintf("%d", lentWeek)
		if weekData, ok := d.GospelCycle.Lenten[weekKey]; ok {
			if reading, ok := weekData[weekday]; ok {
				return &reading
			}
		}
	}

	return nil
}

// combineReadings merges feast and cycle readings according to feast rank.
func combineReadings(cycle *models.DayReadings, feast *models.DayReadings, feasts []models.Feast) []models.DayReadings {
	if feast == nil && cycle == nil {
		return nil
	}
	if feast == nil {
		return []models.DayReadings{*cycle}
	}
	if cycle == nil {
		return []models.DayReadings{*feast}
	}

	// Check if any feast is a great feast — great feasts replace cycle readings
	hasGreatFeast := false
	for _, f := range feasts {
		if f.Rank == models.RankGreat {
			hasGreatFeast = true
			break
		}
	}

	if hasGreatFeast {
		return []models.DayReadings{*feast}
	}

	// Minor/major feasts: show both cycle and feast readings
	return []models.DayReadings{*cycle, *feast}
}
