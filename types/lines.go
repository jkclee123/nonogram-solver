package types

import "fmt"

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

// Print prints all lines in a readable format
func (l *Lines) Print() {
	fmt.Printf("Lines{Total: %d lines\n", len(l.Lines))
	for lineID, line := range l.Lines {
		fmt.Printf("  %s[%d]: ", lineID.Direction.String(), lineID.Index)
		fmt.Printf("Line{Size: %d, Blocks: [", line.Size)
		for i, block := range line.Blocks {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("Block{ColorID: %d, Size: %d}", block.ColorID, block.Size)
		}
		fmt.Println("]}")
	}
	fmt.Println("}")
}
