package index

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewKeywordIndex(t *testing.T) {
	idx := NewKeywordIndex()

	err := idx.Index(0, "test")
	assert.Nil(t, err)

	err = idx.Index(1, "test2")
	assert.Nil(t, err)

	r, err := idx.TermQuery("test")
	assert.Nil(t, err)
	assert.Equal(t, KeywordResult{0: 1}, r)

	r, err = idx.TermsQuery([]string{"test", "test2"})
	assert.Nil(t, err)
	assert.Equal(t, KeywordResult{0: 1, 1: 1}, r)

	s := idx.Stats()
	assert.Equal(t, IdxStats{TermCount: 2}, s)

	// analyzer error
	err = idx.Index(1, 1)
	assert.Equal(t, errors.New("expecting string or []string"), err)

	// cumulative count
	idx = NewKeywordIndex()
	_ = idx.Index(0, []string{"a", "a", "a"})
	assert.Equal(t, 3, idx.Terms[0][0])
}

func TestKeywordResult_Docs(t *testing.T) {
	res := KeywordResult{}
	var empty []int
	assert.Equal(t, empty, res.Docs())

	res = KeywordResult{1: 1}
	assert.Equal(t, []int{1}, res.Docs())
}

func TestIndexKeyword_TermsAgg(t *testing.T) {
	idx := NewKeywordIndex()

	err := idx.Index(0, "test")
	assert.Nil(t, err)

	have, err := idx.TermsAgg()
	assert.Nil(t, err)

	assert.Equal(t, KeywordResult{0: 1}, have)
}
