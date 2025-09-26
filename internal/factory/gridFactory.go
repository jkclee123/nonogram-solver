package factory

import (
	"math/big"
	"nonogram-solver/internal/types"
)

// CreateGridFromClues creates a Grid from extracted clues
func CreateGridFromClues(clues map[types.LineID][]types.ClueItem, width, height int, colorMap map[int]string) types.Grid {
	grid := types.Grid{
		Rows: make([]*types.Line, height),
		Cols: make([]*types.Line, width),
	}

	// Create rows
	for row := 0; row < height; row++ {
		lineID := types.LineID{Direction: types.Row, Index: row}
		clueList := clues[lineID]
		grid.Rows[row] = &types.Line{
			ID:        lineID,
			Direction: types.Row,
			Length:    width,
			Clues:     clueList,
			Facts: &types.Facts{
				FilledByColor: make(map[int]*types.Bitset),
				EmptyMask:     types.NewBitset(big.NewInt(0)),
			},
			Combinations: NewCombinationsProvider(clueList, width),
		}
	}

	// Create columns
	for col := 0; col < width; col++ {
		lineID := types.LineID{Direction: types.Column, Index: col}
		clueList := clues[lineID]
		grid.Cols[col] = &types.Line{
			ID:        lineID,
			Direction: types.Column,
			Length:    height,
			Clues:     clueList,
			Facts: &types.Facts{
				FilledByColor: make(map[int]*types.Bitset),
				EmptyMask:     types.NewBitset(big.NewInt(0)),
			},
			Combinations: NewCombinationsProvider(clueList, height),
		}
	}

	return grid
}
