package analyser

import (
	"bytes"
	"errors"
	"github.com/richardjennings/invertedindex/test"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestFullTextAnalyser_Analyse(t *testing.T) {
	a := FullTextAnalyser{}

	// with string
	r, err := a.Analyse("a b c")
	assert.Equal(t, []string{"a", "b", "c"}, r)
	assert.Nil(t, err)

	// with io.ReadCloser
	r, err = a.Analyse(ioutil.NopCloser(bytes.NewBufferString("d e f")))
	assert.Equal(t, []string{"d", "e", "f"}, r)
	assert.Nil(t, err)

	// with unsuported type
	r, err = a.Analyse(1)
	assert.Nil(t, r)
	assert.Equal(t, errors.New("string or io.ReadCloser type required"), err)

	// with failed io.ReadCloser read
	r, err = a.Analyse(test.NewErrReadCloser())
	assert.Nil(t, r)
	assert.Equal(t, errors.New("test error"), err)
}

func TestKeywordAnalyser_Analyse(t *testing.T) {
	a := KeywordAnalyser{}

	// with string
	r, err := a.Analyse("a b c")
	assert.Equal(t, []string{"a b c"}, r)
	assert.Nil(t, err)

	// with []string
	r, err = a.Analyse([]string{"d e f"})
	assert.Equal(t, []string{"d e f"}, r)
	assert.Nil(t, err)

	// with unsupported type
	r, err = a.Analyse(1)
	assert.Nil(t, r)
	assert.Equal(t, errors.New("expecting string or []string"), err)
}
