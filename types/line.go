package types

import (
	"fmt"
	"math/big"
)

type Line struct {
	Combinations map[int][]*big.Int
	Facts        map[int]*big.Int
}

// Print prints the line information in a readable format
func (l *Line) Print() {
	fmt.Printf("Line{Combinations: %d, Facts: %d", len(l.Combinations), len(l.Facts))

	if len(l.Combinations) > 0 {
		fmt.Printf(", Combinations: {")
		first := true
		for colorID, combos := range l.Combinations {
			if !first {
				fmt.Printf(", ")
			}
			fmt.Printf("%d: %d combos", colorID, len(combos))
			first = false
		}
		fmt.Printf("}")
	}

	if len(l.Facts) > 0 {
		fmt.Printf(", Facts: {")
		first := true
		for colorID, fact := range l.Facts {
			if !first {
				fmt.Printf(", ")
			}
			fmt.Printf("%d: %s", colorID, fact.String())
			first = false
		}
		fmt.Printf("}")
	}

	fmt.Println("}")
}
