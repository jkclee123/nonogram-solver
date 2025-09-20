package types

type Direction int

const (
	Row Direction = iota
	Column
)

func (d Direction) String() string {
	switch d {
	case Row:
		return "Row"
	case Column:
		return "Column"
	default:
		return "Unknown"
	}
}
