package fetcher

import (
	"bytes"
	"io"
	"os"
	"testing"

	"nonogram-solver/types"
)

// captureOutput captures stdout during the execution of a function
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// TestFetchNonogramData77943 tests fetching data for nonogram ID 77943
func TestFetchNonogramData77943(t *testing.T) {
	data, err := FetchNonogramData("77943")
	if err != nil {
		t.Fatalf("Failed to fetch nonogram data for ID 77943: %v", err)
	}

	// Verify basic structure
	if data == nil {
		t.Fatal("NonogramData is nil")
	}

	if data.Width <= 0 {
		t.Errorf("Invalid width: %d", data.Width)
	}

	if data.Height <= 0 {
		t.Errorf("Invalid height: %d", data.Height)
	}

	if data.Clues == nil {
		t.Fatal("Clues map is nil")
	}

	// Verify we have clues for all rows and columns
	expectedRowCount := data.Height
	expectedColCount := data.Width

	rowCount := 0
	colCount := 0

	for lineID := range data.Clues {
		switch lineID.Direction {
		case types.Row:
			rowCount++
		case types.Column:
			colCount++
		}
	}

	if rowCount != expectedRowCount {
		t.Errorf("Expected %d rows, got %d", expectedRowCount, rowCount)
	}

	if colCount != expectedColCount {
		t.Errorf("Expected %d columns, got %d", expectedColCount, colCount)
	}

	// Test complete print output
	output := captureOutput(func() {
		data.Print()
	})

	// Compare with expected output
	if output != expectedOutput77943 {
		t.Errorf("Print output for nonogram 77943 does not match expected output")
		t.Logf("Expected:\n%s", expectedOutput77943)
		t.Logf("Actual:\n%s", output)
	}

	t.Logf("Successfully fetched and validated nonogram 77943 (%dx%d)", data.Width, data.Height)
}

// TestFetchNonogramData77891 tests fetching data for nonogram ID 77891
func TestFetchNonogramData77891(t *testing.T) {
	data, err := FetchNonogramData("77891")
	if err != nil {
		t.Fatalf("Failed to fetch nonogram data for ID 77891: %v", err)
	}

	// Verify basic structure
	if data == nil {
		t.Fatal("NonogramData is nil")
	}

	if data.Width <= 0 {
		t.Errorf("Invalid width: %d", data.Width)
	}

	if data.Height <= 0 {
		t.Errorf("Invalid height: %d", data.Height)
	}

	if data.Clues == nil {
		t.Fatal("Clues map is nil")
	}

	// Verify we have clues for all rows and columns
	expectedRowCount := data.Height
	expectedColCount := data.Width

	rowCount := 0
	colCount := 0

	for lineID := range data.Clues {
		switch lineID.Direction {
		case types.Row:
			rowCount++
		case types.Column:
			colCount++
		}
	}

	if rowCount != expectedRowCount {
		t.Errorf("Expected %d rows, got %d", expectedRowCount, rowCount)
	}

	if colCount != expectedColCount {
		t.Errorf("Expected %d columns, got %d", expectedColCount, colCount)
	}

	// Test complete print output
	output := captureOutput(func() {
		data.Print()
	})

	// Compare with expected output
	if output != expectedOutput77891 {
		t.Errorf("Print output for nonogram 77891 does not match expected output")
		t.Logf("Expected:\n%s", expectedOutput77891)
		t.Logf("Actual:\n%s", output)
	}

	t.Logf("Successfully fetched and validated nonogram 77891 (%dx%d)", data.Width, data.Height)
}

// TestFetchNonogramDataInvalidID tests error handling for invalid IDs
func TestFetchNonogramDataInvalidID(t *testing.T) {
	_, err := FetchNonogramData("")
	if err == nil {
		t.Error("Expected error for empty ID, got nil")
	}

	_, err = FetchNonogramData("invalid")
	if err == nil {
		t.Error("Expected error for invalid ID, got nil")
	}
}

// Expected output for nonogram 77943
const expectedOutput77943 = `
=== ROW CLUES ===
Row 1: (1,2)
Row 2: (1,4)
Row 3: (1,4), (1,3)
Row 4: (1,1), (1,4), (1,1)
Row 5: (1,1), (1,1)
Row 6: (1,3), (1,1), (1,3)
Row 7: (1,1), (1,2), (1,1), (1,2)
Row 8: (1,1), (1,1)
Row 9: (1,1), (1,1)
Row 10: (1,6)

=== COLUMN CLUES ===
Col 1: (1,3)
Col 2: (1,1), (1,4)
Col 3: (1,1), (1,1), (1,1)
Col 4: (1,3), (1,1), (1,1)
Col 5: (1,4), (1,1), (1,1)
Col 6: (1,2), (1,1), (1,1), (1,1)
Col 7: (1,3), (1,1), (1,1)
Col 8: (1,1), (1,1), (1,1)
Col 9: (1,1), (1,4)
Col 10: (1,4)

=== COLOR MAP ===
Color 1: #000000

Grid size: 10x10
`

// Expected output for nonogram 77891
const expectedOutput77891 = `
=== ROW CLUES ===
Row 1: (3,2), (1,1)
Row 2: (1,1), (3,1), (4,1), (3,1), (1,1)
Row 3: (3,4)
Row 4: (3,5)
Row 5: (1,1), (3,5)
Row 6: (1,1), (3,6)
Row 7: (2,1), (2,1), (1,1), (3,2)
Row 8: (2,2), (1,2)

=== COLUMN CLUES ===
Col 1: (1,1), (2,2)
Col 2: (3,3), (2,1)
Col 3: (3,1), (4,1), (3,2), (1,1), (2,1)
Col 4: (1,1), (3,4), (1,1), (1,1)
Col 5: (1,1), (3,4), (1,2)
Col 6: (3,4)
Col 7: (3,4)
Col 8: (3,2)
Col 9: (3,1)
Col 10: (3,1)

=== COLOR MAP ===
Color 1: #B5A61F
Color 2: #1FB54C
Color 3: #EBD728
Color 4: #000000

Grid size: 10x8
`
