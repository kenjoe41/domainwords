package wordprocess

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/edsrzf/mmap-go" // Assuming you are using mmap-go package
	"github.com/kenjoe41/domainwords/pkg/domainwords"
)

// ProcessChunksConcurrently processes wordlist chunks using a worker pool.
func ProcessChunksConcurrently(tempWordsFiles []os.File, outputChan chan string, depth uint, logger *log.Logger) error {
	var wg sync.WaitGroup
	chunkChan := make(chan []string, len(tempWordsFiles)) // Buffered channel for chunks

	// Load chunks into the channel concurrently
	go func() {
		defer close(chunkChan) // Ensure the channel is closed after loading
		for _, tempChunkFile := range tempWordsFiles {
			chunkedWords, err := readLinesWithMmap(tempChunkFile.Name())
			if err != nil {
				logger.Printf("Error reading lines from file: %v", err)
				_ = os.RemoveAll(tempChunkFile.Name()) // Clean up on error
				continue
			}

			logger.Fatalf("We have %d words, %v", len(chunkedWords), chunkedWords)
			chunkChan <- chunkedWords

			// Clean up temp file after processing
			err = os.RemoveAll(tempChunkFile.Name())
			if err != nil {
				logger.Printf("Error removing temp file: %v", err)
			}
		}
	}()

	// Worker pool to process chunks concurrently
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for chunk := range chunkChan {
				processChunk(chunk, depth, outputChan)
			}
		}()
	}

	wg.Wait() // Wait for all workers to finish
	return nil
}

// processChunk handles permutations for a given chunk and sends results to outputChan.
func processChunk(originalWords []string, depth uint, outputChan chan string) {
	// Send original words to the output channel
	for _, word := range originalWords {
		outputChan <- word
	}

	// Generate permutations up to the given depth
	permutateWords(originalWords, depth, outputChan)
}

// permutateWords generates permutations up to a given depth and sends them to outputChan using a strings.Builder.
func permutateWords(originalWords []string, depth uint, outputChan chan string) {
	permutatedWords := originalWords
	for d := uint(1); d < depth; d++ {
		var newPermutatedWords []string

		for _, permutatedWord := range permutatedWords {
			for _, word := range originalWords {
				var builder strings.Builder

				// Preallocate memory for the combined length of word + "." + permutatedWord
				builder.Grow(len(word) + 1 + len(permutatedWord))

				// Construct the new word
				builder.WriteString(word)
				builder.WriteString(".")
				builder.WriteString(permutatedWord)

				newWord := builder.String()
				newPermutatedWords = append(newPermutatedWords, newWord)
				outputChan <- newWord
			}
		}

		permutatedWords = newPermutatedWords // Move to next depth
	}
}

// readLinesWithMmap reads a file using memory-mapped I/O for efficient access.
func readLinesWithMmap(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}
	defer file.Close()

	// Use mmap to map the file into memory
	data, err := mmap.Map(file, mmap.RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("error mapping file: %w", err)
	}
	defer data.Unmap()

	// Split the file content into words (assuming new-line delimited words)
	return splitIntoWords(string(data)), nil
}

// splitIntoWords splits the input data into words by new lines (or other delimiters).
func splitIntoWords(data string) []string {
	// Assuming the words are newline-separated.
	words := strings.Split(data, "\n")
	var cleanWords []string
	for _, word := range words {
		if domainwords.IsCleanWord(word) {
			cleanWords = append(cleanWords, word)
		}
	}
	return cleanWords
}
