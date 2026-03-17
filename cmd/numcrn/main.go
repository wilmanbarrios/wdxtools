// numcrn — abbreviate numbers for humans.
//
// Port of Laravel Framework's Number::abbreviate().
// Original: https://github.com/laravel/framework/blob/12.x/src/Illuminate/Support/Number.php
// Laravel is created by Taylor Otwell and licensed under the MIT License.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/wilmanbarrios/wdxtools/pkg/number"
)

func main() {
	precision := flag.Int("p", 0, "minimum decimal places (precision)")
	flag.IntVar(precision, "precision", 0, "minimum decimal places (precision)")

	maxPrecision := flag.Int("m", -1, "maximum decimal places (max-precision, -1 = unset)")
	flag.IntVar(maxPrecision, "max-precision", -1, "maximum decimal places (max-precision, -1 = unset)")

	long := flag.Bool("l", false, "use long form (thousand, million, ...)")
	flag.BoolVar(long, "long", false, "use long form (thousand, million, ...)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: numcrn [flags] <number>\n\nAbbreviate numbers for humans.\n\nFlags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n  numcrn 489939          → 490K\n  numcrn -m 2 489939     → 489.94K\n  numcrn -l 1000000      → 1 million\n  echo 1000 | numcrn     → 1K\n")
	}

	flag.Parse()

	var opts []number.Option
	if *precision != 0 {
		opts = append(opts, number.WithPrecision(*precision))
	}
	if *maxPrecision >= 0 {
		opts = append(opts, number.WithMaxPrecision(*maxPrecision))
	}

	fn := number.Abbreviate
	if *long {
		fn = number.ForHumans
	}

	args := flag.Args()
	if len(args) > 0 {
		// Process arguments
		for _, arg := range args {
			if err := processInput(arg, fn, opts); err != nil {
				fmt.Fprintf(os.Stderr, "numcrn: %v\n", err)
				os.Exit(1)
			}
		}
		return
	}

	// Check stdin
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			if err := processInput(line, fn, opts); err != nil {
				fmt.Fprintf(os.Stderr, "numcrn: %v\n", err)
				os.Exit(1)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "numcrn: reading stdin: %v\n", err)
			os.Exit(1)
		}
		return
	}

	flag.Usage()
	os.Exit(1)
}

func processInput(input string, fn func(float64, ...number.Option) string, opts []number.Option) error {
	n, err := strconv.ParseFloat(strings.TrimSpace(input), 64)
	if err != nil {
		return fmt.Errorf("invalid number: %q", input)
	}
	fmt.Println(fn(n, opts...))
	return nil
}
