package basically

import (
	"fmt"
	"sort"

	"github.com/algao1/basically/textrank"
)

// A Config represents a setting that changes the summarization process.
// For example, it may configure a custom token filter:
// 		topSents := textrank.Summarize("...", length, textrank.WithCustomFilter(filter))
type Config func(cfgs *Configs)

// Configs controls the summarization process.
type Configs struct {
	filter     textrank.TokenFilter // Default whitelists tokens with noun and verb tags.
	similarity textrank.Similarity  // Default uses the DefaultSimilarity similarity function.

	quotations   bool // Default enables merging sentences within quotations.
	conjunctions bool // Default removes conjunctions from beginning of sentences (not implemented).

	focus     bool    // Default enables a sentence (focus/bias) to be used for ranking sentence scores.
	focusSent string  // Default uses the first sentence of the given text as focus/bias.
	threshold float64 // Default sets the similarity threshold to 0.65 as recommended in Biased TextRank.
}

// WithCustomFilter dictates a custom function for whitelisting/backlisting tokens.
func WithCustomFilter(filter textrank.TokenFilter) Config {
	return func(cfgs *Configs) { cfgs.filter = filter }
}

// WithCustomSimilarity dictates a custom function for calculating sentence similarity.
func WithCustomSimilarity(similarity textrank.Similarity) Config {
	return func(cfgs *Configs) { cfgs.similarity = similarity }
}

// WithoutMergeQuotations disables merging sentences within quotations.
func WithoutMergeQuotations() Config {
	return func(cfgs *Configs) { cfgs.quotations = false }
}

// WithoutFocus disables the use of a focus for ranking sentence scores.
func WithoutFocus() Config {
	return func(cfgs *Configs) { cfgs.focus = false }
}

// WithFocusSentence enables a focus (bias) sentence to be set for ranking sentence scores.
func WithFocusSentence(sent string) Config {
	return func(cfgs *Configs) {
		cfgs.focus = true
		cfgs.focusSent = sent
	}
}

// WithCustomThreshold sets the similarity threshold as per the specification.
// Lower threshold values correspond with sparse graphs, and higher threshold values
// correspond with dense graphs.
func WithCustomThreshold(threshold float64) Config {
	return func(cfgs *Configs) {
		cfgs.threshold = threshold
	}
}

// Summarize returns a summary of the given text with the top relevant phrases according
// to the user-specified configuration.
// For example, summarizing using a custom token filter:
// 		topSents := textrank.Summarize("...", length, textrank.WithCustomFilter(filter))
func Summarize(text string, len int, cfgs ...Config) ([]textrank.Node, error) {
	g := &textrank.Graph{}

	configs := Configs{
		quotations: true,
		focus:      true,
		threshold:  0.65,
	}

	for _, applyConfig := range cfgs {
		applyConfig(&configs)
	}

	sim := configs.similarity
	if sim == nil {
		sim = textrank.DefaultSimilarity
	}

	filter := configs.filter
	if filter == nil {
		filter = textrank.NVFilter
	}

	// Adds nodes and edges representing sentences and inter-sentence relationships respectively.
	acclen := g.Build(text, sim, filter, configs.quotations, configs.focus, configs.focusSent)
	g.Connect(sim, filter, configs.threshold)

	if acclen < len {
		return nil, fmt.Errorf("text is too short")
	}

	// Applies the TextRank algorithm to the graph constructed previously.
	textrank.Rank(g, 15)

	// Sorts the ranked nodes first by score, then by sentence order in the original text.
	nodes := g.Nodes()
	sort.SliceStable(nodes, func(i, j int) bool { return nodes[i].Score > nodes[j].Score })
	topSentences := nodes[:len]
	sort.SliceStable(topSentences, func(i, j int) bool { return topSentences[i].Order < topSentences[j].Order })

	return topSentences, nil
}
