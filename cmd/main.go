package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/algao1/basically"
)

func main() {
	start := time.Now()
	summaryLen, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal()
	}
	file := os.Args[2]

	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	topSentences, err := basically.Summarize(string(data), summaryLen)
	if err != nil {
		log.Fatal(err)
	}

	var result string
	for _, node := range topSentences {
		result += node.Sentence
		fmt.Printf("[%.2f] %s\n", node.Score, node.Sentence)
	}

	elapsed := time.Since(start)
	log.Printf("Function took %s", elapsed)
}
