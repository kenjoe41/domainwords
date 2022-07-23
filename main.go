package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/kenjoe41/domainwords/pkg/domainwords"
	"github.com/kenjoe41/domainwords/pkg/options"
)

var (
	words []string
)

func main() {

	flags := options.ScanFlags()

	if flags.Wordlist == "" {
		// No file provided,
		// TODO: Read from the commandline

		fmt.Fprintln(os.Stderr, "[-] No wordlist provided.")
		options.Usage()
		os.Exit(1)
	}

	var err error
	words, err = domainwords.ReadingLines(flags.Wordlist)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Sort words and remove dups
	words = domainwords.RemoveDuplicateStr(words)
	depth := domainwords.ConfigureDepth(flags.Level)

	outputChan := make(chan string, 1024)

	var outputWG sync.WaitGroup
	outputWG.Add(1)
	go func() {
		defer outputWG.Done()

		for permWord := range outputChan {
			fmt.Println(permWord)
		}

	}()

	// With depth=1 or few words < chuckSize, we have all the different iterations there is.
	if depth == 1 || len(words) < int(flags.ChunkSize) {
		flags.Iterations = 1
	}

	for iter := 0; iter < int(flags.Iterations); iter++ {

		chaoticSlice := domainwords.ChaoticShuffle(words)

		// Divide chaoticSlice into chunks and permutate each chunk.
		// TODO: Test where we should keep this in memory of write it out to temp-files
		dvWordSlices := domainwords.ChunkSlice(chaoticSlice, int(flags.ChunkSize))

		for _, dvWordSlice := range dvWordSlices {

			domainwords.HandleWords(dvWordSlice, depth, outputChan)

		}
	}

	close(outputChan)
	outputWG.Wait()

}
