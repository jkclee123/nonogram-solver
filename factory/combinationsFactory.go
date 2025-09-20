package factory

import (
	"math/big"

	"nonogram-solver/types"
)

// GenerateCombinations returns the set of bitmask combinations for the i-th block
// described by clues within a line of given size. It honors color adjacency rules:
// adjacent blocks with the same color require at least one empty cell between them;
// blocks with different colors can be adjacent without an empty cell.
func GenerateCombinations(clues []types.ClueItem, size uint8, i int) []*big.Int {
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

// GenerateCombinationsForLines populates combinations for every line in the collection.
func GenerateCombinationsForLines(lines *types.Lines) {
	if lines == nil {
		return
	}
	for id, line := range lines.GetLines() {
		GenerateCombinationsForLine(&line)
		lines.SetLine(id, line)
	}
}

// GenerateCombinationsForLine computes, for each block in the line, the set of unique
// bitmask combinations representing all positions that block can occupy across all
// valid placements of the entire line, honoring color adjacency rules.
func GenerateCombinationsForLine(line *types.Line) {
	if line == nil || line.Size == 0 || len(line.Blocks) == 0 {
		return
	}

	numBlocks := len(line.Blocks)

	// Minimum required gap between adjacent blocks based on color equality.
	// If two adjacent blocks share the same color, they require at least one empty cell.
	// If colors differ, no empty cell is required (they can be adjacent).
	minGaps := make([]int, maxInt(numBlocks-1, 0))
	for i := 0; i < len(minGaps); i++ {
		if types.IsSameColor(line, uint8(i), uint8(i+1)) {
			minGaps[i] = 1
		} else {
			minGaps[i] = 0
		}
	}

	// Compute minimal occupied width and available free slots to distribute.
	sizesSum := 0
	for i := 0; i < numBlocks; i++ {
		sizesSum += int(line.Blocks[i].Size)
	}
	minGapSum := 0
	for _, g := range minGaps {
		minGapSum += g
	}

	lineSize := int(line.Size)
	minRequired := sizesSum + minGapSum
	if minRequired > lineSize {
		// Impossible line; leave combinations empty.
		for i := 0; i < numBlocks; i++ {
			line.Blocks[i].Combinations = []*big.Int{}
		}
		return
	}

	freeSlots := lineSize - minRequired

	// Precompute minimal start for each block (with zero extra spaces allocated).
	// startMin[bi] = sum_{j<bi} size[j] + sum_{j<bi} minGaps[j]
	startMin := make([]int, numBlocks)
	accum := 0
	for bi := 0; bi < numBlocks; bi++ {
		startMin[bi] = accum
		accum += int(line.Blocks[bi].Size)
		if bi < numBlocks-1 {
			accum += minGaps[bi]
		}
	}

	// For each block, the set of possible starts is a contiguous range
	// from startMin to startMin+freeSlots. Generate masks by shifting.
	for bi := 0; bi < numBlocks; bi++ {
		var combinations []*big.Int
		runLen := int(line.Blocks[bi].Size)
		minStart := startMin[bi]
		maxStart := minStart + freeSlots
		for s := minStart; s <= maxStart; s++ {
			combinations = append(combinations, makeRunMask(runLen, s))
		}
		line.Blocks[bi].Combinations = combinations
	}
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
