package query

import (
	"errors"
	"github.com/richardjennings/invertedindex/inverted"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseReq(t *testing.T) {
	tcases := []struct {
		name string
		body string
		want *inverted.SearchRequest
		err  error
	}{
		{
			"malformed json request",
			`{`,
			nil,
			errors.New("unexpected end of JSON input"),
		},
		{
			"invalid json request",
			`true`,
			nil,
			errors.New("unknown type"),
		},
		{
			"Query not map",
			`{"query":true}`,
			nil,
			errors.New("expected map"),
		},
		{
			"Query key not map",
			`{"query":{"a":"true"}}`,
			nil,
			errors.New("expected map"),
		},
		{
			"Query key empty map",
			`{"query":{"a":{}}}`,
			nil,
			errors.New("no Field/value specified"),
		},
		{
			"invalid Query component",
			`{"query":{"notathing":{"s":"s"}}}`,
			nil,
			errors.New("unknown key"),
		},
		/*{
			"invalid bool Query not map",
			`{"query":{"bool":{"notathing":"ok"}}}`,
			nil,
			errors.New("expected map"),
		},*/
		{
			"invalid bool Query type",
			`{"query":{"bool":{"notathing":{"term":{"field":"value"}}}}}`,
			nil,
			errors.New("unknown key"),
		},
		{
			"invalid term Query type",
			`{"query":{"term":{"true": true}}}`,
			nil,
			errors.New("expected string"),
		},
		{
			"invalid terms Query value not array",
			`{"query":{"terms":{"field": "notarray"}}}`,
			nil,
			errors.New("invalid value for terms"),
		},
		{
			"invalid Match Query map missing Query",
			`{"query":{"match":{"field":{"true":"true"}}}}`,
			nil,
			errors.New("missing query"),
		},
		{
			"invalid Match Query map Query not string",
			`{"query":{"match":{"field":{"query":true}}}}`,
			nil,
			errors.New("expected string"),
		},
		{
			"invalid Match Query not map not string",
			`{"query":{"match":{"field":true}}}`,
			nil,
			errors.New("invalid"),
		},
		{
			"invalid match_phrase Query missing Query",
			`{"query":{"match_phrase":{"field":{}}}}`,
			nil,
			errors.New("missing query"),
		},
		{
			"invalid match_phrase Query Query not string",
			`{"query":{"match_phrase":{"field":{"query": true}}}}`,
			nil,
			errors.New("expected string"),
		},
		{
			"invalid match_phrase Query not map",
			`{"query":{"match_phrase":{"field":true}}}`,
			nil,
			errors.New("invalid"),
		},
		{
			"invalid match_phrase_prefix Query missing Query",
			`{"query":{"match_phrase_prefix":{"field":{}}}}`,
			nil,
			errors.New("unknown key"),
		},
		{
			"invalid match_phrase_prefix Query Query not string",
			`{"query":{"match_phrase_prefix":{"field":{"query":true}}}}`,
			nil,
			errors.New("unknown key"),
		},
		{
			"invalid match_phrase_prefix Query not map",
			`{"query":{"match_phrase_prefix":{"field":true}}}`,
			nil,
			errors.New("unknown key"),
		},
		{
			"invalid multi_match Query missing Fields declaration",
			`{"query":{"multi_match":{"query":"test"}}}`,
			nil,
			errors.New("missing field declaration"),
		},
		{
			"invalid multi_match Query invalid Fields value type",
			`{"query":{"multi_match":{"fields":[1]}}}`,
			nil,
			errors.New("expecting string"),
		},
		{
			"invalid multi_match Query invalid Fields not array",
			`{"query":{"multi_match":{"fields":{"a":"b"}}}}`,
			nil,
			errors.New("fields must be an array"),
		},
		{
			"invalid multi_match Query missing Query",
			`{"query":{"multi_match":{"fields":["a","b"]}}}`,
			nil,
			errors.New("missing query declaration"),
		},
		{
			"invalid multi_match Query Query not string",
			`{"query":{"multi_match":{"fields":["a","b"], "query": true}}}`,
			nil,
			errors.New("query should be a string"),
		},
		{
			"bool must term",
			`{"query":{"bool":{"must":{"term":{"a":"b"}}}}}`,
			&inverted.SearchRequest{
				Query: &inverted.Query{
					BoolMust: []*inverted.BoolMustQuery{
						{
							Query: &inverted.Query{
								Leaf: &inverted.TermQuery{
									Field: "a",
									Term:  "b",
								},
							},
						},
					},
				},
			},
			nil,
		},
		{
			"bool must multiple",
			`{
				"query":{
					"bool":{
						"must":[
							{"term":{"a":"b"}},
							{"term":{"c":"d"}}
						]
					}
				}
			}`,
			&inverted.SearchRequest{
				Query: &inverted.Query{
					BoolMust: []*inverted.BoolMustQuery{
						{
							Query: &inverted.Query{
								Leaf: &inverted.TermQuery{
									Field: "a",
									Term:  "b",
								},
							},
						},
						{
							Query: &inverted.Query{
								Leaf: &inverted.TermQuery{
									Field: "c",
									Term:  "d",
								},
							},
						},
					},
				},
			},
			nil,
		},
		{
			"bool must not Match",
			`{"query":{"bool":{"must_not":{"match":{"a":"b"}}}}}`,
			&inverted.SearchRequest{
				Query: &inverted.Query{
					BoolMustNot: []*inverted.BoolMustNotQuery{
						{
							Query: &inverted.Query{
								Leaf: &inverted.MatchQuery{
									Field: "a",
									Term:  "b",
								},
							},
						},
					},
				},
			},
			nil,
		},
		{
			"bool must not multiple",
			`{
				"query":{
					"bool":{
						"must_not":[
							{"match":{"a":"b"}},
							{"match":{"b":"c"}}
						]
					}
				}
			}`,
			&inverted.SearchRequest{
				Query: &inverted.Query{
					BoolMustNot: []*inverted.BoolMustNotQuery{
						{
							Query: &inverted.Query{
								Leaf: &inverted.MatchQuery{
									Field: "a",
									Term:  "b",
								},
							},
						},
						{
							Query: &inverted.Query{
								Leaf: &inverted.MatchQuery{
									Field: "b",
									Term:  "c",
								},
							},
						},
					},
				},
			},
			nil,
		},
		{
			"bool should match_phrase",
			`{"query":{"bool":{"should":{"match_phrase":{"a":"b c"}}}}}`,
			&inverted.SearchRequest{
				Query: &inverted.Query{
					BoolShould: []*inverted.BoolShouldQuery{
						{
							Query: &inverted.Query{
								Leaf: &inverted.MatchPhraseQuery{
									Field: "a",
									Term:  "b c",
								},
							},
						},
					},
				},
			},
			nil,
		},
		{
			"bool should multiple",
			`{
				"query":{
					"bool":{
						"should": [
							{"match_phrase":{"a":"b c"}},
							{"match":{"a": "b"}}
						]
					}
				}
			}`,
			&inverted.SearchRequest{
				Query: &inverted.Query{
					BoolShould: []*inverted.BoolShouldQuery{
						{
							Query: &inverted.Query{
								Leaf: &inverted.MatchPhraseQuery{
									Field: "a",
									Term:  "b c",
								},
							},
						},
						{
							Query: &inverted.Query{
								Leaf: &inverted.MatchQuery{
									Field: "a",
									Term:  "b",
								},
							},
						},
					},
				},
			},
			nil,
		},
		{
			"bool filter terms",
			`{"query":{"bool":{"filter":{"terms":{"a":["b", "c"]}}}}}`,
			&inverted.SearchRequest{
				Query: &inverted.Query{
					BoolFilter: []*inverted.BoolFilterQuery{
						{
							Query: &inverted.Query{
								Leaf: &inverted.TermsQuery{
									Field: "a",
									Terms: []string{"b", "c"},
								},
							},
						},
					},
				},
			},
			nil,
		},
		{
			"bool filter multiple",
			`{
				"query": {
    				"bool": {
      					"filter": [
							{ "term": { "color": "red"   }},
        					{ "term": { "brand": "gucci" }}
      					]
    				}
  				}
			}`,
			&inverted.SearchRequest{
				Query: &inverted.Query{
					BoolFilter: []*inverted.BoolFilterQuery{
						{
							Query: &inverted.Query{
								Leaf: &inverted.TermQuery{
									Field: "color",
									Term:  "red",
								},
							},
						},
						{
							Query: &inverted.Query{
								Leaf: &inverted.TermQuery{
									Field: "brand",
									Term:  "gucci",
								},
							},
						},
					},
				},
			},
			nil,
		},
		{
			"term",
			`{"query":{"term":{"a":"b"}}}`,
			&inverted.SearchRequest{
				Query: &inverted.Query{
					Leaf: &inverted.TermQuery{
						Field: "a",
						Term:  "b",
					},
				},
			},
			nil,
		},
		{
			"terms",
			`{"query":{"terms":{"a":["b", "c"]}}}`,
			&inverted.SearchRequest{
				Query: &inverted.Query{
					Leaf: &inverted.TermsQuery{
						Field: "a",
						Terms: []string{"b", "c"},
					},
				},
			},
			nil,
		},
		{
			"Match inline",
			`{"query":{"match":{"field": "query"}}}`,
			&inverted.SearchRequest{
				Query: &inverted.Query{
					Leaf: &inverted.MatchQuery{
						Field: "field",
						Term:  "query",
					},
				},
			},
			nil,
		},
		{
			"Match composite",
			`{"query":{"match":{"field": {"query": "thequery"}}}}`,
			&inverted.SearchRequest{
				Query: &inverted.Query{
					Leaf: &inverted.MatchQuery{
						Field: "field",
						Term:  "thequery",
					},
				},
			},
			nil,
		},
		{
			"match_phrase inline",
			`{"query":{"match_phrase":{"field": "query"}}}`,
			&inverted.SearchRequest{
				Query: &inverted.Query{
					Leaf: &inverted.MatchPhraseQuery{
						Field: "field",
						Term:  "query",
					},
				},
			},
			nil,
		},
		{
			"match_phrase composite",
			`{"query":{"match_phrase":{"field": {"query": "thequery"}}}}`,
			&inverted.SearchRequest{
				Query: &inverted.Query{
					Leaf: &inverted.MatchPhraseQuery{
						Field: "field",
						Term:  "thequery",
					},
				},
			},
			nil,
		},
		{
			"multi_match",
			`{"query":{"multi_match":{"fields":["a","b"],"query": "thequery"}}}`,
			&inverted.SearchRequest{
				Query: &inverted.Query{
					Leaf: &inverted.MultiMatchQuery{
						Fields: []string{"a", "b"},
						Term:   "thequery",
					},
				},
			},
			nil,
		},
	}

	for _, tcase := range tcases {
		have, err := ParseReq([]byte(tcase.body))
		if tcase.err != nil {
			assert.Equal(t, tcase.err.Error(), err.Error(), tcase.name)
		}
		assert.Equal(t, tcase.want, have, tcase.name)
	}
}
