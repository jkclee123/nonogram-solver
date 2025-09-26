package grid

import "nonogram-solver/types"

// OrthogonalLineIndex returns the opposite line id for a given cell index.
// For rows, index is column; for columns, index is row.
func OrthogonalLineIndex(id types.LineID, pos int) types.LineID {
	if id.Direction == types.Row {
		return types.LineID{Direction: types.Column, Index: pos}
	}
	return types.LineID{Direction: types.Row, Index: pos}
}
