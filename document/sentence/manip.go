package sentence

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/algao1/basically"
)

// RemoveConj removes the first coordinating conjunction
// (for, and, nor, etc.) from the sentence.
func RemoveConj(s *basically.Sentence) {
	// Sanity check to ensure that the string is sufficiently long.
	if len(s.Tokens) < 2 {
		return
	}

	if s.Tokens[0].Tag == "CC" {
		idx := strings.Index(s.Raw, s.Tokens[0].Text)
		len := utf8.RuneCountInString(s.Tokens[0].Text) + 1
		s.Raw = Capitalize(SubStr(s.Raw, idx+len, -1))
		s.Tokens = s.Tokens[1:]
	}
}

// Capitalize capitalizes the first letter in a string.
func Capitalize(str string) string {
	var upperStr string
	srunes := []rune(str)
	for idx := range srunes {
		if idx == 0 {
			srunes[idx] = unicode.ToUpper(srunes[idx])
		}
		upperStr += string(srunes[idx])
	}
	return upperStr
}

// SubStr returns the substring of s between start and end.
func SubStr(str string, start, end int) string {
	counter, startIdx := 0, 0
	for idx := range str {
		if counter == start {
			startIdx = idx
		}
		if counter == end {
			return str[startIdx:idx]
		}
		counter++
	}
	return str[startIdx:]
}
