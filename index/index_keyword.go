package index

import (
	"github.com/richardjennings/invertedindex/analyser"
	"sort"
)

const Keyword = "keyword"
const Text = "text"

type IndexKeyword struct {
	Analyser  analyser.KeywordAnalyser
	TermIndex map[string]int
	Terms     []map[int]int
}

type KeywordResult map[int]int

func (k KeywordResult) Docs() []int {
	var docs []int
	for d := range k {
		docs = append(docs, d)
	}
	sort.Ints(docs)
	return docs
}

// NewTextIndex creates a new index struct
func NewKeywordIndex() *IndexKeyword {
	index := IndexKeyword{}
	index.TermIndex = make(map[string]int)
	return &index
}

func (idx *IndexKeyword) Stats() (stats IdxStats) {
	stats.TermCount = len(idx.TermIndex)
	return stats
}

func (idx *IndexKeyword) Index(docId int, content interface{}) error {
	terms, err := idx.Analyser.Analyse(content)
	if err != nil {
		return err
	}
	for _, term := range terms {
		termId, ok := idx.TermIndex[term]
		if !ok {
			idx.TermIndex[term] = len(idx.Terms)
			idx.Terms = append(idx.Terms, map[int]int{docId: 1})
		} else {
			idx.Terms[termId][docId]++
		}
	}

	return nil
}

func (idx *IndexKeyword) TermQuery(query string) (KeywordResult, error) {
	termId, ok := idx.TermIndex[query]
	if !ok {
		return nil, nil
	}
	result := KeywordResult{}
	for docId := range idx.Terms[termId] {
		result[docId] = idx.Terms[termId][docId]
	}
	return result, nil
}

func (idx *IndexKeyword) TermsQuery(query []string) (KeywordResult, error) {
	result := make(KeywordResult)
	for _, term := range query {
		r, err := idx.TermQuery(term)
		if err != nil {
			return nil, err
		}
		for docId, v := range r {
			result[docId] = v
		}
	}
	return result, nil
}

func (idx *IndexKeyword) TermsAgg() (KeywordResult, error) {
	res := KeywordResult{}
	for t, d := range idx.Terms {
		res[t] = len(d)
	}
	return res, nil
}
