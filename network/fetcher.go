package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"nonogram-solver/types"
)

// Constants for URL patterns and data structure indices
const (
	// URL patterns for fetching nonogram data
	baseURLPattern1 = "https://www.nonograms.org/nonograms/i/%s"
	baseURLPattern2 = "https://www.nonograms.org/nonograms2/i/%s"

	// Data structure indices based on JavaScript implementation
	widthDataIndex      = 1
	heightDataIndex     = 2
	colorsDataIndex     = 3
	gridDataIndex       = 4
	colorOffset         = 5
	gridDataStartOffset = 2

	// Array indices within data structures
	dataValueIndex      = 0
	dataOffsetIndex     = 1
	dataModulusIndex    = 2
	dataMultiplierIndex = 3

	// Color calculation constants
	colorBaseValue = 256
)

// FetchPage retrieves the HTML content for a nonogram by ID
// It tries multiple URL patterns to accommodate different nonogram formats
func FetchPage(nonogramID string) ([]byte, error) {
	if nonogramID == "" {
		return nil, fmt.Errorf("nonogramID cannot be empty")
	}

	// Try both URL patterns in case the nonogram uses a different format
	urls := []string{
		fmt.Sprintf(baseURLPattern1, nonogramID),
		fmt.Sprintf(baseURLPattern2, nonogramID),
	}

	client := &http.Client{}

	for _, url := range urls {
		resp, err := client.Get(url)
		if err != nil {
			continue // Try next URL if this one fails
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read response body from %s: %w", url, err)
			}
			return body, nil
		}
	}

	return nil, fmt.Errorf("nonogram with ID %s not found on either URL pattern", nonogramID)
}

// ParseNonogramData extracts and decodes the nonogram clues from the HTML
func ParseNonogramData(htmlContent []byte) (*types.NonogramData, error) {
	if len(htmlContent) == 0 {
		return nil, fmt.Errorf("HTML content is empty")
	}

	htmlStr := string(htmlContent)

	// Extract the JavaScript variable 'd' containing the nonogram data
	dataRegex := regexp.MustCompile(`var d=(\[.*?\]);`)
	matches := dataRegex.FindStringSubmatch(htmlStr)
	if len(matches) < 2 {
		return nil, fmt.Errorf("could not find nonogram data variable 'd' in HTML")
	}

	// Parse the JSON array containing the encoded nonogram data
	var rawData [][]int
	if err := json.Unmarshal([]byte(matches[1]), &rawData); err != nil {
		return nil, fmt.Errorf("failed to parse nonogram data JSON: %w", err)
	}

	if len(rawData) == 0 {
		return nil, fmt.Errorf("parsed nonogram data is empty")
	}

	// Decode the raw data into structured nonogram format
	return decodeNonogramData(rawData)
}

// decodeNonogramData decodes the raw nonogram data from JavaScript format into structured data
func decodeNonogramData(rawData [][]int) (*types.NonogramData, error) {
	if len(rawData) < 4 {
		return nil, fmt.Errorf("insufficient raw data for nonogram decoding (need at least 4 arrays)")
	}

	// Calculate grid dimensions and color count
	width, height, numColors, err := calculateDimensions(rawData)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate dimensions: %w", err)
	}

	fmt.Printf("Grid dimensions: %dx%d, Colors: %d\n", width, height, numColors)

	// Initialize empty grid
	grid := initializeGrid(width, height)

	// Decode color information
	colorMap, err := decodeColorData(rawData, numColors)
	if err != nil {
		return nil, fmt.Errorf("failed to decode color data: %w", err)
	}

	// Decode grid cell data
	if err := decodeGridCells(rawData, grid, width, height, numColors); err != nil {
		return nil, fmt.Errorf("failed to decode grid cells: %w", err)
	}

	// Extract clues from the populated grid
	rowClues, colClues, err := extractAllClues(grid, width, height)
	if err != nil {
		return nil, fmt.Errorf("failed to extract clues: %w", err)
	}

	// Combine row and column clues into a single slice
	allClues := make([]types.LineClue, 0, len(rowClues)+len(colClues))
	allClues = append(allClues, rowClues...)
	allClues = append(allClues, colClues...)

	return &types.NonogramData{
		Clues:    allClues,
		Width:    width,
		Height:   height,
		ColorMap: colorMap,
	}, nil
}

// calculateDimensions extracts width, height, and color count from raw data
func calculateDimensions(rawData [][]int) (width, height, numColors int, err error) {
	if len(rawData) <= colorsDataIndex {
		return 0, 0, 0, fmt.Errorf("insufficient data arrays for dimension calculation")
	}

	// Calculate width using the formula from JavaScript: D = d[1][0]%d[1][3] + d[1][1]%d[1][3] - d[1][2]%d[1][3]
	widthData := rawData[widthDataIndex]
	if len(widthData) < 4 {
		return 0, 0, 0, fmt.Errorf("width data array too short")
	}
	width = widthData[dataValueIndex]%widthData[dataMultiplierIndex] +
		widthData[dataOffsetIndex]%widthData[dataMultiplierIndex] -
		widthData[dataModulusIndex]%widthData[dataMultiplierIndex]

	// Calculate height using the formula from JavaScript: C = d[2][0]%d[2][3] + d[2][1]%d[2][3] - d[2][2]%d[2][3]
	heightData := rawData[heightDataIndex]
	if len(heightData) < 4 {
		return 0, 0, 0, fmt.Errorf("height data array too short")
	}
	height = heightData[dataValueIndex]%heightData[dataMultiplierIndex] +
		heightData[dataOffsetIndex]%heightData[dataMultiplierIndex] -
		heightData[dataModulusIndex]%heightData[dataMultiplierIndex]

	// Calculate number of colors using the formula from JavaScript: Aa = d[3][0]%d[3][3] + d[3][1]%d[3][3] - d[3][2]%d[3][3]
	colorsData := rawData[colorsDataIndex]
	if len(colorsData) < 4 {
		return 0, 0, 0, fmt.Errorf("colors data array too short")
	}
	numColors = colorsData[dataValueIndex]%colorsData[dataMultiplierIndex] +
		colorsData[dataOffsetIndex]%colorsData[dataMultiplierIndex] -
		colorsData[dataModulusIndex]%colorsData[dataMultiplierIndex]

	// Validate calculated dimensions
	if width <= 0 || height <= 0 {
		return 0, 0, 0, fmt.Errorf("invalid grid dimensions: %dx%d", width, height)
	}
	if numColors <= 0 {
		return 0, 0, 0, fmt.Errorf("invalid number of colors: %d", numColors)
	}

	return width, height, numColors, nil
}

// initializeGrid creates an empty 2D grid with the specified dimensions
func initializeGrid(width, height int) [][]int {
	grid := make([][]int, height)
	for i := range grid {
		grid[i] = make([]int, width)
	}
	return grid
}

// decodeColorData extracts color information from the raw data
func decodeColorData(rawData [][]int, numColors int) (map[int]string, error) {
	if len(rawData) < colorOffset+numColors {
		return nil, fmt.Errorf("insufficient data for %d colors", numColors)
	}

	colorMap := make(map[int]string)
	colorBaseData := rawData[gridDataIndex]

	for i := 0; i < numColors; i++ {
		colorIndex := colorOffset + i
		if colorIndex >= len(rawData) {
			break
		}

		colorData := rawData[colorIndex]
		if len(colorData) < 4 {
			continue
		}

		// Decode color value based on JavaScript implementation
		colorValueOffset := colorData[dataValueIndex] - colorBaseData[dataOffsetIndex]
		colorValue := colorValueOffset + colorBaseValue

		// Convert RGB value to hex format
		// Color value is typically in the format 0xRRGGBB
		r := (colorValue >> 16) & 0xFF
		g := (colorValue >> 8) & 0xFF
		b := colorValue & 0xFF

		hexColor := fmt.Sprintf("#%02X%02X%02X", r, g, b)
		colorMap[i+1] = hexColor // Color IDs start from 1
	}

	return colorMap, nil
}

// decodeGridCells populates the grid with cell data from the raw encoded format
func decodeGridCells(rawData [][]int, grid [][]int, width, height, numColors int) error {
	gridDataStart := colorOffset + numColors
	if gridDataStart >= len(rawData) {
		return fmt.Errorf("no grid data found after color data")
	}

	// Extract grid metadata
	gridMetadata := rawData[gridDataStart]
	if len(gridMetadata) < 4 {
		return fmt.Errorf("grid metadata array too short")
	}

	// Calculate number of grid data entries
	gridDataCount := gridMetadata[dataValueIndex]%gridMetadata[dataMultiplierIndex]*(gridMetadata[dataValueIndex]%gridMetadata[dataMultiplierIndex]) +
		gridMetadata[dataOffsetIndex]%gridMetadata[dataMultiplierIndex]*2 +
		gridMetadata[dataModulusIndex]%gridMetadata[dataMultiplierIndex]

	// Extract grid offset data
	gridOffsetData := rawData[gridDataStart+1]
	if len(gridOffsetData) < 4 {
		return fmt.Errorf("grid offset data array too short")
	}

	// Decode each grid cell entry
	for i := 0; i < gridDataCount; i++ {
		dataIndex := gridDataStart + gridDataStartOffset + i
		if dataIndex >= len(rawData) {
			break
		}

		cellData := rawData[dataIndex]
		if len(cellData) < 4 {
			continue
		}

		// Calculate actual grid position and color
		startCol := cellData[dataValueIndex] - gridOffsetData[dataValueIndex] - 1
		colSpan := cellData[dataOffsetIndex] - gridOffsetData[dataOffsetIndex]
		endCol := startCol + colSpan - 1
		color := cellData[dataModulusIndex] - gridOffsetData[dataModulusIndex]
		row := cellData[dataMultiplierIndex] - gridOffsetData[dataMultiplierIndex] - 1

		// Fill the grid cells within bounds
		if row >= 0 && row < height {
			for col := startCol; col <= endCol && col >= 0 && col < width; col++ {
				grid[row][col] = color
			}
		}
	}

	return nil
}

// extractAllClues generates row and column clues from the populated grid
func extractAllClues(grid [][]int, width, height int) ([]types.LineClue, []types.LineClue, error) {
	if len(grid) != height {
		return nil, nil, fmt.Errorf("grid height mismatch: expected %d, got %d", height, len(grid))
	}

	// Extract row clues
	rowClues := make([]types.LineClue, height)
	for row := 0; row < height; row++ {
		if len(grid[row]) != width {
			return nil, nil, fmt.Errorf("grid row %d width mismatch: expected %d, got %d", row, width, len(grid[row]))
		}
		clues := extractCluesFromRow(grid[row])
		rowClues[row] = types.LineClue{
			Clues: clues,
			LineID: types.LineID{
				Direction: types.Row,
				Index:     uint8(row),
			},
		}
	}

	// Extract column clues
	colClues := make([]types.LineClue, width)
	for col := 0; col < width; col++ {
		column := make([]int, height)
		for row := 0; row < height; row++ {
			column[row] = grid[row][col]
		}
		clues := extractCluesFromRow(column)
		colClues[col] = types.LineClue{
			Clues: clues,
			LineID: types.LineID{
				Direction: types.Column,
				Index:     uint8(col),
			},
		}
	}

	return rowClues, colClues, nil
}

// extractCluesFromRow extracts clue numbers from a row or column of the grid
// This function analyzes consecutive blocks of the same color and creates clue items
func extractCluesFromRow(row []int) []types.ClueItem {
	var clues []types.ClueItem

	currentColor := 0
	currentCount := 0

	for _, cell := range row {
		if cell > 0 { // Cell contains a color (non-empty)
			if cell == currentColor {
				// Continue counting consecutive cells of the same color
				currentCount++
			} else {
				// Color changed, save previous block if it exists
				if currentCount > 0 {
					clues = append(clues, types.ClueItem{
						ColorID:   uint8(currentColor),
						BlockSize: uint8(currentCount),
					})
				}
				// Start new block
				currentColor = cell
				currentCount = 1
			}
		} else { // Empty cell (background)
			if currentCount > 0 {
				// End current block
				clues = append(clues, types.ClueItem{
					ColorID:   uint8(currentColor),
					BlockSize: uint8(currentCount),
				})
				// Reset for next block
				currentCount = 0
				currentColor = 0
			}
		}
	}

	// Add the final block if one was in progress
	if currentCount > 0 {
		clues = append(clues, types.ClueItem{
			ColorID:   uint8(currentColor),
			BlockSize: uint8(currentCount),
		})
	}

	return clues
}

// FetchNonogramData fetches and parses the nonogram data for the given ID
// This is the main entry point for retrieving nonogram data from the website
func FetchNonogramData(nonogramID string) (*types.NonogramData, error) {
	if nonogramID == "" {
		return nil, fmt.Errorf("nonogramID cannot be empty")
	}

	htmlContent, err := FetchPage(nonogramID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page for nonogram %s: %w", nonogramID, err)
	}

	nonogramData, err := ParseNonogramData(htmlContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse nonogram data for %s: %w", nonogramID, err)
	}

	return nonogramData, nil
}
