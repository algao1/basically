package textrank

import (
	"math"
	"strings"

	"github.com/algao1/basically/internal/prose"
	"github.com/surgebase/porter2"
)

// A TokenFilter represents a (black/white) filter applied to tokens before similarity calculations.
type TokenFilter func(*prose.Token) bool

// A Similarity computes the similarity of two sentences (nodes) after applying the token filter.
type Similarity func(n1, n2 []*prose.Token, filter TokenFilter) float64

// NVFilter is a filter that whitelists tokens with n(oun) and v(erb) tags.
func NVFilter(tok *prose.Token) bool {
	return tok.Tag == "NN" || tok.Tag == "NNP" || tok.Tag == "NNPS" || tok.Tag == "NNS" ||
		tok.Tag == "VB" || tok.Tag == "VBD" || tok.Tag == "VBG" || tok.Tag == "VBN" || tok.Tag == "VBP" || tok.Tag == "VBZ" || tok.Tag == "MD"
}

// NVAAFilter is a filter that whitelists tokens with n(oun), v(erb), a(djective) and a(dverb) tokens.
func NVAAFilter(tok *prose.Token) bool {
	return tok.Tag == "NN" || tok.Tag == "NNP" || tok.Tag == "NNPS" || tok.Tag == "NNS" ||
		tok.Tag == "VB" || tok.Tag == "VBD" || tok.Tag == "VBG" || tok.Tag == "VBN" || tok.Tag == "VBP" || tok.Tag == "VBZ" || tok.Tag == "MD" ||
		tok.Tag == "JJ" || tok.Tag == "JJR" || tok.Tag == "JJS" ||
		tok.Tag == "RB" || tok.Tag == "RBR" || tok.Tag == "RBS" || tok.Tag == "RP"
}

// DefaultSimilarity is the default similarity implementation used in TextRank.
// Tokens are normalized (converted to lowercase and stemmed), before calculating similarity.
func DefaultSimilarity(n1, n2 []*prose.Token, filter TokenFilter) float64 {
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
