package types

import "fmt"

type NonogramData struct {
	Clues    []LineClue
	Width    int
	Height   int
	ColorMap map[int]string
}

// Print displays the parsed nonogram clues in a readable format
func (nd *NonogramData) Print() {
	// Separate row and column clues
	rowClues := make([]LineClue, 0)
	columnClues := make([]LineClue, 0)

	for _, clue := range nd.Clues {
		switch clue.LineID.Direction {
		case Row:
			rowClues = append(rowClues, clue)
		case Column:
			columnClues = append(columnClues, clue)
		}
	}

	fmt.Printf("\n=== ROW CLUES ===\n")
	for _, row := range rowClues {
		fmt.Printf("Row %d: ", row.LineID.Index+1)
		for j, clue := range row.Clues {
			if j > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("(%d,%d)", clue.ColorID, clue.BlockSize)
		}
		fmt.Printf("\n")
	}

	fmt.Printf("\n=== COLUMN CLUES ===\n")
	for _, col := range columnClues {
		fmt.Printf("Col %d: ", col.LineID.Index+1)
		for j, clue := range col.Clues {
			if j > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("(%d,%d)", clue.ColorID, clue.BlockSize)
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

	fmt.Printf("\nGrid size: %dx%d\n", len(columnClues), len(rowClues))
}
