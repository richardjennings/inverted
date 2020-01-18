package inverted

import "github.com/richardjennings/invertedindex/index"

type CompoundQuery interface {
	Run(i *index.Index) (index.Result, error)
}

type Query struct {
	Leaf        LeafQuery
	BoolMust    []*BoolMustQuery
	BoolShould  []*BoolShouldQuery
	BoolMustNot []*BoolMustNotQuery
	BoolFilter  []*BoolFilterQuery
}

// bool compound
type (
	BoolMustQuery    struct{ *Query }
	BoolShouldQuery  struct{ *Query }
	BoolMustNotQuery struct{ *Query }
	BoolFilterQuery  struct{ *Query }
)

/**
"aggregations" : {
    "<aggregation_name>" : {
        "<aggregation_type>" : {
            <aggregation_body>
        },
        ["aggregations" : { [<sub_aggregation>]* } ]
    }
    [,"<aggregation_name_2>" : { ... } ]*
}
*/
// { "agg": {
type Aggregation struct {
	// child buckets
	Aggregations map[string]*Aggregation
	// aggregations at this level
	Aggs   []Agg
	Filter *Query
}

type Agg interface {
}

type TermsAgg struct {
	Field string
}

func (q Query) Run(i *index.Index) (index.Result, error) {

	if q.Leaf != nil {
		return q.Leaf.Query(i)
	}

	result := QueryResult{}

	// Use this operator for clauses that must appear in the matching documents.
	if q.BoolMust != nil {
		for j, query := range q.BoolMust {
			// AND multiple

			r, err := query.Run(i)
			if err != nil {
				return nil, err
			}
			var t struct{}

			// again need r.Docs() to return a map
			res := make(map[int]struct{})
			for _, d := range r.Docs() {
				res[d] = t
			}

			if j == 0 {
				for d := range res {
					result[d] = t
				}
			} else {
				for d := range result {
					if _, ok := res[d]; !ok {
						delete(result, d)
					}
				}
			}
			//result = r.Docs()

		}

		//if result == nil {
		//	result = make(Result)
		//} else {
		// filter ?
		//}
	}

	// Use this operator for clauses that should appear in the matching documents.
	// For a BooleanQuery with no MUST clauses one or more SHOULD clauses must match
	// a document for the BooleanQuery to match.
	if q.BoolShould != nil {
		for _, query := range q.BoolShould {

			_, err := query.Run(i)
			if err != nil {
				return nil, err
			}
			// @todo results in r, increase score
		}
	}

	// Use this operator for clauses that must not appear in the matching documents.
	// Note that it is not possible to search for queries that only consist of a MUST_NOT clause.
	// These clauses do not contribute to the score of documents.
	if q.BoolMustNot != nil {
		for _, query := range q.BoolMustNot {

			r, err := query.Run(i)
			if err != nil {
				return nil, err
			}
			for _, d := range r.Docs() {
				delete(result, d)
			}
		}
	}

	// Like MUST except that these clauses do not participate in scoring.
	if q.BoolFilter != nil {
		for _, query := range q.BoolFilter {

			r, err := query.Run(i)
			if err != nil {
				return nil, err
			}
			if result == nil {
				var t struct{}
				for _, d := range r.Docs() {
					result[d] = t
				}
			} else {
				// remove any results that are not in filter
				// @todo maybe use map[int]struct{} instead of []int for docs?
				docs := make(map[int]struct{})
				var t struct{}
				for _, d := range r.Docs() {
					docs[d] = t
				}

				//remove all results not in docs
				for d := range result {
					if _, ok := docs[d]; !ok {
						delete(result, d)
					}
				}
			}
		}
	}
	return result, nil
}
