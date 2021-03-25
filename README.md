# basically

`basically` is a Go implementation of the TextRank and Biased TextRank algorithm built upon [`prose`](https://github.com/jdkato/prose). It provides fully unsupervised methods for keyword extraction and focused text summarization, along with additional quality of life features over the original implementations.

## Methods

First, the document is parsed into its constituent sentences and words using a sentence segmenter and tokenizer. Sentiment values are assigned to individual sentences, and tokens are annotated with part of speech tags.

For **keyword extraction**, all words that pass the syntactic filter are added to a undirected, weighted graph, and an edge is added between words that co-occur within a window of $N$ words. The edge weight is set to be inversely proportional to the distance between the words. Each vertex is assigned an initial score of 1, and the following ranking algorithm is run on the graph

$$TR(V_i) = (1 - d) + d \times \sum_{V_j \in In(V_i)} \frac{w_{ji}}{\sum_{V_k \in Out(v_j)} w_{jk}} TR(V_j)$$

During post-processing, adjacent keywords are collapsed into a multi-word keyword, and the top keywords are then extracted.

For **sentence extraction**, every sentence is added to a undirected, weighted graph, with an edge between sentences that share *common content*. The edge weight is set simply as the number of common tokens between the lexical representations of the two sentences. Each vertex is also assigned an initial score of 1, and a bias score based on the focus text, before the following ranking algorithm is run on the graph

$$BTR(V_i) = Bias \times (1 - d) + d \times \sum_{V_j \in In(V_i)} \frac{w_{ji}}{\sum_{V_k \in Out(v_j)} w_{jk}} BTR(V_j)$$

The top weighted sentences are then selected and sorted in chronological order to form a summary.

Further information on the two algorithms can be found [here](https://web.eecs.umich.edu/~mihalcea/papers/mihalcea.emnlp04.pdf) and [here](https://arxiv.org/pdf/2011.01026.pdf).

## Installation

```console
go get https://github.com/algao1/basically
``` 

## Usage

```Go
// Instantiate a document for every text.
doc, err := document.Create(text, &btrank.BiasedTextRank{}, &trank.KWTextRank{}, &parser.Parser{})
if err != nil {
	log.Fatal(err)
}

// Summarize the document into 7 sentences with respect to a focus sentence.
sents, err := document.Summarize(7, focus)
if err != nil {
	log.Fatal(err)
}

for _, sent := range sents {
	fmt.Printf("[%.2f, %.2f] %s\n", sum.Score, sum.Sentiment, sum.Raw)
}

// Highlight the top 7 keywords in the document, with multi-word keywords enabled.
words, err := document.Highlight(7, true)
if err != nil {
	log.Fatal(err)
}

for _, word := range words {
	fmt.Println(word.Weight, word.Word)
}
```

Optionally, we can also specify configurations such as retaining conjunctions at the beginning of sentences for our summary

```Go
doc, err := document.Create(text, &btrank.BiasedTextRank{}, &trank.KWTextRank{}, &parser.Parser{}, document.WithConjunctions())
```

## Things I Learned

This project was started to better familiarize myself with Go, and some best practices

* How to [structure](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1) your application
* How to idiomatically [handle errors](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)
* How to [style](https://github.com/uber-go/guide/blob/master/style.md#specifying-slice-capacity) your code

## Next Steps

Currently the project is more or less complete, with no major foreseeable updates. However, I'll be periodically updating the library as things come to mind.