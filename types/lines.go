package types

import (
	"fmt"
	"sort"
)

type Lines struct {
	Lines map[LineID]Line
}

func (l *Lines) GetLine(lineID LineID) Line {
	return l.Lines[lineID]
}

func (l *Lines) GetLines() map[LineID]Line {
	return l.Lines
}

func (l *Lines) SetLine(lineID LineID, line Line) {
	l.Lines[lineID] = line
}

// Print prints all lines in a readable format, organized by direction and sorted by index
func (l *Lines) Print() {
	fmt.Printf("Lines{Total: %d lines}\n", len(l.Lines))

	// Separate rows and columns
	var rows []LineID
	var columns []LineID

	for lineID := range l.Lines {
		if lineID.Direction == Row {
			rows = append(rows, lineID)
		} else {
			columns = append(columns, lineID)
		}
	}

	// Sort by index
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Index < rows[j].Index
	})
	sort.Slice(columns, func(i, j int) bool {
		return columns[i].Index < columns[j].Index
	})

	// Print rows section
	if len(rows) > 0 {
		fmt.Printf("\nRows (%d):\n", len(rows))
		for _, lineID := range rows {
			fmt.Printf("  Row[%d]: ", lineID.Index)
			line := l.Lines[lineID]
			line.Print()
		}
	}

	// Print columns section
	if len(columns) > 0 {
		fmt.Printf("\nColumns (%d):\n", len(columns))
		for _, lineID := range columns {
			fmt.Printf("  Column[%d]: ", lineID.Index)
			line := l.Lines[lineID]
			line.Print()
		}
	}
	fmt.Println()
}
