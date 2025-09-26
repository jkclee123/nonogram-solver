## Nonogram Solver Architecture

### High-level goal
Design a clean, concurrent project structure for a nonogram solver centered on `Line` operations: overlap and crossReference, with lazy combination generation and a worker/queue model.

### Proposed project structure (Go)
- `cmd/nonogram-solver/`
  - Entrypoint `main.go` (CLI wiring, flags, IO)
- `internal/solver/`
  - `solver.go` (top-level `Solve` orchestrator, convergence loop)
  - `work.go` (work item types, queue, scheduling policy)
  - `worker_pool.go` (pool lifecycle, draining, concurrency)
- `internal/line/`
  - `line.go` (`Line` operations: `Overlap`, `CrossReference`, `FactsDelta`)
  - `facts.go` (bitsets for filled-by-color, empty, helpers) [future]
- `internal/combinatorics/`
  - `provider.go` (lazy generator interface and cache) [future: implementation]
- `internal/grid/`
  - `grid.go` (grid model helpers, orthogonal lookup)
- `internal/factory/`
  - `gridFactory.go` (build `Grid` and `Line` structures from clues)
  - `combinationsFactory.go` (combination generation strategy)
- `internal/network/`
  - `fetcher.go` (fetch and parse puzzle, now returns `Grid` directly)
- `types/`
  - `direction.go`, `lineId.go`, `line.go`, `clueItem.go`, `grid.go` (shared types)
- `test/`
  - Unit tests by package (`*_test.go`) and integration harness

### Core domain models
- `Line`
  - `id LineID`, `direction Direction`, `length int`
  - `clues []ClueItem`
  - `facts map[int]*big.Int` (per-color facts; separate mask for empties if needed)
  - `combinations map[int][]*big.Int` (lazily generated per color)
- `Facts`
  - `filledByColor map[int]*Bitset` (bits 1 = must be that color)
  - `emptyMask *Bitset` (bits 1 = must be empty)
- `Grid`
  - `lines map[LineID]Line`, `width int`, `height int`, `colorMap map[int]string`
- `CombinationsProvider` (optional future abstraction)
  - Cache: `generated map[int]bool`
  - Store: `combosByColor map[int][]*big.Int`
  - API: `Get(color int) []*big.Int`

### Bitset conventions (using `math/big.Int`)
- One bitset represents a line-length vector.
- For one combination of a given color: bit 1 = cell is that color in that combination.
- Overlap per color: intersect across that color’s combos → must-fill for that color.
- Empties: union across all colors’ combos; positions NOT in the union are “must be empty”.
- CrossReference filtering: Given fact bitset `F` and combination `C`, drop `C` if `(F | C) != C` (conflict). Apply for both fills and empties.

### Operations
- `Overlap(line, color)`
  - Step 1 (empties): union across all colors’ combos; complement → new empties.
  - Step 2 (fills): intersect across combos for `color` → new must-fill for that color.
  - Return `FactsDelta` and `changed` flag.
- `CrossReference(line, color)`
  - Remove combos that contradict current facts using bitwise checks.
  - If combos shrink, return `changed` and schedule re-overlap.

### Work model and scheduling
- `WorkItem`
  - `type` ∈ {Overlap, CrossReference}
  - `lineID LineID`, `color int`
- Queue
  - Buffered `chan WorkItem`, de-duplication by `(type,lineID,color)`
- Worker
  - Overlap:
    - Run for `color` (and possibly a special empties pass across all colors).
    - If facts changed → enqueue crossReference for orthogonal lines at impacted positions.
  - CrossReference:
    - Filter combos by facts for `color`.
    - If combos changed → enqueue overlap for `(lineID,color)`.
- Convergence
  - Continue until queue drains and no changes, or max-iterations reached.

### Lazy combination generation
- `Get(color)`
  - If cached, return.
  - Else generate per line given clues and size, respecting same-color gap rules; store and mark generated.
- Generation strategy in `internal/factory/combinationsFactory.go`
  - Full-line arrangements; project per-color bitmasks; cache by (clues,size) key.

### Grid coordination
- Overlap sets facts on a `line`.
- For each changed position `i` on a row, enqueue CrossReference for column `i` (and vice versa).
- When a position becomes empty, schedule crossRef for all colors on the orthogonal line at that index.
- When a position is filled with color `C`, schedule crossRef for color `C` on the orthogonal line at that index.

### Public API
- `solver.Solve(ctx, grid types.Grid, opts solver.Options) (types.Grid, error)`
  - Options: worker count, logging, deterministic mode, early-stop.

### Testing strategy
- `internal/factory/combinationsFactory`: unit tests for generator vs known clues
- `internal/line`: tests for `Overlap` and `CrossReference` on small synthetic lines
- `internal/solver`: end-to-end on tiny puzzles
- Fuzz/property tests: after convergence, no remaining combination contradicts facts

### Incremental implementation order
1) Bitset helpers (if needed)
2) Facts and `Line` operations surface
3) CombinationsProvider or continue using factory generator
4) Implement `Overlap`
5) Implement `CrossReference`
6) Work queue + worker pool
7) Grid propagation wiring
8) Solver orchestration
9) Tests from small to larger
