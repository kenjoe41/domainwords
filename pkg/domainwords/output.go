package domainwords

import (
	"fmt"
	"os"
	"sort"
	"sync"
)

var (
	outputWG sync.WaitGroup
)

func HandleOutput(outputChan <-chan string, filePath string) {
	outputWG.Add(1)

	// Open the file for writing
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file: %v\n", err)
		return
	}

	go func() {
		var results []string // Slice to store results

		defer func() {
			// Close the file when the goroutine exits
			if err := file.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "error closing file: %v\n", err)
			}
			outputWG.Done()

		}()

		for permWord := range outputChan {
			// Print to console
			fmt.Println(permWord)

			// Write to file
			if _, err := fmt.Fprintln(file, permWord); err != nil {
				fmt.Fprintf(os.Stderr, "error writing to file: %v\n", err)
				return // Stop writing to file if an error occurs
			}
			// Store result in the slice
			results = append(results, permWord)
		}

		sort.Strings(results)
		// fmt.Println("Preparing to sync to Github.")
		if err := syncResultsToGitHub(results); err != nil {
			fmt.Fprintf(os.Stderr, "error syncing file to GitHub: %v\n", err)
		}
	}()

}

func WaitOutputCompletion() {
	outputWG.Wait()
}
