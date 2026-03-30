// Package timediff provides human-readable time difference formatting.
//
// The DiffForHumans function is a port of Carbon's CarbonInterface::diffForHumans().
// Original source: https://github.com/briannesbitt/Carbon/blob/master/src/Carbon/Traits/Difference.php
// Carbon is created by Brian Nesbitt and licensed under the MIT License.
package timediff

import "time"

// Syntax controls how the output is framed (absolute vs relative, and to what).
type Syntax int

const (
	// SyntaxRelativeAuto auto-detects: if no "other" time is given, uses
	// SyntaxRelativeToNow; otherwise uses SyntaxRelativeToOther.
	SyntaxRelativeAuto Syntax = 0
	// SyntaxAbsolute outputs bare duration with no temporal modifier.
	SyntaxAbsolute Syntax = 1
	// SyntaxRelativeToNow appends "ago" or "from now".
	SyntaxRelativeToNow Syntax = 2
	// SyntaxRelativeToOther appends "before" or "after".
	SyntaxRelativeToOther Syntax = 3
)

// Options is a bitmask controlling formatting details.
type Options int

const (
	// NoZeroDiff shows "1 second" instead of "0 seconds" when the diff is zero.
	NoZeroDiff Options = 1 << iota
	// JustNow shows "just now" for zero diffs in relative-to-now mode.
	JustNow
	// OneDayWords uses "yesterday"/"tomorrow" for 1-day diffs.
	OneDayWords
	// TwoDayWords uses "before yesterday"/"after tomorrow" for 2-day diffs.
	TwoDayWords
	// SequentialPartsOnly stops collecting parts at the first zero-value gap.
	SequentialPartsOnly
)

// Unit identifies a time unit, used with WithSkip to exclude units from output.
type Unit int

const (
	UnitYear Unit = iota
	UnitMonth
	UnitWeek
	UnitDay
	UnitHour
	UnitMinute
	UnitSecond
)

// Option configures the behavior of DiffForHumans.
type Option func(*config)

type config struct {
	other       *time.Time
	nowOverride *time.Time // overrides time.Now() without affecting syntax auto-detection
	syntax      Syntax
	short       bool
	parts       int
	options     Options
	skip        []Unit
}

func defaultConfig() config {
	return config{
		syntax: SyntaxRelativeAuto,
		parts:  1,
	}
}

// WithOther sets the comparison time. Without this, DiffForHumans compares to time.Now().
func WithOther(t time.Time) Option {
	return func(c *config) {
		c.other = &t
	}
}

// WithSyntax sets the output syntax mode.
func WithSyntax(s Syntax) Option {
	return func(c *config) {
		c.syntax = s
	}
}

// WithShort enables abbreviated unit names (e.g. "2h" instead of "2 hours").
func WithShort(short bool) Option {
	return func(c *config) {
		c.short = short
	}
}

// WithParts sets the maximum number of time units to display (1-7).
func WithParts(n int) Option {
	return func(c *config) {
		if n < 1 {
			n = 1
		}
		if n > 7 {
			n = 7
		}
		c.parts = n
	}
}

// WithOptions sets formatting option flags (bitmask).
func WithOptions(o Options) Option {
	return func(c *config) {
		c.options = o
	}
}

// WithSkip excludes the given units from the output, cascading their values
// to the next smaller unit.
func WithSkip(units ...Unit) Option {
	return func(c *config) {
		c.skip = units
	}
}

// withNowOverride overrides time.Now() for deterministic testing without
// affecting syntax auto-detection (unlike WithOther which changes the mode).
func withNowOverride(t time.Time) Option {
	return func(c *config) {
		c.nowOverride = &t
	}
}
