package document

import (
	"fmt"
	"sort"

	"github.com/algao1/basically"
	"github.com/algao1/basically/document/sentence"
)

// A Config represents a setting that changes the summarization process.
// For example, it may configure a custom token filter:
// 		doc := document.Create(..., document.WithCustomFilter(filter))
type Config func(cfgs *Configs)

// Configs control the summarization process.
type Configs struct {
	sfilter      basically.TokenFilter // Default uses matcher.NVFilter.
	kwfilter     basically.TokenFilter // Default uses matcher.NVNSFilter.
	similarity   basically.Similarity  // Default uses sentence.NVFilter.
	quotations   bool                  // Default disables merging of sentences within quotations.
	conjunctions bool                  // Default removes conjunctions from the beginning of sentences.
	focus        bool                  // Default uses the first sentence as focus if a focus sentence is not provided.
	threshold    float64               // Default sets the similarity threshold to 0.65 as recommended in Biased TextRank.
}

// WithCustomSFilter allows for a custom (black/white) token filter to be set
// for sentence extraction.
func WithCustomSFilter(filter basically.TokenFilter) Config {
	return func(cfgs *Configs) { cfgs.sfilter = filter }
}

// WithCustomKWFilter allows for a custom (black/white) token filter to be set
// for keyword extraction.
func WithCustomKWFilter(filter basically.TokenFilter) Config {
	return func(cfgs *Configs) { cfgs.kwfilter = filter }
}

// WithCustomSimilarity allows for a custom similarity function to be set.
func WithCustomSimilarity(similarity basically.Similarity) Config {
	return func(cfgs *Configs) { cfgs.similarity = similarity }
}

// WithConjunctions disables removing conjunctions from the start of sentences.
func WithConjunctions() Config {
	return func(cfgs *Configs) { cfgs.conjunctions = true }
}

// WithCustomThreshold sets the similarity threshold as per the specification.
// Lower threshold values correspond with sparser graphs, and higher threshold values
// correspond with denser graphs.
func WithCustomThreshold(threshold float64) Config {
	return func(cfgs *Configs) { cfgs.threshold = threshold }
}

// WithoutMergeQuotations disables merging sentences within quotations.
func WithoutMergeQuotations() Config {
	return func(cfgs *Configs) { cfgs.quotations = false }
}

// WithoutFocus disables the use of a focus for ranking sentence scores.
func WithoutFocus() Config {
	return func(cfgs *Configs) { cfgs.focus = false }
}

// Document is an implementation of basically.Document.
type Document struct {
	// Configurations and dependency injection.
	Configs     *Configs
	Summarizer  basically.Summarizer
	Highlighter basically.Highlighter
	Parser      basically.Parser
	// Parsed sentences, and words.
	Sentences []*basically.Sentence
	Words     []*basically.Token
}

func Create(text string, s basically.Summarizer, h basically.Highlighter,
	p basically.Parser, cfgs ...Config) (basically.Document, error) {
	// Initializes and applies the configurations.
	// The threshold is set based on the results from https://www.aclweb.org/anthology/P04-3020.pdf.
	m := sentence.CreateMatcher()
	configs := Configs{
		sfilter:      m.NVFilter,
		kwfilter:     m.NVNSFilter,
		similarity:   sentence.DefaultSimilarity,
		conjunctions: false,
		focus:        true,
		threshold:    0.65,
	}
	for _, applyConfig := range cfgs {
		applyConfig(&configs)
	}

	// Parses the document into sentences and words.
	sents, words, err := p.ParseDocument(text, configs.quotations)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", "unable to parse document", err)
	}

	// Create and return the document.
	doc := &Document{
		Configs:     &configs,
		Summarizer:  s,
		Highlighter: h,
		Parser:      p,
		Sentences:   sents,
		Words:       words,
	}
	return doc, nil
}

// Summarize returns a summary of given length corresponding to the top relevant phrases.
// A focus string may be provided to adjust the summary contents.
func (doc *Document) Summarize(length int, raw string) ([]*basically.Sentence, error) {
	// Sanity check to ensure that the given text is sufficiently large.
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
	doc.Summarizer.Initialize(doc.Sentences, doc.Configs.similarity, doc.Configs.sfilter, focus, doc.Configs.threshold)
	doc.Summarizer.Rank(5)

	// Sorts the ranked sentences first by score, then by their sentence order in the original text.
	sort.SliceStable(doc.Sentences, func(i, j int) bool { return doc.Sentences[i].Score > doc.Sentences[j].Score })
	sort.SliceStable(doc.Sentences[:length], func(i, j int) bool { return doc.Sentences[i].Order < doc.Sentences[j].Order })

	// Handle conjunctions at the beginning of sentences.
	if !doc.Configs.conjunctions {
		for _, sent := range doc.Sentences[:length] {
			sentence.RemoveConj(sent)
		}
	}

	return doc.Sentences[:length], nil
}

// Highlight returns a list of the keywords in the document.
func (doc *Document) Highlight(length int, merge bool) ([]*basically.Keyword, error) {
	// Initialize the highlighter, and apply the ranking algorithm.
	doc.Highlighter.Initialize(doc.Words, doc.Configs.kwfilter, 2)
	doc.Highlighter.Rank(25)
	return doc.Highlighter.Highlight(length, merge)
}
