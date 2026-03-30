// diffh — human-readable time differences.
//
// Port of Carbon's CarbonInterface::diffForHumans().
// Original: https://github.com/briannesbitt/Carbon/blob/master/src/Carbon/Traits/Difference.php
// Carbon is created by Brian Nesbitt and licensed under the MIT License.
package main

import (
	"bufio"
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
	time.RFC3339,
	time.RFC3339Nano,
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

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: diffh [flags] <date> [other-date]\n\nHuman-readable time differences.\n\nFlags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n  diffh 2024-01-15              → 1 year ago\n  diffh -s 2024-01-15            → 1y ago\n  diffh -p 3 2024-01-15          → 1 year, 2 months, and 14 days ago\n  diffh -a 2024-01-15 2025-01-15 → 1 year\n  diffh 1774820647               → 4 hours ago (unix timestamp)\n  echo 2024-01-15 | diffh        → 1 year ago\n")
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

	w := bufio.NewWriter(os.Stdout)

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
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			from, err := parseDate(line)
			if err != nil {
				w.Flush()
				fmt.Fprintf(os.Stderr, "diffh: %v\n", err)
				os.Exit(1)
			}
			result := timediff.DiffForHumans(from, opts...)
			io.WriteString(w, result)
			io.WriteString(w, "\n")
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

// parseDate tries to interpret s as a date/time. It checks for Unix timestamps
// first (pure integers), then tries common date formats.
func parseDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)

	// Try Unix timestamp (pure integer).
	if ts, err := strconv.ParseInt(s, 10, 64); err == nil {
		return time.Unix(ts, 0), nil
	}

	// Try known date formats.
	for _, layout := range dateFormats {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid date: %q", s)
}
