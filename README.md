# wdxtools

[![CI](https://github.com/wilmanbarrios/wdxtools/actions/workflows/ci.yml/badge.svg)](https://github.com/wilmanbarrios/wdxtools/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Personal CLI utilities written in Go. Ports of useful tools found online that lack binary distributions.

## Utilities

### numcrn

Abbreviate numbers for humans. Port of [Laravel's `Number::abbreviate()`](https://github.com/laravel/framework/blob/12.x/src/Illuminate/Support/Number.php) with 100% feature parity.

```
numcrn [flags] <number>

Flags:
  -p, --precision int        minimum decimal places (default 0)
  -m, --max-precision int    maximum decimal places (default: unset)
  -l, --long                 use long form (thousand, million, ...)
```

**Examples:**

```bash
numcrn 489939              # 490K
numcrn -m 2 489939         # 489.94K
numcrn -p 2 1000000        # 1.00M
numcrn -l 1000000          # 1 million
echo -1000 | numcrn        # -1K
numcrn 1000000000000000000 # 1KQ
```

### diffh

Human-readable time differences. Port of [Carbon's `diffForHumans()`](https://github.com/briannesbitt/Carbon/blob/master/src/Carbon/Traits/Difference.php) with 100% feature parity.

```
diffh [flags] <date> [other-date]

Flags:
  -s, --short         use short form (2h, 3d, ...)
  -a, --absolute      no temporal modifier (no ago/from now)
  -p, --parts int     time units to show (default 1, max 7)
  -j, --just-now      show "just now" for zero diffs
  -1, --one-day       use yesterday/tomorrow for 1-day diffs
  -2, --two-day       use "before yesterday"/"after tomorrow"
  -S, --sequential    only show sequential non-zero units
  -n, --no-zero       show "1 second" instead of "0 seconds"
```

**Examples:**

```bash
diffh 2024-01-15                    # 1 year ago
diffh -s 2024-01-15                 # 1y ago
diffh -p 3 2024-01-15              # 1 year, 2 months, and 14 days ago
diffh -a 2024-01-15 2025-01-15     # 1 year
diffh 1774820647                    # 4 hours ago (unix timestamp)
echo 2024-01-15 | diffh            # 1 year ago
```

## Installation

### Homebrew

```bash
brew install wilmanbarrios/wdxtools/numcrn
brew install wilmanbarrios/wdxtools/diffh
```

### GitHub Releases

Download the latest binary from [Releases](https://github.com/wilmanbarrios/wdxtools/releases).

### Go install

```bash
go install github.com/wilmanbarrios/wdxtools/cmd/numcrn@latest
go install github.com/wilmanbarrios/wdxtools/cmd/diffh@latest
```

### Build from source (Docker)

```bash
make build
./bin/numcrn 489939
```

## Credits

- **numcrn**: Port of [`Number::abbreviate()`](https://github.com/laravel/framework/blob/12.x/src/Illuminate/Support/Number.php) from [Laravel Framework](https://laravel.com) by [Taylor Otwell](https://github.com/taylorotwell), licensed under the [MIT License](https://github.com/laravel/framework/blob/12.x/LICENSE.md).
- **diffh**: Port of [`diffForHumans()`](https://github.com/briannesbitt/Carbon/blob/master/src/Carbon/Traits/Difference.php) from [Carbon](https://carbon.nesbot.com) by [Brian Nesbitt](https://github.com/briannesbitt), licensed under the [MIT License](https://github.com/briannesbitt/Carbon/blob/master/LICENSE).

## License

MIT — see [LICENSE](LICENSE).
