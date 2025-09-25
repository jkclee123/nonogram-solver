package types

import (
	"fmt"
	"sort"
)

type Grid struct {
	Lines    map[LineID]Line
	Width    int
	Height   int
	ColorMap map[int]string
}

func (g *Grid) Print() {
	fmt.Printf("\n=== GRID INFO ===\n")
	fmt.Printf("Dimensions: %dx%d\n", g.Width, g.Height)
	fmt.Printf("Total lines: %d\n", len(g.Lines))

	// Count rows and columns
	rowCount := 0
	colCount := 0
	for lineID := range g.Lines {
		if lineID.Direction == Row {
			rowCount++
		} else {
			colCount++
		}
	}
	fmt.Printf("Rows: %d, Columns: %d\n", rowCount, colCount)

	fmt.Printf("\n=== COLOR MAP ===\n")
	if len(g.ColorMap) > 0 {
		// Collect and sort color IDs
		var colorIDs []int
		for colorID := range g.ColorMap {
			colorIDs = append(colorIDs, colorID)
		}
		sort.Ints(colorIDs)

		// Print colors in sorted order
		for _, colorID := range colorIDs {
			hexColor := g.ColorMap[colorID]
			fmt.Printf("Color %d: %s\n", colorID, hexColor)
		}
	} else {
		fmt.Printf("No colors found\n")
	}

	fmt.Printf("\n=== LINE DETAILS ===\n")
	// Collect and sort row line IDs
	var rowLineIDs []LineID
	var colLineIDs []LineID
	for lineID := range g.Lines {
		if lineID.Direction == Row {
			rowLineIDs = append(rowLineIDs, lineID)
		} else {
			colLineIDs = append(colLineIDs, lineID)
		}
	}

	sort.Slice(rowLineIDs, func(i, j int) bool {
		return rowLineIDs[i].Index < rowLineIDs[j].Index
	})
	sort.Slice(colLineIDs, func(i, j int) bool {
		return colLineIDs[i].Index < colLineIDs[j].Index
	})

	// Print row details
	fmt.Printf("Rows:\n")
	for _, lineID := range rowLineIDs {
		line := g.Lines[lineID]
		totalCombos := 0
		for _, combos := range line.Combinations {
			totalCombos += len(combos)
		}
		fmt.Printf("  Row %d: %d combinations, %d facts\n", lineID.Index, totalCombos, len(line.Facts))
	}

	// Print column details
	fmt.Printf("Columns:\n")
	for _, lineID := range colLineIDs {
		line := g.Lines[lineID]
		totalCombos := 0
		for _, combos := range line.Combinations {
			totalCombos += len(combos)
		}
		fmt.Printf("  Col %d: %d combinations, %d facts\n", lineID.Index, totalCombos, len(line.Facts))
	}
}
