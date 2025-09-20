package factory

import (
	"sync"

	"nonogram-solver/types"
)

// CreateLines takes NonogramData and creates Lines containing all row and column lines
func CreateLines(data types.NonogramData) types.Lines {
	totalLines := len(data.RowClues) + len(data.ColumnClues)

	// Use goroutines when total line count is 50 or more
	if totalLines >= 50 {
		return createLinesWithGoroutines(data)
	}

	return createLinesSingleThreaded(data)
}

// createLinesSingleThreaded creates lines using single-threaded approach
func createLinesSingleThreaded(data types.NonogramData) types.Lines {
	lines := types.Lines{
		Lines: make(map[types.LineID]types.Line),
	}

	numRows := len(data.RowClues)
	numCols := len(data.ColumnClues)

	// Create row lines
	for i := range numRows {
		lineID := types.LineID{
			Direction: types.Row,
			Index:     uint8(i),
		}

		line := createLineFromClues(data.RowClues[i], numCols)
		lines.SetLine(lineID, line)
	}

	// Create column lines
	for i := range numCols {
		lineID := types.LineID{
			Direction: types.Column,
			Index:     uint8(i),
		}

		line := createLineFromClues(data.ColumnClues[i], numRows)
		lines.SetLine(lineID, line)
	}

	return lines
}

// createLinesWithGoroutines creates lines using goroutines for parallel processing
func createLinesWithGoroutines(data types.NonogramData) types.Lines {
	lines := types.Lines{
		Lines: make(map[types.LineID]types.Line),
	}

	numRows := len(data.RowClues)
	numCols := len(data.ColumnClues)

	var wg sync.WaitGroup
	var mu sync.Mutex

	// Create row lines in parallel
	for i := range numRows {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			lineID := types.LineID{
				Direction: types.Row,
				Index:     uint8(index),
			}

			line := createLineFromClues(data.RowClues[index], numCols)

			mu.Lock()
			lines.SetLine(lineID, line)
			mu.Unlock()
		}(i)
	}

	// Create column lines in parallel
	for i := range numCols {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			lineID := types.LineID{
				Direction: types.Column,
				Index:     uint8(index),
			}

			line := createLineFromClues(data.ColumnClues[index], numRows)

			mu.Lock()
			lines.SetLine(lineID, line)
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	return lines
}

// createLineFromClues creates a Line from a slice of ClueItems and the line size
func createLineFromClues(clues []types.ClueItem, size int) types.Line {
	blocks := make([]types.Block, len(clues))

	for i, clue := range clues {
		blocks[i] = types.Block{
			ColorID:      clue.ColorID,
			Size:         clue.BlockSize,
			Combinations: GenerateCombinations(clues, size, i),
		}
	}

	return types.Line{
		Blocks: blocks,
		Size:   size,
	}
}
