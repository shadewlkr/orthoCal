# orthoCal

A command-line Greek Orthodox liturgical calendar that displays daily fasting rules, feast days, saints, scripture readings, and quotes.

## Installation

```bash
go build -o orthoCal
```

Requires Go 1.16+ (uses embedded data files).

## Usage

```
orthoCal [options]
```

### Options

| Flag | Description |
|------|-------------|
| `-date YYYY-MM-DD` | Display information for a specific date (default: today) |
| `-simple` | One-line output suitable for scripts, status bars, or shell prompts |
| `-month` | Display a monthly calendar grid |

### Examples

**Today's liturgical information:**
```bash
./orthoCal
```

**Specific date:**
```bash
./orthoCal -date 2026-04-12
```

**One-liner for status bar:**
```bash
./orthoCal --simple
# Output: Thu Feb 5 | 🟢 No Fast | Mk 1:29-35
```

**Monthly calendar grid:**
```bash
./orthoCal --month
./orthoCal --month -date 2026-04-01
```

## Output Sections

### Default View

The default view displays a formatted box with:

- **Date** — Current or specified date
- **Feasts** — Great, major, or minor feast days with Greek names
- **Saints** — Commemorated saints for the day
- **Fasting** — Fasting level with description and reason
- **Scripture Readings** — Daily Epistle and Gospel citations
- **Quote** — Daily quote from Church Fathers

### Fasting Indicators

| Icon | Level | Description |
|------|-------|-------------|
| 🔴 | Strict | No meat, dairy, fish, oil, or wine |
| 🟠 | Oil & Wine | Oil and wine permitted |
| 🟡 | Fish | Fish, oil, and wine permitted |
| 🟡 | Dairy & Fish | Dairy and fish permitted (Cheesefare) |
| 🟢 | No Fast | No fasting restrictions |

### Simple Output Format

```
Day Mon DD | 🟢 Fast Level | Feast/Saint | Gospel Citation
```

Example:
```
Sun Apr 12 | 🟢 No Fast | ✦ Pascha (Resurrection of Christ) | Jn 1:1-17
```

## Scripture Readings

Readings follow the Orthodox lectionary cycle:

- **Paschal Season** (Pascha to Pentecost): John series
- **After Pentecost**: Matthew series, then Luke series with Lukan Jump
- **Great Lent weekdays**: Mark series
- **Feast days**: Override or supplement cycle readings

The lectionary data is embedded and computed algorithmically, so readings are accurate for any year without external API calls.

## Data Sources

All liturgical data is embedded in the binary:

- `fixed_feasts.json` — Fixed-date feasts (Nativity, Theophany, etc.)
- `moveable_feasts.json` — Pascha-relative feasts (Palm Sunday, Pentecost, etc.)
- `saints.json` — Daily saint commemorations
- `fasting_rules.json` — Fasting periods and rules
- `quotes.json` — Church Father quotes
- `epistle_cycle.json` — Weekly epistle readings
- `gospel_cycle.json` — Gospel series (John, Matthew, Luke, Lenten)
- `feast_readings.json` — Feast-specific scripture readings

## License

MIT
