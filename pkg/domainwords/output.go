package domainwords

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

var (
	outputWG sync.WaitGroup
)

// HandleOutput manages the output of the words being processed.
func HandleOutput(outputChan <-chan string, filePath string, sync bool) {
	outputWG.Add(1)

	var file *os.File
	var isFileOutput = false

	if filePath != "" {

		isFileOutput = true

		// Get the directory from the file path
		dir := filepath.Dir(filePath)

		// Create the directory (and all parent directories) if they don't exist
		err := os.MkdirAll(dir, os.ModePerm) // os.ModePerm is 0777 by default
		if err != nil {
			log.Fatalf("Error creating directories: %v", err)
			return
		}

		// Open the file for writing
		file, err = os.Create(filePath)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
			return
		}
	}

	go func() {
		var results []string

		defer func() {
			// Close the file when the goroutine exits
			if isFileOutput {
				if err := file.Close(); err != nil {
					log.Printf("Error closing file: %v", err)
				}
			}

			outputWG.Done()
		}()

		for permWord := range outputChan {
			// Print to console (can be removed if not needed)
			fmt.Println(permWord)

			// Write to file
			if isFileOutput {
				if _, err := fmt.Fprintln(file, permWord); err != nil {
					log.Printf("Error writing to file: %v", err)
					return
				}
			}

			results = append(results, permWord)
		}

		// Sort and sync results to GitHub
		if sync {
			sort.Strings(results)
			if err := SyncResultsToGitHub(results); err != nil {
				log.Printf("Error syncing results to GitHub: %v", err)
			}
		}
	}()
}

// WaitOutputCompletion waits for output to be fully processed.
func WaitOutputCompletion() {
	outputWG.Wait()
}
