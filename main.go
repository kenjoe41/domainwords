package main

import (
	"bufio"
	"fmt"
	"io"
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

		chaoticSlice := domainwords.ChaoticShuffle(words)

		// Write Shuffled words to Temp file
		tempWordsFile, err := domainwords.WriteTempFile(chaoticSlice)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		counter := 0
		var chunkedwords []string

		tfile, err := os.Open(tempWordsFile.Name())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer tfile.Close()

		filereader := bufio.NewReader(tfile)

		// Awfully inefficient, still waiting for beter ideas.
		for {
			word, _, err := filereader.ReadLine()
			if err != nil {
				if err == io.EOF {
					if chunkedwords != nil {

						domainwords.HandleWords(chunkedwords, depth, outputChan)

					}

					break
				}

				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			counter++
			chunkedwords = append(chunkedwords, string(word))

			if counter == int(flags.ChunkSize) {
				domainwords.HandleWords(chunkedwords, depth, outputChan)

				// reset counter
				counter = 0
				// Reset chunk
				chunkedwords = []string{}
			}

		}

		os.RemoveAll(tempWordsFile.Name())

	}

	close(outputChan)
	outputWG.Wait()

}
