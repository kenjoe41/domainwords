package domainwords

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"sync"
)

// ReadLines reads a file line by line into a string slice.
func ReadLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// Pre-allocate memory for lines slice (optional, assuming a rough estimate of 1000 lines)
	lines := make([]string, 0, 1000) // TODO: Dynamically preallocate memory for each slice.

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, strings.ToLower(line))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// RemoveDuplicateStr removes duplicate strings using sync.Map for concurrency-safe operations.
func RemoveDuplicateStr(wordsChan <-chan string) []string {
	var seen sync.Map
	var result []string
	var mu sync.Mutex

	for word := range wordsChan {
		if _, ok := seen.LoadOrStore(word, struct{}{}); !ok {
			mu.Lock()
			result = append(result, word)
			mu.Unlock()
		}
	}

	return result
}

// isCleanWord checks if a word is valid (no single characters or symbols at the start or end).
func IsCleanWord(word string) bool {
	// Check if the word is empty
	if len(word) == 0 {
		return false
	}

	// Check if it's just one character and a symbol
	if len(word) == 1 {
		if isSymbol(word) {
			return false
		}
	}

	// Check if the word's prefix or suffix is a symbol
	if isSymbol(string(word[0])) || isSymbol(string(word[len(word)-1])) {
		return false
	}

	return true
}

// isSymbol checks if a string is a symbol.
func isSymbol(s string) bool {
	return !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(s)
}

// ConfigureDepth ensures depth is between 1 and 5.
func ConfigureDepth(depth uint) uint {
	if depth > 5 {
		fmt.Fprintln(os.Stderr, "[-] The maximum depth is 5. Adjusting.")
		return 5
	} else if depth < 1 {
		fmt.Fprintln(os.Stderr, "[-] The minimum depth is 1. Adjusting.")
		return 1
	}
	return depth
}

// ChunkSlice splits a slice into smaller chunks.
func ChunkSlice(wordsSlice []string, chunkSize int) [][]string {
	var chunks [][]string
	for len(wordsSlice) > 0 {
		if len(wordsSlice) < chunkSize {
			chunkSize = len(wordsSlice)
		}
		chunks = append(chunks, wordsSlice[:chunkSize])
		wordsSlice = wordsSlice[chunkSize:]
	}
	return chunks
}

// ChaoticShuffle shuffles a slice of strings randomly.
func ChaoticShuffle(wordsSlice []string) []string {
	rand.Shuffle(len(wordsSlice), func(i, j int) {
		wordsSlice[i], wordsSlice[j] = wordsSlice[j], wordsSlice[i]
	})
	return wordsSlice
}

// WriteTempChunks writes chunks to temporary files and returns the list of file handlers.
func WriteTempChunks(chunks [][]string) []os.File {
	var chunkFiles []os.File
	tmpDir := os.TempDir()

	// Clean up any residue
	_ = os.RemoveAll(tmpDir + "/domainwords*")

	for _, chunk := range chunks {
		chunkFile, err := writeTempFile(chunk, tmpDir)
		if err == nil {
			chunkFiles = append(chunkFiles, chunkFile)
		}
	}

	return chunkFiles
}

// writeTempFile writes a chunk of words to a temp file.
func writeTempFile(words []string, tmpDir string) (os.File, error) {
	tempFile, err := os.CreateTemp(tmpDir, "domainwords")
	if err != nil {
		return os.File{}, err
	}
	defer tempFile.Close()

	writer := bufio.NewWriter(tempFile)
	for _, word := range words {
		_, _ = writer.WriteString(word + "\n")
	}
	writer.Flush()
	return *tempFile, nil
}
