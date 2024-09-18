package main

import (
	"log"
	"os"

	"github.com/kenjoe41/domainwords/pkg/domainwords"
	"github.com/kenjoe41/domainwords/pkg/options"
)

var (
	words []string
)

func main() {
	// Create a logger
	logger := log.New(os.Stderr, "[DomainWords] ", log.LstdFlags)

	// Parsing command-line flags
	flags := options.ScanFlags()

	// Handle input (stdin or wordlist)
	var err error
	words, err = options.HandleInput(flags)
	if err != nil {
		logger.Fatalf("Error handling input: %v", err)
	}

	// Channel for processing output
	outputChan := make(chan string, 10240)

	// Move the outputFilePath to flags
	outputFilePath := flags.OutputFile
	domainwords.HandleOutput(outputChan, outputFilePath, flags.Sync)

	// Sort words and remove duplicates
	words = domainwords.RemoveDuplicateStr(words)

	// Configure depth based on flags
	depth := domainwords.ConfigureDepth(flags.Level)

	// Check if iterations should be minimized
	if depth == 1 || len(words) < int(flags.ChunkSize) {
		flags.Iterations = 1
	}

	for iter := 0; iter < int(flags.Iterations); iter++ {
		// Shuffle the words
		chaoticSlice := domainwords.ChaoticShuffle(words)

		// Break into chunks
		tempChunks := domainwords.ChunkSlice(chaoticSlice, int(flags.ChunkSize))

		// Write shuffled words to temp files
		tempWordsFiles := domainwords.WriteTempChunks(tempChunks)
		if len(tempWordsFiles) == 0 {
			logger.Fatalf("Errors while writing chunk files.")
		}

		// Process each chunk
		for _, tempChunkFile := range tempWordsFiles {
			chunkedWords, err := domainwords.ReadLines(tempChunkFile.Name())
			if err != nil {
				logger.Printf("Error reading lines from file: %v", err)
				os.RemoveAll(tempChunkFile.Name())
				continue
			}

			// Process the chunked words and write to output channel
			domainwords.HandleWords(chunkedWords, depth, outputChan)

			// Clean up temp files
			err = os.RemoveAll(tempChunkFile.Name())
			if err != nil {
				logger.Printf("Error removing temp file: %v", err)
			}
		}
	}

	// Clean up remaining temp files and close output channel
	os.RemoveAll(os.TempDir() + "/domainwords*")
	close(outputChan)

	// Wait for output to be fully written
	domainwords.WaitOutputCompletion()

	logger.Println("Processing complete.")
}
