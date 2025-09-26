package types

import "fmt"

// Grid represents a nonogram grid with rows and columns of Lines
type Grid struct {
	Rows []*Line
	Cols []*Line
}

// Orthogonal returns the orthogonal line and index for the given line at the specified position.
// For a row line at index i, returns the column line at index i.
// For a column line at index i, returns the row line at index i.
func (g *Grid) Orthogonal(lineID LineID, index int) (LineID, int) {
	switch lineID.Direction {
	case Row:
		return LineID{Direction: Column, Index: index}, lineID.Index
	case Column:
		return LineID{Direction: Row, Index: index}, lineID.Index
	default:
		panic("invalid direction")
	}
}

// Width returns the width of the grid (number of columns)
func (g *Grid) Width() int {
	return len(g.Cols)
}

// Height returns the height of the grid (number of rows)
func (g *Grid) Height() int {
	return len(g.Rows)
}

// Print prints a simple representation of the grid
func (g *Grid) Print() {
	for _, row := range g.Rows {
		for i := 0; i < row.Length; i++ {
			// Simple representation - could be enhanced based on facts
			if row.Facts != nil && row.Facts.EmptyMask.Bit(i) == 1 {
				fmt.Print(".")
			} else {
				fmt.Print("?")
			}
		}
		fmt.Println()
	}
}
