package test

import (
	"math/big"
	"reflect"
	"testing"

	"nonogram-solver/internal/factory"
	"nonogram-solver/types"
)

func TestGenerateColorCombinations(t *testing.T) {
	tests := []struct {
		name     string
		clues    []types.ClueItem
		size     int
		colorID  int
		expected []*big.Int
	}{
		{
			name: "example 1: (1,1), (2,2), Size: 8, ColorID: 1",
			clues: []types.ClueItem{
				{ColorID: 1, Clue: 1},
				{ColorID: 2, Clue: 2},
			},
			size:    8,
			colorID: 1,
			expected: []*big.Int{
				big.NewInt(128),
				big.NewInt(64),
				big.NewInt(32),
				big.NewInt(16),
				big.NewInt(8),
				big.NewInt(4),
			},
		},
		{
			name: "example 1: (1,1), (2,2), Size: 8, ColorID: 2",
			clues: []types.ClueItem{
				{ColorID: 1, Clue: 1},
				{ColorID: 2, Clue: 2},
			},
			size:    8,
			colorID: 2,
			expected: []*big.Int{
				big.NewInt(96),
				big.NewInt(48),
				big.NewInt(24),
				big.NewInt(12),
				big.NewInt(6),
				big.NewInt(3),
			},
		},
		{
			name: "example 2: (1,1), (3,4), (1,2), Size: 8, ColorID 1",
			clues: []types.ClueItem{
				{ColorID: 1, Clue: 1},
				{ColorID: 3, Clue: 4},
				{ColorID: 1, Clue: 2},
			},
			size:    8,
			colorID: 1,
			expected: []*big.Int{
				big.NewInt(134),
				big.NewInt(131),
				big.NewInt(67),
			},
		},
		{
			name: "example 2: (1,1), (3,4), (1,2), Size: 8, ColorID 3",
			clues: []types.ClueItem{
				{ColorID: 1, Clue: 1},
				{ColorID: 3, Clue: 4},
				{ColorID: 1, Clue: 2},
			},
			size:    8,
			colorID: 3,
			expected: []*big.Int{
				big.NewInt(120),
				big.NewInt(60),
			},
		},
		{
			name: "example 3: (1,4), (1,3), Size: 10, ColorID 1",
			clues: []types.ClueItem{
				{ColorID: 1, Clue: 4},
				{ColorID: 1, Clue: 3},
			},
			size:    10,
			colorID: 1,
			expected: []*big.Int{
				big.NewInt(988),
				big.NewInt(974),
				big.NewInt(967),
				big.NewInt(494),
				big.NewInt(487),
				big.NewInt(247),
			},
		},
		{
			name: "example 4: (1,2), (1,1), (1,1), (1,1), Size: 10, ColorID 1",
			clues: []types.ClueItem{
				{ColorID: 1, Clue: 2},
				{ColorID: 1, Clue: 1},
				{ColorID: 1, Clue: 1},
				{ColorID: 1, Clue: 1},
			},
			size:    10,
			colorID: 1,
			expected: []*big.Int{
				big.NewInt(852),
				big.NewInt(850),
				big.NewInt(849),
				big.NewInt(842),
				big.NewInt(841),
				big.NewInt(837),
				big.NewInt(810),
				big.NewInt(809),
				big.NewInt(805),
				big.NewInt(789),
				big.NewInt(426),
				big.NewInt(425),
				big.NewInt(421),
				big.NewInt(405),
				big.NewInt(213),
			},
		},
		{
			name: "example 5: (3,1), (4,1), (3,2), (1,1), (2,1), Size: 8, ColorID 1",
			clues: []types.ClueItem{
				{ColorID: 3, Clue: 1},
				{ColorID: 4, Clue: 1},
				{ColorID: 3, Clue: 2},
				{ColorID: 1, Clue: 1},
				{ColorID: 2, Clue: 1},
			},
			size:    8,
			colorID: 1,
			expected: []*big.Int{
				big.NewInt(8),
				big.NewInt(4),
				big.NewInt(2),
			},
		},
		{
			name: "example 5: (3,1), (4,1), (3,2), (1,1), (2,1), Size: 8, ColorID: 2",
			clues: []types.ClueItem{
				{ColorID: 3, Clue: 1},
				{ColorID: 4, Clue: 1},
				{ColorID: 3, Clue: 2},
				{ColorID: 1, Clue: 1},
				{ColorID: 2, Clue: 1},
			},
			size:    8,
			colorID: 2,
			expected: []*big.Int{
				big.NewInt(4),
				big.NewInt(2),
				big.NewInt(1),
			},
		},
		{
			name: "example 5: (3,1), (4,1), (3,2), (1,1), (2,1), Size: 8, ColorID: 3",
			clues: []types.ClueItem{
				{ColorID: 3, Clue: 1},
				{ColorID: 4, Clue: 1},
				{ColorID: 3, Clue: 2},
				{ColorID: 1, Clue: 1},
				{ColorID: 2, Clue: 1},
			},
			size:    8,
			colorID: 3,
			expected: []*big.Int{
				big.NewInt(176),
				big.NewInt(152),
				big.NewInt(140),
				big.NewInt(88),
				big.NewInt(76),
				big.NewInt(44),
			},
		},
		{
			name: "example 5: (3,1), (4,1), (3,2), (1,1), (2,1), Size: 8, ColorID: 4",
			clues: []types.ClueItem{
				{ColorID: 3, Clue: 1},
				{ColorID: 4, Clue: 1},
				{ColorID: 3, Clue: 2},
				{ColorID: 1, Clue: 1},
				{ColorID: 2, Clue: 1},
			},
			size:    8,
			colorID: 4,
			expected: []*big.Int{
				big.NewInt(64),
				big.NewInt(32),
				big.NewInt(16),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := factory.GenerateColorCombinations(tt.clues, tt.size)
			colorCombinations := result[tt.colorID]

			if !reflect.DeepEqual(colorCombinations, tt.expected) {
				t.Errorf("GenerateColorCombinations() for color %d = %v, want %v", tt.colorID, colorCombinations, tt.expected)
			}
		})
	}
}
