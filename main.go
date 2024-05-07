package main

import (
	"fmt"
	"os"

	"github.com/kenjoe41/domainwords/pkg/domainwords"
	"github.com/kenjoe41/domainwords/pkg/options"
)

var (
	words []string
)

func main() {

	flags := options.ScanFlags()

	var err error
	words, err = options.HandleInput(flags)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	outputChan := make(chan string, 10240)
	outputFilePath := "test/output.txt" // TODO: move to configuration file or flags
	domainwords.HandleOutput(outputChan, outputFilePath)

	// Sort words and remove dups
	words = domainwords.RemoveDuplicateStr(words)

	depth := domainwords.ConfigureDepth(flags.Level)

	// With depth=1 or few words < chuckSize, we have all the different iterations there is.
	if depth == 1 || len(words) < int(flags.ChunkSize) {
		flags.Iterations = 1
	}

	for iter := 0; iter < int(flags.Iterations); iter++ {

		// Shuffle and randomise the words
		chaoticSlice := domainwords.ChaoticShuffle(words)

		// Break up into chunks
		tempchunks := domainwords.ChunkSlice(chaoticSlice, int(flags.ChunkSize))

		// Write Shuffled words to Temp file
		tempWordsFiles := domainwords.WriteTempChunks(tempchunks)
		if len(tempWordsFiles) == 0 {
			fmt.Fprintln(os.Stderr, "Errors while writing chunk files.")
			os.Exit(1)
		}

		// Lets work on one chunk at a time.
		for _, tempChunkFile := range tempWordsFiles {

			chunkedwords, err := domainwords.ReadLines(tempChunkFile.Name())
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.RemoveAll(tempChunkFile.Name())
				continue

			}

			domainwords.HandleWords(chunkedwords, depth, outputChan)

			err = os.RemoveAll(tempChunkFile.Name())
			if err != nil {
				continue
			}
		}

	}
	os.RemoveAll(os.TempDir() + "/domainwords*")

	close(outputChan)
	domainwords.WaitOutputCompletion() // Wait for the output goroutine to finish

}
