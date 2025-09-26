package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	network "nonogram-solver/internal/network"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Provide a nonogram ID\n")
		return
	}
	nonogramID := os.Args[1]
	grid, err := network.FetchGrid(nonogramID)
	if err != nil {
		fmt.Printf("Error fetching nonogram grid: %v\n", err)
		return
	}

	// Print the parsed data
	// nonogramData.Print()

	// Force garbage collection to get accurate memory stats
	runtime.GC()
	var memStatsBefore runtime.MemStats
	runtime.ReadMemStats(&memStatsBefore)

	start := time.Now()
	// Grid is already constructed by the fetcher; keep timing to report end-to-end
	elapsed := time.Since(start)

	// Force garbage collection and get final memory stats
	runtime.GC()
	var memStatsAfter runtime.MemStats
	runtime.ReadMemStats(&memStatsAfter)

	// grid.Print()
	fmt.Printf("Grid created %dx%d \n", grid.Width, grid.Height)
	fmt.Printf("Grid creation completed in %v\n", elapsed)
	fmt.Printf("Memory usage: %.2f MB (allocated), %.2f MB (total allocated)\n",
		float64(memStatsAfter.Alloc-memStatsBefore.Alloc)/1024/1024,
		float64(memStatsAfter.TotalAlloc-memStatsBefore.TotalAlloc)/1024/1024)

}
