package timediff

import "time"

// Interval holds the decomposed calendar difference between two times.
type Interval struct {
	Years   int
	Months  int
	Weeks   int
	Days    int // days excluding weeks
	Hours   int
	Minutes int
	Seconds int
	Invert  bool // true when "from" is after "to" (past direction)
}

// Pre-computed days per month (index 1-12). February defaults to 28;
// leap years are handled by isLeapYear. Avoids time.Date allocation.
var daysPerMonth = [13]int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

// Diff computes the calendar interval between two times.
// This replicates PHP DateInterval's behavior: each component is computed
// via calendar arithmetic (not total-seconds conversion), so "Jan 1 to Feb 1"
// is always 1 month regardless of how many days January has. Month-end
// clipping matches PHP's diff() semantics.
func Diff(from, to time.Time) Interval {
	invert := from.After(to)
	if invert {
		from, to = to, from
	}

	y1, m1, d1 := from.Date()
	y2, m2, d2 := to.Date()
	h1, min1, s1 := from.Clock()
	h2, min2, s2 := to.Clock()

	// Time components with borrowing.
	seconds := s2 - s1
	minutes := min2 - min1
	hours := h2 - h1

	if seconds < 0 {
		seconds += 60
		minutes--
	}
	if minutes < 0 {
		minutes += 60
		hours--
	}
	if hours < 0 {
		hours += 24
		d2-- // borrow a day from the date
	}

	// If d2 became 0 (or negative), roll back one month.
	if d2 <= 0 {
		m2--
		if m2 <= 0 {
			m2 = 12
			y2--
		}
		d2 += daysIn(y2, m2)
	}

	// Month/year: total months between the two adjusted dates.
	totalMonths := (y2-y1)*12 + int(m2) - int(m1)

	// Day remainder: d2 - d1, handling month-boundary clipping.
	days := d2 - d1
	if days < 0 {
		totalMonths--

		prevM := m2 - 1
		prevY := y2
		if prevM <= 0 {
			prevM = 12
			prevY--
		}
		dim := daysIn(prevY, prevM)

		// Clip from's day to the max day of that month (PHP behavior).
		clipped := d1
		if clipped > dim {
			clipped = dim
		}
		days = dim - clipped + d2
	}

	years := totalMonths / 12
	months := totalMonths % 12
	weeks := days / 7
	days = days % 7

	// Suppress unused warnings for h1, min1, s1.
	_ = h1
	_ = min1
	_ = s1

	return Interval{
		Years:   years,
		Months:  months,
		Weeks:   weeks,
		Days:    days,
		Hours:   hours,
		Minutes: minutes,
		Seconds: seconds,
		Invert:  invert,
	}
}

// daysIn returns the number of days in the given month using a pre-computed
// table. No time.Date allocation needed.
func daysIn(year int, month time.Month) int {
	if month == 2 && isLeapYear(year) {
		return 29
	}
	return daysPerMonth[month]
}

// isLeapYear reports whether year is a leap year.
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// values returns the interval components as an array aligned with Unit constants.
func (iv Interval) values() [7]int {
	return [7]int{iv.Years, iv.Months, iv.Weeks, iv.Days, iv.Hours, iv.Minutes, iv.Seconds}
}

// IsZero reports whether all components are zero.
func (iv Interval) IsZero() bool {
	return iv.Years == 0 && iv.Months == 0 && iv.Weeks == 0 &&
		iv.Days == 0 && iv.Hours == 0 && iv.Minutes == 0 && iv.Seconds == 0
}
