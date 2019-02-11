package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

func main() {
	mvs := MarkVShaney{}
	for _, filename := range os.Args[1:] {
		if err := parseInput(filename, mvs); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	// Output 5 paragraphs.
	// TODO(camh-): Make number of paragraphs a CLI option
	for i := 0; i < 5; i++ {
		fmt.Println(strings.Join(mvs.GetDocument(), " "), "\n")
	}
}

// parseInput parses in input file for paragraphs and then parses
// those paragraphs with parseParagraph. Returns an error if the given
// filename could not be opened.
func parseInput(filename string, mvs MarkVShaney) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(scanParagraph)
	for scanner.Scan() {
		parseParagraph(scanner.Text(), mvs)
	}
	return nil
}

// parseParagraph parses a single paragraph for words and feeds the words
// and their prefixes into a MarkVShaney.
func parseParagraph(p string, mvs MarkVShaney) {
	prefix := Initial
	scanner := bufio.NewScanner(strings.NewReader(p))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Text()
		mvs.Add(prefix, word)
		prefix.Shift(word)
	}
}

// scanParagraph is a bufio.SplitFunc for a scanner to split input on paragraph
// boundaries: "\n\n".
func scanParagraph(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		// All done. At EOF and no unprocessed input
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte("\n\n")); i >= 0 {
		// Return up to \n\n and consume input after it
		return i + 2, data[0:i], nil
	}
	if atEOF {
		// Return remainder of buffer. EOF terminates a paragraph
		return len(data), data, nil
	}
	// Not at EOF, but no paragraph break. ask the scanner for more input
	return 0, nil, nil
}
