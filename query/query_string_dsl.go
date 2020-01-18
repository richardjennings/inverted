package query

import (
	"errors"
	"github.com/richardjennings/invertedindex/inverted"
	"net/url"
)

// Field => Query type => Query string
type CompositeQueryConf struct {
	Typ   int
	Query string
}

type CompositeQuery map[string]CompositeQueryConf

const (
	SingleTerm = iota
	Phrase
)

// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-query-string-query.html#query-string-syntax

// where the status Field contains active
// status:active

// where the author Field contains the exact phrase "john smith" (depending on tokenizer)
// author:"John Smith"

func ParseReqFromQuery(cfg url.Values) (*inverted.SearchRequest, error) {
	q := cfg.Get("q")
	if q == "" {
		return nil, errors.New("missing q in url Query params")
	}
	qs := []byte(q)
	var field string
	var term string
	var cont []byte
	var typ int

	for i, ch := range qs {
		if ch != ':' {
			field += string(ch)
		} else {
			cont = qs[i+1:]
			break
		}
	}

	if len(cont) == 0 {
		return nil, errors.New("only supporting Field:term")
	}

	typ = SingleTerm
	for i, ch := range cont {
		if ch == '"' {
			typ = Phrase
			continue
		}
		if i == len(cont)-1 && typ == Phrase && ch != '"' {
			return nil, errors.New(`unterminated "`)
		}
		term += string(ch)
	}

	switch typ {
	case Phrase:
		return &inverted.SearchRequest{Query: &inverted.Query{Leaf: &inverted.MatchPhraseQuery{Field: field, Term: term}}}, nil
	default:
		return &inverted.SearchRequest{Query: &inverted.Query{Leaf: &inverted.MatchQuery{Field: field, Term: term}}}, nil
	}
}
