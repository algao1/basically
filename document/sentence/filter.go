package sentence

import "github.com/algao1/basically"

// NVFilter is a filter that whitelists tokens with n(oun) and v(erb) tags.
func NVFilter(tok *basically.Token) bool {
	return IsNoun(tok.Tag) || IsVerb(tok.Tag)
}

// NVAAFilter is a filter that whitelists tokens with n(oun), v(erb), a(djective) and a(dverb) tokens.
func NVAAFilter(tok *basically.Token) bool {
	return IsNoun(tok.Tag) || IsVerb(tok.Tag) || IsAdj(tok.Tag) || IsAdv(tok.Tag)
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
