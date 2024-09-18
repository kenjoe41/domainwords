package options

import flag "github.com/spf13/pflag"

type Options struct {
	Wordlist   string
	Level      uint
	ChunkSize  uint
	Iterations uint
	OutputFile string
	Sync       bool
}

func ScanFlags() Options {
	wordlistPtr := flag.StringP("wordlist", "w", "", "File containing list of words. Or cat wordlist.txt | domainwords")
	levelPtr := flag.UintP("level", "l", 3, "Level of Permutations to do (1-5).")
	chunkSizePtr := flag.UintP("chunk", "c", 20000, "Chunk size per slice.")
	iterationsPtr := flag.UintP("iterations", "i", 10, "Number of Iterations of shuffling, chunking, and permutation [For BIG wordlists].")
	outputFilePtr := flag.StringP("output", "o", "", "Output filename path.")
	syncPtr := flag.BoolP("sync", "", false, "Sync File to Github.")

	flag.Parse()

	return Options{
		Wordlist:   *wordlistPtr,
		Level:      *levelPtr,
		ChunkSize:  *chunkSizePtr,
		Iterations: *iterationsPtr,
		OutputFile: *outputFilePtr,
		Sync:       *syncPtr,
	}
}

func Usage() {
	flag.Usage()
}
