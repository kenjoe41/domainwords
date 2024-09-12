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
	"strings"

	"github.com/google/go-github/v64/github"
	"golang.org/x/oauth2"
)

var (
	owner = os.Getenv("GITHUB_USERNAME")
	token = os.Getenv("GITHUB_TOKEN")
	repo  = "workflows"
	path  = "permutations.txt"
)

// HandleWords manages word permutations up to the specified depth.
func HandleWords(originalWords []string, depth uint, outputChan chan string) {
	for _, word := range originalWords {
		outputChan <- word
	}

	permutatedWords := originalWords
	for depth > 1 {
		permutatedWords = permutateWords(permutatedWords, originalWords, outputChan)
		depth--
	}
}

// permutateWords generates permutations and sends them to the output channel.
func permutateWords(permutatedWords, originalWords []string, outputChan chan string) []string {
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

// ReadLines reads a file line by line into a string slice.
func ReadLines(filename string) ([]string, error) {
	var lines []string
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

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

// RemoveDuplicateStr removes duplicate strings from a slice and returns only clean words.
func RemoveDuplicateStr(strSlice []string) []string {
	seen := make(map[string]struct{})
	var result []string

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

// isCleanWord checks if a word is valid (no single characters or symbols at the start or end).
func isCleanWord(word string) bool {
	if len(word) == 1 && isSymbol(word) {
		return false
	}
	if isSymbol(string(word[0])) || isSymbol(string(word[len(word)-1])) {
		return false
	}
	return true
}

// isSymbol checks if a string is a symbol.
func isSymbol(s string) bool {
	return !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(s)
}

// ConfigureDepth ensures depth is between 1 and 5.
func ConfigureDepth(depth uint) uint {
	if depth > 5 {
		fmt.Fprintln(os.Stderr, "[-] The maximum depth is 5. Adjusting.")
		return 5
	} else if depth < 1 {
		fmt.Fprintln(os.Stderr, "[-] The minimum depth is 1. Adjusting.")
		return 1
	}
	return depth
}

// ChunkSlice splits a slice into smaller chunks.
func ChunkSlice(wordsSlice []string, chunkSize int) [][]string {
	var chunks [][]string
	for len(wordsSlice) > 0 {
		if len(wordsSlice) < chunkSize {
			chunkSize = len(wordsSlice)
		}
		chunks = append(chunks, wordsSlice[:chunkSize])
		wordsSlice = wordsSlice[chunkSize:]
	}
	return chunks
}

// ChaoticShuffle shuffles a slice of strings randomly.
func ChaoticShuffle(wordsSlice []string) []string {
	rand.Shuffle(len(wordsSlice), func(i, j int) {
		wordsSlice[i], wordsSlice[j] = wordsSlice[j], wordsSlice[i]
	})
	return wordsSlice
}

// WriteTempChunks writes chunks to temporary files and returns the list of file handlers.
func WriteTempChunks(chunks [][]string) []os.File {
	var chunkFiles []os.File
	tmpDir := os.TempDir()

	// Clean up any residue
	_ = os.RemoveAll(tmpDir + "/domainwords*")

	for _, chunk := range chunks {
		chunkFile, err := writeTempFile(chunk, tmpDir)
		if err == nil {
			chunkFiles = append(chunkFiles, chunkFile)
		}
	}

	return chunkFiles
}

// writeTempFile writes a chunk of words to a temp file.
func writeTempFile(words []string, tmpDir string) (os.File, error) {
	tempFile, err := os.CreateTemp(tmpDir, "domainwords")
	if err != nil {
		return os.File{}, err
	}
	defer tempFile.Close()

	writer := bufio.NewWriter(tempFile)
	for _, word := range words {
		_, _ = writer.WriteString(word + "\n")
	}
	writer.Flush()
	return *tempFile, nil
}

// SyncResultsToGitHub syncs the results to a GitHub repository.
func SyncResultsToGitHub(results []string) error {

	if owner == "" || token == "" {
		return fmt.Errorf("no Github Username or Token supplied. Hint: populate the GITHUB_USERNAME and GITHUB_TOKEN environment variables")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Retrieve existing content
	content, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		return err
	}

	newContent := strings.Join(results, "\n")

	hasher := sha1.New()
	hasher.Write([]byte(newContent))
	newContentSHA := fmt.Sprintf("%x", hasher.Sum(nil))

	if content != nil && *content.SHA == newContentSHA {
		return nil
	}

	// Create or update the file in the repository
	_, _, err = client.Repositories.UpdateFile(ctx, owner, repo, path, &github.RepositoryContentFileOptions{
		Message:   github.String("Update file"),
		Content:   []byte(newContent),
		SHA:       content.SHA,
		Committer: &github.CommitAuthor{Name: github.String("kenjoe41"), Email: github.String("kenjoe41.nafuti@gmail.com")},
	})

	return err
}
