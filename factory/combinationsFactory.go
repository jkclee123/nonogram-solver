package factory

import (
	"math/big"

	"nonogram-solver/types"
)

// GenerateCombinations returns the set of bitmask combinations for the i-th block
// described by clues within a line of given size. It honors color adjacency rules:
// adjacent blocks with the same color require at least one empty cell between them;
// blocks with different colors can be adjacent without an empty cell.
func GenerateCombinations(clues []types.ClueItem, size int, i int) []*big.Int {
	if size == 0 || len(clues) == 0 || i < 0 || i >= len(clues) {
		return []*big.Int{}
	}

	numBlocks := len(clues)

	// Minimum required gap between adjacent blocks based on color equality.
	minGaps := make([]int, maxInt(numBlocks-1, 0))
	for gi := 0; gi < len(minGaps); gi++ {
		if clues[gi].ColorID == clues[gi+1].ColorID {
			minGaps[gi] = 1
		} else {
			minGaps[gi] = 0
		}
	}

	// Compute minimal occupied width and available free slots to distribute.
	sizesSum := 0
	for bi := 0; bi < numBlocks; bi++ {
		sizesSum += int(clues[bi].BlockSize)
	}
	minGapSum := 0
	for _, g := range minGaps {
		minGapSum += g
	}

	lineSize := int(size)
	minRequired := sizesSum + minGapSum
	if minRequired > lineSize {
		return []*big.Int{}
	}

	freeSlots := lineSize - minRequired

	// Precompute minimal start for each block (with zero extra spaces allocated).
	startMin := make([]int, numBlocks)
	accum := 0
	for bi := 0; bi < numBlocks; bi++ {
		startMin[bi] = accum
		accum += int(clues[bi].BlockSize)
		if bi < numBlocks-1 {
			accum += minGaps[bi]
		}
	}

	// Generate masks for the i-th block by shifting across the free slots.
	var result []*big.Int
	runLen := int(clues[i].BlockSize)
	minStart := startMin[i]
	maxStart := minStart + freeSlots
	for s := minStart; s <= maxStart; s++ {
		result = append(result, makeRunMask(runLen, s))
	}
	return result
}

// enumerateDistributions generates all non-negative integer vectors of length n
// summing to total, invoking cb for each vector.
// enumerateDistributions removed: no longer needed with shift-based generation

// makeRunMask creates a bitmask with a contiguous run of 1s of length runLen
// starting at bit position start (LSB at position 0).
func makeRunMask(runLen int, start int) *big.Int {
	if runLen <= 0 {
		return big.NewInt(0)
	}
	one := big.NewInt(1)
	mask := new(big.Int).Lsh(one, uint(runLen))
	mask.Sub(mask, one)
	mask.Lsh(mask, uint(start))
	return mask
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
