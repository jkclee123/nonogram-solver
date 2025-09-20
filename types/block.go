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
		for e := b.Combinations.Front(); e != nil && count < 2; e = e.Next() {
			if count > 0 {
				fmt.Print(", ")
			}
			combination := e.Value.(*big.Int)
			fmt.Printf("%032b", combination)
			count++
		}
		if b.Combinations.Len() > 2 {
			fmt.Print("...")
		}
		fmt.Print("]")
	} else {
		fmt.Print(", Combinations: none")
	}
	fmt.Println("}")
}
