package output

import (
	"bufio"
	"fmt"
	"os"
)

// HandleOutput writes processed words to a file or stdout
func HandleOutput(outputChan <-chan string, outputFile string) error {
	var writer *bufio.Writer

	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			return err
		}
		defer file.Close()
		writer = bufio.NewWriter(file)
	} else {
		writer = bufio.NewWriter(os.Stdout)
	}

	for word := range outputChan {
		_, err := fmt.Fprintln(writer, word)
		if err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}
