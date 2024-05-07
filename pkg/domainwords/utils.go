package domainwords

import (
	"bufio"
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"math/rand"
	"os"
	"regexp"

	"github.com/google/go-github/v61/github"
	"golang.org/x/oauth2"

	"strings"
)

var (
	owner = os.Getenv("GITHUB_USERNAME")
	token = os.Getenv("GITHUB_TOKEN")
	repo  = "workflows"
	path  = "permutations.txt"
)

func HandleWords(originalWords []string, depth uint, outputChan chan string) {

	for _, word := range originalWords {
		outputChan <- word
	}

	permutatedWords := originalWords
	for ; depth > 1; depth-- {
		permutatedWords = permutateWords(permutatedWords, originalWords, outputChan)
	}

}
func permutateWords(permutatedWords []string, originalWords []string, outputChan chan string) []string {
	var newPermutatedWords []string

	for _, permutatedWord := range permutatedWords {
		for _, word := range originalWords {
			newWord := word + "." + permutatedWord
			newPermutatedWords = append(newPermutatedWords, newWord)

			outputChan <- newWord
		}

	}
	return newPermutatedWords
}

func ReadLines(filename string) ([]string, error) {
	var lines []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Use buffered reader for better performance
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return lines, err
		}
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, strings.ToLower(line))
		}
	}

	return lines, nil
}

func RemoveDuplicateStr(strSlice []string) []string {
	// Create a set to store unique strings
	seen := make(map[string]struct{})
	result := make([]string, 0, len(strSlice))

	// Iterate over the slice and add each unique string to the set
	for _, str := range strSlice {
		if isCleanWord(str) {
			if _, ok := seen[str]; !ok {
				seen[str] = struct{}{}
				result = append(result, str)
			}
		}
	}

	return result
}

func isCleanWord(word string) bool {
	// Check if its just one character and a symbol, or word's pre- / suffix is a symbol.
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
	if depth > 5 {
		fmt.Fprintln(os.Stderr, "[-] The maximum is 5. Configuring")
		auxDepth = 5
	} else if depth < 1 {
		fmt.Fprintln(os.Stderr, "[-] The minimum is 1. Configuring")
		auxDepth = 1
	}
	return auxDepth
}

func ChunkSlice(wordsSlice []string, chunkSize int) [][]string {
	var chunks [][]string
	for {
		if len(wordsSlice) == 0 {
			break
		}

		if len(wordsSlice) < chunkSize {
			chunkSize = len(wordsSlice)
		}

		chunks = append(chunks, wordsSlice[0:chunkSize])
		wordsSlice = wordsSlice[chunkSize:]
	}

	return chunks
}

func ChaoticShuffle(wordsSlice []string) []string {

	rand.Shuffle(len(wordsSlice), func(i, j int) {
		wordsSlice[i], wordsSlice[j] = wordsSlice[j], wordsSlice[i]
	})

	return wordsSlice
}

func WriteTempChunks(chunks [][]string) []os.File {
	var chunkfiles []os.File
	tmpDir := os.TempDir()

	// Clean up any residue
	os.RemoveAll(tmpDir + "/domainwords*")

	for _, chunk := range chunks {
		chunkfile, err := writeTempFile(chunk, tmpDir)
		if err != nil {
			continue
		}
		chunkfiles = append(chunkfiles, chunkfile)
	}

	return chunkfiles

}

func writeTempFile(words []string, tmpDir string) (os.File, error) {
	tempfile, err := os.CreateTemp(tmpDir, "domainwords")
	if err != nil {
		return *tempfile, err
	}

	wordsfile, err := os.OpenFile(tempfile.Name(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return *wordsfile, err
	}

	datawriter := bufio.NewWriter(wordsfile)

	for _, data := range words {
		_, _ = datawriter.WriteString(data + "\n")
	}

	datawriter.Flush()
	wordsfile.Close()

	return *wordsfile, nil
}

// SyncResultsToGitHub syncs the provided results to a GitHub repository.
func syncResultsToGitHub(results []string) error {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// Retrieve existing content
	content, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {

		return err
	}

	// Convert results to string
	newContent := strings.Join(results, "\n")

	// Calculate SHA256 hash of the new content
	hasher := sha1.New()
	hasher.Write([]byte(newContent))
	newContentSHAString := fmt.Sprintf("%x", hasher.Sum(nil))

	// If the content exists and the SHA matches, no need to update
	if content != nil && *content.SHA == newContentSHAString {
		fmt.Println("same same")
		return nil
	}

	// Create or update the file
	_, _, err = client.Repositories.UpdateFile(ctx, owner, repo, path, &github.RepositoryContentFileOptions{
		Message:   github.String("Update file"),
		Content:   []byte(newContent),
		SHA:       content.SHA,
		Committer: &github.CommitAuthor{Name: github.String("kenjoe41"), Email: github.String("kenjoe41.nafuti@gmail.com")},
	})
	if err != nil {
		return err
	}

	return nil
}
