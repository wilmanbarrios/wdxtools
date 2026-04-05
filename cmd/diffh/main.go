// diffh — human-readable time differences.
//
// Port of Carbon's CarbonInterface::diffForHumans().
// Original: https://github.com/briannesbitt/Carbon/blob/master/src/Carbon/Traits/Difference.php
// Carbon is created by Brian Nesbitt and licensed under the MIT License.
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/wilmanbarrios/wdxtools/pkg/timediff"
)

// Supported date formats, tried in order.
var dateFormats = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02",
}

// version is set at build time via -ldflags.
var version = "dev"

func main() {
	versionFlag := flag.Bool("v", false, "print version and exit")
	flag.BoolVar(versionFlag, "version", false, "print version and exit")

	short := flag.Bool("s", false, "use short form (2h, 3d, ...)")
	flag.BoolVar(short, "short", false, "use short form (2h, 3d, ...)")

	absolute := flag.Bool("a", false, "absolute mode (no ago/from now)")
	flag.BoolVar(absolute, "absolute", false, "absolute mode (no ago/from now)")

	parts := flag.Int("p", 1, "time units to show (1-7)")
	flag.IntVar(parts, "parts", 1, "time units to show (1-7)")

	justNow := flag.Bool("j", false, "show \"just now\" for zero diffs")
	flag.BoolVar(justNow, "just-now", false, "show \"just now\" for zero diffs")

	oneDay := flag.Bool("1", false, "use yesterday/tomorrow for 1-day diffs")
	flag.BoolVar(oneDay, "one-day", false, "use yesterday/tomorrow for 1-day diffs")

	twoDay := flag.Bool("2", false, "use \"before yesterday\"/\"after tomorrow\"")
	flag.BoolVar(twoDay, "two-day", false, "use \"before yesterday\"/\"after tomorrow\"")

	sequential := flag.Bool("S", false, "only show sequential non-zero units")
	flag.BoolVar(sequential, "sequential", false, "only show sequential non-zero units")

	noZero := flag.Bool("n", false, "show \"1 second\" instead of \"0 seconds\"")
	flag.BoolVar(noZero, "no-zero", false, "show \"1 second\" instead of \"0 seconds\"")

	formatTpl := flag.String("f", "", "output template for batch mode ($1..$N=fields, $d=diff)")
	flag.StringVar(formatTpl, "format", "", "output template for batch mode ($1..$N=fields, $d=diff)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: diffh [flags] <date> [other-date]\n\nHuman-readable time differences.\n\nFlags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n  diffh 2024-01-15              → 1 year ago\n  diffh -s 2024-01-15            → 1y ago\n  diffh -p 3 2024-01-15          → 1 year, 2 months, and 14 days ago\n  diffh -a 2024-01-15 2025-01-15 → 1 year\n  diffh 1774820647               → 4 hours ago (unix timestamp)\n  echo 2024-01-15 | diffh        → 1 year ago\n  printf '%%s\\n' ts1 ts2 | diffh  → one line per input (batch mode)\n\nTemplate mode (-f): $1..$N = input fields, $d = time diff\n  ... | diffh -f '$2 ($d)'       → second field + diff in parens\n  ... | diffh -f '$1\\t$d'        → first field + tab + diff\n  ... | diffh -f '$d — $1 $2'    → diff followed by two fields\n")
	}

	flag.Parse()

	if *versionFlag {
		fmt.Println("diffh " + version)
		return
	}

	// Build options.
	var opts []timediff.Option

	if *short {
		opts = append(opts, timediff.WithShort(true))
	}
	if *absolute {
		opts = append(opts, timediff.WithSyntax(timediff.SyntaxAbsolute))
	}
	if *parts != 1 {
		opts = append(opts, timediff.WithParts(*parts))
	}

	var optFlags timediff.Options
	if *justNow {
		optFlags |= timediff.JustNow
	}
	if *oneDay {
		optFlags |= timediff.OneDayWords
	}
	if *twoDay {
		optFlags |= timediff.TwoDayWords
	}
	if *sequential {
		optFlags |= timediff.SequentialPartsOnly
	}
	if *noZero {
		optFlags |= timediff.NoZeroDiff
	}
	if optFlags != 0 {
		opts = append(opts, timediff.WithOptions(optFlags))
	}

	w := bufio.NewWriterSize(os.Stdout, 32768)

	args := flag.Args()
	if len(args) > 0 {
		from, err := parseDate(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "diffh: %v\n", err)
			os.Exit(1)
		}

		// Optional second positional argument as "other" date.
		if len(args) >= 2 {
			other, err := parseDate(args[1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "diffh: %v\n", err)
				os.Exit(1)
			}
			opts = append(opts, timediff.WithOther(other))
		}

		result := timediff.DiffForHumans(from, opts...)
		io.WriteString(w, result)
		io.WriteString(w, "\n")
		w.Flush()
		return
	}

	// Check stdin.
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		f := timediff.NewFormatter(opts...)
		var buf [128]byte
		scanner := bufio.NewScanner(os.Stdin)

		if *formatTpl != "" {
			tpl := unescapeTemplate(*formatTpl)
			var tplBuf [256]byte
			for scanner.Scan() {
				var fields [16][]byte
				nf := splitFields(scanner.Bytes(), &fields)
				if nf == 0 {
					continue
				}
				from, err := parseDateBytes(fields[0])
				if err != nil {
					w.Flush()
					fmt.Fprintf(os.Stderr, "diffh: %v\n", err)
					os.Exit(1)
				}
				diff := f.AppendFormat(buf[:0], from)
				out := expandTemplate(tplBuf[:0], tpl, fields[:nf], diff)
				w.Write(out)
				w.WriteByte('\n')
			}
		} else {
			for scanner.Scan() {
				line := bytes.TrimSpace(scanner.Bytes())
				if len(line) == 0 {
					continue
				}
				from, err := parseDateBytes(line)
				if err != nil {
					w.Flush()
					fmt.Fprintf(os.Stderr, "diffh: %v\n", err)
					os.Exit(1)
				}
				b := f.AppendFormat(buf[:0], from)
				b = append(b, '\n')
				w.Write(b)
			}
		}

		if err := scanner.Err(); err != nil {
			w.Flush()
			fmt.Fprintf(os.Stderr, "diffh: reading stdin: %v\n", err)
			os.Exit(1)
		}
		w.Flush()
		return
	}

	flag.Usage()
	os.Exit(1)
}

// parseDateBytes tries to interpret b as a date/time without string allocation
// for the common case of Unix timestamps (pure integers).
func parseDateBytes(b []byte) (time.Time, error) {
	if ts, ok := parseIntBytes(b); ok {
		return time.Unix(ts, 0).UTC(), nil
	}
	return parseDateFormat(string(b))
}

// parseIntBytes parses a decimal integer from b without allocating a string.
func parseIntBytes(b []byte) (int64, bool) {
	if len(b) == 0 {
		return 0, false
	}
	neg := b[0] == '-'
	i := 0
	if neg || b[0] == '+' {
		i = 1
	}
	if i >= len(b) {
		return 0, false
	}
	var n int64
	for ; i < len(b); i++ {
		c := b[i] - '0'
		if c > 9 {
			return 0, false
		}
		n = n*10 + int64(c)
	}
	if neg {
		n = -n
	}
	return n, true
}

// parseDate tries to interpret s as a date/time. It checks for Unix timestamps
// first (pure integers), then tries common date formats. Used for CLI arguments
// where input hasn't been pre-validated.
func parseDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)

	// Try Unix timestamp (pure integer).
	if ts, err := strconv.ParseInt(s, 10, 64); err == nil {
		return time.Unix(ts, 0).UTC(), nil
	}

	return parseDateFormat(s)
}

// parseDateFormat parses a date string using a length+character discriminator
// to select the correct format on the first try, avoiding failed time.Parse
// attempts that allocate internally.
func parseDateFormat(s string) (time.Time, error) {
	n := len(s)

	// "2006-01-02" (10 chars)
	if n == 10 {
		if t, err := time.Parse("2006-01-02", s); err == nil {
			return t, nil
		}
	}

	// Longer formats share the "2006-01-02" prefix; discriminate on s[10].
	if n > 10 {
		switch s[10] {
		case 'T':
			// RFC3339 variants: "2006-01-02T15:04:05Z07:00" or "2006-01-02T15:04:05"
			if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
				return t, nil
			}
			if t, err := time.Parse(time.RFC3339, s); err == nil {
				return t, nil
			}
			if t, err := time.Parse("2006-01-02T15:04:05", s); err == nil {
				return t, nil
			}
		case ' ':
			// "2006-01-02 15:04:05"
			if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
				return t, nil
			}
		}
	}

	// Fallback: try all formats for unexpected inputs.
	for _, layout := range dateFormats {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid date: %q", s)
}

// splitFields splits b on whitespace into the caller-supplied array, returning
// the number of fields found. No heap allocation.
func splitFields(b []byte, dst *[16][]byte) int {
	n := 0
	i := 0
	for i < len(b) && n < len(dst) {
		// Skip whitespace.
		for i < len(b) && (b[i] == ' ' || b[i] == '\t') {
			i++
		}
		if i >= len(b) {
			break
		}
		// Start of field.
		start := i
		for i < len(b) && b[i] != ' ' && b[i] != '\t' {
			i++
		}
		dst[n] = b[start:i]
		n++
	}
	return n
}

// expandTemplate replaces $d with the formatted diff, and $1..$N with input
// fields. Unknown $ sequences are kept as-is.
func expandTemplate(out []byte, tpl string, fields [][]byte, diff []byte) []byte {
	for i := 0; i < len(tpl); i++ {
		if tpl[i] != '$' || i+1 >= len(tpl) {
			out = append(out, tpl[i])
			continue
		}
		if tpl[i+1] == 'd' {
			out = append(out, diff...)
			i++
			continue
		}
		j := i + 1
		for j < len(tpl) && tpl[j] >= '0' && tpl[j] <= '9' {
			j++
		}
		if j > i+1 {
			n := 0
			for k := i + 1; k < j; k++ {
				n = n*10 + int(tpl[k]-'0')
			}
			if n >= 1 && n <= len(fields) {
				out = append(out, fields[n-1]...)
			}
			i = j - 1
			continue
		}
		out = append(out, '$')
	}
	return out
}

// unescapeTemplate interprets common escape sequences in the template string.
func unescapeTemplate(s string) string {
	s = strings.ReplaceAll(s, "\\033", "\033")
	s = strings.ReplaceAll(s, "\\e", "\033")
	s = strings.ReplaceAll(s, "\\n", "\n")
	s = strings.ReplaceAll(s, "\\t", "\t")
	return s
}
