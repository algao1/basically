package parser

import (
	"strings"

	"github.com/algao1/basically"
	"github.com/algao1/basically/internal/prose"
	"github.com/jonreiter/govader"
)

type Parser struct {
	sentTokenizer *prose.PunktSentenceTokenizer
	wordTokenizer *prose.IterTokenizer
	tagger        *prose.PerceptronTagger
	analyzer      *govader.SentimentIntensityAnalyzer
}

var _ basically.Parser = (*Parser)(nil)

// Create initializes the tokenizers, tagger, and sentiment analyzer.
// Token classification is disabled for performance speed-up.
func Create() *Parser {
	return &Parser{
		sentTokenizer: prose.NewPunktSentenceTokenizer(),
		wordTokenizer: prose.NewIterTokenizer(),
		tagger:        prose.DefaultModel(true, false).Tagger,
		analyzer:      govader.NewSentimentIntensityAnalyzer(),
	}
}

// ParseDocument parses a document into sentences and tokens.
// The result contains additional information such as sentence sentiment,
// and POS-tags for tokens.
func (p *Parser) ParseDocument(doc string, quote bool) ([]*basically.Sentence, []*basically.Token, error) {
	sents := p.sentTokenizer.Segment(doc)
	retSents := make([]*basically.Sentence, 0, len(sents))
	retTokens := make([]*basically.Token, 0, len(sents)*15)

	tokCounter := 0
	for idx, sent := range sents {
		tokens := p.wordTokenizer.Tokenize(sent.Text)
		tokens = p.tagger.Tag(tokens)

		// Convert struct from []*prose.Token to []*basically.Token.
		btokens := make([]*basically.Token, 0, len(tokens))
		for _, tok := range tokens {
			btok := basically.Token{Tag: tok.Tag, Text: strings.ToLower(tok.Text), Order: tokCounter}
			btokens = append(btokens, &btok)
			tokCounter++
		}

		// Analyzes sentence sentiment.
		sentiment := p.analyzer.PolarityScores(sent.Text).Compound

		retSents = append(retSents, &basically.Sentence{
			Raw:       sent.Text,
			Tokens:    btokens,
			Sentiment: sentiment,
			Bias:      1.0,
			Order:     idx,
		})

		retTokens = append(retTokens, btokens...)
	}

	return retSents, retTokens, nil
}
