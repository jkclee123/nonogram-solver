package test

import (
	"testing"

	"nonogram-solver/internal/factory"
	"nonogram-solver/internal/types"
)

func TestCombinationsProvider(t *testing.T) {
	// Test data: simple line with one color clue
	clues := []types.ClueItem{
		{ColorID: 1, Clue: 2},
	}
	size := 5

	provider := factory.NewCombinationsProvider(clues, size)

	// Test getting combinations for color 1
	combos, err := provider.Get(1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(combos) == 0 {
		t.Fatal("Expected at least one combination, got none")
	}

	// Test caching - second call should return same result
	combos2, err := provider.Get(1)
	if err != nil {
		t.Fatalf("Expected no error on second call, got %v", err)
	}

	if len(combos) != len(combos2) {
		t.Fatalf("Expected same number of combinations, got %d vs %d", len(combos), len(combos2))
	}

	// Test different color (should return empty for color 2)
	combos3, err := provider.Get(2)
	if err != nil {
		t.Fatalf("Expected no error for color 2, got %v", err)
	}

	if len(combos3) != 1 {
		t.Fatalf("Expected exactly one empty combination for color 2, got %d", len(combos3))
	}
}
