package types

import (
	"fmt"
	"math/big"
)

type Block struct {
	ColorID      uint8
	Size         uint8
	Combinations []*big.Int
}

// Print prints the block information in a readable format
func (b *Block) Print() {
	fmt.Printf("Block{ColorID: %d, Size: %d", b.ColorID, b.Size)
	if len(b.Combinations) > 0 {
		fmt.Printf(", Combinations: %d total", len(b.Combinations))
		// Show first few combinations
		count := 0
		fmt.Print(" [")
		for _, combination := range b.Combinations {
			if count > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%032b", combination)
			count++
		}
		fmt.Print("]")
	} else {
		fmt.Print(", Combinations: none")
	}
	fmt.Println("}")
}

// BitwiseOr performs bitwise OR on all combinations within this block.
// Returns a big.Int where each bit represents the OR result across all combinations.
func (b *Block) BitwiseOr() *big.Int {
	result := big.NewInt(0)
	for _, combination := range b.Combinations {
		result.Or(result, combination)
	}
	return result
}

// BitwiseAnd performs bitwise AND on all combinations within this block.
// Returns a big.Int where each bit represents the AND result across all combinations.
func (b *Block) BitwiseAnd() *big.Int {
	result := big.NewInt(0)
	for _, combination := range b.Combinations {
		result.And(result, combination)
	}
	return result
}

// PrintWithWidth prints the block and its combinations, padding bitmasks to the given width.
func (b *Block) PrintWithWidth(width int) {
	fmt.Printf("Block{ColorID: %d, Size: %d", b.ColorID, b.Size)
	if len(b.Combinations) > 0 {
		fmt.Printf(", Combinations: %d total", len(b.Combinations))
		count := 0
		fmt.Print(" [")
		for _, combination := range b.Combinations {
			if count > 0 {
				fmt.Print(", ")
			}
			// Render bits left-to-right with cell index == bit index.
			// Position 0 (leftmost) corresponds to bit 0, etc.
			w := width
			builder := make([]byte, w)
			for pos := 0; pos < w; pos++ {
				if combination.Bit(pos) == 1 {
					builder[pos] = '1'
				} else {
					builder[pos] = '0'
				}
			}
			fmt.Print(string(builder))
			count++
		}
		fmt.Print("]")
	} else {
		fmt.Print(", Combinations: none")
	}
	fmt.Print("}")
}
