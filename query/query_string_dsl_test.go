package query

import (
	"github.com/richardjennings/invertedindex/inverted"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestQsDSLFactory(t *testing.T) {
	have, err := ParseReqFromQuery(url.Values{})
	assert.Nil(t, have)
	assert.Errorf(t, err, "missing q in url Query params")

	have, err = ParseReqFromQuery(url.Values{"q": {"Field:value"}})
	assert.Nil(t, err)
	assert.Equal(t, &inverted.SearchRequest{Query: &inverted.Query{Leaf: &inverted.MatchQuery{Field: "Field", Term: "value"}}}, have)

	have, err = ParseReqFromQuery(url.Values{"q": {`Field:"value"`}})
	assert.Equal(t, &inverted.SearchRequest{Query: &inverted.Query{Leaf: &inverted.MatchPhraseQuery{Field: "Field", Term: "value"}}}, have)
	assert.Nil(t, err)

	have, err = ParseReqFromQuery(url.Values{"q": {"value"}})
	assert.Nil(t, have)
	assert.Errorf(t, err, "only supporting Field:term")

	have, err = ParseReqFromQuery(url.Values{"q": {`a:"a"s`}})
	assert.Nil(t, have)
	assert.Errorf(t, err, `" expected last"`)

	have, err = ParseReqFromQuery(url.Values{"q": {`a:"a`}})
	assert.Nil(t, have)
	assert.Errorf(t, err, `" expected last"`)
}
