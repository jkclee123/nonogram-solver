package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Provide a nonogram ID\n")
		return
	}
	nonogramID := os.Args[1]

	// Fetch the nonogram data
	data, err := FetchPage(nonogramID)
	if err != nil {
		fmt.Printf("Error fetching nonogram data: %v\n", err)
		return
	}

	fmt.Printf("Successfully fetched nonogram data %s \n", data)
	// TODO: Parse the data into NonogramData struct
}
