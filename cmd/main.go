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
)

func main() {
	start := time.Now()

	sumlen, _ := strconv.Atoi(os.Args[1])
	file := os.Args[2]
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	focus := ""
	if len(os.Args) > 3 {
		focus = os.Args[3]
	}

	p := &parser.Parser{}
	s := &btrank.BiasedTextRank{}

	document, err := document.Create(string(data), s, nil, p)
	if err != nil {
		log.Fatal(err)
	}

	sums, err := document.Summarize(sumlen, focus)
	if err != nil {
		log.Fatal(err)
	}

	for _, sum := range sums {
		fmt.Printf("[%.2f, %.2f]\n", sum.Score, sum.Sentiment)
		fmt.Println(sum.Raw)
	}

	elapsed := time.Since(start)
	log.Printf("Function took %s", elapsed)
}
