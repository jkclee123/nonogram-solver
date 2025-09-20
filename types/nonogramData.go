package types

import "fmt"

type NonogramData struct {
	RowClues    [][]ClueItem
	ColumnClues [][]ClueItem
}

// Print displays the parsed nonogram clues in a readable format
func (nd *NonogramData) Print() {
	fmt.Printf("\n=== ROW CLUES ===\n")
	for i, row := range nd.RowClues {
		fmt.Printf("Row %d: ", i+1)
		for j, clue := range row {
			if j > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("(%d,%d)", clue.Color, clue.BlockSize)
		}
		fmt.Printf("\n")
	}

	fmt.Printf("\n=== COLUMN CLUES ===\n")
	for i, col := range nd.ColumnClues {
		fmt.Printf("Col %d: ", i+1)
		for j, clue := range col {
			if j > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("(%d,%d)", clue.Color, clue.BlockSize)
		}
		fmt.Printf("\n")
	}

	fmt.Printf("\nGrid size: %dx%d\n", len(nd.ColumnClues), len(nd.RowClues))
}
