package timediff

import (
	"strconv"
	"time"
)

// unitInfo holds the English strings for a single time unit.
type unitInfo struct {
	singular string // "year"
	plural   string // "years"
	short    string // "y"
}

// Unit string tables indexed by Unit constants. No map — direct array lookup.
var unitTable = [7]unitInfo{
	{"year", "years", "y"},
	{"month", "months", "mo"},
	{"week", "weeks", "w"},
	{"day", "days", "d"},
	{"hour", "hours", "h"},
	{"minute", "minutes", "m"},
	{"second", "seconds", "s"},
}

// Cascade factors: how many of unit[i+1] fit in one of unit[i].
// year→month=12, month→week=4, week→day=7, day→hour=24, hour→min=60, min→sec=60
var cascadeFactors = [6]int{12, 4, 7, 24, 60, 60}

// DiffForHumans returns a human-readable string describing the time difference
// between from and now (or a specified other time via WithOther).
func DiffForHumans(from time.Time, opts ...Option) string {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}

	// Determine comparison target.
	to := time.Now()
	if cfg.nowOverride != nil {
		to = *cfg.nowOverride
	}
	hasOther := false
	if cfg.other != nil {
		to = *cfg.other
		hasOther = true
	}

	// Auto-detect syntax.
	syntax := cfg.syntax
	if syntax == SyntaxRelativeAuto {
		if hasOther {
			syntax = SyntaxRelativeToOther
		} else {
			syntax = SyntaxRelativeToNow
		}
	}

	iv := Diff(from, to)
	return formatInterval(iv, syntax, cfg.short, cfg.parts, cfg.options, cfg.skip)
}

// formatInterval is the core formatting engine (mirrors CarbonInterval::forHumans).
// Builds the entire output string in a single stack-allocated byte buffer to
// minimize heap allocations.
func formatInterval(iv Interval, syntax Syntax, short bool, parts int, opts Options, skip []Unit) string {
	absolute := syntax == SyntaxAbsolute
	relativeToNow := syntax == SyntaxRelativeToNow

	vals := iv.values()

	// Apply skip: cascade skipped unit values to the next smaller unit.
	if len(skip) > 0 {
		skipped := makeSkipSet(skip)
		for i := 0; i < len(vals)-1; i++ {
			if skipped[i] && vals[i] != 0 {
				vals[i+1] += vals[i] * cascadeFactors[i]
				vals[i] = 0
			}
		}
		if skipped[len(vals)-1] {
			vals[len(vals)-1] = 0
		}
	}

	// Collect non-zero parts into a stack-allocated array (max 7 parts).
	type collected struct {
		count int
		unit  int
	}
	var result [7]collected
	resultLen := 0
	var lastUnit int

	for i, v := range vals {
		if v > 0 {
			result[resultLen] = collected{v, i}
			resultLen++
		} else if opts&SequentialPartsOnly != 0 && resultLen > 0 {
			break
		}
		if resultLen >= parts {
			break
		}
		lastUnit = i
	}

	// Handle zero diff.
	if resultLen == 0 {
		if relativeToNow && opts&JustNow != 0 {
			return "just now"
		}
		fallbackUnit := lastUnit
		if fallbackUnit < int(UnitSecond) {
			fallbackUnit = int(UnitSecond)
		}
		count := 0
		if opts&NoZeroDiff != 0 {
			count = 1
		}
		// Build: "{count} {unit}" + suffix in one buffer.
		var buf [64]byte
		b := appendUnit(buf[:0], count, fallbackUnit, short)
		b = appendSuffix(b, iv.Invert, absolute, relativeToNow)
		return string(b)
	}

	// Special day words (only for single-part day results).
	if parts == 1 && resultLen == 1 && result[0].unit == int(UnitDay) && !absolute {
		if relativeToNow {
			if result[0].count == 1 && opts&OneDayWords != 0 {
				if !iv.Invert {
					return "yesterday"
				}
				return "tomorrow"
			}
			if result[0].count == 2 && opts&TwoDayWords != 0 {
				if !iv.Invert {
					return "before yesterday"
				}
				return "after tomorrow"
			}
		}
	}

	// Build the entire output in a single byte buffer:
	// "{part1}, {part2}, and {partN} ago"
	var buf [128]byte
	b := buf[:0]

	for i := 0; i < resultLen; i++ {
		if i > 0 {
			if i == resultLen-1 {
				if resultLen == 2 {
					b = append(b, " and "...)
				} else {
					b = append(b, ", and "...)
				}
			} else {
				b = append(b, ", "...)
			}
		}
		b = appendUnit(b, result[i].count, result[i].unit, short)
	}

	b = appendSuffix(b, iv.Invert, absolute, relativeToNow)
	return string(b)
}

// appendUnit appends a formatted "{count}{unit}" or "{count} {unit}" to b.
// Uses strconv.AppendInt to avoid Itoa allocation.
func appendUnit(b []byte, count, unitIdx int, short bool) []byte {
	b = strconv.AppendInt(b, int64(count), 10)
	info := unitTable[unitIdx]
	if short {
		return append(b, info.short...)
	}
	b = append(b, ' ')
	if count == 1 {
		return append(b, info.singular...)
	}
	return append(b, info.plural...)
}

// appendSuffix appends the temporal modifier (" ago", " from now", etc.) to b.
func appendSuffix(b []byte, invert, absolute, relativeToNow bool) []byte {
	if absolute {
		return b
	}
	if relativeToNow {
		if !invert {
			return append(b, " ago"...)
		}
		return append(b, " from now"...)
	}
	if !invert {
		return append(b, " before"...)
	}
	return append(b, " after"...)
}

// joinParts joins formatted strings with Oxford comma style:
// 1 part: "X", 2 parts: "X and Y", 3+: "X, Y, and Z".
func joinParts(parts []string) string {
	switch len(parts) {
	case 0:
		return ""
	case 1:
		return parts[0]
	case 2:
		return parts[0] + " and " + parts[1]
	default:
		var b []byte
		for i, p := range parts {
			if i > 0 {
				if i == len(parts)-1 {
					b = append(b, ", and "...)
				} else {
					b = append(b, ", "...)
				}
			}
			b = append(b, p...)
		}
		return string(b)
	}
}

// formatUnit formats a single count+unit pair. Used by tests.
func formatUnit(count, unitIdx int, short bool) string {
	var buf [20]byte
	b := appendUnit(buf[:0], count, unitIdx, short)
	return string(b)
}

// makeSkipSet converts a skip slice to a fixed-size boolean array for O(1) lookup.
func makeSkipSet(skip []Unit) [7]bool {
	var set [7]bool
	for _, u := range skip {
		if u >= 0 && int(u) < len(set) {
			set[u] = true
		}
	}
	return set
}
