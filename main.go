package main

import (
	"fmt"
	"os"

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
	nonogramData, err := network.FetchNonogramData(nonogramID)
	if err != nil {
		fmt.Printf("Error fetching nonogram data: %v\n", err)
		return
	}

	lines := factory.CreateLines(*nonogramData)
	factory.GenerateCombinationsForLines(&lines)
	nonogramData.Print()
	fmt.Println("\n=== GENERATED LINES ===")
	lines.Print()
}
