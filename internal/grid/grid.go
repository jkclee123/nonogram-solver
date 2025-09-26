package grid

import (
	"fmt"
	"nonogram-solver/internal/types"
)

// GridOperations provides additional operations on grids beyond the basic types
type GridOperations struct {
	Grid *types.Grid
}

// NewGridOperations creates a new GridOperations instance
func NewGridOperations(grid *types.Grid) *GridOperations {
	return &GridOperations{Grid: grid}
}

// GetLine returns the line for the given LineID
func (g *GridOperations) GetLine(lineID types.LineID) *types.Line {
	switch lineID.Direction {
	case types.Row:
		if lineID.Index >= 0 && lineID.Index < len(g.Grid.Rows) {
			return g.Grid.Rows[lineID.Index]
		}
	case types.Column:
		if lineID.Index >= 0 && lineID.Index < len(g.Grid.Cols) {
			return g.Grid.Cols[lineID.Index]
		}
	}
	return nil
}

// GetOrthogonalLines returns all lines that intersect with the given position
func (g *GridOperations) GetOrthogonalLines(lineID types.LineID, index int) []*types.Line {
	var lines []*types.Line
	switch lineID.Direction {
	case types.Row:
		// For a row, return the column at that index
		if index >= 0 && index < len(g.Grid.Cols) {
			lines = append(lines, g.Grid.Cols[index])
		}
	case types.Column:
		// For a column, return the row at that index
		if index >= 0 && index < len(g.Grid.Rows) {
			lines = append(lines, g.Grid.Rows[index])
		}
	}
	return lines
}

// ValidateGrid checks if the grid has consistent dimensions
func (g *GridOperations) ValidateGrid() error {
	if len(g.Grid.Rows) == 0 || len(g.Grid.Cols) == 0 {
		return fmt.Errorf("grid must have at least one row and one column")
	}

	// Check that all rows have the same length as the number of columns
	expectedRowLength := len(g.Grid.Cols)
	for i, row := range g.Grid.Rows {
		if row.Length != expectedRowLength {
			return fmt.Errorf("row %d has length %d, expected %d", i, row.Length, expectedRowLength)
		}
	}

	// Check that all columns have the same length as the number of rows
	expectedColLength := len(g.Grid.Rows)
	for i, col := range g.Grid.Cols {
		if col.Length != expectedColLength {
			return fmt.Errorf("column %d has length %d, expected %d", i, col.Length, expectedColLength)
		}
	}

	return nil
}
