package factory

import (
	"math/big"
	"sort"

	"nonogram-solver/types"
)

func GenerateColorCombinations(clues []types.ClueItem, size int, colorID int) []*big.Int {
	if size <= 0 || len(clues) == 0 {
		return []*big.Int{}
	}

	// Generate combinations only for the requested color using backtracking.
	// We traverse all clues (to respect order and same-color gap rules),
	// but we only accumulate and deduplicate the mask for the target color.
	results := make([]*big.Int, 0)
	seen := make(map[string]struct{})

	n := len(clues)

	// Precompute minimal suffix requirement for clues[i:]
	suffix := make([]int, n+1)
	suffix[n] = 0
	for i := n - 1; i >= 0; i-- {
		gap := 0
		if i < n-1 && clues[i].ColorID == clues[i+1].ColorID {
			gap = 1
		}
		suffix[i] = clues[i].Clue + gap + suffix[i+1]
	}

	// Current accumulated mask for the target color only
	currentMask := big.NewInt(0)

	var backtrack func(i, currentPos int)
	backtrack = func(i, currentPos int) {
		if i == n {
			key := string(currentMask.Bytes())
			if _, ok := seen[key]; !ok {
				seen[key] = struct{}{}
				results = append(results, new(big.Int).Set(currentMask))
			}
			return
		}

		clue := clues[i]
		// Respect gap if same color as previous
		minStart := currentPos
		if i > 0 && clues[i-1].ColorID == clue.ColorID {
			minStart++
		}

		// Must leave room for the rest
		maxStart := size - clue.Clue - suffix[i+1]
		if maxStart < minStart {
			return
		}

		for start := minStart; start <= maxStart; start++ {
			block := runMask(size, start, clue.Clue)

			// Only OR/AND the mask if this clue matches the requested color
			if clue.ColorID == colorID {
				currentMask.Or(currentMask, block)
			}

			nextPos := start + clue.Clue
			backtrack(i+1, nextPos)

			if clue.ColorID == colorID {
				currentMask.AndNot(currentMask, block)
			}
		}
	}

	backtrack(0, 0)

	// Sort results in descending order to match expected output
	sort.Slice(results, func(i, j int) bool { return results[i].Cmp(results[j]) > 0 })

	return results
}

// runMask returns a mask with a contiguous run of 'length' ones starting at 'start' (0-based from left).
// Bit mapping: leftmost cell -> most significant bit (bit index size-1), rightmost -> bit 0.
func runMask(size, start, length int) *big.Int {
	// Efficient contiguous run using shifts: ((1<<length)-1) << (size-start-length)
	if length <= 0 {
		return big.NewInt(0)
	}
	ones := new(big.Int).Lsh(big.NewInt(1), uint(length))
	ones.Sub(ones, big.NewInt(1))
	shift := uint(size - start - length)
	return new(big.Int).Lsh(ones, shift)
}
