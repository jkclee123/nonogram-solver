package line

import (
	"math/big"

	"nonogram-solver/types"
)

// FactsDelta captures changes applied to a line's facts
type FactsDelta struct {
	FilledByColor map[int]*big.Int // colorID -> mask of new fills
	Empties       *big.Int         // mask of new empties
}

// Overlap computes forced empties and fills for a given color.
// Returns whether any facts changed and the delta.
func Overlap(l *types.Line, colorID int) (changed bool, delta FactsDelta) {
	// stub
	return false, FactsDelta{}
}

// CrossReference filters combinations based on known facts.
// Returns whether combinations changed.
func CrossReference(l *types.Line, colorID int) (changed bool) {
	// stub
	return false
}
