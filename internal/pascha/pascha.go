package pascha

import "time"

// Compute returns the date of Orthodox Pascha (Easter) for the given year.
// Uses the Meeus Julian Easter algorithm and converts to the Gregorian calendar.
func Compute(year int) time.Time {
	// Meeus algorithm for Julian Easter
	a := year % 4
	b := year % 7
	c := year % 19
	d := (19*c + 15) % 30
	e := (2*a + 4*b - d + 34) % 7
	month := (d + e + 114) / 31
	day := ((d + e + 114) % 31) + 1

	// This gives the Julian calendar date. Convert to Gregorian.
	// The offset between Julian and Gregorian calendars:
	// For years 1900-2099, the offset is 13 days.
	// General formula: century = year/100; offset = century - century/4 - 2
	century := year / 100
	offset := century - century/4 - 2

	julian := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	gregorian := julian.AddDate(0, 0, offset)
	return gregorian
}
