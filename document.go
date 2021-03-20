package basically

// A Token represents an individual token of text such as a word or punctuation
// symbol.
type Token struct {
	Tag   string // The token's part-of-speech tag.
	Text  string // The token's actual content.
	Label string // The token's IOB label.
}

type Sentence struct {
	Raw       string
	Tokens    []*Token
	Sentiment float64
	Order     int
}

type Document struct {
	Sentences []*Sentence
	Words     []*Token
}

type Parser interface {
	ParseDocument(doc string, quote bool) ([]*Sentence, []*Token, error)
}

type SentenceGraph interface {
}
