package main

import (
	"fmt"
	"os"

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

	for _, word := range words {
		fmt.Println(word)
	}
}