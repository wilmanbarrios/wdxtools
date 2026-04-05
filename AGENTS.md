# wdxtools

Everyday formatting tools for the command line. Lightning-fast Go ports of popular library functions.

## Project principles
- Performance first: use strconv over fmt, pre-computed tables, minimize allocations
- 100% feature parity with source implementations
- Zero external dependencies in `pkg/` libraries (stdlib only)
- CLI layer uses cobra for robust subcommand handling
- All builds via Docker (no Go required on host)
- Credit original authors in source headers and README

## Structure
- `cmd/wdxtools/main.go` — single binary entry point (busybox pattern for backward compat)
- `internal/cmd/` — cobra command definitions (one file per subcommand)
- `pkg/<domain>/` — reusable library per utility (stdlib only, no CLI dependencies)
- Dockerfile: multi-stage (builder, runtime)
- Makefile: `build`, `test`, `bench` targets (all via Docker)

## Adding a new utility
1. Create `pkg/<domain>/` with library + tests
2. Create `internal/cmd/<name>.go` with cobra command (register via `init()`)
3. Update README.md

## Distribution
- GoReleaser + GitHub Actions on tag push (v*)
- Homebrew tap: wilmanbarrios/homebrew-tap (with symlinks per subcommand)
- go install support: `go install github.com/wilmanbarrios/wdxtools/cmd/wdxtools@latest`

## Testing
- `make test` — runs all tests via Docker
- `make bench` — runs benchmarks via Docker
- Table-driven tests mirroring original implementations
