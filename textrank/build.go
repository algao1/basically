package textrank

import (
	"strings"

	"github.com/algao1/basically/internal/prose"
)

// Build constructs the underlying undirected graph using the given text, and configurations.
func (g *Graph) Build(text string, similarity Similarity, filter TokenFilter,
	mergeQuo bool, focus bool, focusSent string) int {
	// Initialize the tokenizers and tagger.
	sentTokenizer := prose.NewPunktSentenceTokenizer()
	wordTokenizer := prose.NewIterTokenizer()
	tagger := prose.DefaultModel(true, false)

	sents := sentTokenizer.Segment(text)

	// Merge sentences within quotations.
	if mergeQuo {
		sents = mergeQuotations(sents)
	}

	// Setup focus sentence to calculate bias score.
	var ftokens []*prose.Token
	if len(focusSent) == 0 {
		focusSent = sents[0].Text
	}
	if focus {
		ftokens = wordTokenizer.Tokenize(focusSent)
		ftokens = tagger.Tagger.Tag(ftokens)
	}

	// Process and add the sentences as nodes.
	for idx, sent := range sents {
		tokens := wordTokenizer.Tokenize(sent.Text)
		tokens = tagger.Tagger.Tag(tokens)

		bias := 1.0
		if focus {
			bias = similarity(ftokens, tokens, filter)
		}
		bias = bias * (float64(len(sents)) - float64(idx) + 1) / float64(len(sents))

		g.AddNode(&Node{
			Sentence: sent.Text,
			Score:    0.15,
			Bias:     bias,
			Order:    idx,
			Tokens:   tokens,
		})
	}

	// Returns the number of sentences in the text.
	return len(sents)
}

// Connect constructs edges between nodes (sentences) using the given similarity function,
// and token filter.
func (g *Graph) Connect(similarity Similarity, filter TokenFilter, threshold float64) {
	for i := range g.nodes {
		for j := 0; j < i; j++ {
			sim := similarity(g.nodes[i].Tokens, g.nodes[j].Tokens, filter)
			g.edges[i][j] = 0
			if sim > threshold {
				g.edges[i][j] = sim
			}
		}
	}
}

// mergeQuotations merges sentences within quotations.
func mergeQuotations(sents []prose.Sentence) []prose.Sentence {
	var quote bool
	var j int

	for i := 0; i < len(sents); i++ {
		if quote {
			sents[j].Text += " " + prim(sents[i].Text)
		} else {
			sents[j].Text = prim(sents[i].Text)
		}

		if countQuotes(sents[i].Text)%2 == 1 {
			quote = !quote
		}

		if !quote {
			j++
		}
	}
	return sents[:j]
}

// prim removes the '\n' character from text, and trims extra whitespace.
func prim(str string) string {
	return strings.TrimSpace(strings.ReplaceAll(str, "\n", ""))
}

// countQuotes returns the number of quotation marks within the text.
func countQuotes(str string) int {
	return strings.Count(str, "“") + strings.Count(str, "”") +
		strings.Count(str, "″") + strings.Count(str, "\"")
}
