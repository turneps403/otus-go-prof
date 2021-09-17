package hw03frequencyanalysis

import (
	"regexp"
	"sort"
)

// total O(N + MlogM + k) ~ O(NlogN)
func TopK(s string, k int) []string {
	// split by words, O(N)
	re := regexp.MustCompile(`\s+`)
	rawWords := re.Split(s, -1)

	words := make([]string, 0, len(rawWords))
	for _, w := range rawWords {
		if len(w) > 0 {
			words = append(words, w)
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
