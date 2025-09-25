package types

import (
)

type Grid struct {
	Lines map[LineID]Line
	Width int
	Height int
	ColorMap map[int]string
}
