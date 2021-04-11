package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/algao1/basically/btrank"
	"github.com/algao1/basically/document"
	"github.com/algao1/basically/parser"
	"github.com/algao1/basically/trank"
)

func main() {
	start := time.Now()

	sumlen, _ := strconv.Atoi(os.Args[1])
	files := os.Args[2:]

	p := parser.Create()
	s := &btrank.BiasedTextRank{}
	kwtr := &trank.KWTextRank{}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		doc, err := document.Create(string(data), s, kwtr, p)
		if err != nil {
			log.Fatal(err)
		}

		sums, err := doc.Summarize(sumlen, 0.3, "")
		if err != nil {
			log.Fatal(err)
		}

		orig, redu := doc.Characters()
		fmt.Println("Reduced by:", (1-float64(redu)/float64(orig))*100)

		for _, sum := range sums {
			fmt.Printf("[%.2f, %.2f]\n", sum.Score, sum.Sentiment)
			fmt.Println(sum.Raw)
		}

		kws, err := doc.Highlight(20, true)
		if err != nil {
			log.Fatal(err)
		}

		for _, kw := range kws {
			fmt.Println(kw.Weight, kw.Word)
		}
	}

	elapsed := time.Since(start)
	log.Printf("Function took %s", elapsed)
}
