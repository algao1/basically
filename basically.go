package basically

// A Document represents a given text, and is responsible for
// handling the summarization and keyword extraction process.
type Document interface {
	Summarize(length int, threshold float64, focus string) ([]*Sentence, error)
	Highlight(length int, merge bool) ([]*Keyword, error)
	Characters() (int, int)
}

// A Parser is responsible for parsing and tokenizing a document
// into strings and words. A Parser also performs additional tasks
// such as POS-tagging and sentiment analysis.
type Parser interface {
	ParseDocument(doc string, quote bool) ([]*Sentence, []*Token, error)
}

// A Summarizer is responsible for extracting key sentences from a
// document.
type Summarizer interface {
	Initialize(sents []*Sentence, similar Similarity, filter TokenFilter,
		focusString *Sentence, threshold float64)
	Rank(iters int)
}

// A Highlighter is responsible for extracting key words from a document.
type Highlighter interface {
	Initialize(tokens []*Token, filter TokenFilter, window int)
	Rank(iters int)
	Highlight(length int, merge bool) ([]*Keyword, error)
}

// A TokenFilter represents a (black/white) filter applied to tokens before similarity calculations.
type TokenFilter func(*Token) bool

// A Similarity computes the similarity of two sentences after applying the token filter.
type Similarity func(n1, n2 []*Token, filter TokenFilter) float64

// A Token represents an individual token of text such as a word or punctuation
// symbol.
type Token struct {
	Tag   string // The token's part-of-speech tag.
	Text  string // The token's actual content.
	Order int    // The token's order in the text.
}

// A Keyword is the keyword belonging to a highlighted document.
// A Keyword contains the raw word, and its associated weight.
type Keyword struct {
	Word   string  // Raw keyword.
	Weight float64 // Weight of the keyword.
}

// A Sentence represents an individual sentence within the text.
type Sentence struct {
	Raw       string   // Raw sentence string.
	Tokens    []*Token // Tokenized sentence.
	Sentiment float64  // Sentiment score.
	Score     float64  // Score (weight) of the sentence.
	Bias      float64  // Bias assigned to the sentence for ranking.
	Order     int      // The sentence's order in the text.
}
