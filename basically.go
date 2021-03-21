package basically

// A Document represents a given text, and is responsible for
// handling the summarization and keyword extraction process.
type Document interface {
	Summarize(length int, focus string) ([]*Sentence, error)
	// Extract() ([]*Token, error)
}

// A Parser ...
type Parser interface {
	ParseDocument(doc string, quote bool) ([]*Sentence, []*Token, error)
}

// A Summarizer ...
type Summarizer interface {
	Initialize(sents []*Sentence, similar Similarity, filter TokenFilter,
		focusString *Sentence, threshold float64)
	Rank(iters int)
}

// A Highlighter ...
type Highlighter interface {
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
	Label string // The token's IOB label.
}

// A Sentence represents an individual sentence within the text.
type Sentence struct {
	Raw       string
	Tokens    []*Token
	Sentiment float64
	Score     float64
	Bias      float64
	Order     int
}
