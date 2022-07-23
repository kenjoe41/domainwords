package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/kenjoe41/domainwords/pkg/domainwords"
	"github.com/kenjoe41/domainwords/pkg/options"
)

// TODO: Take Already known subdomains file, Hello goSubsWordlist

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

	outputChan := make(chan string, 1024)

	var outputWG sync.WaitGroup
	outputWG.Add(1)
	go func() {

		for permWord := range outputChan {
			fmt.Println(permWord)
		}
		outputWG.Done()
	}()

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

			chunkedwords, err := domainwords.ReadingLines(tempChunkFile.Name())
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.RemoveAll(tempChunkFile.Name())
				continue

			}

			domainwords.HandleWords(chunkedwords, depth, outputChan)

			os.RemoveAll(tempChunkFile.Name())
		}

	}

	close(outputChan)
	outputWG.Wait()

}
