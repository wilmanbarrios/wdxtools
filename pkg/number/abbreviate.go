// Package number provides numeric formatting utilities.
//
// The Abbreviate function is a port of Laravel Framework's Number::abbreviate().
// Original source: https://github.com/laravel/framework/blob/12.x/src/Illuminate/Support/Number.php
// Laravel is created by Taylor Otwell and licensed under the MIT License.
package number

// No external imports needed — math.Log10/Pow replaced with tables and comparisons.

// Option configures the behavior of Abbreviate and ForHumans.
type Option func(*config)

type config struct {
	precision    int
	maxPrecision int
}

func defaultConfig() config {
	return config{
		precision:    0,
		maxPrecision: -1, // unset
	}
}

// WithPrecision sets the minimum number of decimal places.
func WithPrecision(p int) Option {
	return func(c *config) {
		c.precision = p
	}
}

// WithMaxPrecision sets the maximum number of decimal places.
// Trailing zeros between precision and maxPrecision are removed.
func WithMaxPrecision(p int) Option {
	return func(c *config) {
		c.maxPrecision = p
	}
}

// Pre-computed powers of 10 for common magnitudes.
var powers = [6]float64{1, 1e3, 1e6, 1e9, 1e12, 1e15}

// Suffix arrays indexed by displayExponent/3 for O(1) lookup (no map hashing).
var shortSuffixes = [6]string{"", "K", "M", "B", "T", "Q"}
var longSuffixes = [6]string{"", " thousand", " million", " billion", " trillion", " quadrillion"}

// Abbreviate formats a number into a short human-readable string.
// Examples: 1000 → "1K", 1500000 → "2M", 489939 with MaxPrecision(2) → "489.94K"
func Abbreviate(number float64, opts ...Option) string {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	if number < 0 {
		return "-" + summarize(-number, cfg.precision, cfg.maxPrecision, &shortSuffixes)
	}
	return summarize(number, cfg.precision, cfg.maxPrecision, &shortSuffixes)
}

// ForHumans formats a number into a long human-readable string.
// Examples: 1000 → "1 thousand", 1000000 → "1 million"
func ForHumans(number float64, opts ...Option) string {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	if number < 0 {
		return "-" + summarize(-number, cfg.precision, cfg.maxPrecision, &longSuffixes)
	}
	return summarize(number, cfg.precision, cfg.maxPrecision, &longSuffixes)
}

// summarize is the core formatting logic, ported from Laravel's Number::summarize().
// Uses comparison-based exponent detection (no math.Log10) and pre-computed power
// table (no math.Pow) for maximum throughput.
func summarize(number float64, precision, maxPrecision int, units *[6]string) string {
	if number == 0.0 {
		return formatNumber(0, precision, maxPrecision)
	}

	// For numbers >= 1Q (1e15), recurse once: divide by 1e15 and append quadrillion suffix.
	if number >= 1e15 {
		return summarize(number/1e15, precision, maxPrecision, units) + units[5]
	}

	// Determine display exponent via comparisons instead of math.Log10+math.Floor.
	// At most 4 float comparisons (single UCOMISD instructions), exact at IEEE 754 boundaries.
	var displayExponent int
	switch {
	case number >= 1e12:
		displayExponent = 12
	case number >= 1e9:
		displayExponent = 9
	case number >= 1e6:
		displayExponent = 6
	case number >= 1e3:
		displayExponent = 3
	default:
		return formatNumber(number, precision, maxPrecision)
	}

	value := number / powers[displayExponent/3]
	return formatNumber(value, precision, maxPrecision) + units[displayExponent/3]
}
