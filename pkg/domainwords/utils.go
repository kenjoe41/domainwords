package domainwords

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func HandleWords(originalWords []string, depth uint) {
	for _, word := range originalWords {
		fmt.Println(word) // TODO: Chan this
	}

	permutatedWords := originalWords
	for ; depth > 1; depth-- {
		permutatedWords = permutateWords(permutatedWords, originalWords)
	}
}

func permutateWords(permutatedWords []string, originalWords []string) []string {
	var newPermutatedWords []string

	for _, permutatedWord := range permutatedWords {
		for _, word := range originalWords {
			newWord := word + "." + permutatedWord
			newPermutatedWords = append(newPermutatedWords, newWord)

			fmt.Println(newWord) // TODO: Chan this
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
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
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
