# wdxtools

Personal CLI utilities written in Go. Ports of useful tools found online that lack binary distributions.

## Project principles
- Performance first: use strconv over fmt, pre-computed tables, minimize allocations
- 100% feature parity with source implementations
- Zero external dependencies (stdlib only)
- All builds via Docker (no Go required on host)
- Each utility is independent: own cmd/ entry + own pkg/ library
- Credit original authors in source headers and README

## Structure
- `cmd/<name>/main.go` — CLI entry point per utility
- `pkg/<domain>/` — reusable library per utility
- Dockerfile: multi-stage (builder, tester, runtime)
- Makefile: `build`, `test`, `bench` targets (all via Docker)

## Adding a new utility
1. Create `cmd/<name>/main.go`
2. Create `pkg/<domain>/` with library + tests
3. Add build target to Dockerfile and Makefile
4. Add binary to `.goreleaser.yml`
5. Update README.md

## Distribution
- GoReleaser + GitHub Actions on tag push (v*)
- Homebrew tap: wilmanbarrios/homebrew-wdxtools
- go install support

## Testing
- `make test` — runs all tests via Docker
- `make bench` — runs benchmarks via Docker
- Table-driven tests mirroring original implementations
