package textrank

import (
	"fmt"

	"github.com/algao1/basically/internal/prose"
)

// A Graph is a representation of the text containing nodes and undirected edges.
// Nodes represent individual sentences, and edges being the connections between them.
type Graph struct {
	nodes []*Node
	edges [][]float64
}

// A Node is a representation of a sentence (lexical unit).
// Also contains the individual words/symbols of the sentence as tokens.
type Node struct {
	Sentence string
	Score    float64
	Bias     float64
	Order    int
	Tokens   []*prose.Token
}

// AddNode adds a node to the undirected graph.
func (g *Graph) AddNode(n *Node) {
	g.nodes = append(g.nodes, n)
	g.edges = append(g.edges, make([]float64, len(g.nodes)))
}

// Node returns a copy of node at the selected index.
// Returns an error if the index is out of bounds.
func (g *Graph) Node(n int) (Node, error) {
	if n >= len(g.nodes) {
		return Node{}, fmt.Errorf("index %d out of bounds", n)
	}
	return *g.nodes[n], nil
}

// Nodes returns a copy of the nodes in the graph.
func (g *Graph) Nodes() []Node {
	retNodes := make([]Node, len(g.nodes))
	for idx, node := range g.nodes {
		retNodes[idx] = *node
	}
	return retNodes
}

// SetEdge sets the weight of an undirected edge between nodes
// at index n1 and n2.
func (g *Graph) SetEdge(n1, n2 int, weight float64) error {
	if n1 < n2 {
		n1, n2 = n2, n1
	}
	if n1 >= len(g.nodes) {
		return fmt.Errorf("index %d out of bounds", n2)
	}
	g.edges[n1][n2] = weight
	return nil
}

// Edge returns the weight of the edge (n1,n2).
func (g *Graph) Edge(n1, n2 int) (float64, error) {
	if n1 < n2 {
		n1, n2 = n2, n1
	}
	if n1 >= len(g.nodes) {
		return -1, fmt.Errorf("index %d out of bounds", n2)
	}
	return g.edges[n1][n2], nil
}

// Display is a helper function that draws the graph
// and displays it in a user-friendly fashion.
func (g *Graph) Display() {
	for y := 0; y < len(g.edges); y++ {
		for x := 0; x < len(g.edges); x++ {
			if y > x {
				fmt.Printf("     ")
			} else {
				fmt.Printf("%3.2f ", g.edges[x][y])
			}
		}
		fmt.Println()
	}
}
