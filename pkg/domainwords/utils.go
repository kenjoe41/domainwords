package domainwords

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

func HandleWords(originalWords []string, depth uint, outputChan chan string) {

	for _, word := range originalWords {
		outputChan <- word
	}

	permutatedWords := originalWords
	for ; depth > 1; depth-- {
		permutatedWords = permutateWords(permutatedWords, originalWords, outputChan)
	}

}
func permutateWords(permutatedWords []string, originalWords []string, outputChan chan string) []string {
	var newPermutatedWords []string

	for _, permutatedWord := range permutatedWords {
		for _, word := range originalWords {
			newWord := word + "." + permutatedWord
			newPermutatedWords = append(newPermutatedWords, newWord)

			outputChan <- newWord
		}

	}
	return newPermutatedWords
}

func ReadingLines(filename string) ([]string, error) {
	// Credits to https://github.com/j3ssie/str-replace
	var result []string
	file, err := os.Open(filename)
	if err != nil {
		return result, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		val := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if val == "" {
			continue
		}
		result = append(result, val)
	}

	if err := scanner.Err(); err != nil {
		return result, err
	}
	return result, nil
}

func RemoveDuplicateStr(strSlice []string) []string {

	sort.Strings(strSlice)

	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {

		if !isCleanWord(item) {
			continue
		}

		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func isCleanWord(word string) bool {
	// Check if its just one character and a symbol
	if len(word) == 1 {
		if isSymbol(word) {
			return false
		}
	}

	if isSymbol(string(word[0])) || isSymbol(string(word[len(word)-1])) {
		return false
	}
	return true
}

func isSymbol(s string) bool {
	isStringOrNumber := regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(s)

	return !isStringOrNumber
}

func ConfigureDepth(depth uint) uint {
	// Logic from https://github.com/Josue87/gotator
	auxDepth := depth
	if depth > 5 {
		fmt.Fprintln(os.Stderr, "[-] The maximum is 5. Configuring")
		auxDepth = 5
	} else if depth < 1 {
		fmt.Fprintln(os.Stderr, "[-] The minimum is 1. Configuring")
		auxDepth = 1
	}
	return auxDepth
}

func ChunkSlice(wordsSlice []string, chunkSize int) [][]string {
	var chunks [][]string
	for {
		if len(wordsSlice) == 0 {
			break
		}

		if len(wordsSlice) < chunkSize {
			chunkSize = len(wordsSlice)
		}

		chunks = append(chunks, wordsSlice[0:chunkSize])
		wordsSlice = wordsSlice[chunkSize:]
	}

	return chunks
}

func ChaoticShuffle(wordsSlice []string) []string {
	rand.Seed(time.Now().UnixNano())

	rand.Shuffle(len(wordsSlice), func(i, j int) {
		wordsSlice[i], wordsSlice[j] = wordsSlice[j], wordsSlice[i]
	})

	return wordsSlice
}

func WriteTempFile(words []string) (os.File, error) {
	tmpDir := os.TempDir()
	tempfile, err := ioutil.TempFile(tmpDir, "domainwords")
	if err != nil {
		return *tempfile, err
	}

	wordsfile, err := os.OpenFile(tempfile.Name(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return *wordsfile, err
	}

	datawriter := bufio.NewWriter(wordsfile)

	for _, data := range words {
		_, _ = datawriter.WriteString(data + "\n")
	}

	datawriter.Flush()
	wordsfile.Close()

	return *wordsfile, nil
}
