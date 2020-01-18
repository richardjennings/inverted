package query

import (
	"encoding/json"
	"errors"
	"github.com/richardjennings/invertedindex/inverted"
)

// JSON request parsing
func ParseReq(q []byte) (*inverted.SearchRequest, error) {

	req := &inverted.SearchRequest{}
	var i interface{}
	err := json.Unmarshal(q, &i)
	if err != nil {
		return nil, err
	}
	switch a := i.(type) {
	case map[string]interface{}:

		if q, ok := a["query"]; ok {
			req.Query, err = parseQuery(q)
			if err != nil {
				return nil, err
			}
		}

		if agg, ok := a["aggs"]; ok {
			req.Agg, err = parseAggregation(agg)
			if err != nil {
				return nil, err
			}
		}
	default:
		return nil, errors.New("unknown type")
	}
	return req, nil
}

// expect a map[string]interface{}
func mapStrI(a interface{}) (map[string]interface{}, error) {
	switch a.(type) {
	case map[string]interface{}:
		return a.(map[string]interface{}), nil
	default:
		return nil, errors.New("expected map")
	}
}

// JSON query parsing
func parseQuery(i interface{}) (*inverted.Query, error) {
	a, err := mapStrI(i)
	if err != nil {
		return nil, err
	}
	for k, v := range a {
		j, err := mapStrI(v)
		if err != nil {
			return nil, err
		}
		if len(j) == 0 {
			return nil, errors.New("no Field/value specified")
		}
		switch k {
		case "bool":
			q, err := parseBoolQuery(j)
			if err != nil {
				return nil, err
			}
			return q, nil
		case "term":
			q, err := parseTermQuery(j)
			if err != nil {
				return nil, err
			}
			return &inverted.Query{Leaf: q}, nil
		case "terms":
			q, err := parseTermsQuery(j)
			if err != nil {
				return nil, err
			}
			return &inverted.Query{Leaf: q}, nil
		case "match":
			q, err := parseMatchQuery(j)
			if err != nil {
				return nil, err
			}
			return &inverted.Query{Leaf: q}, nil
		case "match_phrase":
			q, err := parseMatchPhraseQuery(j)
			if err != nil {
				return nil, err
			}
			return &inverted.Query{Leaf: q}, nil
		case "multi_match":
			q, err := parseMultiMatchQuery(j)
			if err != nil {
				return nil, err
			}
			return &inverted.Query{Leaf: q}, nil
		default:
			return nil, errors.New("unknown key")
		}
	}
	return nil, errors.New("empty")
}

// JSON match_phrase query parsing
func parseMatchPhraseQuery(a map[string]interface{}) (*inverted.MatchPhraseQuery, error) {
	q := inverted.MatchPhraseQuery{}
	for k, v := range a {
		q.Field = k
		switch j := v.(type) {
		case map[string]interface{}:
			query, ok := j["query"]
			if !ok {
				return nil, errors.New("missing query")
			}
			switch query.(type) {
			case string:
				q.Term = query.(string)
			default:
				return nil, errors.New("expected string")
			}
		case string:
			q.Term = j
		default:
			return nil, errors.New("invalid")
		}
	}
	return &q, nil
}

// JSON multi_match query parsing
func parseMultiMatchQuery(a map[string]interface{}) (*inverted.MultiMatchQuery, error) {
	q := inverted.MultiMatchQuery{}
	fields, ok := a["fields"]
	if !ok {
		return nil, errors.New("missing field declaration")
	}
	switch f := fields.(type) {
	case []interface{}:
		for _, j := range f {
			switch j.(type) {
			case string:
				q.Fields = append(q.Fields, j.(string))
			default:
				return nil, errors.New("expecting string")
			}
		}
	default:
		return nil, errors.New("fields must be an array")
	}
	query, ok := a["query"]
	if !ok {
		return nil, errors.New("missing query declaration")
	}
	switch query.(type) {
	case string:
		q.Term = query.(string)
	default:
		return nil, errors.New("query should be a string")
	}
	return &q, nil
}

// JSON match query parsing
func parseMatchQuery(a map[string]interface{}) (*inverted.MatchQuery, error) {
	q := inverted.MatchQuery{}
	for k, v := range a {
		q.Field = k
		switch j := v.(type) {
		case map[string]interface{}:
			query, ok := j["query"]
			if !ok {
				return nil, errors.New("missing query")
			}
			switch query.(type) {
			case string:
				q.Term = query.(string)
			default:
				return nil, errors.New("expected string")
			}
		case string:
			q.Term = j
		default:
			return nil, errors.New("invalid")
		}
	}
	return &q, nil
}

// JSON term query parsing
func parseTermQuery(a map[string]interface{}) (*inverted.TermQuery, error) {
	q := inverted.TermQuery{}
	for k, v := range a {
		q.Field = k
		switch v.(type) {
		case string:
			q.Term = v.(string)
		default:
			return nil, errors.New("expected string")
		}
	}
	return &q, nil
}

// JSON terms query parsing
func parseTermsQuery(a map[string]interface{}) (*inverted.TermsQuery, error) {
	q := inverted.TermsQuery{}
	for k, v := range a {
		switch terms := v.(type) {
		case []interface{}:
			for _, t := range terms {
				q.Terms = append(q.Terms, t.(string))
			}
		default:
			return nil, errors.New("invalid value for terms")
		}
		q.Field = k
	}
	return &q, nil
}

func parseBoolMust(v interface{}, q *inverted.Query) error {
	var err error
	b := &inverted.BoolMustQuery{}
	b.Query, err = parseQuery(v)
	if err != nil {
		return err
	}
	q.BoolMust = append(q.BoolMust, b)
	return nil
}

func parseBoolFilter(v interface{}, q *inverted.Query) error {
	var err error
	b := &inverted.BoolFilterQuery{}
	b.Query, err = parseQuery(v)
	if err != nil {
		return err
	}
	q.BoolFilter = append(q.BoolFilter, b)
	return nil
}

func parseBoolMustNot(v interface{}, q *inverted.Query) error {
	var err error
	b := &inverted.BoolMustNotQuery{}
	b.Query, err = parseQuery(v)
	if err != nil {
		return err
	}
	q.BoolMustNot = append(q.BoolMustNot, b)
	return nil
}

func parseBoolShould(v interface{}, q *inverted.Query) error {
	var err error
	b := &inverted.BoolShouldQuery{}
	b.Query, err = parseQuery(v)
	if err != nil {
		return err
	}
	q.BoolShould = append(q.BoolShould, b)
	return nil
}

// JSON bool query parsing
func parseBoolQuery(a map[string]interface{}) (*inverted.Query, error) {
	var err error
	var q inverted.Query
	for k, v := range a {
		switch k {
		case "must":
			switch v.(type) {
			// eg. "must":[{"....
			case []interface{}:
				v = v.([]interface{})
				for _, fq := range v.([]interface{}) {
					err = parseBoolMust(fq, &q)
					if err != nil {
						return nil, err
					}
				}
			// eg. "must":{"..
			case map[string]interface{}:
				err = parseBoolMust(v, &q)
				if err != nil {
					return nil, err
				}
			default:
				return nil, errors.New("unexpected type")
			}

		case "filter":
			switch v.(type) {
			// eg. "filter":[{"....
			case []interface{}:
				v = v.([]interface{})
				for _, fq := range v.([]interface{}) {
					err = parseBoolFilter(fq, &q)
					if err != nil {
						return nil, err
					}
				}
			// eg. "filter":{"..
			case map[string]interface{}:
				err = parseBoolFilter(v, &q)
				if err != nil {
					return nil, err
				}
			default:
				return nil, errors.New("unexpected type")
			}

		case "must_not":
			switch v.(type) {
			// eg. "must_not":[{"....
			case []interface{}:
				v = v.([]interface{})
				for _, fq := range v.([]interface{}) {
					err = parseBoolMustNot(fq, &q)
					if err != nil {
						return nil, err
					}
				}
			// eg. "must_not":{"..
			case map[string]interface{}:
				err = parseBoolMustNot(v, &q)
				if err != nil {
					return nil, err
				}
			default:
				return nil, errors.New("unexpected type")
			}

		case "should":
			switch v.(type) {
			// eg. "should":[{"....
			case []interface{}:
				v = v.([]interface{})
				for _, fq := range v.([]interface{}) {
					err = parseBoolShould(fq, &q)
					if err != nil {
						return nil, err
					}
				}
			// eg. "should":{"..
			case map[string]interface{}:
				err = parseBoolShould(v, &q)
				if err != nil {
					return nil, err
				}
			default:
				return nil, errors.New("unexpected type")
			}
		default:
			return nil, errors.New("unknown key")
		}
	}
	return &q, nil
}
