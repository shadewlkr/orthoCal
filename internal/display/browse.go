package display

import (
	"fmt"
	"greekOrtho/internal/models"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

const (
	reverseVideo = "\033[7m"
	clearScreen  = "\033[2J\033[H"
)

// Browse runs an interactive calendar browser starting at startDate.
// getDayInfo is called to fetch liturgical data for any date the user navigates to.
func Browse(getDayInfo func(time.Time) models.DayInfo, startDate time.Time) error {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("cannot enter raw terminal mode: %w", err)
	}
	defer term.Restore(fd, oldState)

	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	selected := startDate

	buf := make([]byte, 3)
	for {
		screen := renderBrowseScreen(getDayInfo, selected, today)
		fmt.Fprint(os.Stdout, screen)

		n, err := os.Stdin.Read(buf)
		if err != nil {
			break
		}

		switch {
		case n == 1 && buf[0] == 'q', n == 1 && buf[0] == 0x03:
			fmt.Fprint(os.Stdout, clearScreen)
			return nil
		case n == 1 && buf[0] == 't':
			selected = today
		case n == 1 && buf[0] == 'n':
			selected = nextMonth(selected)
		case n == 1 && buf[0] == 'p':
			selected = prevMonth(selected)
		case n == 3 && buf[0] == 0x1b && buf[1] == '[':
			switch buf[2] {
			case 'A': // up
				selected = selected.AddDate(0, 0, -7)
			case 'B': // down
				selected = selected.AddDate(0, 0, 7)
			case 'C': // right
				selected = selected.AddDate(0, 0, 1)
			case 'D': // left
				selected = selected.AddDate(0, 0, -1)
			}
		}
	}
	return nil
}

func nextMonth(d time.Time) time.Time {
	year, month, day := d.Date()
	next := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	lastDay := next.AddDate(0, 1, -1).Day()
	if day > lastDay {
		day = lastDay
	}
	return time.Date(next.Year(), next.Month(), day, 0, 0, 0, 0, time.UTC)
}

func prevMonth(d time.Time) time.Time {
	year, month, day := d.Date()
	prev := time.Date(year, month-1, 1, 0, 0, 0, 0, time.UTC)
	lastDay := prev.AddDate(0, 1, -1).Day()
	if day > lastDay {
		day = lastDay
	}
	return time.Date(prev.Year(), prev.Month(), day, 0, 0, 0, 0, time.UTC)
}

func renderBrowseScreen(getDayInfo func(time.Time) models.DayInfo, selected, today time.Time) string {
	var sb strings.Builder
	sb.WriteString(clearScreen)

	sb.WriteString(renderBrowseMonth(getDayInfo, selected, today))
	sb.WriteString("\r\n")
	sb.WriteString(dim + strings.Repeat(divHoriz, 60) + reset + "\r\n")
	sb.WriteString("\r\n")

	info := getDayInfo(selected)
	sb.WriteString(renderBrowseDayInfo(info))

	sb.WriteString("\r\n")
	sb.WriteString(dim + strings.Repeat(divHoriz, 60) + reset + "\r\n")
	sb.WriteString(dimWhite + " ← → ↑ ↓ Navigate   n/p Month   t Today   q Quit" + reset + "\r\n")

	return sb.String()
}

func renderBrowseMonth(getDayInfo func(time.Time) models.DayInfo, selected, today time.Time) string {
	var sb strings.Builder

	year, month, _ := selected.Date()
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	daysInMonth := firstOfMonth.AddDate(0, 1, -1).Day()
	startWeekday := int(firstOfMonth.Weekday())

	// Pre-fetch day info for the whole month
	dayInfos := make(map[int]models.DayInfo, daysInMonth)
	for i := 1; i <= daysInMonth; i++ {
		d := time.Date(year, month, i, 0, 0, 0, 0, time.UTC)
		dayInfos[i] = getDayInfo(d)
	}

	title := fmt.Sprintf("☦  %s %d", month, year)
	sb.WriteString("\r\n")
	sb.WriteString(" " + boldGold + title + reset + "\r\n")
	sb.WriteString("\r\n")

	header := " "
	for _, d := range []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"} {
		header += fmt.Sprintf(" %-7s", d)
	}
	sb.WriteString(bold + header + reset + "\r\n")

	dayNum := 1
	for row := 0; dayNum <= daysInMonth; row++ {
		rowStr := " "
		for col := 0; col < 7; col++ {
			if row == 0 && col < startWeekday {
				rowStr += strings.Repeat(" ", cellWidth)
				continue
			}
			if dayNum > daysInMonth {
				rowStr += strings.Repeat(" ", cellWidth)
				continue
			}

			info := dayInfos[dayNum]
			fastColor, _ := fastingStyle(info.FastingLevel)

			numStr := fmt.Sprintf("%d", dayNum)
			hasFeast := len(info.Feasts) > 0
			if hasFeast {
				numStr += "✦"
			}

			isSelected := selected.Day() == dayNum
			isToday := today.Year() == year &&
				today.Month() == month &&
				today.Day() == dayNum

			cell := ""
			switch {
			case isSelected && isToday:
				cell = reverseVideo + bold + underline + fastColor + numStr + reset
			case isSelected:
				cell = reverseVideo + bold + fastColor + numStr + reset
			case isToday:
				cell = bold + underline + fastColor + numStr + reset
			default:
				cell = fastColor + numStr + reset
			}

			// Pad to cell width accounting for multi-byte ✦ (3 bytes in UTF-8)
			visualLen := len(numStr)
			if hasFeast {
				visualLen = visualLen - 2 // ✦ is 3 bytes but ~1 visual char
			}
			padding := cellWidth - visualLen
			if padding < 1 {
				padding = 1
			}
			cell += strings.Repeat(" ", padding)

			rowStr += cell
			dayNum++
		}
		sb.WriteString(rowStr + "\r\n")
	}

	sb.WriteString("\r\n")
	legend := " 🔴 Strict  🟠 Oil/Wine  🟡 Fish  🟢 No Fast  ✦ Feast"
	sb.WriteString(legend + "\r\n")

	return sb.String()
}

func renderBrowseDayInfo(info models.DayInfo) string {
	var sb strings.Builder

	sb.WriteString(" " + boldWhite + info.Date.Format("Monday, January 2, 2006") + reset + "\r\n")
	sb.WriteString("\r\n")

	// Feasts
	if len(info.Feasts) > 0 {
		for _, f := range info.Feasts {
			sb.WriteString(" " + boldGold + "✦ " + f.Name + reset + "\r\n")
			if f.Rank != "" {
				sb.WriteString("   " + yellow + rankDisplay(f.Rank) + reset + "\r\n")
			}
			if f.GreekName != "" {
				sb.WriteString("   " + dimWhite + f.GreekName + reset + "\r\n")
			}
		}
		sb.WriteString("\r\n")
	}

	// Saints
	if len(info.Saints) > 0 {
		sb.WriteString(" " + boldCyan + "Saints Commemorated" + reset + "\r\n")
		for _, s := range info.Saints {
			title := ""
			if s.Title != "" {
				title = dimWhite + " — " + s.Title + reset
			}
			sb.WriteString("   " + cyan + "• " + s.Name + title + reset + "\r\n")
		}
		sb.WriteString("\r\n")
	}

	// Fasting
	fastColor, fastIcon := fastingStyle(info.FastingLevel)
	sb.WriteString(" " + bold + fastIcon + " Fasting" + reset + "\r\n")
	sb.WriteString("   " + fastColor + FastingDescription(info.FastingLevel) + reset + "\r\n")
	if info.FastingReason != "" {
		sb.WriteString("   " + dimWhite + info.FastingReason + reset + "\r\n")
	}
	sb.WriteString("\r\n")

	// Scripture Readings
	if len(info.Readings) > 0 {
		sb.WriteString(" " + bold + "📖 Scripture Readings" + reset + "\r\n")
		for _, r := range info.Readings {
			if r.Epistle != nil {
				sb.WriteString("   " + blue + "Epistle: " + r.Epistle.Book + " " + r.Epistle.Passage + reset + "\r\n")
			}
			if r.Gospel != nil {
				sb.WriteString("   " + blue + "Gospel:  " + r.Gospel.Book + " " + r.Gospel.Passage + reset + "\r\n")
			}
		}
		sb.WriteString("\r\n")
	}

	// Quote
	sb.WriteString(" " + bold + "✼ Quote of the Day" + reset + "\r\n")
	maxWidth := 56
	words := strings.Fields(info.Quote.Text)
	if len(words) > 0 {
		lines := wrapWords(words, maxWidth)
		for i, l := range lines {
			if i == 0 {
				sb.WriteString("   " + italic + magenta + "\"" + l + reset + "\r\n")
			} else if i == len(lines)-1 {
				sb.WriteString("   " + italic + magenta + " " + l + "\"" + reset + "\r\n")
			} else {
				sb.WriteString("   " + italic + magenta + " " + l + reset + "\r\n")
			}
		}
		attribution := "     " + dimWhite + "— " + info.Quote.Author
		if info.Quote.Source != "" {
			attribution += ", " + info.Quote.Source
		}
		attribution += reset
		sb.WriteString(attribution + "\r\n")
	}

	return sb.String()
}
