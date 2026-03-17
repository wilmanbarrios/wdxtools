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

## Installation

### Homebrew

```bash
brew install wilmanbarrios/wdxtools/numcrn
```

### GitHub Releases

Download the latest binary from [Releases](https://github.com/wilmanbarrios/wdxtools/releases).

### Go install

```bash
go install github.com/wilmanbarrios/wdxtools/cmd/numcrn@latest
```

### Build from source (Docker)

```bash
make build
./bin/numcrn 489939
```

## Credits

- **numcrn**: Port of [`Number::abbreviate()`](https://github.com/laravel/framework/blob/12.x/src/Illuminate/Support/Number.php) from [Laravel Framework](https://laravel.com) by [Taylor Otwell](https://github.com/taylorotwell), licensed under the [MIT License](https://github.com/laravel/framework/blob/12.x/LICENSE.md).

## License

MIT — see [LICENSE](LICENSE).
