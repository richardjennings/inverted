package inverted

import (
	"errors"
	"github.com/richardjennings/invertedindex/index"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatchQuery_Query(t *testing.T) {
	cidx, err := index.NewIndex(map[string]map[string]string{"test": {"type": "text"}, "test2": {"type": "keyword"}})
	assert.Nil(t, err)
	err = cidx.Index("0", map[string]interface{}{"test": "some content"})
	assert.Nil(t, err)

	// query result
	q := MatchQuery{"test", "some"}
	r, err := q.Query(cidx)
	assert.Nil(t, err)
	assert.Equal(t, index.TermFreqResult{0: {1}}, r)

	// error index not exists
	q = MatchQuery{"notexists", "some"}
	r, err = q.Query(cidx)
	assert.Nil(t, r)
	assert.Equal(t, errors.New("field not found"), err)

	// error field does not support match queries
	q = MatchQuery{"test2", "some"}
	r, err = q.Query(cidx)
	assert.Nil(t, r)
	assert.Equal(t, errors.New("field does not support match queries"), err)
}

func TestMatchPhraseQuery_Query(t *testing.T) {
	cidx, err := index.NewIndex(map[string]map[string]string{"test": {"type": "text"}, "test2": {"type": "keyword"}})
	assert.Nil(t, err)
	err = cidx.Index("0", map[string]interface{}{"test": "there is some content"})
	assert.Nil(t, err)

	// query result
	q := MatchPhraseQuery{"test", "is some"}
	r, err := q.Query(cidx)
	assert.Nil(t, err)
	assert.Equal(t, index.PostingResult{0: {1}}, r)

	// error index not exists
	q = MatchPhraseQuery{"notexists", "is some"}
	r, err = q.Query(cidx)
	assert.Nil(t, r)
	assert.Equal(t, errors.New("field not found"), err)

	// error field does not support match queries
	q = MatchPhraseQuery{"test2", "is some"}
	r, err = q.Query(cidx)
	assert.Nil(t, r)
	assert.Equal(t, errors.New("field does not support match phrase queries"), err)
}

func TestMultiMatchQuery_Query(t *testing.T) {
	cidx, err := index.NewIndex(
		map[string]map[string]string{
			"test":  {"type": "text"},
			"test2": {"type": "keyword"},
			"test3": {"type": "text"},
		},
	)
	assert.Nil(t, err)

	err = cidx.Index(
		"0",
		map[string]interface{}{
			"test":  "there is some",
			"test2": "keyword",
			"test3": "and some more text",
		},
	)
	assert.Nil(t, err)

	err = cidx.Index(
		"1",
		map[string]interface{}{
			"test":  "there is some content",
			"test2": "keyword",
			"test3": "and some more text content",
		},
	)
	assert.Nil(t, err)

	// query result
	q := MultiMatchQuery{[]string{"test", "test3"}, "content"}
	r, err := q.Query(cidx)
	assert.Nil(t, err)
	assert.Equal(t, QueryResult{1: struct{}{}}, r)

	// error index not exists
	q = MultiMatchQuery{[]string{"notexists"}, "is some"}
	r, err = q.Query(cidx)
	assert.Nil(t, r)
	assert.Equal(t, errors.New("field not found"), err)

	// error field does not support match queries
	q = MultiMatchQuery{[]string{"test2"}, "is some"}
	r, err = q.Query(cidx)
	assert.Nil(t, r)
	assert.Equal(t, errors.New("field does not support multi match queries"), err)
}

func TestTermQuery_Query(t *testing.T) {
	cidx, err := index.NewIndex(map[string]map[string]string{"test": {"type": "text"}, "test2": {"type": "keyword"}})
	assert.Nil(t, err)
	err = cidx.Index("0", map[string]interface{}{"test2": "a keyword"})
	assert.Nil(t, err)

	// query result
	q := TermQuery{"test2", "a keyword"}
	r, err := q.Query(cidx)
	assert.Nil(t, err)
	assert.Equal(t, index.KeywordResult{0: 1}, r)

	// error index not exists
	q = TermQuery{"notexists", "some"}
	r, err = q.Query(cidx)
	assert.Nil(t, r)
	assert.Equal(t, errors.New("field not found"), err)

	// error field does not support match queries
	q = TermQuery{"test", "some"}
	r, err = q.Query(cidx)
	assert.Nil(t, r)
	assert.Equal(t, errors.New("field does not support term queries"), err)
}

func TestTermsQuery_Query(t *testing.T) {
	cidx, err := index.NewIndex(map[string]map[string]string{"test": {"type": "text"}, "test2": {"type": "keyword"}})
	assert.Nil(t, err)
	err = cidx.Index("0", map[string]interface{}{"test2": "keyword1"})
	assert.Nil(t, err)
	err = cidx.Index("1", map[string]interface{}{"test2": "keyword2"})
	assert.Nil(t, err)
	err = cidx.Index("2", map[string]interface{}{"test2": "keyword3"})
	assert.Nil(t, err)

	// query result
	q := TermsQuery{"test2", []string{"keyword1", "keyword2"}}
	r, err := q.Query(cidx)
	assert.Nil(t, err)
	assert.Equal(t, index.KeywordResult{0: 1, 1: 1}, r)

	// error index not exists
	q = TermsQuery{"notexists", []string{"some"}}
	r, err = q.Query(cidx)
	assert.Nil(t, r)
	assert.Equal(t, errors.New("field not found"), err)

	// error field does not support match queries
	q = TermsQuery{"test", []string{"some"}}
	r, err = q.Query(cidx)
	assert.Nil(t, r)
	assert.Equal(t, errors.New("field does not support terms queries"), err)
}

func TestQueryResult_Docs(t *testing.T) {
	var s struct{}
	q := QueryResult{1: s, 2: s, 3: s, 4: s}
	assert.Equal(t, []int{1, 2, 3, 4}, q.Docs())
}
