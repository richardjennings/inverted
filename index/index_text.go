package index

import (
	"github.com/richardjennings/invertedindex/analyser"
	"sort"
)

type IndexText struct {
	TermIndex map[string]int
	Terms     []map[int]map[int]int
	Analyser  analyser.FullTextAnalyser
}

// NewTextIndex creates a new index struct
func NewTextIndex() *IndexText {
	index := IndexText{}
	index.TermIndex = make(map[string]int)
	return &index
}

type TermFreqResult map[int][]int

func (p TermFreqResult) Docs() []int {
	var docs []int
	for d := range p {
		docs = append(docs, d)
	}
	// until score sort by id for testing consistency
	sort.Ints(docs)
	return docs
}

type PostingResult map[int][]int

func (p PostingResult) Docs() []int {
	var docs []int
	for d := range p {
		docs = append(docs, d)
	}
	// until score sort by id for testing consistency
	sort.Ints(docs)
	return docs
}

// Stats provides some information about
// the terms and documents in the inverted index
func (idx *IndexText) Stats() (stats IdxStats) {
	stats.TermCount = len(idx.TermIndex)
	return stats
}

// IndexDocument adds a document to the inverted index
func (idx *IndexText) Index(docId int, content interface{}) error {
	terms, err := idx.Analyser.Analyse(content)
	if err != nil {
		return err
	}
	pretid := -1

	for j := 0; j < len(terms); j++ {

		// look up term id
		tid, ok := idx.TermIndex[terms[j]]
		if !ok {
			tid = len(idx.TermIndex)
			idx.TermIndex[terms[j]] = tid
			d := make(map[int]map[int]int)
			idx.Terms = append(idx.Terms, d)
		}

		// update previous term with next (this) term
		if pretid != -1 {
			idx.Terms[pretid][docId][j-1] = tid
		}

		if _, ok := idx.Terms[tid][docId]; !ok {
			idx.Terms[tid][docId] = make(map[int]int)
		}

		// set term doc pos with placeholder next tid
		idx.Terms[tid][docId][j] = 0
		pretid = tid
	}

	return nil
}

// MatchQuery looks up a terms in the inverted index and returns
// the associated documents and position data
func (idx IndexText) MatchQuery(query string) (TermFreqResult, error) {
	//nilResult := make(map[int]map[int]int)
	result := make(TermFreqResult)

	// tokenize the query
	terms, err := idx.Analyser.Analyse(query)
	if err != nil {
		return nil, err
	}

	for i, term := range terms {
		// find term in index
		termId, ok := idx.TermIndex[term]
		if !ok {
			// if it was AND, return nil result
			//return result, nil
			// but default or atm so continue
			continue
		}

		for docId, postings := range idx.Terms[termId] {
			if len(result[docId]) < i {
				// pad 0 counts for previous terms that did not match
				for j := 0; j < i; j++ {
					result[docId] = append(result[docId], 0)
				}
			}
			result[docId] = append(result[docId], len(postings))
		}
	}

	// return result
	return result, nil
}

// phraseQuery uses position data in the inverted index
// to find documents that have a sequence of terms
// where the positions of terms in order is incremental
func (idx IndexText) PhraseQuery(query string) (PostingResult, error) {
	// initialize result
	//result := make(map[int][]int)
	result := make(PostingResult)

	// tokenize query
	terms, err := idx.Analyser.Analyse(query)
	if err != nil {
		return nil, err
	}
	lenTerms := len(terms)

	// look up termIds from index
	var termIds = make(map[int]int)
	for i, t := range terms {
		id, ok := idx.TermIndex[t]
		if !ok {
			return result, nil
		}
		termIds[i] = id
	}

OUTER:
	// iterate all document matches for term 0
	for docId, postings := range idx.Terms[termIds[0]] {
		// check match for document against all terms
		for i := 1; i < lenTerms; i++ {
			_, ok := idx.Terms[termIds[i]][docId]
			if !ok {
				continue OUTER
			}
		}
	POSTING:
		// now check that for this document, there is a position match for each term posting
		for posting, ntid := range postings {
			// use the next term info to skip if needed
			if ntid != termIds[1] {
				continue
			}

			for i := 1; i < lenTerms; i++ {
				// check for term in position
				ntid, ok := idx.Terms[termIds[i]][docId][posting+i]
				if !ok {
					continue POSTING
				}

				// if last iteration, we have a result
				if i == lenTerms-1 {

					result[docId] = append(result[docId], posting)
					continue
				}

				// check next term can match
				if ntid != termIds[i+1] {
					continue POSTING
				}
			}
		}
	}

	for _, r := range result {
		sort.Ints(r)
	}

	return result, nil
}
