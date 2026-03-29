package number

import (
	"math"
	"strconv"
	"strings"
)

// Pre-computed powers of 10 for precision lookups (indices 0-10).
// Avoids math.Pow calls for the common case of small precision values.
var pow10 = [11]float64{1, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10}

// formatNumber formats a float64 with the given precision and maxPrecision.
//
// precision controls the minimum number of decimal places.
// maxPrecision controls the maximum number of decimal places (trailing zeros removed).
// If maxPrecision is -1, it is unset and precision is used exactly.
//
// Uses math.Round for half-away-from-zero rounding (PHP parity).
func formatNumber(number float64, precision, maxPrecision int) string {
	if maxPrecision >= 0 {
		return formatWithMaxPrecision(number, precision, maxPrecision)
	}
	return formatWithPrecision(number, precision)
}

// formatWithPrecision formats number with exactly `precision` decimal places.
func formatWithPrecision(number float64, precision int) string {
	if precision <= 0 {
		return strconv.FormatInt(int64(math.Round(number)), 10)
	}

	// Pre-round with half-away-from-zero (PHP parity)
	factor := pow10Factor(precision)
	rounded := math.Round(number*factor) / factor

	return strconv.FormatFloat(rounded, 'f', precision, 64)
}

// formatWithMaxPrecision formats number with at most `maxPrecision` decimal
// places and at least `precision` decimal places, removing trailing zeros
// between precision and maxPrecision.
func formatWithMaxPrecision(number float64, precision, maxPrecision int) string {
	// Pre-round to maxPrecision digits (half-away-from-zero)
	factor := pow10Factor(maxPrecision)
	rounded := math.Round(number*factor) / factor

	s := strconv.FormatFloat(rounded, 'f', maxPrecision, 64)

	// Trim trailing zeros, but keep at least `precision` decimals
	if dotIdx := strings.IndexByte(s, '.'); dotIdx >= 0 {
		minLen := dotIdx + 1 + precision // e.g. "1." + precision digits
		if precision <= 0 {
			minLen = dotIdx // no dot needed
		}

		// Remove trailing zeros from the right
		end := len(s)
		for end > minLen && s[end-1] == '0' {
			end--
		}
		// Remove trailing dot if no decimals left
		if end > 0 && s[end-1] == '.' {
			if precision <= 0 {
				end--
			}
		}
		s = s[:end]
	}

	return s
}

// pow10Factor returns 10^n using a lookup table for n in [0,10],
// falling back to math.Pow for larger values.
func pow10Factor(n int) float64 {
	if n >= 0 && n < len(pow10) {
		return pow10[n]
	}
	return math.Pow(10, float64(n))
}
