package options

import (
	"bufio"
	"os"
	"strings"

	"github.com/kenjoe41/domainwords/pkg/domainwords"
)

func HandleInput(flags Options) ([]string, error) {

	var words []string

	if flags.Wordlist == "" {
		// No file provided,

		// Check for stdin input
		stat, err := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return nil, err
		}

		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			words = append(words, strings.TrimSpace(sc.Text()))
		}

		// check there were no errors reading stdin (unlikely)
		if err := sc.Err(); err != nil {
			return nil, err
		}

	} else {

		var err error
		words, err = domainwords.ReadingLines(flags.Wordlist)
		if err != nil {
			return nil, err
		}
	}

	return words, nil
}
