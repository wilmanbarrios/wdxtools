package timediff

import (
	"testing"
	"time"
)

// fixed returns a deterministic time in UTC for testing.
func fixed(year int, month time.Month, day, hour, min, sec int) time.Time {
	return time.Date(year, month, day, hour, min, sec, 0, time.UTC)
}

func TestDiffInterval(t *testing.T) {
	tests := []struct {
		name string
		from time.Time
		to   time.Time
		want Interval
	}{
		{
			"same time",
			fixed(2024, 1, 15, 10, 30, 0),
			fixed(2024, 1, 15, 10, 30, 0),
			Interval{},
		},
		{
			"1 second",
			fixed(2024, 1, 1, 0, 0, 0),
			fixed(2024, 1, 1, 0, 0, 1),
			Interval{Seconds: 1},
		},
		{
			"1 minute",
			fixed(2024, 1, 1, 0, 0, 0),
			fixed(2024, 1, 1, 0, 1, 0),
			Interval{Minutes: 1},
		},
		{
			"1 hour",
			fixed(2024, 1, 1, 0, 0, 0),
			fixed(2024, 1, 1, 1, 0, 0),
			Interval{Hours: 1},
		},
		{
			"1 day",
			fixed(2024, 1, 1, 0, 0, 0),
			fixed(2024, 1, 2, 0, 0, 0),
			Interval{Days: 1},
		},
		{
			"1 week",
			fixed(2024, 1, 1, 0, 0, 0),
			fixed(2024, 1, 8, 0, 0, 0),
			Interval{Weeks: 1},
		},
		{
			"1 month",
			fixed(2024, 1, 1, 0, 0, 0),
			fixed(2024, 2, 1, 0, 0, 0),
			Interval{Months: 1},
		},
		{
			"1 year",
			fixed(2024, 1, 1, 0, 0, 0),
			fixed(2025, 1, 1, 0, 0, 0),
			Interval{Years: 1},
		},
		{
			"complex interval",
			fixed(2024, 1, 15, 10, 30, 0),
			fixed(2025, 3, 20, 14, 45, 30),
			Interval{Years: 1, Months: 2, Days: 5, Hours: 4, Minutes: 15, Seconds: 30},
		},
		{
			"inverted (from after to)",
			fixed(2025, 6, 1, 0, 0, 0),
			fixed(2024, 1, 1, 0, 0, 0),
			Interval{Years: 1, Months: 5, Invert: true},
		},
		{
			"leap year feb to march",
			fixed(2024, 2, 28, 0, 0, 0),
			fixed(2024, 3, 1, 0, 0, 0),
			Interval{Days: 2},
		},
		{
			"non-leap year feb to march",
			fixed(2023, 2, 28, 0, 0, 0),
			fixed(2023, 3, 1, 0, 0, 0),
			Interval{Days: 1},
		},
		{
			"month boundary Jan31 to Mar1",
			fixed(2024, 1, 31, 0, 0, 0),
			fixed(2024, 3, 1, 0, 0, 0),
			Interval{Months: 1, Days: 1},
		},
		{
			"2 weeks 3 days",
			fixed(2024, 1, 1, 0, 0, 0),
			fixed(2024, 1, 18, 0, 0, 0),
			Interval{Weeks: 2, Days: 3},
		},
		{
			"seconds borrow",
			fixed(2024, 1, 1, 0, 0, 45),
			fixed(2024, 1, 1, 0, 1, 15),
			Interval{Seconds: 30},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Diff(tt.from, tt.to)
			if got != tt.want {
				t.Errorf("Diff() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestDiffForHumans(t *testing.T) {
	now := fixed(2025, 3, 29, 12, 0, 0)

	tests := []struct {
		name     string
		from     time.Time
		opts     []Option
		expected string
	}{
		// Relative to now (default) — past
		{"1 second ago", fixed(2025, 3, 29, 11, 59, 59), nil, "1 second ago"},
		{"30 seconds ago", fixed(2025, 3, 29, 11, 59, 30), nil, "30 seconds ago"},
		{"1 minute ago", fixed(2025, 3, 29, 11, 59, 0), nil, "1 minute ago"},
		{"5 minutes ago", fixed(2025, 3, 29, 11, 55, 0), nil, "5 minutes ago"},
		{"1 hour ago", fixed(2025, 3, 29, 11, 0, 0), nil, "1 hour ago"},
		{"2 hours ago", fixed(2025, 3, 29, 10, 0, 0), nil, "2 hours ago"},
		{"1 day ago", fixed(2025, 3, 28, 12, 0, 0), nil, "1 day ago"},
		{"1 week ago", fixed(2025, 3, 22, 12, 0, 0), nil, "1 week ago"},
		{"1 month ago", fixed(2025, 2, 28, 12, 0, 0), nil, "1 month ago"},
		{"1 year ago", fixed(2024, 3, 29, 12, 0, 0), nil, "1 year ago"},

		// Relative to now — future
		{"1 second from now", fixed(2025, 3, 29, 12, 0, 1), nil, "1 second from now"},
		{"1 hour from now", fixed(2025, 3, 29, 13, 0, 0), nil, "1 hour from now"},
		{"1 day from now", fixed(2025, 3, 30, 12, 0, 0), nil, "1 day from now"},
		{"1 year from now", fixed(2026, 3, 29, 12, 0, 0), nil, "1 year from now"},

		// Absolute mode
		{"absolute 2 hours", fixed(2025, 3, 29, 10, 0, 0),
			[]Option{WithSyntax(SyntaxAbsolute)}, "2 hours"},
		{"absolute 1 year", fixed(2024, 3, 29, 12, 0, 0),
			[]Option{WithSyntax(SyntaxAbsolute)}, "1 year"},

		// Relative to other
		{"before other", fixed(2024, 6, 1, 0, 0, 0),
			[]Option{WithOther(fixed(2025, 6, 1, 0, 0, 0))}, "1 year before"},
		{"after other", fixed(2025, 6, 1, 0, 0, 0),
			[]Option{WithOther(fixed(2024, 6, 1, 0, 0, 0))}, "1 year after"},

		// Short mode
		{"short 2h ago", fixed(2025, 3, 29, 10, 0, 0),
			[]Option{WithShort(true)}, "2h ago"},
		{"short 1y ago", fixed(2024, 3, 29, 12, 0, 0),
			[]Option{WithShort(true)}, "1y ago"},
		{"short absolute", fixed(2025, 3, 29, 10, 0, 0),
			[]Option{WithShort(true), WithSyntax(SyntaxAbsolute)}, "2h"},

		// Multiple parts
		{"2 parts", fixed(2025, 3, 29, 9, 30, 0),
			[]Option{WithParts(2)}, "2 hours and 30 minutes ago"},
		{"3 parts", fixed(2024, 1, 15, 10, 30, 0),
			[]Option{WithParts(3)}, "1 year, 2 months, and 2 weeks ago"},
		{"3 parts short", fixed(2024, 1, 15, 10, 30, 0),
			[]Option{WithParts(3), WithShort(true)}, "1y, 2mo, and 2w ago"},

		// Zero diff
		{"zero diff default", now, nil, "0 seconds ago"},
		{"zero diff no-zero", now,
			[]Option{WithOptions(NoZeroDiff)}, "1 second ago"},
		{"zero diff just-now", now,
			[]Option{WithOptions(JustNow)}, "just now"},
		{"zero diff just-now absolute", now,
			[]Option{WithOptions(JustNow), WithSyntax(SyntaxAbsolute)}, "0 seconds"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := append([]Option{withNowOverride(now)}, tt.opts...)
			got := DiffForHumans(tt.from, opts...)
			if got != tt.expected {
				t.Errorf("DiffForHumans() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestDiffForHumansOptions(t *testing.T) {
	now := fixed(2025, 3, 29, 12, 0, 0)

	tests := []struct {
		name     string
		from     time.Time
		opts     []Option
		expected string
	}{
		// OneDayWords
		{"yesterday", fixed(2025, 3, 28, 12, 0, 0),
			[]Option{WithOptions(OneDayWords)}, "yesterday"},
		{"tomorrow", fixed(2025, 3, 30, 12, 0, 0),
			[]Option{WithOptions(OneDayWords)}, "tomorrow"},
		{"one day words only for 1 day", fixed(2025, 3, 27, 12, 0, 0),
			[]Option{WithOptions(OneDayWords)}, "2 days ago"},

		// TwoDayWords
		{"before yesterday", fixed(2025, 3, 27, 12, 0, 0),
			[]Option{WithOptions(TwoDayWords)}, "before yesterday"},
		{"after tomorrow", fixed(2025, 3, 31, 12, 0, 0),
			[]Option{WithOptions(TwoDayWords)}, "after tomorrow"},
		{"two day words only for 2 days", fixed(2025, 3, 26, 12, 0, 0),
			[]Option{WithOptions(TwoDayWords)}, "3 days ago"},

		// Combined one + two day words
		{"combined yesterday", fixed(2025, 3, 28, 12, 0, 0),
			[]Option{WithOptions(OneDayWords | TwoDayWords)}, "yesterday"},
		{"combined before yesterday", fixed(2025, 3, 27, 12, 0, 0),
			[]Option{WithOptions(OneDayWords | TwoDayWords)}, "before yesterday"},

		// SequentialPartsOnly
		{"sequential stops at gap", fixed(2024, 3, 29, 10, 0, 0),
			[]Option{WithParts(7), WithOptions(SequentialPartsOnly)}, "1 year ago"},
		{"sequential with contiguous", fixed(2025, 3, 29, 9, 30, 15),
			[]Option{WithParts(7), WithOptions(SequentialPartsOnly)}, "2 hours, 29 minutes, and 45 seconds ago"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := append([]Option{withNowOverride(now)}, tt.opts...)
			got := DiffForHumans(tt.from, opts...)
			if got != tt.expected {
				t.Errorf("DiffForHumans() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestDiffForHumansSkip(t *testing.T) {
	now := fixed(2025, 3, 29, 12, 0, 0)

	tests := []struct {
		name     string
		from     time.Time
		opts     []Option
		expected string
	}{
		{"skip weeks", fixed(2025, 3, 12, 12, 0, 0),
			[]Option{WithParts(2), WithSkip(UnitWeek)},
			"17 days ago"},
		{"skip hours cascade to minutes", fixed(2025, 3, 29, 9, 30, 0),
			[]Option{WithParts(2), WithSkip(UnitHour)},
			"150 minutes ago"},
		{"skip year cascade to months", fixed(2023, 3, 29, 12, 0, 0),
			[]Option{WithParts(2), WithSkip(UnitYear)},
			"24 months ago"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := append([]Option{withNowOverride(now)}, tt.opts...)
			got := DiffForHumans(tt.from, opts...)
			if got != tt.expected {
				t.Errorf("DiffForHumans() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestAppendUnit(t *testing.T) {
	tests := []struct {
		name    string
		count   int
		unitIdx int
		short   bool
		expect  string
	}{
		{"1 year", 1, int(UnitYear), false, "1 year"},
		{"2 years", 2, int(UnitYear), false, "2 years"},
		{"0 seconds", 0, int(UnitSecond), false, "0 seconds"},
		{"1 second", 1, int(UnitSecond), false, "1 second"},
		{"5h short", 5, int(UnitHour), true, "5h"},
		{"1mo short", 1, int(UnitMonth), true, "1mo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf [20]byte
			got := string(appendUnit(buf[:0], tt.count, tt.unitIdx, tt.short))
			if got != tt.expect {
				t.Errorf("appendUnit() = %q, want %q", got, tt.expect)
			}
		})
	}
}

// Formatter tests

func TestFormatter(t *testing.T) {
	now := fixed(2025, 3, 29, 12, 0, 0)

	tests := []struct {
		name     string
		from     time.Time
		opts     []Option
		expected string
	}{
		{"1 hour ago", fixed(2025, 3, 29, 11, 0, 0), nil, "1 hour ago"},
		{"1 day ago", fixed(2025, 3, 28, 12, 0, 0), nil, "1 day ago"},
		{"just now", now, []Option{WithOptions(JustNow)}, "just now"},
		{"short", fixed(2025, 3, 29, 10, 0, 0), []Option{WithShort(true)}, "2h ago"},
		{"multi-part", fixed(2025, 3, 29, 9, 30, 0), []Option{WithParts(2)}, "2 hours and 30 minutes ago"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := append([]Option{withNowOverride(now)}, tt.opts...)
			f := NewFormatter(opts...)
			got := f.Format(tt.from)
			if got != tt.expected {
				t.Errorf("Formatter.Format() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestFormatterAppendFormat(t *testing.T) {
	now := fixed(2025, 3, 29, 12, 0, 0)
	f := NewFormatter(withNowOverride(now))

	var buf [128]byte
	b := f.AppendFormat(buf[:0], fixed(2025, 3, 29, 11, 0, 0))

	got := string(b)
	if got != "1 hour ago" {
		t.Errorf("AppendFormat() = %q, want %q", got, "1 hour ago")
	}

	// Verify it actually appends to existing content.
	prefix := []byte("prefix: ")
	b = f.AppendFormat(prefix, fixed(2025, 3, 28, 12, 0, 0))
	got = string(b)
	if got != "prefix: 1 day ago" {
		t.Errorf("AppendFormat(prefix) = %q, want %q", got, "prefix: 1 day ago")
	}
}

// Benchmarks

func BenchmarkDiffForHumans(b *testing.B) {
	from := fixed(2024, 6, 15, 10, 30, 0)
	now := fixed(2025, 3, 29, 12, 0, 0)
	for i := 0; i < b.N; i++ {
		DiffForHumans(from, withNowOverride(now))
	}
}

func BenchmarkDiffForHumansMultiPart(b *testing.B) {
	from := fixed(2024, 1, 15, 10, 30, 0)
	now := fixed(2025, 3, 29, 12, 0, 0)
	for i := 0; i < b.N; i++ {
		DiffForHumans(from, withNowOverride(now), WithParts(3))
	}
}

func BenchmarkDiffForHumansShort(b *testing.B) {
	from := fixed(2024, 6, 15, 10, 30, 0)
	now := fixed(2025, 3, 29, 12, 0, 0)
	for i := 0; i < b.N; i++ {
		DiffForHumans(from, withNowOverride(now), WithShort(true))
	}
}

func BenchmarkDiffInterval(b *testing.B) {
	from := fixed(2024, 1, 15, 10, 30, 0)
	to := fixed(2025, 3, 29, 12, 0, 0)
	for i := 0; i < b.N; i++ {
		Diff(from, to)
	}
}

func BenchmarkFormatter(b *testing.B) {
	from := fixed(2024, 6, 15, 10, 30, 0)
	now := fixed(2025, 3, 29, 12, 0, 0)
	f := NewFormatter(withNowOverride(now))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Format(from)
	}
}

func BenchmarkFormatterAppend(b *testing.B) {
	from := fixed(2024, 6, 15, 10, 30, 0)
	now := fixed(2025, 3, 29, 12, 0, 0)
	f := NewFormatter(withNowOverride(now))
	var buf [128]byte
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.AppendFormat(buf[:0], from)
	}
}

func BenchmarkDiffForHumansWith(b *testing.B) {
	from := fixed(2024, 6, 15, 10, 30, 0)
	now := fixed(2025, 3, 29, 12, 0, 0)
	for i := 0; i < b.N; i++ {
		DiffForHumansWith(from, now, SyntaxRelativeToNow, false, 1, 0, nil)
	}
}

func BenchmarkDiffForHumansWithMultiPart(b *testing.B) {
	from := fixed(2024, 1, 15, 10, 30, 0)
	now := fixed(2025, 3, 29, 12, 0, 0)
	for i := 0; i < b.N; i++ {
		DiffForHumansWith(from, now, SyntaxRelativeToNow, false, 3, 0, nil)
	}
}

func BenchmarkDiffIntervalUTC(b *testing.B) {
	from := fixed(2024, 1, 15, 10, 30, 0).UTC()
	to := fixed(2025, 3, 29, 12, 0, 0).UTC()
	for i := 0; i < b.N; i++ {
		Diff(from, to)
	}
}

// Fuzz test

func FuzzDiffForHumans(f *testing.F) {
	f.Add(int64(0))
	f.Add(int64(1704067200))  // 2024-01-01
	f.Add(int64(1774820647))  // tmux-style timestamp
	f.Add(int64(-86400))      // negative
	f.Add(int64(2000000000))  // far future
	f.Add(int64(946684800))   // 2000-01-01

	now := fixed(2025, 3, 29, 12, 0, 0)

	f.Fuzz(func(t *testing.T, ts int64) {
		from := time.Unix(ts, 0)
		result := DiffForHumans(from, withNowOverride(now))
		if len(result) == 0 {
			t.Errorf("DiffForHumans(Unix(%d)) returned empty string", ts)
		}
	})
}
