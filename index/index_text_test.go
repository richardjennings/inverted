package index

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestIndexText_Stats(t *testing.T) {
	type tcase struct {
		documents []string
		termCount int
	}
	tcases := []tcase{
		// test an empty index
		{[]string{}, 0},
		// test index with single document
		{[]string{"1 2 3"}, 3},
		// test index with multiple documents
		{[]string{"1 2 3", "1 2 4"}, 4},
	}
	for _, tcase := range tcases {
		index := NewTextIndex()
		if len(tcase.documents) > 0 {
			for docId, content := range tcase.documents {
				err := index.Index(docId, content)
				assert.Nil(t, err)
			}
		}
		want := IdxStats{
			TermCount: tcase.termCount,
		}
		have := index.Stats()
		assert.Equal(t, want, have)
	}
}

// test postings for given terms in a given document
// have the correct positions when indexed
func TestIndexTextDocumentTermPostings(t *testing.T) {
	index := NewTextIndex()
	err := index.Index(0, "1 2 3")
	assert.Nil(t, err)
	have := index.Terms[0][0]
	want := map[int]int{0: 1}
	assert.True(t, reflect.DeepEqual(want, have))
	have = index.Terms[1][0]
	want = map[int]int{1: 2}
	assert.True(t, reflect.DeepEqual(want, have))
	have = index.Terms[2][0]
	want = map[int]int{2: 0}
	assert.True(t, reflect.DeepEqual(want, have))
}

func TestIndexText_MatchQuery(t *testing.T) {
	type tcase struct {
		name      string
		documents []string
		query     string
		want      Result
		err       error
	}

	tcases := []tcase{
		{
			"test single term match query against single document",
			[]string{"mary had a little lamb"},
			"little",
			TermFreqResult{0: {1}},
			nil,
		},
		{
			"test multiple term match query against single document",
			[]string{"mary had a little little lamb"},
			"little lamb",
			TermFreqResult{0: {2, 1}},
			nil,
		},
		{
			"test multiple terms against more than one document",
			[]string{
				"once upon a time in a land far far away",
				"mary had a little lamb",
			},
			"once a",
			TermFreqResult{0: {1, 2}, 1: {0, 1}},
			nil,
		},
		{
			"test query term not in index",
			[]string{"a b c d"},
			"e",
			TermFreqResult{},
			nil,
		},
	}
	for _, tcase := range tcases {
		idx := NewTextIndex()
		for docId, content := range tcase.documents {
			err := idx.Index(docId, content)
			assert.Nil(t, err)
		}
		have, err := idx.MatchQuery(tcase.query)
		assert.Equal(t, tcase.want, have, tcase.name)
		assert.Equal(t, tcase.err, err)
	}
}

func TestIndexText_PhraseQuery(t *testing.T) {
	type tcase struct {
		name      string
		documents []string
		query     string
		want      PostingResult
	}

	tcases := []tcase{
		{
			"single document with single match",
			[]string{"once upon a time in a land far far away"},
			"a land far",
			PostingResult{0: {5}},
		},
		{
			"more than one document with multiple matches",
			[]string{
				"i like roast dinners",
				"i like all food",
			},
			"i like",
			PostingResult{0: {0}, 1: {0}},
		},
		{
			"missing terms return early",
			[]string{"a b c"},
			"e f",
			PostingResult{},
		},
		{
			"next term missing returns early",
			[]string{"d a b c"},
			"b c d",
			PostingResult{},
		},
		{
			"next term can match",
			[]string{"d a b c f"},
			"b c d",
			PostingResult{},
		},
		{
			"missing term in one document, out of order in other",
			[]string{
				"a b c",
				"1 2 3 d c",
			},
			"a c d",
			PostingResult{},
		},
		{
			"query can match the same document more than once",
			[]string{
				"i like roast dinners i like lots of things",
				"i like all food",
			},
			"i like",
			PostingResult{0: {0, 4}, 1: {0}},
		},
	}

	for _, tcase := range tcases {
		idx := NewTextIndex()
		for docId, content := range tcase.documents {
			err := idx.Index(docId, content)
			assert.Nil(t, err, tcase.name)
		}
		have, err := idx.PhraseQuery(tcase.query)
		assert.Nil(t, err)
		// sort result postings so comparison works
		// do not want to sort postings in query funcs (yet)
		for k := range have {
			sort.Ints(have[k])
		}
		assert.True(t, reflect.DeepEqual(tcase.want, have), tcase.name)
	}
}

func BenchmarkIndexText_MatchQuery(b *testing.B) {

	file, err := os.Open("../../test/corpus/the-comedy-of-errors.txt")
	if err != nil {
		panic(err)
	}
	cidx, err := NewIndex(map[string]map[string]string{"txt": {"type": "text"}})
	if err != nil {
		b.Error(err)
	}
	err = cidx.Index("the-comedy-of-errors", map[string]interface{}{"txt": ioutil.NopCloser(file)})
	assert.Nil(b, err)

	// run benchmark
	for n := 0; n < b.N; n++ {
		_, err = cidx.Idxs["txt"].(Match).MatchQuery("in")
		assert.Nil(b, err)
	}
}

func BenchmarkIndexText_PhraseQuery(b *testing.B) {

	file, err := os.Open("../../test/corpus/the-comedy-of-errors.txt")
	if err != nil {
		panic(err)
	}
	cidx, err := NewIndex(map[string]map[string]string{"txt": {"type": "text"}})
	if err != nil {
		b.Error(err)
	}
	err = cidx.Index("the-comedy-of-errors", map[string]interface{}{"txt": ioutil.NopCloser(file)})
	assert.Nil(b, err)

	benchmarks := []string{
		"in this",
		"against accepting unsolicited donations",
		"By rushing in their houses, bearing thence",
	}

	// run benchmarks
	for _, q := range benchmarks {
		b.Run(q, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_, err := cidx.Idxs["txt"].(Phrase).PhraseQuery(q)
				assert.Nil(b, err)
			}
		})
	}
}

func TestTermFreqResult_Docs(t *testing.T) {
	p := TermFreqResult{3: {2, 6, 3}, 2: {1, 7, 9}, 1: {3, 1, 2}}
	assert.Equal(t, []int{1, 2, 3}, p.Docs())
}

func TestPostingResult_Docs(t *testing.T) {
	p := PostingResult{3: {2, 6, 3}, 2: {1, 7, 9}, 1: {3, 1, 2}}
	assert.Equal(t, []int{1, 2, 3}, p.Docs())
}
