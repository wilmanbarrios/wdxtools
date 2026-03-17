package number

import (
	"testing"
)

func TestAbbreviate(t *testing.T) {
	tests := []struct {
		name     string
		number   float64
		opts     []Option
		expected string
	}{
		// Basic cases from Laravel
		{"1", 1, nil, "1"},
		{"10", 10, nil, "10"},
		{"100", 100, nil, "100"},
		{"1000", 1000, nil, "1K"},
		{"1000 precision 1", 1000, []Option{WithPrecision(1)}, "1.0K"},
		{"489939", 489939, nil, "490K"},
		{"489939 maxPrecision 2", 489939, []Option{WithMaxPrecision(2)}, "489.94K"},
		{"1M", 1e6, nil, "1M"},
		{"1M precision 2", 1e6, []Option{WithPrecision(2)}, "1.00M"},
		{"1.5M", 1.5e6, nil, "2M"},
		{"1B", 1e9, nil, "1B"},
		{"1T", 1e12, nil, "1T"},
		{"1Q", 1e15, nil, "1Q"},

		// Zero
		{"zero", 0, nil, "0"},

		// Negative numbers
		{"-1000", -1000, nil, "-1K"},
		{"-1M", -1e6, nil, "-1M"},
		{"-489939", -489939, nil, "-490K"},

		// Very large: >= 1e15, recursive
		{"1KQ", 1e18, nil, "1KQ"},

		// Small numbers (no abbreviation)
		{"999", 999, nil, "999"},
		{"50", 50, nil, "50"},

		// Precision and maxPrecision combos
		{"1234 maxPrecision 2", 1234, []Option{WithMaxPrecision(2)}, "1.23K"},
		{"1200 maxPrecision 2", 1200, []Option{WithMaxPrecision(2)}, "1.2K"},
		{"1000 maxPrecision 2", 1000, []Option{WithMaxPrecision(2)}, "1K"},
		{"1000 precision 2 maxPrecision 4", 1000, []Option{WithPrecision(2), WithMaxPrecision(4)}, "1.00K"},
		{"999950 maxPrecision 1", 999950, []Option{WithMaxPrecision(1)}, "1000K"},

		// Edge: number just under 1000
		{"999.5", 999.5, nil, "1000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Abbreviate(tt.number, tt.opts...)
			if got != tt.expected {
				t.Errorf("Abbreviate(%v) = %q, want %q", tt.number, got, tt.expected)
			}
		})
	}
}

func TestForHumans(t *testing.T) {
	tests := []struct {
		name     string
		number   float64
		opts     []Option
		expected string
	}{
		{"1000", 1000, nil, "1 thousand"},
		{"1M", 1e6, nil, "1 million"},
		{"1B", 1e9, nil, "1 billion"},
		{"1T", 1e12, nil, "1 trillion"},
		{"1Q", 1e15, nil, "1 quadrillion"},
		{"1KQ", 1e18, nil, "1 thousand quadrillion"},
		{"-1000", -1000, nil, "-1 thousand"},
		{"489939 maxPrecision 2", 489939, []Option{WithMaxPrecision(2)}, "489.94 thousand"},
		{"zero", 0, nil, "0"},
		{"100", 100, nil, "100"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ForHumans(tt.number, tt.opts...)
			if got != tt.expected {
				t.Errorf("ForHumans(%v) = %q, want %q", tt.number, got, tt.expected)
			}
		})
	}
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		name         string
		number       float64
		precision    int
		maxPrecision int
		expected     string
	}{
		{"integer", 42, 0, -1, "42"},
		{"precision 2", 42, 2, -1, "42.00"},
		{"maxPrecision 2 no trail", 42, 0, 2, "42"},
		{"maxPrecision 2 with decimals", 42.123, 0, 2, "42.12"},
		{"maxPrecision 2 one decimal", 42.1, 0, 2, "42.1"},
		{"precision 2 maxPrecision 4", 42.12345, 2, 4, "42.1235"},
		{"precision 2 maxPrecision 4 short", 42.1, 2, 4, "42.10"},
		{"half round up", 2.5, 0, -1, "3"},
		{"negative half round", -2.5, 0, -1, "-3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatNumber(tt.number, tt.precision, tt.maxPrecision)
			if got != tt.expected {
				t.Errorf("formatNumber(%v, %d, %d) = %q, want %q",
					tt.number, tt.precision, tt.maxPrecision, got, tt.expected)
			}
		})
	}
}

func BenchmarkAbbreviate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Abbreviate(489939, WithMaxPrecision(2))
	}
}

func BenchmarkAbbreviateLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Abbreviate(1e18)
	}
}

func BenchmarkForHumans(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ForHumans(489939, WithMaxPrecision(2))
	}
}

func FuzzAbbreviate(f *testing.F) {
	f.Add(0.0)
	f.Add(1.0)
	f.Add(-1.0)
	f.Add(1000.0)
	f.Add(489939.0)
	f.Add(1e6)
	f.Add(1e9)
	f.Add(1e12)
	f.Add(1e15)
	f.Add(1e18)
	f.Add(-489939.0)

	f.Fuzz(func(t *testing.T, n float64) {
		// Should not panic
		result := Abbreviate(n)
		if len(result) == 0 {
			t.Errorf("Abbreviate(%v) returned empty string", n)
		}
	})
}
