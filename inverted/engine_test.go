package inverted

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEngine_NewIndex(t *testing.T) {
	e := New()

	// test empty list when no indexes
	assert.Equal(t, []string{}, e.IndexList())

	// test new index is listed
	_, err := e.NewIndex("test", nil)
	assert.Nil(t, err)
	assert.Equal(t, []string{"test"}, e.IndexList())

	// test multiple index list
	_, err = e.NewIndex("test2", nil)
	assert.Nil(t, err)

	assert.Contains(t, e.IndexList(), "test")
	assert.Contains(t, e.IndexList(), "test2")

	// test index already exists
	s, err := e.NewIndex("test", nil)
	assert.Error(t, err, IndexAlreadyExists)
	assert.Nil(t, s)

	// new index unknown field type
	s, err = e.NewIndex("ksksks", map[string]map[string]string{"a": {"type": "notexists"}})
	assert.Error(t, err, "unknown field type")
	assert.Nil(t, s)
}

func TestEngine_GetIndex(t *testing.T) {
	// index does not exist
	e := New()
	idx, err := e.GetIndex("a")
	assert.Nil(t, idx)
	assert.Equal(t, errors.New(IndexNotFound), err)
}

func TestEngine_IndexStats(t *testing.T) {
	// index does not exist
	e := New()
	idx, err := e.IndexStats("a")
	assert.Nil(t, idx)
	assert.Equal(t, errors.New(IndexNotFound), err)
}

func TestEngine_Index(t *testing.T) {
	// index does not exist
	e := New()
	err := e.Index("a", "b", map[string]interface{}{"a": "b"})
	assert.Equal(t, errors.New(IndexNotFound), err)
}

/*
func TestEngine_IndexStats_IndexContent(t *testing.T) {
	e := New()

	// test index does not exist
	s, err := e.IndexStats("test")
	assert.Nil(t, s)
	assert.EqualError(t, err, IndexNotFound)

	// composite index field not exist
	_, err = e.NewIndex("test33", nil)
	assert.Nil(t, err)
	buf := bytes.NewBufferString("a")
	err = e.IndexContent("test33", "a", "a", ioutil.NopCloser(buf))
	assert.Errorf(t, err, IndexNotFound)

	// test existing index with no documents
	_, err = e.NewIndex("test", map[string]map[string]string{"field": {"type": "text"}})
	assert.Nil(t, err)
	s, err = e.IndexStats("test")
	assert.Nil(t, err)
	assert.Equal(t, map[string]index.Stats{"field": {0, 0}}, s)

	// test index with 1 document
	b := ioutil.NopCloser(bytes.NewBufferString("a b c"))
	err = e.IndexContent("test", "a", "field", b)
	assert.Nil(t, err)

	s, err = e.IndexStats("test")
	assert.Nil(t, err)
	assert.Equal(t, map[string]index.Stats{"field": index.Stats{1, 3}}, s)

	// test index with 2 documents
	b = ioutil.NopCloser(bytes.NewBufferString("a b c d"))
	err = e.IndexContent("test", "b", "field", b)
	assert.Nil(t, err)
	s, err = e.IndexStats("test")
	assert.Nil(t, err)
	assert.Equal(t, map[string]index.Stats{"field": {2, 4}}, s)

	// test index content non existent index
	err = e.IndexContent("a", "a", "field", b)
	assert.Errorf(t, err, IndexNotFound)

	// index content reader error
	err = e.IndexContent("test", "a", "field", util.NewErrReadCloser())
	assert.Errorf(t, err, "test error")
}
*/

func TestEngine_DeleteIndex(t *testing.T) {
	e := New()

	// test delete non existing index err
	err := e.DeleteIndex("woops")
	assert.EqualError(t, err, IndexNotFound)

	// create index and delete
	_, err = e.NewIndex("deleteme", nil)
	assert.Nil(t, err)
	err = e.DeleteIndex("deleteme")
	assert.Nil(t, err)
	assert.Equal(t, []string{}, e.IndexList())
}

func TestEngine_Query(t *testing.T) {
	e := New()
	_, err := e.NewIndex(
		"films",
		map[string]map[string]string{
			"title":       {"type": "text"},
			"description": {"type": "text"},
			"genre":       {"type": "keyword"},
		},
	)
	assert.Nil(t, err)
	for i, content := range []map[string]interface{}{
		{
			"title":       "The Shawshank Redemption",
			"description": "The Shawshank Redemption is about",
			"genre":       "drama",
		},
		{
			"title":       "The Godfather",
			"description": "The Godfather is about",
			"genre":       "crime",
		},
		{
			"title":       "The Godfather, Part II",
			"description": "The Godfather, Part II is about",
			"genre":       "crime",
		},
		{
			"title":       "The Dark Knight",
			"description": "The Dark Knight is about",
			"genre":       "thriller",
		},
		{
			"title":       "12 Angry Men",
			"description": "12 Angry Men is about",
			"genre":       "crime",
		},
		{
			"title":       "Schindler's List",
			"description": "Schindler's List is about",
			"genre":       "biography",
		},
		{
			"title":       "Pulp Fiction",
			"description": "Pulp Fiction is about",
			"genre":       "crime",
		},
		{
			"title":       "The Lord of the Rings: The Return of the King",
			"description": "The Return of the King is about",
			"genre":       "action",
		},
		{
			"title":       "The Good, the Bad, and the Ugly",
			"description": "The Good, the Bad, and the Ugly is about",
			"genre":       "western",
		},
		{
			"title":       "Fight Club",
			"description": "Fight Club is about",
			"genre":       "drama",
		},
		{
			"title":       "The Lord of the Rings: The Fellowship of the Ring",
			"description": "The Fellowship of the Ring is ",
			"genre":       "action",
		},
		{
			"title":       "Forrest Gump",
			"description": "Forrest Gump is about",
			"genre":       "romance",
		},
	} {
		err := e.Index("films", string(i+1), content)
		assert.Nil(t, err)
	}

	for _, tcase := range []struct {
		name  string
		index string
		req   *SearchRequest
		want  map[string][]int
		err   error
	}{
		{
			"missing query",
			"films",
			&SearchRequest{},
			map[string][]int{},
			nil,
		},
		{
			"index does not exist",
			"ssss",
			&SearchRequest{},
			nil,
			errors.New(IndexNotFound),
		},
		{
			"match single term",
			"films",
			&SearchRequest{Query: &Query{Leaf: &MatchQuery{Field: "title", Term: "Godfather"}}},
			map[string][]int{"hits": {1}},
			nil,
		},
		{
			"match query multiple terms",
			"films",
			&SearchRequest{Query: &Query{Leaf: &MatchQuery{Field: "title", Term: "The"}}},
			map[string][]int{"hits": {0, 1, 2, 3, 7, 8, 10}},
			nil,
		},
		{
			"match phrase multiple terms",
			"films",
			&SearchRequest{Query: &Query{Leaf: &MatchPhraseQuery{Field: "title", Term: "The Lord"}}},
			map[string][]int{"hits": {7, 10}},
			nil,
		},
		{
			"multi match multiple terms",
			"films",
			&SearchRequest{Query: &Query{Leaf: &MultiMatchQuery{Fields: []string{"title", "description"}, Term: "of"}}},
			map[string][]int{"hits": {7, 10}},
			nil,
		},
		{
			"term query",
			"films",
			&SearchRequest{Query: &Query{Leaf: &TermQuery{Term: "action", Field: "genre"}}},
			map[string][]int{"hits": {7, 10}},
			nil,
		},
		{
			"terms query",
			"films",
			&SearchRequest{Query: &Query{Leaf: &TermsQuery{Terms: []string{"action", "western"}, Field: "genre"}}},
			map[string][]int{"hits": {7, 8, 10}},
			nil,
		},
		{
			"bool must match",
			"films",
			&SearchRequest{Query: &Query{BoolMust: []*BoolMustQuery{{Query: &Query{Leaf: &MatchQuery{Field: "title", Term: "The"}}}}}},
			map[string][]int{"hits": {0, 1, 2, 3, 7, 8, 10}},
			nil,
		},
		{
			"must match title The and must not match title Lord",
			"films",
			&SearchRequest{
				Query: &Query{
					BoolMust: []*BoolMustQuery{
						{
							Query: &Query{
								Leaf: &MatchQuery{Field: "title", Term: "The"},
							},
						},
					},
					BoolMustNot: []*BoolMustNotQuery{
						{
							Query: &Query{
								Leaf: &MatchQuery{Field: "title", Term: "Lord"},
							},
						},
					},
				},
			},
			map[string][]int{"hits": {0, 1, 2, 3, 8}},
			nil,
		},
		{
			"must match title The and must not match title Lord and filter genres crime or drama or thriller",
			"films",
			&SearchRequest{
				Query: &Query{
					BoolMust: []*BoolMustQuery{
						{
							Query: &Query{
								Leaf: &MatchQuery{Field: "title", Term: "The"},
							},
						},
					},
					BoolMustNot: []*BoolMustNotQuery{
						{
							Query: &Query{
								Leaf: &MatchQuery{Field: "title", Term: "Lord"},
							},
						},
					},
					BoolFilter: []*BoolFilterQuery{
						{
							Query: &Query{
								Leaf: &TermsQuery{Field: "genre", Terms: []string{"crime", "drama", "thriller"}},
							},
						},
					},
				},
			},
			map[string][]int{"hits": {0, 1, 2, 3}},
			nil,
		},
	} {
		r, err := e.Search(tcase.index, tcase.req)
		assert.Equal(t, tcase.want, r, tcase.name)
		assert.Equal(t, tcase.err, err)
	}
}
