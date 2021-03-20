package parser

import (
	"github.com/algao1/basically"
	"github.com/algao1/basically/internal/prose"
	"github.com/jonreiter/govader"
)

type Parser struct{}

var _ basically.Parser = (*Parser)(nil)

// ParseDocument parses a document into sentences and tokens.
// The result contains additional information such as sentence sentiment,
// and POS-tags for tokens.
func (p *Parser) ParseDocument(doc string, quote bool) ([]*basically.Sentence, []*basically.Token, error) {
	// Initialize tokenizers, taggers, and sentiment analyzers.
	sentTokenizer := prose.NewPunktSentenceTokenizer()
	wordTokenizer := prose.NewIterTokenizer()
	// Disable token classification for performance speed-up.
	tagger := prose.DefaultModel(true, false)
	analyzer := govader.NewSentimentIntensityAnalyzer()

	sents := sentTokenizer.Segment(doc)
	retSents := make([]*basically.Sentence, 0, len(sents))
	retTokens := make([]*basically.Token, 0, len(sents)*15)

	for idx, sent := range sents {
		tokens := wordTokenizer.Tokenize(sent.Text)
		tokens = tagger.Tagger.Tag(tokens)

		// Convert struct from []*prose.Token to []*basically.Token.
		btokens := make([]*basically.Token, 0, len(tokens))
		for _, tok := range tokens {
			btok := basically.Token(*tok)
			btokens = append(btokens, &btok)
		}

		sentiment := analyzer.PolarityScores(sent.Text).Compound

		retSents = append(retSents, &basically.Sentence{
			Raw:       sent.Text,
			Tokens:    btokens,
			Sentiment: sentiment,
			Order:     idx,
		})
		retTokens = append(retTokens, btokens...)
	}

	return retSents, retTokens, nil
}
