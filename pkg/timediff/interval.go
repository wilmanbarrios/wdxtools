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
	h1, _, s1 := from.Clock()
	h2, _, s2 := to.Clock()
	_, min1, _ := from.Clock()
	_, min2, _ := to.Clock()

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
		d2 += daysInMonth(y2, m2)
	}

	// Month/year: total months between the two adjusted dates.
	totalMonths := (y2-y1)*12 + int(m2) - int(m1)

	// Day remainder: d2 - d1, handling month-boundary clipping.
	days := d2 - d1
	if days < 0 {
		// We overshot: back off one month and count remaining days.
		totalMonths--

		// Find the month we land in after totalMonths from (y1, m1).
		prevM := m2 - 1
		prevY := y2
		if prevM <= 0 {
			prevM = 12
			prevY--
		}
		dim := daysInMonth(prevY, prevM)

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

// daysInMonth returns the number of days in the given month.
func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// values returns the interval components as a slice aligned with Unit constants.
func (iv Interval) values() [7]int {
	return [7]int{iv.Years, iv.Months, iv.Weeks, iv.Days, iv.Hours, iv.Minutes, iv.Seconds}
}

// IsZero reports whether all components are zero.
func (iv Interval) IsZero() bool {
	return iv.Years == 0 && iv.Months == 0 && iv.Weeks == 0 &&
		iv.Days == 0 && iv.Hours == 0 && iv.Minutes == 0 && iv.Seconds == 0
}
