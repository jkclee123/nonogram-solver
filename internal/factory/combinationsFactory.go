package factory

import (
	"math/big"
	"runtime"
	"sync"

	"nonogram-solver/internal/types"
)

// CombinationsProviderImpl implements combinatorics.CombinationsProvider
type CombinationsProviderImpl struct {
	clues         []types.ClueItem
	size          int
	generated     map[int]bool
	combosByColor map[int][]*types.Bitset
	mu            sync.RWMutex
}

// NewCombinationsProvider creates a new lazy combinations provider
func NewCombinationsProvider(clues []types.ClueItem, size int) *CombinationsProviderImpl {
	return &CombinationsProviderImpl{
		clues:         clues,
		size:          size,
		generated:     make(map[int]bool),
		combosByColor: make(map[int][]*types.Bitset),
	}
}

// Get returns combinations for the specified color, generating them lazily if needed
func (cp *CombinationsProviderImpl) Get(color int) ([]*types.Bitset, error) {
	cp.mu.RLock()
	if cp.generated[color] {
		result := cp.combosByColor[color]
		cp.mu.RUnlock()
		return result, nil
	}
	cp.mu.RUnlock()

	// Generate combinations outside of read lock
	rawCombos := GenerateColorCombinations(cp.clues, cp.size, color)

	cp.mu.Lock()
	defer cp.mu.Unlock()

	// Double-check in case another goroutine generated it while we were waiting
	if cp.generated[color] {
		return cp.combosByColor[color], nil
	}

	// Convert []*big.Int to []*types.Bitset
	bitsets := make([]*types.Bitset, len(rawCombos))
	for i, combo := range rawCombos {
		bitsets[i] = types.NewBitset(combo)
	}

	cp.combosByColor[color] = bitsets
	cp.generated[color] = true

	return bitsets, nil
}

// GenerateColorCombinations enumerates combinations for a single color by
// projecting the multi-color clue line onto only the target color and treating
// other-color clues as fixed-length separators. This dramatically reduces the
// search space compared to enumerating all colors.
func GenerateColorCombinations(clues []types.ClueItem, size int, colorID int) []*big.Int {
	if size <= 0 || len(clues) == 0 {
		return []*big.Int{}
	}

	// Quick feasibility: minimal required cells across entire line
	// equals sum(lengths) + gaps between adjacent same-color clues.
	minRequired := 0
	for i := 0; i < len(clues); i++ {
		minRequired += clues[i].Clue
		if i > 0 && clues[i-1].ColorID == clues[i].ColorID {
			minRequired++
		}
	}
	if minRequired > size {
		return []*big.Int{}
	}

	// Collect target-color blocks (lengths and original indices)
	targetIdx := make([]int, 0)
	targetLen := make([]int, 0)
	for i := range clues {
		if clues[i].ColorID == colorID {
			targetIdx = append(targetIdx, i)
			targetLen = append(targetLen, clues[i].Clue)
		}
	}

	// If no target-color clues, there is exactly one mask: all zeros (if feasible)
	if len(targetIdx) == 0 {
		return []*big.Int{big.NewInt(0)}
	}

	m := len(targetIdx)

	// Compute required prefix, inter-target separators, and suffix
	// - prefix: minimal cells occupied before the first target block
	// - sep[k]: minimal cells between target block k and k+1
	// - suffix: minimal cells after the last target block
	prefix := 0
	for i := 0; i < targetIdx[0]; i++ {
		prefix += clues[i].Clue
		if i > 0 && clues[i-1].ColorID == clues[i].ColorID {
			prefix++
		}
	}

	sep := make([]int, m-1)
	for k := 0; k < m-1; k++ {
		a := targetIdx[k]
		b := targetIdx[k+1]
		if b == a+1 {
			// Adjacent target clues of the same color require a 1-cell gap
			sep[k] = 1
			continue
		}
		req := 0
		for i := a + 1; i < b; i++ {
			req += clues[i].Clue
			if i > a+1 && clues[i-1].ColorID == clues[i].ColorID {
				req++
			}
		}
		sep[k] = req
	}

	suffix := 0
	lastIdx := targetIdx[m-1]
	for i := lastIdx + 1; i < len(clues); i++ {
		suffix += clues[i].Clue
		if i > lastIdx+1 && clues[i-1].ColorID == clues[i].ColorID {
			suffix++
		}
	}

	// Earliest starts for each target block via forward pass
	earliest := make([]int, m)
	earliest[0] = prefix
	for k := 1; k < m; k++ {
		earliest[k] = earliest[k-1] + targetLen[k-1] + sep[k-1]
	}

	// Minimal tail requirement from k to end (inclusive)
	tailMin := make([]int, m+1)
	tailMin[m] = 0
	for k := m - 1; k >= 0; k-- {
		if k == m-1 {
			tailMin[k] = targetLen[k]
		} else {
			tailMin[k] = targetLen[k] + sep[k] + tailMin[k+1]
		}
	}

	// Latest starts for each target block via backward feasibility
	latest := make([]int, m)
	for k := 0; k < m; k++ {
		latest[k] = size - suffix - tailMin[k]
		if latest[k] < earliest[k] {
			// No feasible placement
			return []*big.Int{}
		}
	}

	// Helper to check room after placing block k at start s
	canPlace := func(k, s int) bool {
		// s + minimal cells from k to end must fit before size - suffix
		consumed := s - earliest[k] // shift relative doesn't matter for minimal tail
		_ = consumed                // not used; keep for clarity
		return s+tailMin[k] <= size-suffix
	}

	// DFS over target blocks only. Iterate start from min to max to yield
	// masks in descending numeric order (leftmost bits first).
	var dfs func(k int, prevStart int, mask *big.Int, out *[]*big.Int)
	dfs = func(k int, prevStart int, mask *big.Int, out *[]*big.Int) {
		if k == m {
			*out = append(*out, new(big.Int).Set(mask))
			return
		}

		minStart := earliest[k]
		if k > 0 {
			minStart = prevStart + targetLen[k-1] + sep[k-1]
			if minStart < earliest[k] {
				minStart = earliest[k]
			}
		}
		maxStart := latest[k]
		if maxStart < minStart {
			return
		}

		for s := minStart; s <= maxStart; s++ {
			if !canPlace(k, s) {
				continue
			}
			block := runMask(size, s, targetLen[k])
			mask.Or(mask, block)
			nextPrev := s
			dfs(k+1, nextPrev, mask, out)
			mask.AndNot(mask, block)
		}
	}

	// Top-level parallelization over the first block's start positions,
	// merging results deterministically in ascending start order to preserve
	// the final descending numeric ordering of masks.
	min0, max0 := earliest[0], latest[0]
	if min0 > max0 {
		return []*big.Int{}
	}
	choices := max0 - min0 + 1

	// Small ranges: run single-threaded for lower overhead
	if choices <= 3 {
		result := make([]*big.Int, 0)
		mask := big.NewInt(0)
		for s := min0; s <= max0; s++ {
			if !canPlace(0, s) {
				continue
			}
			block := runMask(size, s, targetLen[0])
			mask.Or(mask, block)
			dfs(1, s, mask, &result)
			mask.AndNot(mask, block)
		}
		return result
	}

	// Parallel path
	workers := runtime.GOMAXPROCS(0)
	if workers < 4 {
		workers = 4
	}
	if workers > choices {
		workers = choices
	}
	sem := make(chan struct{}, workers)
	parts := make([][]*big.Int, choices)
	var wg sync.WaitGroup

	for idx := 0; idx < choices; idx++ {
		s := min0 + idx
		if !canPlace(0, s) {
			// keep empty slice
			continue
		}
		wg.Add(1)
		sem <- struct{}{}
		go func(localIdx, start int) {
			defer wg.Done()
			defer func() { <-sem }()
			mask := big.NewInt(0)
			block := runMask(size, start, targetLen[0])
			mask.Or(mask, block)
			local := make([]*big.Int, 0)
			dfs(1, start, mask, &local)
			mask.AndNot(mask, block)
			parts[localIdx] = local
		}(idx, s)
	}
	wg.Wait()

	// Merge in ascending start order to maintain descending numeric order overall
	result := make([]*big.Int, 0)
	for i := 0; i < choices; i++ {
		if len(parts[i]) == 0 {
			continue
		}
		result = append(result, parts[i]...)
	}
	return result
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
