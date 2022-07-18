package options

import "flag"

type Options struct {
	Wordlist string
	Level    uint
}

func ScanFlags() Options {
	wordlistPtr := flag.String("w", "", "File containing list of words.")
	levelPtr := flag.Uint("l", 3, "Level of Permutations to do (1-5).")

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
