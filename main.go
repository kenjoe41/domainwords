package main

import (
	"bufio"
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
	isStdin := false

	if flags.Wordlist == "" {
		// No file provided,
		isStdin = true

		// Check for stdin input
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			fmt.Fprintln(os.Stderr, "No words/text detected. Hint: cat wordlist.txt | domainwords")
			os.Exit(1)
		}

	} else {

		var err error
		words, err = domainwords.ReadingLines(flags.Wordlist)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	if isStdin {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			words = append(words, sc.Text())
		}

		// check there were no errors reading stdin (unlikely)
		if err := sc.Err(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	outputChan := make(chan string, 1024)

	var outputWG sync.WaitGroup
	outputWG.Add(1)
	go func() {
		defer outputWG.Done()

		for permWord := range outputChan {
			fmt.Println(permWord)
		}

	}()

	// Sort words and remove dups
	words = domainwords.RemoveDuplicateStr(words)
	depth := domainwords.ConfigureDepth(flags.Level)

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
