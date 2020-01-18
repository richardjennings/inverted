package query

import (
	"github.com/richardjennings/invertedindex/inverted"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseAgg(t *testing.T) {
	for _, tcase := range []struct {
		name string
		body string
		want *inverted.SearchRequest
		err  error
	}{
		{
			"simple filter aggregation",
			`{"aggs":{"bucket_name":{"filter":{"bool":{"must":{"term":{"test":"a"}}}}}}}`,
			&inverted.SearchRequest{
				Agg: &inverted.Aggregation{
					Aggregations: map[string]*inverted.Aggregation{
						"bucket_name": {
							Filter: &inverted.Query{
								BoolMust: []*inverted.BoolMustQuery{
									{
										Query: &inverted.Query{
											Leaf: &inverted.TermQuery{
												Field: "test",
												Term:  "a",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			nil,
		},
		{
			"simple term aggregation",
			`{"aggs":{"bucket_name":{"terms":{"field":"test"}}}}`,
			&inverted.SearchRequest{
				Agg: &inverted.Aggregation{
					Aggregations: map[string]*inverted.Aggregation{
						"bucket_name": {
							Aggs: []inverted.Agg{
								inverted.TermsAgg{Field: "test"},
							},
						},
					},
				},
			},
			nil,
		},
		{
			// select * where test = a group by b
			"term aggregation with bool filter",
			`
		{
			"aggs":{
				"bucket_name":{
					"filter":{
						"bool":{
							"must":{
								"term":{
									"test": "a"
								}
							}
						}
					},
					"aggs":{
						"terms":{
							"field":"b"
						}
					}
				}
			}
		}`,
			&inverted.SearchRequest{
				Agg: &inverted.Aggregation{
					Aggregations: map[string]*inverted.Aggregation{
						"bucket_name": {
							Filter: &inverted.Query{
								BoolMust: []*inverted.BoolMustQuery{
									{
										Query: &inverted.Query{
											Leaf: &inverted.TermQuery{
												Field: "test",
												Term:  "a",
											},
										},
									},
								},
							},
							Aggs: []inverted.Agg{
								inverted.TermsAgg{Field: "b"},
							},
						},
					},
				},
			},
			nil,
		},
		{
			"multi level term aggregations with filters",
			`
{
	"aggs":{
		"bucket_name":{
			"filter":{
				"bool":{
					"must":{
						"term":{
							"test": "a"
						}
					}
				}
			},
			"terms": {
				"field": "c"
			},
			"aggs":{
				"some-name": {
					"filter": {
						"bool": {
							"must":{
								"term":{
									"test": "b"
								}
							}
						}
					},
					"terms": {
						"field": "d"
					},
					"aggs": {
						"another-name": {
							"filter": {
								"bool": {
									"must":{
										"term":{
											"test2": "c"
										}
									}
								}
							},
							"terms":{
								"field":"test"
							}
						}
					}
				}
			}
		}
		
	}
}`,
			&inverted.SearchRequest{
				Agg: &inverted.Aggregation{
					Aggregations: map[string]*inverted.Aggregation{
						"bucket_name": {
							Filter: &inverted.Query{
								BoolMust: []*inverted.BoolMustQuery{
									{
										Query: &inverted.Query{
											Leaf: &inverted.TermQuery{
												Field: "test",
												Term:  "a",
											},
										},
									},
								},
							},
							Aggs: []inverted.Agg{
								inverted.TermsAgg{Field: "c"},
							},
							Aggregations: map[string]*inverted.Aggregation{
								"some-name": {
									Filter: &inverted.Query{
										BoolMust: []*inverted.BoolMustQuery{
											{
												Query: &inverted.Query{
													Leaf: &inverted.TermQuery{
														Field: "test",
														Term:  "b",
													},
												},
											},
										},
									},
									Aggs: []inverted.Agg{
										inverted.TermsAgg{Field: "d"},
									},
									Aggregations: map[string]*inverted.Aggregation{
										"another-name": {
											Filter: &inverted.Query{
												BoolMust: []*inverted.BoolMustQuery{
													{
														Query: &inverted.Query{
															Leaf: &inverted.TermQuery{
																Field: "test2",
																Term:  "c",
															},
														},
													},
												},
											},
											Aggs: []inverted.Agg{
												inverted.TermsAgg{Field: "test"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			nil,
		},
	} {
		have, err := ParseReq([]byte(tcase.body))
		assert.Equal(t, tcase.want, have, tcase.name)
		if tcase.err != nil {
			assert.Equal(t, tcase.err.Error(), err.Error(), tcase.name)
		} else {
			assert.Nil(t, err)
		}
	}
}
