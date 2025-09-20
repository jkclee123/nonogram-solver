package factory

import (
	"container/list"
	"math/big"
	"sort"

	"nonogram-solver/types"
)

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
			line.Blocks[i].Combinations = list.New()
		}
		return
	}

	freeSlots := lineSize - minRequired
	numGaps := numBlocks + 1 // left pad, inter-gaps, right pad

	// Prepare unique mask sets per block.
	perBlockMasks := make([]map[string]*big.Int, numBlocks)
	for i := 0; i < numBlocks; i++ {
		perBlockMasks[i] = make(map[string]*big.Int)
	}

	// Enumerate distributions of extra spaces across gaps.
	enumerateDistributions(numGaps, freeSlots, func(extras []int) {
		// Compute start positions from extras and minGaps
		starts := make([]int, numBlocks)
		pos := extras[0]
		for bi := 0; bi < numBlocks; bi++ {
			starts[bi] = pos
			pos += int(line.Blocks[bi].Size)
			if bi < numBlocks-1 {
				pos += minGaps[bi] + extras[bi+1]
			}
		}

		// Record masks for each block at these starts
		for bi := 0; bi < numBlocks; bi++ {
			mask := makeRunMask(int(line.Blocks[bi].Size), starts[bi])
			key := mask.Text(10)
			if _, exists := perBlockMasks[bi][key]; !exists {
				perBlockMasks[bi][key] = mask
			}
		}
	})

	// Transfer unique masks into each block's list, sorted ascending for determinism.
	for bi := 0; bi < numBlocks; bi++ {
		masks := make([]*big.Int, 0, len(perBlockMasks[bi]))
		for _, m := range perBlockMasks[bi] {
			masks = append(masks, m)
		}
		sort.Slice(masks, func(i, j int) bool { return masks[i].Cmp(masks[j]) < 0 })

		lst := list.New()
		for _, m := range masks {
			// Store a copy to avoid accidental mutations later
			lst.PushBack(new(big.Int).Set(m))
		}
		line.Blocks[bi].Combinations = lst
	}
}

// enumerateDistributions generates all non-negative integer vectors of length n
// summing to total, invoking cb for each vector.
func enumerateDistributions(n int, total int, cb func([]int)) {
	current := make([]int, n)
	var dfs func(idx int, remaining int)
	dfs = func(idx int, remaining int) {
		if idx == n-1 {
			current[idx] = remaining
			cb(current)
			return
		}
		for v := 0; v <= remaining; v++ {
			current[idx] = v
			dfs(idx+1, remaining-v)
		}
	}
	dfs(0, total)
}

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
