package types

import "math/big"

// Facts represents the known facts about a line (bitsets for filled and empty positions)
type Facts struct {
	FilledByColor map[int]*Bitset // bits 1 = must be that color
	EmptyMask     *Bitset         // bits 1 = must be empty
}

// IsKnown returns true if the position is known (either filled or empty)
func (f *Facts) IsKnown(i int) bool {
	if f.EmptyMask.Bit(i) == 1 {
		return true
	}
	for _, bitset := range f.FilledByColor {
		if bitset.Bit(i) == 1 {
			return true
		}
	}
	return false
}

// MarkEmpty marks position i as empty
func (f *Facts) MarkEmpty(i int) {
	f.EmptyMask.SetBit(nil, i, 1)
}

// MarkFilled marks position i as filled with the given color
func (f *Facts) MarkFilled(i int, color int) {
	if f.FilledByColor[color] == nil {
		f.FilledByColor[color] = NewBitset(big.NewInt(0))
	}
	f.FilledByColor[color].SetBit(nil, i, 1)
}
