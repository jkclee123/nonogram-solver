package main

import (
	"fmt"
	"os"
	"time"

	factory "nonogram-solver/factory"
	network "nonogram-solver/network"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Provide a nonogram ID\n")
		return
	}
	nonogramID := os.Args[1]

	// Fetch and parse the nonogram data
	start := time.Now()
	nonogramData, err := network.FetchNonogramData(nonogramID)
	elapsed := time.Since(start)
	fmt.Printf("Fetched nonogram data in %v\n", elapsed)
	if err != nil {
		fmt.Printf("Error fetching nonogram data: %v\n", err)
		return
	}

	start = time.Now()
	lines := factory.CreateLines(*nonogramData)
	elapsed = time.Since(start)
	// fmt.Println("\n=== GENERATED LINES ===")
	// lines.Print()
	fmt.Printf("Created %d lines in %v\n", len(lines.Lines), elapsed)
}
