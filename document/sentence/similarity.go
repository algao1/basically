package sentence

import (
	"math"
	"strings"

	"github.com/algao1/basically"
	"github.com/surgebase/porter2"
)

// DefaultSimilarity is the default similarity implementation used in Biased TextRank.
// Tokens are normalized (converted to lowercase and stemmed), before calculating similarity.
func DefaultSimilarity(n1, n2 []*basically.Token, filter basically.TokenFilter) float64 {
	var ret float64
	l1, l2 := float64(len(n1)), float64(len(n2))
	freqTable := make(map[string]int)

	for _, tok := range append(n1, n2...) {
		norm := porter2.Stem(strings.ToLower(tok.Text))
		if _, ok := freqTable[norm]; ok && filter(tok) {
			ret++
		} else if !ok {
			freqTable[norm]++
		}
	}
	return ret / (math.Log10(l1) + math.Log10(l2))
}
