// Mark V. Shaney
//
// Mark V. Shaney implements a Markov Chain of words, recording a prefix
// (a sequence of words) and the frequency at which certain other words
// follow that prefix in a given corpus. The chain is built up by
// feeding each word in the corpus to the chain with the prefix of words
// that preceeded it.
//
// The caller starts with the empty prefix and adds the first word
// of the corpus. The prefix is shifted, dropping the first word of
// the prefix and adding the current word to the end of it, and then
// continues with the next word in the corpus, using the shifted prefix
// with the next word.
//
// A "document" is generated by walking the chain. The starting point
// is always the empty prefix, so if you only feed in one document, you
// will always start at the same point. Usually a paragraph of text is
// a "document", resetting the prefix at the start of each paragraph.
// This causes the generated document to start at a more random point.
// However, it is up to the caller to decide what the document boundary
// is.
//
// The prefix length is a constant defined in this file. The longer the
// prefix, the more the output document will resemble the input corpus,
// following longer phrases of the corpus. The usual length for this is
// two.

package main

import (
	"math/rand"
)

const PrefixLength = 2

type WordBag map[string]int // frequency-weighted set of words
type Prefix [PrefixLength]string
type MarkVShaney map[Prefix]WordBag

var (
	Initial  = Prefix{}
	Terminal = "" // zero value, not a valid word in the input
)

// Add a word to the Markov Chain for a given Prefix.
func (mvs MarkVShaney) Add(prefix Prefix, word string) {
	bag, ok := mvs[prefix]
	if !ok {
		bag = WordBag{}
		mvs[prefix] = bag
	}
	bag[word]++
}

// GetDocument returns a list of Words generated by walking the chain. A
// document always starts with the `Initial` prefix, and picks a suffix
// the normal way (`WordBag.GetOne`). If only one `Initial` prefix was
// used, the document always starts at the same place. A paragraph of
// text is the most common sort of document, but the document boundaries
// are up to the caller of `MarkVShaney.Add`.
func (mvs MarkVShaney) GetDocument() []string {
	result := []string{}
	prefix := Initial
	for {
		word := mvs.Walk(prefix)
		if word == Terminal {
			break
		}
		result = append(result, word)
		prefix.Shift(word)
	}
	return result
}

// Walk returns a random word from the Markov Chain given a Prefix.
func (mvs MarkVShaney) Walk(prefix Prefix) string {
	if bag, ok := mvs[prefix]; ok {
		return bag.GetOne()
	}
	return Terminal
}

// GetOne returns a random word from the WordBag, weighted by the word frequency.
func (bag WordBag) GetOne() string {
	entry := rand.Intn(bag.Len())
	var sum int
	for word, count := range bag {
		sum += count
		if entry <= sum {
			return word
		}
	}
	// Should never be reached, unless the bag length changes under us
	panic("entry out of range")
}

// Len returns the length of the WordBag, summing all the weights.
func (bag WordBag) Len() (sum int) {
	for _, count := range bag {
		sum += count
	}
	return
}

// Shift puts a new word at the end of a Prefix, shifting down the others.
func (p *Prefix) Shift(word string) {
	for i := 1; i < len(p); i++ {
		p[i-1] = p[i]
	}
	p[len(p)-1] = word
}
