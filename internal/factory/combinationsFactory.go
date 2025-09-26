package factory

import (
	"math/big"
	"sort"
	"sync"

	"nonogram-solver/types"
)

// GenerateCombinations returns the set of bitmask combinations for the i-th block
// described by clues within a line of given size. It honors color adjacency rules:
// adjacent blocks with the same color require at least one empty cell between them;
// blocks with different colors can be adjacent without an empty cell.
var combinationsCache sync.Map // key string -> map[int][]*big.Int

func GenerateColorCombinations(clues []types.ClueItem, size int) map[int][]*big.Int {
	combinations := make(map[int][]*big.Int)

	if size <= 0 || len(clues) == 0 {
		return combinations
	}

	// Cache lookup
	key := cacheKey(clues, size)
	if cached, ok := combinationsCache.Load(key); ok {
		return cached.(map[int][]*big.Int)
	}

	// General case: handle any number of clues with mixed patterns (optimized)
	genGeneralColorCombinations(clues, size, combinations)
	combinationsCache.Store(key, combinations)
	return combinations
}

func cacheKey(clues []types.ClueItem, size int) string {
	// Build a stable key: size|n|c1:len1|c2:len2|...
	// Using small stack buffer via byte slice builder
	b := make([]byte, 0, 32+len(clues)*6)
	b = appendInt(b, size)
	b = append(b, '|')
	b = appendInt(b, len(clues))
	for _, c := range clues {
		b = append(b, '|')
		b = appendInt(b, c.ColorID)
		b = append(b, ':')
		b = appendInt(b, c.Clue)
	}
	return string(b)
}

func appendInt(dst []byte, v int) []byte {
	if v == 0 {
		return append(dst, '0')
	}
	// collect digits in reverse
	var buf [20]byte
	i := len(buf)
	n := v
	for n > 0 {
		i--
		buf[i] = byte('0' + (n % 10))
		n /= 10
	}
	return append(dst, buf[i:]...)
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

// genGeneralColorCombinations handles arbitrary number of clues with mixed color patterns.
// For each color, it finds all possible positions where blocks of that color can be placed,
// considering the constraints from all other blocks.
func genGeneralColorCombinations(clues []types.ClueItem, size int, out map[int][]*big.Int) {
	n := len(clues)
	if n == 0 {
		return
	}

	// Collect unique colors
	colors := make(map[int]struct{})
	for _, c := range clues {
		colors[c.ColorID] = struct{}{}
	}

	// Per-color accumulated masks (mutated in place during backtracking)
	current := make(map[int]*big.Int, len(colors))
	for color := range colors {
		current[color] = big.NewInt(0)
	}

	// Dedup sets per-color using byte representation as key
	seen := make(map[int]map[string]struct{}, len(colors))
	for color := range colors {
		seen[color] = make(map[string]struct{})
	}

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

	var backtrack func(i, currentPos int)
	backtrack = func(i, currentPos int) {
		if i == n {
			// Record masks per color
			for color, m := range current {
				key := string(m.Bytes())
				if _, ok := seen[color][key]; !ok {
					seen[color][key] = struct{}{}
					out[color] = append(out[color], new(big.Int).Set(m))
				}
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

			cm := current[clue.ColorID]
			// Add block bits
			cm.Or(cm, block)
			nextPos := start + clue.Clue

			backtrack(i+1, nextPos)

			// Remove block bits (revert)
			cm.AndNot(cm, block)
		}
	}

	backtrack(0, 0)

	// Sort results per color in descending order (to match expected output)
	for _, masks := range out {
		sort.Slice(masks, func(i, j int) bool { return masks[i].Cmp(masks[j]) > 0 })
	}
}
