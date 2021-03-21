package document

import (
	"fmt"
	"sort"

	"github.com/algao1/basically"
)

// A Config represents a setting that changes the summarization process.
// For example, it may configure a custom token filter:
// 		doc := document.Create(..., document.WithCustomFilter(filter))
type Config func(cfgs *Configs)

// Configs control the summarization process.
type Configs struct {
	quotations   bool    // Default disables merging sentences within quotations.
	conjunctions bool    // TODO: Default removes conjunctions from beginning of sentences.
	focus        bool    // Default uses the first sentence as focus if a focus is not provided.
	threshold    float64 // Default sets the similarity threshold to 0.65 as recommended in Biased TextRank.
}

// WithoutMergeQuotations disables merging sentences within quotations.
func WithoutMergeQuotations() Config {
	return func(cfgs *Configs) { cfgs.quotations = false }
}

// WithoutFocus disables the use of a focus for ranking sentence scores.
func WithoutFocus() Config {
	return func(cfgs *Configs) { cfgs.focus = false }
}

// WithCustomThreshold sets the similarity threshold as per the specification.
// Lower threshold values correspond with sparse graphs, and higher threshold values
// correspond with dense graphs.
func (doc *Document) WithCustomThreshold(threshold float64) Config {
	return func(cfgs *Configs) {
		cfgs.threshold = threshold
	}
}

// Document is an implementation of basically.Document.
type Document struct {
	// Configurations and dependency injection.
	Configs     *Configs
	Similarity  basically.Similarity
	Filter      basically.TokenFilter
	Summarizer  basically.Summarizer
	Highlighter basically.Highlighter
	Parser      basically.Parser
	// Parsed sentences, and words.
	Sentences []*basically.Sentence
	Words     []*basically.Token
}

func Create(text string, s basically.Summarizer, h basically.Highlighter, p basically.Parser,
	fil basically.TokenFilter, sim basically.Similarity, cfgs ...Config) (basically.Document, error) {
	// Initializes and applies the configurations.
	// The threshold is set based on the results from https://www.aclweb.org/anthology/P04-3020.pdf.
	configs := Configs{focus: true, threshold: 0.65}
	for _, applyConfig := range cfgs {
		applyConfig(&configs)
	}

	// Parses the document for sentences and words.
	sents, words, err := p.ParseDocument(text, configs.quotations)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", "unable to parse document", err)
	}

	// Create and return the document.
	doc := &Document{
		Configs:     &configs,
		Filter:      fil,
		Similarity:  sim,
		Summarizer:  s,
		Highlighter: h,
		Parser:      p,
		Sentences:   sents,
		Words:       words,
	}
	return doc, nil
}

// Summarize returns a summary of the given text with the top relevant phrases according
// to the user-specified configuration.
// For example, summarizing using a custom token filter:
// 		topSents := textrank.Summarize("...", length, textrank.WithCustomFilter(filter))
func (doc *Document) Summarize(length int, raw string) ([]*basically.Sentence, error) {
	// Sanity check to see if the text is sufficiently long.
	if length > len(doc.Sentences) {
		return nil, fmt.Errorf("text is too short")
	}

	var focus *basically.Sentence
	if len(raw) > 0 {
		sents, _, err := doc.Parser.ParseDocument(raw, doc.Configs.quotations)
		if err != nil {
			return nil, fmt.Errorf("%q: %w", "unable to parse focus sentence", err)
		}
		if len(sents) > 0 {
			focus = sents[0]
		}
	} else if doc.Configs.focus {
		focus = doc.Sentences[0]
	}

	// Initializes and ranks the sentences.
	doc.Summarizer.Initialize(doc.Sentences, doc.Similarity, doc.Filter, focus, doc.Configs.threshold)
	doc.Summarizer.Rank(5)

	// Sorts the ranked sentences first by score, then by their sentence order in the original text.
	sort.SliceStable(doc.Sentences, func(i, j int) bool { return doc.Sentences[i].Score > doc.Sentences[j].Score })
	sort.SliceStable(doc.Sentences[:length], func(i, j int) bool { return doc.Sentences[i].Order < doc.Sentences[j].Order })

	return doc.Sentences[:length], nil
}
