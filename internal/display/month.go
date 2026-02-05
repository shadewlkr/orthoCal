package display

import (
	"fmt"
	"greekOrtho/internal/models"
	"strings"
	"time"
)

const (
	underline = "\033[4m"
	cellWidth = 8
)

// PrintMonth renders a monthly calendar grid with fasting colors and feast markers.
func PrintMonth(days []models.DayInfo, today time.Time) {
	if len(days) == 0 {
		return
	}

	month := days[0].Date.Month()
	year := days[0].Date.Year()

	fmt.Println()
	fmt.Println(topBorder())

	// Header
	fmt.Println(emptyLine())
	title := fmt.Sprintf("☦  Greek Orthodox Calendar — %s %d", month, year)
	fmt.Println(line(boldGold + "  " + title + reset))
	fmt.Println(emptyLine())

	// Day-of-week header
	fmt.Println(divider())
	fmt.Println(emptyLine())
	header := "  "
	for _, d := range []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"} {
		header += fmt.Sprintf("%-8s", d)
	}
	fmt.Println(line(bold + header + reset))
	fmt.Println(emptyLine())

	// Build lookup maps
	dayInfo := make(map[int]models.DayInfo)
	for _, d := range days {
		dayInfo[d.Date.Day()] = d
	}

	// Determine starting weekday offset
	firstDay := days[0].Date
	startWeekday := int(firstDay.Weekday()) // 0=Sun
	daysInMonth := len(days)

	// Print calendar grid
	dayNum := 1
	for row := 0; dayNum <= daysInMonth; row++ {
		rowStr := "  "
		for col := 0; col < 7; col++ {
			if row == 0 && col < startWeekday {
				rowStr += strings.Repeat(" ", cellWidth)
				continue
			}
			if dayNum > daysInMonth {
				rowStr += strings.Repeat(" ", cellWidth)
				continue
			}

			info := dayInfo[dayNum]
			fastColor, _ := fastingStyle(info.FastingLevel)

			numStr := fmt.Sprintf("%d", dayNum)
			hasFeast := len(info.Feasts) > 0
			if hasFeast {
				numStr += "✦"
			}

			// Highlight today
			isToday := today.Year() == year &&
				today.Month() == month &&
				today.Day() == dayNum

			cell := ""
			if isToday {
				cell = bold + underline + fastColor + numStr + reset
			} else {
				cell = fastColor + numStr + reset
			}

			// Pad to cell width accounting for multi-byte ✦ (3 bytes in UTF-8)
			visualLen := len(numStr)
			if hasFeast {
				visualLen = len(numStr) - 2 // ✦ is 3 bytes but 1 visual char... actually visual width varies
			}
			padding := cellWidth - visualLen
			if padding < 1 {
				padding = 1
			}
			cell += strings.Repeat(" ", padding)

			rowStr += cell
			dayNum++
		}
		fmt.Println(line(rowStr))
	}
	fmt.Println(emptyLine())

	// Legend
	fmt.Println(divider())
	legend := "  🔴 Strict  🟠 Oil/Wine  🟡 Fish  🟢 No Fast  ✦ Feast"
	fmt.Println(line(legend))

	// Feasts this month
	var feasts []string
	for _, d := range days {
		for _, f := range d.Feasts {
			entry := fmt.Sprintf("  %s %d — %s",
				d.Date.Format("Jan"), d.Date.Day(), f.Name)
			feasts = append(feasts, entry)
		}
	}

	if len(feasts) > 0 {
		fmt.Println(divider())
		fmt.Println(line(bold + "  Feasts this month:" + reset))
		for _, f := range feasts {
			fmt.Println(line(yellow + f + reset))
		}
	}

	fmt.Println(bottomBorder())
	fmt.Println()
}
