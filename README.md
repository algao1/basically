# basically

`basically` is a Go implementation of the TextRank and Biased TextRank algorithm built on [`prose`](https://github.com/jdkato/prose). It provides fully unsupervised methods for keyword extraction and focused text summarization, along with some additional quality of life features over the original implementations.

## Methods

First, the document is parsed into its constituent sentences and words using a sentence segmenter and tokenizer. Sentiment values are assigned to individual sentences, and tokens are annotated with part of speech tags.

For **keyword extraction**, all words that pass the syntactic filter are added to a undirected, weighted graph, and an edge is added between words that co-occur within a window of *N* words. The edge weight is set to be inversely proportional to the distance between the words. Each vertex is assigned an initial score of 1, and the following ranking algorithm is run on the graph

<div align="center">
	<img src="https://latex.codecogs.com/svg.latex?TR(V_i)&space;=&space;(1&space;-&space;d)&space;&plus;&space;d&space;\times&space;\sum_{V_j&space;\in&space;In(V_i)}&space;\frac{w_{ji}}{\sum_{V_k&space;\in&space;Out(v_j)}&space;w_{jk}}&space;TR(V_j)" title="TR(V_i) = (1 - d) + d \times \sum_{V_j \in In(V_i)} \frac{w_{ji}}{\sum_{V_k \in Out(v_j)} w_{jk}} TR(V_j)"/>
</div>

During post-processing, adjacent keywords are collapsed into a multi-word keyword, and the top keywords are then extracted.

For **sentence extraction**, every sentence is added to a undirected, weighted graph, with an edge between sentences that share *common content*. The edge weight is set simply as the number of common tokens between the lexical representations of the two sentences. Each vertex is also assigned an initial score of 1, and a bias score based on the focus text, before the following ranking algorithm is run on the graph

<div align="center">
	<img src="https://latex.codecogs.com/svg.latex?BTR(V_i)&space;=&space;Bias&space;\times&space;(1&space;-&space;d)&space;&plus;&space;d&space;\times&space;\sum_{V_j&space;\in&space;In(V_i)}&space;\frac{w_{ji}}{\sum_{V_k&space;\in&space;Out(v_j)}&space;w_{jk}}&space;BTR(V_j)" title="BTR(V_i) = Bias \times (1 - d) + d \times \sum_{V_j \in In(V_i)} \frac{w_{ji}}{\sum_{V_k \in Out(v_j)} w_{jk}} BTR(V_j)"/>
</div>

The top weighted sentences are then selected and sorted in chronological order to form a summary.

Further information on the two algorithms can be found [here](https://web.eecs.umich.edu/~mihalcea/papers/mihalcea.emnlp04.pdf) and [here](https://arxiv.org/pdf/2011.01026.pdf).

## Installation

```shell
go get github.com/algao1/basically
``` 

## Usage

Initialization:

```Go
// Instantiate the summarizer, highlighter, and parser.
s := &btrank.BiasedTextRank{}
h := &trank.KWTextRank{}
p := parser.Create()

// Instantiate a document for every given text.
doc, err := document.Create(text, s, h, p)
if err != nil {
	log.Fatal(err)
}
```

Text Summarization:

```Go
// Summarize the document into 7 sentences, with no threshold value, and with respect to a focus sentence.
sents, err := doc.Summarize(7, 0, focus)
if err != nil {
	log.Fatal(err)
}

for _, sent := range sents {
	fmt.Printf("[%.2f, %.2f] %s\n", sent.Score, sent.Sentiment, sent.Raw)
}
```

Keyword Extraction:

```Go
// Highlight the top 7 keywords in the document, with multi-word keywords enabled.
words, err := doc.Highlight(7, true)
if err != nil {
	log.Fatal(err)
}

for _, word := range words {
	fmt.Println(word.Weight, word.Word)
}
```

Optionally, we can specify configurations such as retaining conjunctions at the beginning of sentences for our summary

```Go
doc, err := document.Create(text, s, h, p, document.WithConjunctions())
```

## Benchmarks

### Text Summarization & Keyword Extraction

Below is a rudimentary comparison of `basically`'s performance against other implementations using news articles from *The Guardian*:

| Library     | Language    | Avg Speed |
| ----------- | ----------- | --------- |
| summa       | Python      | 1.67s     |
| basically   | Go          | 0.48s     |

## Things I Learned

This project was started to better familiarize myself with Go, and some best practices

* How to idiomatically [structure](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1) applications
* How to idiomatically [handle errors](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)
* How to [style and format](https://github.com/uber-go/guide/blob/master/style.md#specifying-slice-capacity) Go code
* How to [test and benchmark](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)

## Next Steps

Currently the project is more or less complete, with no major foreseeable updates. However, I'll be periodically updating the library as things come to mind.
