package factory

import (
	"nonogram-solver/types"
)

// CreateLines takes NonogramData and creates Lines containing all row and column lines
func CreateLines(data types.NonogramData) types.Lines {
	lines := types.Lines{
		Lines: make(map[types.LineID]types.Line),
	}

	numRows := len(data.RowClues)
	numCols := len(data.ColumnClues)

	// Create row lines
	for i := 0; i < numRows; i++ {
		lineID := types.LineID{
			Direction: types.Row,
			Index:     uint8(i),
		}

		line := createLineFromClues(data.RowClues[i], uint8(numCols), i)
		lines.SetLine(lineID, line)
	}

	// Create column lines
	for i := 0; i < numCols; i++ {
		lineID := types.LineID{
			Direction: types.Column,
			Index:     uint8(i),
		}

		line := createLineFromClues(data.ColumnClues[i], uint8(numRows), i)
		lines.SetLine(lineID, line)
	}

	return lines
}

// createLineFromClues creates a Line from a slice of ClueItems and the line size
func createLineFromClues(clues []types.ClueItem, size uint8, i int) types.Line {
	blocks := make([]types.Block, len(clues))

	for bi, clue := range clues {
		blocks[bi] = types.Block{
			ColorID:      clue.ColorID,
			Size:         clue.BlockSize,
			Combinations: GenerateCombinations(clues, size, bi),
		}
	}

	return types.Line{
		Blocks: blocks,
		Size:   size,
	}
}
