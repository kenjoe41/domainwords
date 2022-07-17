package domainwords

import (
	"bufio"
	"os"
	"strings"
)

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
		val := strings.TrimSpace(scanner.Text())
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
