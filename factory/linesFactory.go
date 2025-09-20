package factory

import (
	"sync"

	"nonogram-solver/types"
)

// CreateLines takes NonogramData and creates Lines containing all row and column lines
func CreateLines(data types.NonogramData) types.Lines {
	// Use goroutines when total line count is 50 or more
	if len(data.Clues) >= 50 {
		return createLinesWithGoroutines(data)
	}

	return createLinesSingleThreaded(data)
}

// createLinesSingleThreaded creates lines using single-threaded approach
func createLinesSingleThreaded(data types.NonogramData) types.Lines {
	lines := types.Lines{
		Lines: make(map[types.LineID]types.Line),
	}

	for i := range data.Clues {
		lineID := data.Clues[i].LineID
		if lineID.Direction == types.Row {
			lines.SetLine(lineID, createLineFromClues(data.Clues[i], data.Width))
		} else {
			lines.SetLine(lineID, createLineFromClues(data.Clues[i], data.Height))
		}
	}

	return lines
}

// createLinesWithGoroutines creates lines using goroutines for parallel processing
func createLinesWithGoroutines(data types.NonogramData) types.Lines {
	lines := types.Lines{
		Lines: make(map[types.LineID]types.Line),
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	// Process all clues in parallel goroutines
	for i := range data.Clues {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			clue := data.Clues[index]
			lineID := clue.LineID

			var line types.Line
			if lineID.Direction == types.Row {
				line = createLineFromClues(clue, data.Width)
			} else {
				line = createLineFromClues(clue, data.Height)
			}

			mu.Lock()
			lines.SetLine(lineID, line)
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	return lines
}

// createLineFromClues creates a Line from a LineClue and the line size
func createLineFromClues(lineClue types.LineClue, size int) types.Line {
	clues := lineClue.Clues
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
