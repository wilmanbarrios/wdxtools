// Port of Laravel Framework's Number::abbreviate().
// Original: https://github.com/laravel/framework/blob/12.x/src/Illuminate/Support/Number.php
// Laravel is created by Taylor Otwell and licensed under the MIT License.
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wilmanbarrios/wdxtools/pkg/number"
)

var numcrnCmd = &cobra.Command{
	Use:   "numcrn [flags] <number>",
	Short: "Abbreviate numbers for humans",
	Long:  "Abbreviate numbers for humans (port of Laravel Number::abbreviate).",
	Example: `  wdxtools numcrn 489939          # 490K
  wdxtools numcrn -m 2 489939     # 489.94K
  wdxtools numcrn -l 1000000      # 1 million
  echo 1000 | wdxtools numcrn     # 1K`,
	Args:         cobra.ArbitraryArgs,
	SilenceUsage: true,
	RunE:         runNumcrn,
}

func init() {
	f := numcrnCmd.Flags()
	f.IntP("precision", "p", 0, "minimum decimal places (precision)")
	f.IntP("max-precision", "m", -1, "maximum decimal places (-1 = unset)")
	f.BoolP("long", "l", false, "use long form (thousand, million, ...)")

	rootCmd.AddCommand(numcrnCmd)
}

func runNumcrn(cmd *cobra.Command, args []string) error {
	f := cmd.Flags()
	precision, _ := f.GetInt("precision")
	maxPrecision, _ := f.GetInt("max-precision")
	long, _ := f.GetBool("long")

	var opts []number.Option
	if precision != 0 {
		opts = append(opts, number.WithPrecision(precision))
	}
	if maxPrecision >= 0 {
		opts = append(opts, number.WithMaxPrecision(maxPrecision))
	}

	fn := number.Abbreviate
	if long {
		fn = number.ForHumans
	}

	w := bufio.NewWriterSize(os.Stdout, 32768)

	if len(args) > 0 {
		for _, arg := range args {
			if err := processNumcrnInput(w, arg, fn, opts); err != nil {
				w.Flush()
				return fmt.Errorf("numcrn: %w", err)
			}
		}
		w.Flush()
		return nil
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
			if err := processNumcrnInput(w, line, fn, opts); err != nil {
				w.Flush()
				fmt.Fprintf(os.Stderr, "numcrn: %v\n", err)
				os.Exit(1)
			}
		}
		if err := scanner.Err(); err != nil {
			w.Flush()
			fmt.Fprintf(os.Stderr, "numcrn: reading stdin: %v\n", err)
			os.Exit(1)
		}
		w.Flush()
		return nil
	}

	return cmd.Help()
}

func processNumcrnInput(w io.Writer, input string, fn func(float64, ...number.Option) string, opts []number.Option) error {
	n, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return fmt.Errorf("invalid number: %q", input)
	}
	result := fn(n, opts...)
	io.WriteString(w, result)
	io.WriteString(w, "\n")
	return nil
}
