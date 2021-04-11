package document

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/algao1/basically/btrank"
	"github.com/algao1/basically/parser"
	"github.com/algao1/basically/trank"
)

func BenchmarkSummarize(b *testing.B) {
	file, err := os.Open("../test")
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()

	// Read all files and folders from directory.
	list, _ := file.Readdirnames(0)
	texts := make([]string, len(list))

	// Retrieve text from files.
	for _, name := range list {
		if !strings.Contains(name, ".txt") {
			continue
		}
		data, _ := os.ReadFile("../test/" + name)
		texts = append(texts, string(data))
	}

	b.ResetTimer()

	parser := parser.Create()
	s := &btrank.BiasedTextRank{}
	h := &trank.KWTextRank{}

	for n := 0; n < b.N; n++ {
		for _, text := range texts {
			doc, _ := Create(text, s, h, parser)
			doc.Summarize(5, 0.3, "")
			doc.Highlight(5, true)
		}
	}
}
