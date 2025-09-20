package types

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