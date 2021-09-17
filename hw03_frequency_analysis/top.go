package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var (
	reStringSplitter = regexp.MustCompile(`\s+`)
	// https://stackoverflow.com/questions/1716609/how-to-match-cyrillic-characters-with-a-regular-expression
	reWordSanitizer = regexp.MustCompile(`(^[^\p{L}]+|[^\p{L}]+$)`)
)

// total O(N + MlogM + k) ~ O(NlogN)
func TopK(s string, k int) []string {
	// split by words, O(N)
	rawWords := reStringSplitter.Split(s, -1)

	words := make([]string, 0, len(rawWords))
	for _, w := range rawWords {
		w = string(reWordSanitizer.ReplaceAll([]byte(w), []byte{}))
		if len(w) > 0 {
			words = append(words, strings.ToLower(w))
		}
	}

	// accumulate pattern
	frq := make(map[string]int)
	for _, w := range words {
		frq[w]++
	}

	// collect keys only
	uniq := make([]string, 0)
	for k := range frq {
		uniq = append(uniq, k)
	}

	// O(MlogM) where M is count of unique words
	sort.SliceStable(uniq, func(a, b int) bool {
		if frq[uniq[a]] == frq[uniq[b]] {
			return uniq[a] < uniq[b]
		}
		return frq[uniq[b]] < frq[uniq[a]]
	})

	if k > len(uniq) {
		k = len(uniq)
	}

	return uniq[:k]
}

func Top10(s string) []string {
	return TopK(s, 10)
}
