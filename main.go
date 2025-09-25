package main

import (
	"fmt"
	"os"
	"time"

	network "nonogram-solver/network"
	factory"nonogram-solver/factory"
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

	// Print the parsed data
	// nonogramData.Print()

	grid := factory.CreateGrid(*nonogramData)
	grid.Print()

}
