package main

import (
	"fmt"
	"io"
	"net/http"
)

func FetchPage(nonogramID string) ([]byte, error) {
	// Try both URL patterns
	urls := []string{
		fmt.Sprintf("https://www.nonograms.org/nonograms/i/%s", nonogramID),
		fmt.Sprintf("https://www.nonograms.org/nonograms2/i/%s", nonogramID),
	}

	client := &http.Client{}

	for _, url := range urls {
		resp, err := client.Get(url)
		if err != nil {
			continue // Try next URL
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			// Read the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read response body: %v", err)
			}
			return body, nil
		}
	}

	return nil, fmt.Errorf("nonogram with ID %s not found on either URL pattern", nonogramID)
}
