package types

import (
	"container/list"
)

type Block struct {
	ColorID      uint8
	Size         uint8
	Combinations *list.List
}
