package factory

import (
	"math/big"

	"nonogram-solver/types"
)

// CreateGridFromClues builds a Grid directly from clues, dimensions, and color map.
// It mirrors CreateGrid but avoids the NonogramData wrapper.
func CreateGridFromClues(clues map[types.LineID][]types.ClueItem, width, height int, colorMap map[int]string) types.Grid {
	lines := make(map[types.LineID]types.Line)
	for lineID, cls := range clues {
		var size int
		if lineID.Direction == types.Row {
			size = width
		} else {
			size = height
		}
		lines[lineID] = createLineFromClues(cls)
	}
	return types.Grid{Lines: lines, Width: width, Height: height, ColorMap: colorMap}
}

// createLineFromClues creates a Line from clues and the line size
func createLineFromClues(clues []types.ClueItem) types.Line {
	// Find all unique colors in this line
	colorSet := make(map[int]bool)
	for _, clue := range clues {
		colorSet[clue.ColorID] = true
	}

	// Initialize Combinations and Facts for each color (combinations will be lazily generated)
	combinations := make(map[int][]*big.Int)
	facts := make(map[int]*big.Int)

	for colorID := range colorSet {
		combinations[colorID] = []*big.Int{} // Empty slice, will be populated lazily
		facts[colorID] = big.NewInt(0)
	}

	return types.Line{
		Combinations: combinations,
		Facts:        facts,
	}
}
