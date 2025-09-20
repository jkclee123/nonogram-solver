package types

import (
	"container/list"
	"fmt"
	"math/big"
)

type Block struct {
	ColorID      uint8
	Size         uint8
	Combinations *list.List
}

// Print prints the block information in a readable format
func (b *Block) Print() {
	fmt.Printf("Block{ColorID: %d, Size: %d", b.ColorID, b.Size)
	if b.Combinations != nil && b.Combinations.Len() > 0 {
		fmt.Printf(", Combinations: %d total", b.Combinations.Len())
		// Show first few combinations
		count := 0
		fmt.Print(" [")
		for e := b.Combinations.Front(); e != nil; e = e.Next() {
			if count > 0 {
				fmt.Print(", ")
			}
			combination := e.Value.(*big.Int)
			fmt.Printf("%032b", combination)
			count++
		}
		fmt.Print("]")
	} else {
		fmt.Print(", Combinations: none")
	}
	fmt.Println("}")
}

// PrintWithWidth prints the block and its combinations, padding bitmasks to the given width.
func (b *Block) PrintWithWidth(width uint8) {
	fmt.Printf("Block{ColorID: %d, Size: %d", b.ColorID, b.Size)
	if b.Combinations != nil && b.Combinations.Len() > 0 {
		fmt.Printf(", Combinations: %d total", b.Combinations.Len())
		count := 0
		fmt.Print(" [")
		for e := b.Combinations.Front(); e != nil; e = e.Next() {
			if count > 0 {
				fmt.Print(", ")
			}
			if combination, ok := e.Value.(*big.Int); ok {
				// Render bits left-to-right with cell index == bit index.
				// Position 0 (leftmost) corresponds to bit 0, etc.
				w := int(width)
				builder := make([]byte, w)
				for pos := 0; pos < w; pos++ {
					if combination.Bit(pos) == 1 {
						builder[pos] = '1'
					} else {
						builder[pos] = '0'
					}
				}
				fmt.Print(string(builder))
			} else {
				fmt.Print("<invalid>")
			}
			count++
		}
		fmt.Print("]")
	} else {
		fmt.Print(", Combinations: none")
	}
	fmt.Print("}")
}
