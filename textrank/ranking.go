package textrank

// Rank iterates the TextRank algorithm over the graph, and updates the scores.
func Rank(g *Graph, maxIter int) {
	outWeights := outWeights(g)
	for iter := 0; iter < maxIter; iter++ {
		for x := 0; x < len(g.nodes); x++ {
			var sum float64
			for y := 0; y < len(g.nodes); y++ {
				if outWeights[y] < 1e-4 {
					continue
				}
				edge, _ := g.Edge(y, x)
				sum += edge * (g.nodes[y].Score / outWeights[y])
			}
			g.nodes[x].Score = g.nodes[x].Bias*(1-0.85) + 0.85*sum
		}
	}
}

// outWeights calculates the weights of edges going out of a node.
func outWeights(g *Graph) []float64 {
	weights := make([]float64, 0, len(g.nodes))
	for y := 0; y < len(g.nodes); y++ {
		var sum float64
		for x := 0; x < len(g.nodes); x++ {
			edge, _ := g.Edge(x, y)
			sum += edge
		}
		weights = append(weights, sum)
	}
	return weights
}
