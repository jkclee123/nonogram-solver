# Nonogram Solver Architecture

## High-level Goal
Design a clean, concurrent project structure for a nonogram solver centered on `Line` operations: overlap and crossReference, with lazy combination generation and a worker/queue model.

## Proposed Project Structure (Go)
- `cmd/nonogram-solver/`
  - Entrypoint `main.go` (CLI wiring, flags, IO)
- `internal/solver/`
  - `solver.go` (top-level `Solve` orchestrator, convergence loop)
  - `work.go` (work item types, queue, scheduling policy)
  - `worker_pool.go` (pool lifecycle, draining, concurrency)
- `internal/line/`
  - `line.go` (`Line` model: id, direction, length, clues, facts)
  - `facts.go` (bitsets for filled-by-color, empty, helpers)
  - `overlap.go` (overlap operation)
  - `crossref.go` (crossReference operation)
- `internal/combinatorics/`
  - `combinations.go` (lazy generator interface)
  - `bitset.go` (thin wrapper over `math/big.Int`)
- `internal/grid/`
  - `grid.go` (grid model, row/col indexing, orthogonal lookup)
  - `propagation.go` (map overlap deltas to orthogonal work)
- `internal/factory/`
  - `gridFactory.go` (build `Grid` and `Line` structures from input)
  - `combinationsFactory.go` (wire `CombinationsProvider` strategy and generator implementation)
- `internal/network/`
  - `fetcher.go` (fetch puzzle data if needed)
- `internal/types/` (or keep top-level `types/` you already have)
  - `direction.go`, `lineId.go`, `clueItem.go`, `facts.go`, `grid.go`
  - `line.go` (Line model with references to combinatorics types)
- `test/`
  - Unit tests by package (`*_test.go`) and integration harness

**Note:** You already have `internal/factory`, `internal/network`, and `types/`. We can either:
- migrate `types/` into `internal/types/` to keep APIs internal,
Either is fine; prefer internalizing if you don't need external consumers.

## Core Domain Models
- `Line`
  - `id LineID`, `direction Direction`, `length int`
  - `clues []ClueItem` (colored runs as needed)
  - `facts Facts` (bitsets; see below)
  - `combinations CombinationsProvider` (lazy, per-color)
- `Facts`
  - `filledByColor map[Color]*combinatorics.Bitset` (bits 1 = must be that color)
  - `emptyMask *combinatorics.Bitset` (bits 1 = must be empty)
  - Helpers: `IsKnown(i)`, `MarkFilled(i,color)`, `MarkEmpty(i)`
- `CombinationsProvider` (in `internal/combinatorics`)
  - Cache: `generated map[Color]bool`
  - Store: `combosByColor map[Color][]*Bitset`
  - API: `Get(color) ([]*Bitset, error)`
- `Grid`
  - `rows []*Line`, `cols []*Line`
  - `At(row,col)` gives cell view; `Orthogonal(line, index)` maps to the other axis

## Bitset Conventions (using `math/big.Int`)
- One `Bitset` represents a line-length vector.
- For a single combination of a given color: bit 1 = cell is that color in that combination.
- Overlap per color: intersect across that color's combos to get "must be that color".
- Empties: compute union across all colors' combos; positions not in the union are "must be empty".
- CrossReference filtering: For a fact bitset `F` and a combination bitset `C`, drop `C` if `(F | C) != C` (conflict with known facts). Apply both for empties and filled-by-color (i.e., ensure color-specific 1s are compatible and empty 1s forbid C's 1s).

## Operations
- `Overlap(line, color)`:
  - Step 1 (empties): **Only when all color combinations have been generated** - union across all colors' combos, complement → new empties.
  - Step 2 (fills): intersect across combos for `color` → new must-fill for that color.
  - Return `FactsDelta` and changed flag.
- `CrossReference(line, color)`:
  - Apply known facts to eliminate incompatible combos for `color`.
  - If combos filtered, immediately perform Overlap for `lineID,color` to update facts.

## Work Model and Scheduling
- **Initial Work Queue Seeding**: At solver startup, prioritize lines with the most overlap potential using slack score (lineLength − sum(clues) − (numClues − 1); lower slack = higher priority). Pick top K lines where K = min(32, totalLines/4), configurable via `maxInitialSeeds`. For each selected line, generate all colors' combinations upfront, then enqueue Overlap for all its colors (enabling empties calculation). Subsequent operations use CrossReference; when CrossReference filters combinations, immediately perform Overlap locally. When queue drains without changes, seed the next K lines once. Stop after two batches or when solved/stable.
- `WorkItem`
  - `type` ∈ {Overlap, CrossReference}
  - `lineID LineID`
  - `color Color` (Overlap needs all colors processed over time; CrossRef often uses the impacted color, plus "empty" propagation)
- Queue
  - Buffered `chan WorkItem`
  - De-duplication set to avoid floods (key: `(type,lineID,color)`)
- Worker
  - Pulls a WorkItem
  - For Overlap:
    - Run overlap for the `color` (and optionally a fast empties pass against all colors if you split empties as a special `color`).
    - If facts changed → enqueue CrossReference for all orthogonal lines at impacted positions, for relevant colors.
  - For CrossReference:
    - Filter combos by facts for `color`.
    - If combos changed → immediately perform Overlap for `lineID,color` (no queueing).
- Convergence
  - Continue until queue drains and no changes.
  - Solver returns when stable or solved; add a guard (max iterations) for safety.

## Lazy Combination Generation
- `CombinationsProvider.Get(color)`:
  - If `generated[color]` return cached combos.
  - Else generate with current `facts` (optional pruning using known empties/fills to avoid generating impossible placements), cache, mark generated.
- Generation granularity
  - Multi-color nonograms: color-specific combos are slices of the full-line arrangements; generator must respect the full line's multi-color clues. If needed, have a full-line generator produce arrangements, then project to per-color bitsets and cache.

## Grid Coordination
- Overlap sets facts on `line` positions.
- For each changed position `i` on a row, schedule CrossReference for the column at index `i` (and vice versa).
- When a position becomes "empty", schedule crossRef for all colors of the orthogonal line at that index.
- When a position becomes "filled with color C", schedule crossRef for color C on the orthogonal line at that index, and optionally schedule a quick check for other colors if your model encodes exclusivity constraints via empties elsewhere.

## Public API
- `solver.Solve(ctx, data types.NonogramData, opts solver.Options) (grid.Grid, error)`
  - Options: worker count, logging, deterministic mode, early-stop, etc.

## Testing Strategy
- `internal/factory`: unit tests for combinations generator and provider
- `internal/combinatorics`: unit tests for bitset operations
- `internal/line`: tests for overlap and crossReference on synthetic lines
- `internal/solver`: end-to-end on tiny puzzles (2–5 tests)
- Fuzz tests for crossRef stability (no resurrection of eliminated combos)
- Property: after convergence, no combination remains that contradicts facts

## Incremental Implementation Order
1. Bitset wrapper
2. Facts and Line model
3. CombinationsProvider with minimal generator
4. Overlap
5. CrossReference
6. Work queue + worker pool
7. Grid + propagation wiring
8. Solver orchestration
9. Tests from small to larger