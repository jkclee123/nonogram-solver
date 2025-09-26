package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"
	"time"

	factory "nonogram-solver/internal/factory"
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

	// Minimum number of encoded fields required for decoding
	minEncodedFields = 4

	// Timeout for outbound requests
	requestTimeout = 10 * time.Second
)

var (
	dataRegex  = regexp.MustCompile(`var d=(\[.*?\]);`)
	httpClient = &http.Client{}
)

// FetchPage retrieves the HTML content for a nonogram by ID.
// It tries multiple URL patterns in parallel to accommodate different formats.
func FetchPage(nonogramID string) ([]byte, error) {
	if nonogramID == "" {
		return nil, fmt.Errorf("nonogramID cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	urls := []string{
		fmt.Sprintf(baseURLPattern2, nonogramID),
		fmt.Sprintf(baseURLPattern1, nonogramID),
	}

	resultChan := make(chan fetchResult, len(urls))

	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			fetchURL(ctx, url, resultChan)
		}(url)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		if result.err == nil {
			cancel()
			return result.body, nil
		}
	}

	return nil, fmt.Errorf("nonogram with ID %s not found on either URL pattern", nonogramID)
}

// fetchResult holds the result of a URL fetch attempt
type fetchResult struct {
	body []byte
	err  error
	url  string
}

// fetchURL attempts to fetch a URL and sends the result to the channel
func fetchURL(ctx context.Context, url string, resultChan chan<- fetchResult) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		resultChan <- fetchResult{err: fmt.Errorf("failed to create request for %s: %w", url, err), url: url}
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		resultChan <- fetchResult{err: fmt.Errorf("failed to fetch %s: %w", url, err), url: url}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		resultChan <- fetchResult{err: fmt.Errorf("HTTP %d from %s", resp.StatusCode, url), url: url}
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		resultChan <- fetchResult{err: fmt.Errorf("failed to read response body from %s: %w", url, err), url: url}
		return
	}

	resultChan <- fetchResult{body: body, url: url}
}

// Note: We no longer materialize NonogramData; we directly build Grid.

// calculateDimensions extracts width, height, and color count from raw data
func calculateDimensions(rawData [][]int) (width, height, numColors int, err error) {
	if len(rawData) <= colorsDataIndex {
		return 0, 0, 0, fmt.Errorf("insufficient data arrays for dimension calculation")
	}

	widthData := rawData[widthDataIndex]
	if len(widthData) < minEncodedFields {
		return 0, 0, 0, fmt.Errorf("width data array too short")
	}
	width = decodeDimension(widthData)

	heightData := rawData[heightDataIndex]
	if len(heightData) < minEncodedFields {
		return 0, 0, 0, fmt.Errorf("height data array too short")
	}
	height = decodeDimension(heightData)

	colorsData := rawData[colorsDataIndex]
	if len(colorsData) < minEncodedFields {
		return 0, 0, 0, fmt.Errorf("colors data array too short")
	}
	numColors = decodeDimension(colorsData)

	if width <= 0 || height <= 0 {
		return 0, 0, 0, fmt.Errorf("invalid grid dimensions: %dx%d", width, height)
	}
	if numColors <= 0 {
		return 0, 0, 0, fmt.Errorf("invalid number of colors: %d", numColors)
	}

	return width, height, numColors, nil
}

func decodeDimension(data []int) int {
	return data[dataValueIndex]%data[dataMultiplierIndex] +
		data[dataOffsetIndex]%data[dataMultiplierIndex] -
		data[dataModulusIndex]%data[dataMultiplierIndex]
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
		if len(colorData) < minEncodedFields {
			continue
		}

		colorMap[i+1] = decodeColor(colorBaseData, colorData, numColors)
	}

	return colorMap, nil
}

func decodeColor(colorBaseData, colorData []int, numColors int) string {
	if numColors == 1 {
		return "#000000"
	}

	r := clampColor(colorData[dataValueIndex] - colorBaseData[dataValueIndex])
	g := clampColor(colorData[dataOffsetIndex] - colorBaseData[dataValueIndex])
	b := clampColor(colorData[dataModulusIndex] - colorBaseData[dataMultiplierIndex])

	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

func clampColor(value int) int {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return value
}

// decodeGridCells populates the grid with cell data from the raw encoded format
func decodeGridCells(rawData [][]int, grid [][]int, width, height, numColors int) error {
	gridDataStart := colorOffset + numColors
	if gridDataStart >= len(rawData) {
		return fmt.Errorf("no grid data found after color data")
	}

	gridMetadata := rawData[gridDataStart]
	if len(gridMetadata) < minEncodedFields {
		return fmt.Errorf("grid metadata array too short")
	}

	gridDataCount := calculateGridDataCount(gridMetadata)

	gridOffsetData := rawData[gridDataStart+1]
	if len(gridOffsetData) < minEncodedFields {
		return fmt.Errorf("grid offset data array too short")
	}

	for i := 0; i < gridDataCount; i++ {
		dataIndex := gridDataStart + gridDataStartOffset + i
		if dataIndex >= len(rawData) {
			break
		}

		cellData := rawData[dataIndex]
		if len(cellData) < minEncodedFields {
			continue
		}

		startCol, endCol, color, row := decodeGridCell(cellData, gridOffsetData)
		if row < 0 || row >= height {
			continue
		}

		for col := startCol; col <= endCol; col++ {
			if col < 0 || col >= width {
				continue
			}
			grid[row][col] = color
		}
	}

	return nil
}

func calculateGridDataCount(gridMetadata []int) int {
	encodedValue := gridMetadata[dataValueIndex] % gridMetadata[dataMultiplierIndex]
	return encodedValue*encodedValue +
		(gridMetadata[dataOffsetIndex]%gridMetadata[dataMultiplierIndex])*2 +
		(gridMetadata[dataModulusIndex] % gridMetadata[dataMultiplierIndex])
}

func decodeGridCell(cellData, gridOffsetData []int) (startCol, endCol, color, row int) {
	startCol = cellData[dataValueIndex] - gridOffsetData[dataValueIndex] - 1
	colSpan := cellData[dataOffsetIndex] - gridOffsetData[dataOffsetIndex]
	endCol = startCol + colSpan - 1
	color = cellData[dataModulusIndex] - gridOffsetData[dataModulusIndex]
	row = cellData[dataMultiplierIndex] - gridOffsetData[dataMultiplierIndex] - 1
	return
}

// extractAllClues generates row and column clues from the populated grid
func extractAllClues(grid [][]int, width, height int) (map[types.LineID][]types.ClueItem, error) {
	if len(grid) != height {
		return nil, fmt.Errorf("grid height mismatch: expected %d, got %d", height, len(grid))
	}

	clues := make(map[types.LineID][]types.ClueItem)

	for rowIndex, row := range grid {
		if len(row) != width {
			return nil, fmt.Errorf("grid row %d width mismatch: expected %d, got %d", rowIndex, width, len(row))
		}
		lineID := types.LineID{Direction: types.Row, Index: rowIndex}
		clues[lineID] = extractCluesFromRow(row)
	}

	for col := 0; col < width; col++ {
		column := buildColumn(grid, col)
		lineID := types.LineID{Direction: types.Column, Index: col}
		clues[lineID] = extractCluesFromRow(column)
	}

	return clues, nil
}

func buildColumn(grid [][]int, col int) []int {
	column := make([]int, len(grid))
	for rowIndex := range grid {
		column[rowIndex] = grid[rowIndex][col]
	}
	return column
}

// extractCluesFromRow extracts clue numbers from a row or column of the grid.
// This function analyzes consecutive blocks of the same color and creates clue items.
func extractCluesFromRow(row []int) []types.ClueItem {
	var (
		clues        []types.ClueItem
		currentColor int
		currentCount int
	)

	for _, cell := range row {
		switch {
		case cell <= 0:
			if currentCount == 0 {
				continue
			}
			clues = append(clues, types.ClueItem{ColorID: currentColor, Clue: currentCount})
			currentColor = 0
			currentCount = 0
		case cell == currentColor:
			currentCount++
		default:
			if currentCount > 0 {
				clues = append(clues, types.ClueItem{ColorID: currentColor, Clue: currentCount})
			}
			currentColor = cell
			currentCount = 1
		}
	}

	if currentCount > 0 {
		clues = append(clues, types.ClueItem{ColorID: currentColor, Clue: currentCount})
	}

	return clues
}

// FetchNonogramData removed: use FetchGrid instead.

// FetchGrid fetches, parses, and constructs a Grid directly for the given nonogram ID.
// This bypasses exposing NonogramData to callers by internally converting clues to a Grid.
func FetchGrid(nonogramID string) (types.Grid, error) {
	if nonogramID == "" {
		return types.Grid{}, fmt.Errorf("nonogramID cannot be empty")
	}

	htmlContent, err := FetchPage(nonogramID)
	if err != nil {
		return types.Grid{}, fmt.Errorf("failed to fetch page for nonogram %s: %w", nonogramID, err)
	}

	if len(htmlContent) == 0 {
		return types.Grid{}, fmt.Errorf("HTML content is empty")
	}

	htmlStr := string(htmlContent)
	matches := dataRegex.FindStringSubmatch(htmlStr)
	if len(matches) < 2 {
		return types.Grid{}, fmt.Errorf("could not find nonogram data variable 'd' in HTML")
	}

	var rawData [][]int
	if err := json.Unmarshal([]byte(matches[1]), &rawData); err != nil {
		return types.Grid{}, fmt.Errorf("failed to parse nonogram data JSON: %w", err)
	}
	if len(rawData) == 0 {
		return types.Grid{}, fmt.Errorf("parsed nonogram data is empty")
	}

	width, height, numColors, err := calculateDimensions(rawData)
	if err != nil {
		return types.Grid{}, fmt.Errorf("failed to calculate dimensions: %w", err)
	}

	gridData := initializeGrid(width, height)
	colorMap, err := decodeColorData(rawData, numColors)
	if err != nil {
		return types.Grid{}, fmt.Errorf("failed to decode color data: %w", err)
	}
	if err := decodeGridCells(rawData, gridData, width, height, numColors); err != nil {
		return types.Grid{}, fmt.Errorf("failed to decode grid cells: %w", err)
	}
	clues, err := extractAllClues(gridData, width, height)
	if err != nil {
		return types.Grid{}, fmt.Errorf("failed to extract clues: %w", err)
	}

	grid := factory.CreateGridFromClues(clues, width, height, colorMap)
	return grid, nil
}
