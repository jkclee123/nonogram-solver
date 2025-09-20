package types

type Grid struct {
	Grid      [][]uint8
	Width     uint8
	Height    uint8
	ColorMap  map[uint8]string
}
