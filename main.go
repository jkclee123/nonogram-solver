package main

import (
	"fmt"
	"os"

	fetcher "nonogram-solver/internal"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Provide a nonogram ID\n")
		return
	}
	nonogramID := os.Args[1]

	// Fetch and parse the nonogram data
	data, err := fetcher.FetchNonogramData(nonogramID)
	if err != nil {
		fmt.Printf("Error fetching nonogram data: %v\n", err)
		return
	}

	fmt.Printf("Successfully fetched and parsed nonogram data!\n")
	data.Print()
}
