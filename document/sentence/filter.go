package sentence

import (
	"os"
	"path"
	"runtime"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/algao1/basically"
)

type Matcher struct {
	Stopwords map[string]struct{}
}

func CreateMatcher() *Matcher {
	m := &Matcher{Stopwords: make(map[string]struct{})}
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	words, _ := os.ReadFile(path.Join(path.Dir(filename), "stopwords.txt"))
	for _, word := range strings.Split(string(words), "\n") {
		m.Stopwords[word] = struct{}{}
	}
	return m
}

// NVFilter is a filter that whitelists tokens with (n)oun and (v)erb tags.
func (m *Matcher) NVFilter(tok *basically.Token) bool {
	return AlphaStart(tok.Text) && (IsNoun(tok.Tag) || IsVerb(tok.Tag))
}

// NVNSFilter is a filter that whitelists tokens with (n)oun and (v)erb tags,
// and blacklists tokens that are (s)topwords.
func (m *Matcher) NVNSFilter(tok *basically.Token) bool {
	_, stop := m.Stopwords[tok.Text]
	return !stop && AlphaStart(tok.Text) && (IsNoun(tok.Tag) || IsAdj(tok.Tag))
}

// NVAAFilter is a filter that whitelists tokens with n(oun), v(erb), a(djective) and a(dverb) tokens.
func (m *Matcher) NVAAFilter(tok *basically.Token) bool {
	return AlphaStart(tok.Text) && (IsNoun(tok.Tag) || IsVerb(tok.Tag) || IsAdj(tok.Tag) || IsAdv(tok.Tag))
}

func IsNoun(tag string) bool {
	return tag == "NN" || tag == "NNP" || tag == "NNPS" || tag == "NNS"
}

func IsVerb(tag string) bool {
	return tag == "VB" || tag == "VBD" || tag == "VBG" ||
		tag == "VBN" || tag == "VBP" || tag == "VBZ" || tag == "MD"
}

func IsAdj(tag string) bool {
	return tag == "JJ" || tag == "JJR" || tag == "JJS"
}

func IsAdv(tag string) bool {
	return tag == "RB" || tag == "RBR" || tag == "RBS" || tag == "RP"
}

func AlphaStart(text string) bool {
	r, _ := utf8.DecodeRuneInString(text)
	return unicode.IsLetter(r)
}
