package solver

import (
	"context"

	"nonogram-solver/types"
)

// Options configures Solve behavior
type Options struct {
	WorkerCount int
}

// Solve orchestrates the solving process until convergence.
func Solve(ctx context.Context, grid types.Grid, opts Options) (types.Grid, error) {
	// Stub: return input grid unchanged for now
	return grid, nil
}
