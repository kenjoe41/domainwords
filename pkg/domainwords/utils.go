package domainwords

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

func HandleWords(originalWords []string, depth uint, permutationsChan chan string, outputChan chan string) {
	for _, word := range originalWords {
		// fmt.Println(word) // TODO: Chan this
		permutationsChan <- word
		outputChan <- word
	}

	// permutatedWords := originalWords
	// for ; depth > 1; depth-- {
	// 	permutatedWords = permutateWords(permutatedWords, originalWords, outputChan)
	// }

	permutateWords(originalWords, depth, permutationsChan, outputChan)

	outputChan <- ""

}

func permutateWords(originalWords []string, depth uint, permutationsChan chan string, outputChan chan string) {
	// var newPermutatedWords []string

	for permutatedWord := range permutationsChan {
		if strings.Count(permutatedWord, ".") == int(depth) {
			continue
		}

		for _, word := range originalWords {
			newWord := word + "." + permutatedWord
			// newPermutatedWords = append(newPermutatedWords, newWord)

			outputChan <- newWord
			// permutationsChan <- newWord

			if strings.Count(newWord, ".") < int(depth) {
				permutationsChan <- newWord
			}
		}

	}
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
	if depth > 3 {
		fmt.Fprintln(os.Stderr, "[-] The maximum is 3. Configuring")
		auxDepth = 3
	} else if depth < 1 {
		fmt.Fprintln(os.Stderr, "[-] The minimum is 1. Configuring")
		auxDepth = 1
	}
	return auxDepth
}
