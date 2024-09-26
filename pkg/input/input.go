package input

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// HandleInput reads words from stdin or a file and sends them through a channel
func HandleInput(inputFilePath string) (<-chan string, error) {
	wordChan := make(chan string)

	var reader io.Reader
	if inputFilePath != "" {
		file, err := os.Open(inputFilePath)
		if err != nil {
			return nil, fmt.Errorf("could not open file: %w", err)
		}
		reader = file
	} else {
		reader = os.Stdin
	}

	// Start a goroutine to scan input asynchronously
	go func() {
		defer close(wordChan) // Close the channel once done
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			wordChan <- scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "error reading input: %v", err)
		}
	}()

	return wordChan, nil
}
