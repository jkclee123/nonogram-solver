package factory

import (
	"math/big"
	"sync"

	"nonogram-solver/types"
)

// CreateGrid takes NonogramData and creates Grid containing all row and column lines
func CreateGrid(data types.NonogramData) types.Grid {
	// Use goroutines when total line count is 50 or more
	if len(data.Clues) >= 50 {
		return createGridWithGoroutines(data)
	}

	return createGridSingleThreaded(data)
}

// createGridSingleThreaded creates grid using single-threaded approach
func createGridSingleThreaded(data types.NonogramData) types.Grid {
	lines := make(map[types.LineID]types.Line)

	for lineID, clues := range data.Clues {
		var size int
		if lineID.Direction == types.Row {
			size = data.Width
		} else {
			size = data.Height
		}
		lines[lineID] = createLineFromClues(clues, size)
	}

	return types.Grid{
		Lines:    lines,
		Width:    data.Width,
		Height:   data.Height,
		ColorMap: data.ColorMap,
	}
}

// createGridWithGoroutines creates grid using goroutines for parallel processing
func createGridWithGoroutines(data types.NonogramData) types.Grid {
	lines := make(map[types.LineID]types.Line)

	var wg sync.WaitGroup
	var mu sync.Mutex

	// Process all clues in parallel goroutines
	for lineID, clues := range data.Clues {
		wg.Add(1)
		go func(lid types.LineID, cls []types.ClueItem) {
			defer wg.Done()

			var size int
			if lid.Direction == types.Row {
				size = data.Width
			} else {
				size = data.Height
			}

			line := createLineFromClues(cls, size)

			mu.Lock()
			lines[lid] = line
			mu.Unlock()
		}(lineID, clues)
	}

	wg.Wait()
	return types.Grid{
		Lines:    lines,
		Width:    data.Width,
		Height:   data.Height,
		ColorMap: data.ColorMap,
	}
}

// createLineFromClues creates a Line from clues and the line size
func createLineFromClues(clues []types.ClueItem, size int) types.Line {
	// Find all unique colors in this line
	colorSet := make(map[int]bool)
	for _, clue := range clues {
		colorSet[clue.ColorID] = true
	}

	// Initialize Combinations and Facts for each color
	combinations := make(map[int][]*big.Int)
	facts := make(map[int]*big.Int)

	for colorID := range colorSet {
		combinations = GenerateColorCombinations(clues, size)
		facts[colorID] = big.NewInt(0)
	}

	return types.Line{
		Combinations: combinations,
		Facts:        facts,
	}
}
