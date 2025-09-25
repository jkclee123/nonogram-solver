package types

import (
)

type Lines struct {
	Lines map[LineID]Line
	Width int
	Height int
	ColorMap map[int]string
}
