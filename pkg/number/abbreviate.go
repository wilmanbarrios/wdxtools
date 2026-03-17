// Package number provides numeric formatting utilities.
//
// The Abbreviate function is a port of Laravel Framework's Number::abbreviate().
// Original source: https://github.com/laravel/framework/blob/12.x/src/Illuminate/Support/Number.php
// Laravel is created by Taylor Otwell and licensed under the MIT License.
package number

import (
	"math"
	"strings"
)

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

// Suffix map for abbreviated form.
var shortSuffixes = map[int]string{
	3:  "K",
	6:  "M",
	9:  "B",
	12: "T",
	15: "Q",
}

// Long suffix map for human-readable form.
var longSuffixes = map[int]string{
	3:  " thousand",
	6:  " million",
	9:  " billion",
	12: " trillion",
	15: " quadrillion",
}

// Abbreviate formats a number into a short human-readable string.
// Examples: 1000 → "1K", 1500000 → "2M", 489939 with MaxPrecision(2) → "489.94K"
func Abbreviate(number float64, opts ...Option) string {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	return summarize(number, cfg.precision, cfg.maxPrecision, shortSuffixes)
}

// ForHumans formats a number into a long human-readable string.
// Examples: 1000 → "1 thousand", 1000000 → "1 million"
func ForHumans(number float64, opts ...Option) string {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	return summarize(number, cfg.precision, cfg.maxPrecision, longSuffixes)
}

// summarize is the core formatting logic, ported from Laravel's Number::summarize().
func summarize(number float64, precision, maxPrecision int, units map[int]string) string {
	if number == 0.0 {
		return formatNumber(0, precision, maxPrecision)
	}

	if number < 0 {
		var b strings.Builder
		b.WriteByte('-')
		b.WriteString(summarize(-number, precision, maxPrecision, units))
		return b.String()
	}

	// For numbers >= 1Q (1e15), recurse: divide by 1e15 and append "Q"
	if number >= 1e15 {
		var b strings.Builder
		b.WriteString(summarize(number/1e15, precision, maxPrecision, units))
		b.WriteString(units[15])
		return b.String()
	}

	exponent := int(math.Floor(math.Log10(number)))
	displayExponent := exponent - (exponent % 3)

	if displayExponent > 0 {
		value := number / math.Pow(10, float64(displayExponent))
		suffix := units[displayExponent]

		var b strings.Builder
		b.WriteString(formatNumber(value, precision, maxPrecision))
		b.WriteString(suffix)
		return b.String()
	}

	return formatNumber(number, precision, maxPrecision)
}
