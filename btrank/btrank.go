package btrank

import "github.com/algao1/basically"

// A SGraph is an undirected graph, representing the sentences within a document.
// The nodes represent individual sentences, and edges represent the connection
// between sentences.
type SGraph struct {
	Nodes []*basically.Sentence
	Edges [][]float64
}

// BiasedTextRank implements the Summarizer interface.
type BiasedTextRank struct {
	Graph *SGraph
}

var _ basically.Summarizer = (*BiasedTextRank)(nil)

// Initialize initializes the underlying SGraph by inserting nodes and edges.
func (btr *BiasedTextRank) Initialize(sents []*basically.Sentence, similar basically.Similarity,
	filter basically.TokenFilter, focusString *basically.Sentence, threshold float64) {
	// Instantiate a new SGraph with the appropriate nodes and edges.
	btr.Graph = &SGraph{Nodes: sents, Edges: make([][]float64, len(sents))}

	// Updates the bias value of each node (sentence) if necessary, and
	// constructs the edge slice for each node.
	for idx, sent := range sents {
		if focusString != nil {
			btr.Graph.Nodes[idx].Bias = similar(focusString.Tokens, sent.Tokens, filter)
		}
		btr.Graph.Edges[idx] = make([]float64, idx+1)
	}

	// Constructs the edges between nodes satisfying the threshold,
	// using the given similarity function, and token filter.
	for i := 0; i < len(btr.Graph.Nodes); i++ {
		for j := 0; j < i; j++ {
			sim := similar(btr.Graph.Nodes[i].Tokens, btr.Graph.Nodes[j].Tokens, filter)
			if sim > threshold {
				btr.Graph.Edges[i][j] = sim
			}
		}
	}
}

// edge returns the weight of the edge.
// Since SGraph is symmetric, edge(y, x) == edge(x, y).
func (btr *BiasedTextRank) edge(x, y int) float64 {
	if y > x {
		return btr.Graph.Edges[y][x]
	}
	return btr.Graph.Edges[x][y]
}

// outWeights calculates the weights of outgoing edges.
func (btr *BiasedTextRank) outWeights() []float64 {
	n := len(btr.Graph.Nodes)
	weights := make([]float64, 0, n)

	for x := 0; x < n; x++ {
		var sum float64
		for y := 0; y < n; y++ {
			sum += btr.edge(x, y)
		}
		weights = append(weights, sum)
	}
	return weights
}

// Rank applies the Biased TextRank algorithm on the SGraph for some specified iterations.
func (btr *BiasedTextRank) Rank(iters int) {
	n := len(btr.Graph.Nodes)
	outWeights := btr.outWeights()

	for iter := 0; iter < iters; iter++ {
		for x := 0; x < n; x++ {
			var sum float64

			for y := 0; y < n; y++ {
				// Ignore node if the outWeights are too small.
				if outWeights[y] < 1e-4 {
					continue
				}
				sum += btr.edge(x, y) * (btr.Graph.Nodes[y].Score / outWeights[y])
			}

			btr.Graph.Nodes[x].Score = btr.Graph.Nodes[x].Bias*0.15 + 0.85*sum
		}
	}
}
