package types

import (
	"fmt"
	"sort"
)

type NonogramData struct {
	Clues    map[LineID][]ClueItem
	Width    int
	Height   int
	ColorMap map[int]string
}

// Print displays the parsed nonogram clues in a readable format
func (nd *NonogramData) Print() {
	fmt.Printf("\n=== ROW CLUES ===\n")

	// Collect and sort row line IDs
	var rowLineIDs []LineID
	for lineID := range nd.Clues {
		if lineID.Direction == Row {
			rowLineIDs = append(rowLineIDs, lineID)
		}
	}
	sort.Slice(rowLineIDs, func(i, j int) bool {
		return rowLineIDs[i].Index < rowLineIDs[j].Index
	})

	// Print rows in sorted order
	for _, lineID := range rowLineIDs {
		clues := nd.Clues[lineID]
		fmt.Printf("Row %d: ", lineID.Index+1)
		for j, clue := range clues {
			if j > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("(%d,%d)", clue.ColorID, clue.Clue)
		}
		fmt.Printf("\n")
	}

	fmt.Printf("\n=== COLUMN CLUES ===\n")

	// Collect and sort column line IDs
	var colLineIDs []LineID
	for lineID := range nd.Clues {
		if lineID.Direction == Column {
			colLineIDs = append(colLineIDs, lineID)
		}
	}
	sort.Slice(colLineIDs, func(i, j int) bool {
		return colLineIDs[i].Index < colLineIDs[j].Index
	})

	// Print columns in sorted order
	for _, lineID := range colLineIDs {
		clues := nd.Clues[lineID]
		fmt.Printf("Col %d: ", lineID.Index+1)
		for j, clue := range clues {
			if j > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("(%d,%d)", clue.ColorID, clue.Clue)
		}
		fmt.Printf("\n")
	}

	fmt.Printf("\n=== COLOR MAP ===\n")
	if len(nd.ColorMap) > 0 {
		for colorID, hexColor := range nd.ColorMap {
			fmt.Printf("Color %d: %s\n", colorID, hexColor)
		}
	} else {
		fmt.Printf("No colors found\n")
	}

	fmt.Printf("\nGrid size: %dx%d\n", nd.Width, nd.Height)
}
