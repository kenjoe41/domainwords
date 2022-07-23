package options

import flag "github.com/spf13/pflag"

type Options struct {
	Wordlist   string
	Level      uint
	ChunkSize  uint
	Iterations uint
}

func ScanFlags() Options {
	wordlistPtr := flag.StringP("wordlist", "w", "", "File containing list of words. Or cat wordlist.txt | domainwords")
	levelPtr := flag.UintP("level", "l", 3, "Level of Permutations to do (1-5).")
	ChuckSizePtr := flag.UintP("chuck", "c", 20000, "Chuck size per slice.")
	IterationsPtr := flag.UintP("iterations", "i", 10, "Number of Iterations of shuffling, chunking and permutation [For BIG wordlists].")

	flag.Parse()

	optFlags := Options{
		*wordlistPtr,
		*levelPtr,
		*ChuckSizePtr,
		*IterationsPtr,
	}

	return optFlags
}

func Usage() {
	flag.Usage()
}
