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

	var memStatsBefore runtime.MemStats
	runtime.ReadMemStats(&memStatsBefore)

	start := time.Now()
	grid, _ := network.FetchGrid(nonogramID)
	elapsed := time.Since(start)

	runtime.GC()
	var memStatsAfter runtime.MemStats
	runtime.ReadMemStats(&memStatsAfter)

	grid.Print()

	fmt.Printf("Grid created %dx%d \n", grid.Width(), grid.Height())
	fmt.Printf("Grid creation completed in %v\n", elapsed)
	fmt.Printf("Memory usage: %.2f MB (allocated), %.2f MB (total allocated)\n",
		float64(memStatsAfter.Alloc-memStatsBefore.Alloc)/1024/1024,
		float64(memStatsAfter.TotalAlloc-memStatsBefore.TotalAlloc)/1024/1024)

}
