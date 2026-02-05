package main

import (
	"flag"
	"fmt"
	"greekOrtho/internal/calendar"
	"greekOrtho/internal/data"
	"greekOrtho/internal/display"
	"os"
	"time"
)

func main() {
	dateFlag := flag.String("date", "", "Date to display in YYYY-MM-DD format (defaults to today)")
	flag.Parse()

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
	info := cal.GetDayInfo(date)
	display.PrintDayInfo(info)
}
