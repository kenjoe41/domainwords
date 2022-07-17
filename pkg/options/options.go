package options

import "flag"

type Options struct {
	Wordlist string
	Level    int
}

func ScanFlags() Options {
	wordlistPtr := flag.String("w", "", "File containing list of words.")
	levelPtr := flag.Int("l", 3, "Level of Permutations to do.")

	flag.Parse()

	optFlags := Options{
		*wordlistPtr,
		*levelPtr,
	}

	return optFlags
}

func Usage() {
	flag.Usage()
}