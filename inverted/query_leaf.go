package inverted

import (
	"errors"
	"github.com/richardjennings/invertedindex/index"
	"sort"
)

type LeafQuery interface {
	Query(i *index.Index) (index.Result, error)
}

// full text Leaf
type (
	MatchQuery struct {
		Field string
		Term  string
	}
	MatchPhraseQuery struct {
		Field string
		Term  string
	}
	MultiMatchQuery struct {
		Fields []string
		Term   string
	}
)

// keyword Leaf
type (
	TermQuery struct {
		Field string
		Term  string
	}
	TermsQuery struct {
		Field string
		Terms []string
	}
)

type QueryResult map[int]struct{}

func (q QueryResult) Docs() []int {
	var result []int
	for d := range q {
		result = append(result, d)
	}
	sort.Ints(result)
	return result
}

// full text Leaf interface implementations
func (m MatchQuery) Query(cidx *index.Index) (index.Result, error) {
	idx, err := cidx.GetFieldIdx(m.Field)
	if err != nil {
		return nil, err
	}
	_, ok := idx.(index.Match)
	if !ok {
		return nil, errors.New("field does not support match queries")
	}
	return idx.(index.Match).MatchQuery(m.Term)
}

func (m MatchPhraseQuery) Query(cidx *index.Index) (index.Result, error) {
	idx, err := cidx.GetFieldIdx(m.Field)
	if err != nil {
		return nil, err
	}
	_, ok := idx.(index.Phrase)
	if !ok {
		return nil, errors.New("field does not support match phrase queries")
	}
	return idx.(index.Phrase).PhraseQuery(m.Term)
}

func (m MultiMatchQuery) Query(cidx *index.Index) (index.Result, error) {
	result := QueryResult{}
	for _, field := range m.Fields {
		idx, err := cidx.GetFieldIdx(field)
		if err != nil {
			return nil, err
		}
		_, ok := idx.(index.Match)
		if !ok {
			return nil, errors.New("field does not support multi match queries")
		}
		r, err := idx.(index.Match).MatchQuery(m.Term)
		if err != nil {
			return nil, err
		}
		for _, d := range r.Docs() {
			var t struct{}
			result[d] = t
		}
	}
	return result, nil
}

// keyword Leaf interface implementations
func (m TermQuery) Query(cidx *index.Index) (index.Result, error) {
	idx, err := cidx.GetFieldIdx(m.Field)
	if err != nil {
		return nil, err
	}
	_, ok := idx.(index.Term)
	if !ok {
		return nil, errors.New("field does not support term queries")
	}
	return idx.(index.Term).TermQuery(m.Term)
}

func (m TermsQuery) Query(cidx *index.Index) (index.Result, error) {
	idx, err := cidx.GetFieldIdx(m.Field)
	if err != nil {
		return nil, err
	}
	_, ok := idx.(index.Terms)
	if !ok {
		return nil, errors.New("field does not support terms queries")
	}
	return idx.(index.Terms).TermsQuery(m.Terms)
}
