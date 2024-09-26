package main

import (
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/kenjoe41/domainwords/pkg/domainwords"
	"github.com/kenjoe41/domainwords/pkg/input"
	"github.com/kenjoe41/domainwords/pkg/options"
	"github.com/kenjoe41/domainwords/pkg/output"
	"github.com/kenjoe41/domainwords/pkg/wordprocess"
)

func main() {
	// Create a logger
	logger := log.New(os.Stderr, "[DomainWords] ", log.LstdFlags)

	// Parsing command-line flags
	flags := options.ScanFlags()

	// Handle input (stdin or wordlist)
	var err error
	logger.Println("Reading input from file or stdin...")
	wordChan, err := input.HandleInput(flags.Wordlist)
	if err != nil {
		logger.Fatalf("Error handling input: %v", err)
	}

	// Set up output channel and buffer size dynamically
	outputChan := make(chan string, calculateBufferSize(flags.Wordlist))

	// Handle output in a separate goroutine
	var outputWG sync.WaitGroup
	outputWG.Add(1)
	go func() {
		if err := output.HandleOutput(outputChan, flags.OutputFile); err != nil {
			logger.Fatalf("Error writing output: %v", err)
		}
		outputWG.Done()
	}()

	// Sort words and remove duplicates
	logger.Println("Removing duplicates from input wordlist")
	words := domainwords.RemoveDuplicateStr(wordChan)

	// Configure depth based on flags
	logger.Println("Performing necessary configurations.")
	depth := domainwords.ConfigureDepth(flags.Level)

	// Check if iterations should be minimized
	if depth == 1 || len(words) < int(flags.ChunkSize) {
		flags.Iterations = 1
	}

	for iter := 0; iter < int(flags.Iterations); iter++ {
		logger.Printf("Iteration %d\n", iter+1)

		// Shuffle the words
		logger.Println("Performing Chaotic Shuffle of the input wordlist.")
		chaoticSlice := domainwords.ChaoticShuffle(words)

		// Break into chunks
		logger.Println("Breaking words into chunks.")
		tempChunks := domainwords.ChunkSlice(chaoticSlice, int(flags.ChunkSize))

		//clear chaoticSlice
		chaoticSlice = nil

		// Write shuffled words to temp files
		tempWordsFiles := domainwords.WriteTempChunks(tempChunks)
		if len(tempWordsFiles) == 0 {
			logger.Fatalf("Errors while writing chunk files.")
		}
		// Clear TempChunks
		tempChunks = nil

		// Process each chunk
		logger.Println("Fire up your engines, we start to process the chunks for this iteration, one chunk at a time.")
		wordprocess.ProcessChunksConcurrently(tempWordsFiles, outputChan, flags.Level, logger)

	}
	// Clean up remaining temp files and wait for output channel
	os.RemoveAll(os.TempDir() + "/domainwords*")
	close(outputChan)

	outputWG.Wait()

	logger.Println("Processing complete.")
}

// calculateBufferSize determines buffer size based on file size and memory constraints
func calculateBufferSize(filePath string) int {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	availableMemory := memStats.Sys

	bufferSize := 100 // Default size for stdin

	if filePath != "" {
		fileInfo, err := os.Stat(filePath)
		if err == nil {
			fileSize := fileInfo.Size()
			estimatedWordCount := fileSize / 8

			bufferSize = int(estimatedWordCount)

			if bufferSize > 1000 {
				bufferSize = 1000 // Cap it for large files
			}
		}
	}

	if bufferSize > int(availableMemory/(100*1024)) {
		bufferSize = int(availableMemory / (100 * 1024)) // Ensure it doesn't exceed 1% of available memory
	}

	return bufferSize
}
