package factory

import (
	"math/big"
	"sync"

	"nonogram-solver/types"
)

// CreateGridFromClues builds a Grid directly from clues, dimensions, and color map.
// It mirrors CreateGrid but avoids the NonogramData wrapper.
func CreateGridFromClues(clues map[types.LineID][]types.ClueItem, width, height int, colorMap map[int]string) types.Grid {
	if len(clues) >= 50 {
		return createGridFromCluesWithGoroutines(clues, width, height, colorMap)
	}
	return createGridFromCluesSingleThreaded(clues, width, height, colorMap)
}

func createGridFromCluesSingleThreaded(clues map[types.LineID][]types.ClueItem, width, height int, colorMap map[int]string) types.Grid {
	lines := make(map[types.LineID]types.Line)
	for lineID, cls := range clues {
		var size int
		if lineID.Direction == types.Row {
			size = width
		} else {
			size = height
		}
		lines[lineID] = createLineFromClues(cls, size)
	}
	return types.Grid{Lines: lines, Width: width, Height: height, ColorMap: colorMap}
}

func createGridFromCluesWithGoroutines(clues map[types.LineID][]types.ClueItem, width, height int, colorMap map[int]string) types.Grid {
	lines := make(map[types.LineID]types.Line)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for lineID, cls := range clues {
		wg.Add(1)
		go func(lid types.LineID, c []types.ClueItem) {
			defer wg.Done()
			var size int
			if lid.Direction == types.Row {
				size = width
			} else {
				size = height
			}
			line := createLineFromClues(c, size)
			mu.Lock()
			lines[lid] = line
			mu.Unlock()
		}(lineID, cls)
	}

	wg.Wait()
	return types.Grid{Lines: lines, Width: width, Height: height, ColorMap: colorMap}
}

// createLineFromClues creates a Line from clues and the line size
func createLineFromClues(clues []types.ClueItem, size int) types.Line {
	// Find all unique colors in this line
	colorSet := make(map[int]bool)
	for _, clue := range clues {
		colorSet[clue.ColorID] = true
	}

	// Initialize Combinations and Facts for each color
	combinations := GenerateColorCombinations(clues, size)
	facts := make(map[int]*big.Int)

	for colorID := range colorSet {
		facts[colorID] = big.NewInt(0)
	}

	return types.Line{
		Combinations: combinations,
		Facts:        facts,
	}
}
