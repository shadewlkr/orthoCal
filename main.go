package main

import (
	"flag"
	"fmt"
	"greekOrtho/internal/calendar"
	"greekOrtho/internal/data"
	"greekOrtho/internal/display"
	"greekOrtho/internal/models"
	"os"
	"time"
)

func main() {
	dateFlag := flag.String("date", "", "Date to display in YYYY-MM-DD format (defaults to today)")
	simpleFlag := flag.Bool("simple", false, "One-liner output suitable for piping or status bars")
	monthFlag := flag.Bool("month", false, "Show monthly calendar grid")
	flag.Parse()

	if *simpleFlag && *monthFlag {
		fmt.Fprintf(os.Stderr, "Error: --simple and --month are mutually exclusive\n")
		os.Exit(1)
	}

	var date time.Time
	if *dateFlag != "" {
		var err error
		date, err = time.Parse("2006-01-02", *dateFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid date format %q (use YYYY-MM-DD)\n", *dateFlag)
			os.Exit(1)
		}
	} else {
		date = time.Now()
		// Normalize to midnight UTC for consistent behavior
		date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	}

	d, err := data.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading calendar data: %v\n", err)
		os.Exit(1)
	}

	cal := calendar.New(d)

	switch {
	case *simpleFlag:
		info := cal.GetDayInfo(date)
		display.PrintSimple(info)

	case *monthFlag:
		today := time.Now()
		today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)

		year, month, _ := date.Date()
		firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		daysInMonth := firstOfMonth.AddDate(0, 1, -1).Day()

		days := make([]models.DayInfo, daysInMonth)
		for i := 0; i < daysInMonth; i++ {
			d := firstOfMonth.AddDate(0, 0, i)
			days[i] = cal.GetDayInfo(d)
		}
		display.PrintMonth(days, today)

	default:
		info := cal.GetDayInfo(date)
		display.PrintDayInfo(info)
	}
}
