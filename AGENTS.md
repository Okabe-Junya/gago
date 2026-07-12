# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
make test          # Run all tests (go test ./...)
make lint          # Run golangci-lint (must be installed locally)

go test ./pkg/ga/...                          # Test a single package tree
go test -run TestName ./pkg/ga                # Run a single test
go test -bench=. ./pkg/ga                     # Run benchmarks (ga_benchmark_test.go)
go test -race ./...                           # Race detector

./scripts/run_examples.sh                     # Build and run every example under examples/
go run ./examples/find_max/find_max.go        # Run a single example
```

CI (`.github/workflows/`) runs `make test`, `golangci-lint`, and verifies `go mod tidy` is clean — match those locally before opening a PR.

## Architecture

### Two parallel API surfaces (refactor in progress)

The repo currently exposes the same concepts twice. New code generally lives in the sub-packages; the root `pkg/ga` keeps the old surface for backwards compatibility.

- **Legacy / public entry point**: `pkg/ga/` defines `GA`, `Individual`, `Genotype`, `Phenotype`, `Population`, and the selection/crossover/mutation operators all in one package. Examples and the README use this API.
- **Refactored sub-packages**:
  - `pkg/ga/encoding/` — `Genotype` + the four encodings (Binary / Integer / Real / Permutation) with a stricter typed-error API (`ErrInvalidGenomeType`, `ErrInvalidGenomePosition`, etc.) and `*Unsafe` fast-path accessors.
  - `pkg/ga/population/` — `Individual`, `Phenotype` (note: `Features` here is `map[string]interface{}`, while `pkg/ga.Phenotype.Features` is `[]float64` — they are not interchangeable), `Population`.
  - `pkg/ga/operators/` — selection operators (tournament, roulette, rank, SUS, truncation, Boltzmann, NSGA-II-style multi-objective).
  - `pkg/ga/termination/` — richer termination conditions (composite AND/OR, stagnation, diversity-based) that still implement `ga.TerminationCondition` from the root package, so they plug into the legacy `GA` struct.
- **Bridge**: `pkg/ga/genotype.go` exists only as an adapter between `ga.Genotype` and `encoding.Genotype` (`ConvertToEncodingGenotype` / `ConvertFromEncodingGenotype`). When touching encoding code, decide whether you are working in the legacy or refactored API and keep the adapters in sync.

### How an evolution run is wired together

1. The user constructs a `ga.GA` with operator funcs (`Selection`, `Crossover`, `Mutation` — all of which now take a trailing `*rand.Rand`), rates, generation count, and optionally `Seed`, `EarlyStopping`, `OnGeneration`, `TermCondition`, `ElitismCount`, `AdaptiveParams`, `EnableLogger`.
2. `GA.Initialize(size, initGenotype, evalPhenotype)` builds the initial `Population` (the `initGenotype` callback receives the GA's `*rand.Rand`), fills in defaults (`NumParallelEvals = runtime.NumCPU()`, mutation/crossover rates clamped to (0,1] with defaults 0.1/0.8), validates required operators, and seeds `History` with the initial `Statistics`. The GA's rng is seeded from `Seed` (non-zero) or `time.Now().UnixNano()` (zero).
3. `GA.Evolve(evalPhenotype)` is the main loop and returns `(*Result, error)`. Each generation: select → crossover → mutate → (optionally) snapshot elites → (optionally) `updateAdaptiveParams` from current diversity → parallel evaluate offspring → replace population → reinsert elites → recompute statistics → append to `History` → invoke `OnGeneration` → check `EarlyStopping` then `TermCondition`. Errors propagate up rather than panicking.
4. Best individual is tracked across generations (highest `Phenotype.Fitness` ever seen, not the last generation's best). The captured best is **deep-cloned** on capture so subsequent in-place mutation of the population cannot corrupt it.

### Genome representation

All encodings share a `[]byte` backing array. Real and integer values are normalized to 0–255 byte slots and decoded on access via `GetRealValue` / `GetIntegerValue`. Consequences worth knowing:

- Per-gene precision for real-valued problems is ~8 bits — fine for most demos, possibly too coarse for tight optimization.
- Crossover and mutation operators work directly on bytes, which lets the same operator (e.g. `SinglePointCrossover`) work across encodings.
- `MinValues` / `MaxValues` slices on `Genotype` are required for integer and real encodings; binary and permutation leave them nil.

### Fitness convention

Higher fitness is always better. Every selection operator, `GetBestIndividual`, statistics, convergence checks, and the parallel evaluator's panic fallback (which assigns `-math.MaxFloat64`) all assume this. Minimization problems must be reformulated before being passed to the GA.

### Parallel evaluation

`evaluatePopulationInParallel` uses a worker-pool pattern (`NumParallelEvals` workers, defaulting to `runtime.NumCPU()`) with per-worker panic recovery. A panicking fitness function won't crash the run — that individual gets `-math.MaxFloat64` fitness and an error is logged. Set `NumParallelEvals = 1` to fall back to sequential evaluation (useful for deterministic debugging).

### Logging

`internal/logger` wraps `log/slog`. The logger is opt-in via `GA.EnableLogger = true`; when disabled, `GA.Logger` is `nil` and all `ga.Logger.XYZ(...)` calls are nil-safe by design — preserve that pattern when adding new log sites.

## Conventions

- Linting is configured in `.golangci.yaml` (Go 1.23 target). Test files are excluded from `errcheck`, `noctx`, and `govet` — production code is not.
- Tests live next to the code they cover (`*_test.go`); benchmarks live in `pkg/ga/ga_benchmark_test.go`.
- Do not add a `co-author` / `Co-authored-by` trailer to commits.
- Branch naming and conventional-commit format are defined in the user's global instructions and apply here.

## Things that are easy to get wrong

- `pkg/ga.Phenotype.Features` (`[]float64`) and `pkg/ga/population.Phenotype.Features` (`map[string]interface{}`) are different types. Cloning/conversion code in `pkg/ga/ga.go` (`cloneIndividual`) and `pkg/ga/genotype.go` (the adapters) only handles the legacy `[]float64` form.
- `TermCondition` is evaluated **after** `History` is appended for the generation, so `GenerationCountTermination(n)` stops once `len(History) >= n`, which includes the initial pre-evolution snapshot.
- `EarlyStopping` is checked **before** `TermCondition` on tied generations, so when both would fire the same generation `Result.StopReason` reports the more specific `EarlyStopping` reason.
- The seeded `*rand.Rand` is single-goroutine only; parallel fitness evaluation must not touch it (fitness functions should be deterministic anyway).
- `FitnessStagnationTermination` and `FitnessImprovementTermination` close over mutable state (`bestFitness`, `prevFitness`). They are not safe to reuse across separate `GA` runs — construct a fresh condition per run.
