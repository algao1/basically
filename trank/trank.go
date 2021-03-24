package trank

import (
	"fmt"
	"math"
	"sort"

	"github.com/algao1/basically"
)

// A WGraph is an undirected graph, with nodes representing words within
// a document. Edges represent the connection between words.
type WGraph struct {
	Nodes map[string]float64
	Edges map[string]map[string]int
}

// KWTextRank implements the Highlighter interface.
type KWTextRank struct {
	Graph  *WGraph
	Tokens []*basically.Token
}

var _ basically.Highlighter = (*KWTextRank)(nil)

// Initialize initializes the underlying WGraph by inserting nodes and edges.
func (kwtr *KWTextRank) Initialize(tokens []*basically.Token,
	filter basically.TokenFilter, window int) {
	// Instantiate a new WGraph.
	kwtr.Graph = &WGraph{Nodes: make(map[string]float64), Edges: make(map[string]map[string]int)}
	kwtr.Tokens = tokens

	for i := 0; i < len(tokens); i++ {
		// Insert tokens as nodes if they pass through filter.
		if filter(tokens[i]) {
			text := tokens[i].Text
			kwtr.Graph.Nodes[text] = 1.0
			if _, ok := kwtr.Graph.Edges[text]; !ok {
				kwtr.Graph.Edges[text] = make(map[string]int)
			}

			// Backtracks to add edges.
			if i >= window {
				kwtr.lookBack(tokens[i-window:i+1], filter)
			}
		}
	}
}

// lookBack adds weighted edges to the WGraph if any of the previous tokens within a window
// satisfies the filter.
func (kwtr *KWTextRank) lookBack(tokens []*basically.Token, filter basically.TokenFilter) {
	end := len(tokens) - 1
	t1 := tokens[end]
	for idx, t2 := range tokens[:end] {
		if filter(t2) {
			lt1, lt2 := t1.Text, t2.Text
			// The weight of the edge increases the closer the words are.
			kwtr.Graph.Edges[lt1][lt2] = end / (end - idx)
			kwtr.Graph.Edges[lt2][lt1] = end / (end - idx)
		}
	}
}

// Rank applies the TextRank algorithm on the WGraph for some specified iterations.
func (kwtr *KWTextRank) Rank(iters int) {
	outWeights := kwtr.outWeights()

	for iter := 0; iter < iters; iter++ {
		for to, edges := range kwtr.Graph.Edges {
			var sum float64
			for from := range edges {
				// Ignore outWeights if 0 to avoid Inf or NaN.
				if outWeights[to] == 0 {
					continue
				}
				sum += float64(kwtr.Graph.Edges[to][from]) / outWeights[from] * kwtr.Graph.Nodes[from]
			}
			kwtr.Graph.Nodes[to] = 0.15 + 0.85*sum
		}
	}
}

// outWeights calculates the weights of outgoing edges.
func (kwtr *KWTextRank) outWeights() map[string]float64 {
	n := len(kwtr.Graph.Nodes)
	weights := make(map[string]float64, n)

	// Iterate over the edge sets.
	for t, e := range kwtr.Graph.Edges {
		var sum float64
		// Iterate over the edges/weights in the edge set.
		for _, w := range e {
			sum += float64(w)
		}
		weights[t] = sum
	}

	return weights
}

// Highlight sorts the keywords by weight, and returns the most significant keywords.
// Can optionally be specified to merge keywords together to get multi-word extraction.
func (kwtr *KWTextRank) Highlight(words int, merge bool) ([]*basically.Keyword, error) {
	// Sanity check to ensure that requested word count is below maximum.
	if words > len(kwtr.Graph.Nodes) {
		return nil, fmt.Errorf("unable to highlight document, not enough keywords")
	}

	// If the number of keywords is negative, automatically set to 1/3 of nodes.
	if words < 0 {
		words = len(kwtr.Graph.Nodes) / 3
	}

	// Create a list of keywords.
	kwords := make([]*basically.Keyword, 0, len(kwtr.Graph.Nodes))
	for kw, w := range kwtr.Graph.Nodes {
		kwords = append(kwords, &basically.Keyword{Word: kw, Weight: w})
	}

	sort.Slice(kwords, func(i, j int) bool { return kwords[i].Weight > kwords[j].Weight })
	kwords = kwords[:words]

	// Form multi-word keywords if specified.
	if merge {
		// Create a dictionary to avoid duplicates.
		kwdict := make(map[string]*basically.Keyword)
		for _, kw := range kwords {
			kwdict[kw.Word] = kw
		}

		queue := make([]*basically.Keyword, 0)
		for _, tok := range kwtr.Tokens {
			if kw, ok := kwdict[tok.Text]; ok {
				queue = append(queue, kw)
				continue
			} else if len(queue) > 1 {
				// Append the new keyword if non-duplicate.
				nKW := mergeKW(queue)
				if _, ok := kwdict[nKW.Word]; !ok {
					kwdict[nKW.Word] = nKW
					kwords = append(kwords, nKW)
				}
			}
			queue = nil
		}
	}

	sort.Slice(kwords, func(i, j int) bool { return kwords[i].Weight > kwords[j].Weight })
	return kwords[:words], nil
}

// mergeKW merges a slice of keywords to form one keyword.
func mergeKW(kws []*basically.Keyword) *basically.Keyword {
	word := ""
	max := 0.0
	min := 0.0
	sum := 0.0
	for idx, kw := range kws {
		if idx > 0 {
			word += " "
		}
		word += kw.Word
		max = math.Max(max, kw.Weight)
		min = math.Min(min, kw.Weight)
		sum += kw.Weight
	}

	weight := max + math.Log10(max) - math.Log(min+1)
	return &basically.Keyword{Word: word, Weight: weight}
}
