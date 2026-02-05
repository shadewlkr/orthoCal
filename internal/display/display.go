package display

import (
	"fmt"
	"greekOrtho/internal/models"
	"strings"
)

// ANSI color codes
const (
	reset     = "\033[0m"
	bold      = "\033[1m"
	dim       = "\033[2m"
	italic    = "\033[3m"
	red       = "\033[31m"
	green     = "\033[32m"
	yellow    = "\033[33m"
	blue      = "\033[34m"
	magenta   = "\033[35m"
	cyan      = "\033[36m"
	white     = "\033[37m"
	boldRed   = "\033[1;31m"
	boldGold  = "\033[1;33m"
	boldCyan  = "\033[1;36m"
	boldWhite = "\033[1;37m"
	dimWhite  = "\033[2;37m"
)

const contentWidth = 56

// Box-drawing characters
const (
	topLeft     = "╔"
	topRight    = "╗"
	bottomLeft  = "╚"
	bottomRight = "╝"
	horizontal  = "═"
	vertical    = "║"
	divLeft     = "╟"
	divRight    = "╢"
	divHoriz    = "─"
)

func topBorder() string {
	return dim + topLeft + strings.Repeat(horizontal, contentWidth+2) + topRight + reset
}

func bottomBorder() string {
	return dim + bottomLeft + strings.Repeat(horizontal, contentWidth+2) + bottomRight + reset
}

func divider() string {
	return dim + divLeft + strings.Repeat(divHoriz, contentWidth+2) + divRight + reset
}

func line(content string) string {
	return dim + vertical + reset + " " + content
}

func emptyLine() string {
	return dim + vertical + reset
}

// FastingDescription returns a human-readable description of the fasting level.
func FastingDescription(level models.FastingLevel) string {
	switch level {
	case models.FastingStrict:
		return "Strict Fast (no meat, dairy, fish, oil, or wine)"
	case models.FastingOilWine:
		return "Oil and Wine Permitted (no meat, dairy, or fish)"
	case models.FastingFish:
		return "Fish, Oil, and Wine Permitted (no meat or dairy)"
	case models.FastingDairyFish:
		return "Dairy and Fish Permitted (no meat)"
	case models.FastingNone:
		return "No Fast"
	default:
		return "Unknown"
	}
}

// PrintDayInfo formats and prints the day's liturgical information.
func PrintDayInfo(info models.DayInfo) {
	fmt.Println()
	fmt.Println(topBorder())

	// Header
	fmt.Println(emptyLine())
	fmt.Println(line(boldGold + "  ☦  Greek Orthodox Calendar" + reset))
	fmt.Println(line(boldWhite + "  " + info.Date.Format("Monday, January 2, 2006") + reset))
	fmt.Println(emptyLine())

	// Feasts
	if len(info.Feasts) > 0 {
		fmt.Println(divider())
		fmt.Println(emptyLine())
		for i, f := range info.Feasts {
			fmt.Println(line(boldGold + "  ✦ " + f.Name + reset))
			if f.Rank != "" {
				fmt.Println(line(yellow + "    " + rankDisplay(f.Rank) + reset))
			}
			if f.GreekName != "" {
				fmt.Println(line(dimWhite + "    " + f.GreekName + reset))
			}
			if i < len(info.Feasts)-1 {
				fmt.Println(emptyLine())
			}
		}
		fmt.Println(emptyLine())
	}

	// Saints
	if len(info.Saints) > 0 {
		fmt.Println(divider())
		fmt.Println(emptyLine())
		fmt.Println(line(boldCyan + "  Saints Commemorated" + reset))
		for _, s := range info.Saints {
			title := ""
			if s.Title != "" {
				title = dimWhite + " — " + s.Title + reset
			}
			fmt.Println(line(cyan + "    • " + s.Name + title))
		}
		fmt.Println(emptyLine())
	}

	// Fasting
	fmt.Println(divider())
	fmt.Println(emptyLine())
	fastColor, fastIcon := fastingStyle(info.FastingLevel)
	fmt.Println(line(bold + "  " + fastIcon + " Fasting" + reset))
	fmt.Println(line(fastColor + "    " + FastingDescription(info.FastingLevel) + reset))
	if info.FastingReason != "" {
		fmt.Println(line(dimWhite + "    " + info.FastingReason + reset))
	}
	fmt.Println(emptyLine())

	// Scripture Readings
	if len(info.Readings) > 0 {
		fmt.Println(divider())
		fmt.Println(emptyLine())
		fmt.Println(line(bold + "  📖 Scripture Readings" + reset))
		for _, r := range info.Readings {
			if r.Epistle != nil {
				fmt.Println(line(blue + "    Epistle: " + r.Epistle.Book + " " + r.Epistle.Passage + reset))
			}
			if r.Gospel != nil {
				fmt.Println(line(blue + "    Gospel:  " + r.Gospel.Book + " " + r.Gospel.Passage + reset))
			}
		}
		fmt.Println(emptyLine())
	}

	// Quote
	fmt.Println(divider())
	fmt.Println(emptyLine())
	fmt.Println(line(bold + "  ✼ Quote of the Day" + reset))
	printWrappedQuote(info.Quote)

	fmt.Println(bottomBorder())
	fmt.Println()
}

func fastingStyle(level models.FastingLevel) (string, string) {
	switch level {
	case models.FastingStrict:
		return boldRed, "🔴"
	case models.FastingOilWine:
		return red, "🟠"
	case models.FastingFish:
		return yellow, "🟡"
	case models.FastingDairyFish:
		return yellow, "🟡"
	case models.FastingNone:
		return green, "🟢"
	default:
		return white, "○"
	}
}

func rankDisplay(r models.FeastRank) string {
	switch r {
	case models.RankGreat:
		return "Great Feast"
	case models.RankMajor:
		return "Major Feast"
	case models.RankMinor:
		return "Minor Observance"
	default:
		return string(r)
	}
}

// printWrappedQuote prints a quote with word wrapping and proper indentation.
func printWrappedQuote(q models.Quote) {
	maxWidth := contentWidth - 8
	words := strings.Fields(q.Text)
	if len(words) == 0 {
		return
	}

	lines := wrapWords(words, maxWidth)

	for i, l := range lines {
		if i == 0 {
			fmt.Println(line(italic + magenta + "    \"" + l + reset))
		} else if i == len(lines)-1 {
			fmt.Println(line(italic + magenta + "     " + l + "\"" + reset))
		} else {
			fmt.Println(line(italic + magenta + "     " + l + reset))
		}
	}
	attribution := dimWhite + "      — " + q.Author
	if q.Source != "" {
		attribution += ", " + q.Source
	}
	attribution += reset
	fmt.Println(line(attribution))
	fmt.Println(emptyLine())
}

// shortFastingLabel returns a short plain-text label for the fasting level.
func shortFastingLabel(level models.FastingLevel) string {
	switch level {
	case models.FastingStrict:
		return "Strict"
	case models.FastingOilWine:
		return "Oil & Wine"
	case models.FastingFish:
		return "Fish"
	case models.FastingDairyFish:
		return "Dairy & Fish"
	case models.FastingNone:
		return "No Fast"
	default:
		return "Unknown"
	}
}

// PrintSimple prints a one-liner summary suitable for piping, shell prompts, or status bars.
// Format: Thu Feb 5 | 🟠 Oil & Wine | St. Agatha | Lk 6:17-23
func PrintSimple(info models.DayInfo) {
	_, icon := fastingStyle(info.FastingLevel)
	label := shortFastingLabel(info.FastingLevel)

	parts := []string{
		info.Date.Format("Mon Jan 2"),
		icon + " " + label,
	}

	if len(info.Feasts) > 0 {
		names := make([]string, len(info.Feasts))
		for i, f := range info.Feasts {
			names[i] = "✦ " + f.Name
		}
		parts = append(parts, strings.Join(names, ", "))
	} else if len(info.Saints) > 0 {
		parts = append(parts, info.Saints[0].Name)
	}

	// Append gospel citation
	if len(info.Readings) > 0 && info.Readings[0].Gospel != nil {
		g := info.Readings[0].Gospel
		parts = append(parts, shortBookName(g.Book)+" "+g.Passage)
	}

	fmt.Println(strings.Join(parts, " | "))
}

// shortBookName returns abbreviated book names for compact display.
func shortBookName(book string) string {
	abbrevs := map[string]string{
		"Matthew":         "Mt",
		"Mark":            "Mk",
		"Luke":            "Lk",
		"John":            "Jn",
		"Acts":            "Acts",
		"Romans":          "Rom",
		"1 Corinthians":   "1 Cor",
		"2 Corinthians":   "2 Cor",
		"Galatians":       "Gal",
		"Ephesians":       "Eph",
		"Philippians":     "Phil",
		"Colossians":      "Col",
		"1 Thessalonians": "1 Thess",
		"2 Thessalonians": "2 Thess",
		"1 Timothy":       "1 Tim",
		"2 Timothy":       "2 Tim",
		"Titus":           "Titus",
		"Philemon":        "Phlm",
		"Hebrews":         "Heb",
		"James":           "Jas",
		"1 Peter":         "1 Pet",
		"2 Peter":         "2 Pet",
		"1 John":          "1 Jn",
		"2 John":          "2 Jn",
		"3 John":          "3 Jn",
		"Jude":            "Jude",
		"Revelation":      "Rev",
	}
	if abbrev, ok := abbrevs[book]; ok {
		return abbrev
	}
	return book
}

func wrapWords(words []string, maxWidth int) []string {
	var lines []string
	current := words[0]
	for _, w := range words[1:] {
		if len(current)+1+len(w) > maxWidth {
			lines = append(lines, current)
			current = w
		} else {
			current += " " + w
		}
	}
	lines = append(lines, current)
	return lines
}
