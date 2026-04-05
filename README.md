# wdxtools

[![CI](https://github.com/wilmanbarrios/wdxtools/actions/workflows/ci.yml/badge.svg)](https://github.com/wilmanbarrios/wdxtools/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Everyday formatting tools for the command line. Lightning-fast Go ports of popular library functions.

## Install

```bash
brew install wilmanbarrios/wdxtools/wdxtools
# or
go install github.com/wilmanbarrios/wdxtools/cmd/wdxtools@latest
```

## Utilities

| Command | Description | Ported from |
|---------|-------------|-------------|
| [`wdxtools numcrn`](#numcrn) | Abbreviate numbers for humans | Laravel `Number::abbreviate()` |
| [`wdxtools diffh`](#diffh) | Human-readable time differences | Carbon `diffForHumans()` |

---

## numcrn

Abbreviate numbers for humans. Port of [Laravel's `Number::abbreviate()`](https://github.com/laravel/framework/blob/12.x/src/Illuminate/Support/Number.php) with 100% feature parity.

### Examples

```bash
wdxtools numcrn 489939              # 490K
wdxtools numcrn -m 2 489939         # 489.94K
wdxtools numcrn -p 2 1000000        # 1.00M
wdxtools numcrn -l 1000000          # 1 million
echo -1000 | wdxtools numcrn        # -1K
wdxtools numcrn 1000000000000000000 # 1KQ
```

Run `wdxtools numcrn --help` for all flags.

### Original vs wdxtools

| | Laravel (PHP 8.3) | wdxtools (Go) | Speedup |
|---|---|---|---|
| `Abbreviate(489939)` | 13,603 ns/op | 201 ns/op | **68x** |
| `Abbreviate(1e18)` | 9,616 ns/op | 59 ns/op | **163x** |
| `ForHumans(489939)` | 13,021 ns/op | 203 ns/op | **64x** |
| `Abbreviate(42)` | 9,664 ns/op | 17 ns/op | **568x** |

| | Original (PHP) | wdxtools (Go) |
|---|---|---|
| Runtime | PHP 8.x + Laravel | Static binary |
| Install | `composer require laravel/framework` | `brew install` / single binary |
| Dependencies | Laravel + intl extension | None (stdlib only) |
| Pipe support | No | Yes (`echo 1000 \| wdxtools numcrn`) |

---

## diffh

Human-readable time differences. Port of [Carbon's `diffForHumans()`](https://github.com/briannesbitt/Carbon/blob/master/src/Carbon/Traits/Difference.php) with 100% feature parity.

### Examples

```bash
wdxtools diffh 2024-01-15                    # 1 year ago
wdxtools diffh -s 2024-01-15                 # 1y ago
wdxtools diffh -p 3 2024-01-15              # 1 year, 2 months, and 14 days ago
wdxtools diffh -a 2024-01-15 2025-01-15     # 1 year
wdxtools diffh 1774820647                    # 4 hours ago (unix timestamp)
echo 2024-01-15 | wdxtools diffh            # 1 year ago

# Batch mode with templates (-f)
printf '%s\n' ts1 ts2 | wdxtools diffh               # one line per input
cat log.csv | wdxtools diffh -f '$2 ($d)'            # $1..$N=fields, $d=diff
cat events.tsv | wdxtools diffh -f '$1\t$2\t$d'      # tab-separated output
```

Run `wdxtools diffh --help` for all flags.

### Original vs wdxtools

| | Carbon (PHP 8.3) | wdxtools (Go) | Speedup |
|---|---|---|---|
| `diffForHumans()` | 27,685 ns/op | 167 ns/op | **166x** |
| `diffForHumans(parts: 3)` | 37,653 ns/op | 199 ns/op | **189x** |
| `diffForHumans(short)` | 27,798 ns/op | 168 ns/op | **165x** |

| | Original (PHP) | wdxtools (Go) |
|---|---|---|
| Runtime | PHP 8.x + Carbon | Static binary |
| Install | `composer require nesbot/carbon` | `brew install` / single binary |
| Dependencies | Carbon + intl extension | None (stdlib only) |
| Unix timestamps | No | Yes (`wdxtools diffh 1774820647`) |
| Pipe support | No | Yes (`echo date \| wdxtools diffh`) |

---

## Backward compatibility

Homebrew installs symlinks for each subcommand, so `diffh 2024-01-15` and `numcrn 1000` continue to work as standalone commands.

---

## Benchmark methodology

All benchmarks run inside Docker containers on the same host to ensure a fair comparison. No emulation — both containers run natively on arm64.

| | PHP container | Go container |
|---|---|---|
| Base image | `php:8.3-cli-alpine` | `golang:1.24-alpine` |
| Runtime | PHP 8.3.30 (NTS) | Go 1.24.13 |
| Optimization | OPcache + JIT (tracing mode) | `CGO_ENABLED=0 -ldflags="-s -w"` |
| Architecture | linux/arm64 (native) | linux/arm64 (native) |

**Host:** Apple M2, 8 cores, 16 GB RAM, Docker Desktop 29.2.0 (8 CPUs / 8 GB allocated).

PHP runs with its best performance mode: OPcache enabled and JIT in tracing mode. Go binaries are statically compiled with no CGo. Both containers have access to all 8 cores and the full memory allocation.

Reproduce locally: `make bench` (Go) and `make bench-php` (PHP).

---

## Build from source

All builds run via Docker — no Go toolchain required on host.

```bash
make build        # builds binary to ./bin/wdxtools
make test         # runs all tests
make bench        # runs benchmarks
```

## Credits

| Tool | Original | Author | License |
|------|----------|--------|---------|
| `numcrn` | [`Number::abbreviate()`](https://github.com/laravel/framework/blob/12.x/src/Illuminate/Support/Number.php) | [Taylor Otwell](https://github.com/taylorotwell) (Laravel) | MIT |
| `diffh` | [`diffForHumans()`](https://github.com/briannesbitt/Carbon/blob/master/src/Carbon/Traits/Difference.php) | [Brian Nesbitt](https://github.com/briannesbitt) (Carbon) | MIT |

## License

MIT — see [LICENSE](LICENSE).
