package query

import (
	"errors"
	"github.com/richardjennings/invertedindex/inverted"
)

func parseAgg(i interface{}) ([]inverted.Agg, error) {
	var aggs []inverted.Agg
	a, err := mapStrI(i)
	if err != nil {
		return nil, err
	}

	// terms aggregation config
	if t, ok := a["terms"]; ok {
		a, err = mapStrI(t)
		if err != nil {
			return nil, err
		}
		if f, ok := a["field"]; ok {
			switch f.(type) {
			case string:
				aggs = append(aggs, inverted.TermsAgg{Field: f.(string)})
			default:
				return nil, errors.New("invalid type for field")
			}
		} else {
			return nil, errors.New("invalid terms agg")
		}
	}

	return aggs, nil
}

// JSON aggregation parsing
func parseAggregation(i interface{}) (*inverted.Aggregation, error) {
	aggregation := &inverted.Aggregation{}
	aggregation.Aggregations = make(map[string]*inverted.Aggregation)
	a, err := mapStrI(i)
	if err != nil {
		return nil, err
	}

	for aggName, v := range a {
		cAgg := &inverted.Aggregation{}
		//cAgg.aggregations = make(map[string]*Aggregation)
		aggregation.Aggregations[aggName] = cAgg
		j, err := mapStrI(v)
		if err != nil {
			return nil, err
		}
		if len(j) == 0 {
			return nil, errors.New("no Field/value specified")
		}

		// for filter
		if _, ok := j["filter"]; ok {
			cAgg.Filter, err = parseQuery(j["filter"])
			if err != nil {
				return nil, err
			}
		}

		hasAgg := false
		// for agg config not wrapped in "aggs"
		cAgg.Aggs, err = parseAgg(j)
		if err != nil {
			return nil, err
		}
		if len(cAgg.Aggs) > 0 {
			hasAgg = true
		}

		if _, ok := j["aggs"]; ok {

			// as agg
			if !hasAgg {
				cAgg.Aggs, err = parseAgg(j["aggs"])
				if err != nil {
					return nil, err
				}
				if len(cAgg.Aggs) > 0 {
					// is config not recursive child
					continue
				}
			}
			//

			// as child aggregation
			pl, err := parseAggregation(j["aggs"])
			if err != nil {
				return nil, err
			}
			if len(pl.Aggregations) > 0 {
				cAgg.Aggregations = pl.Aggregations
			}
		}
		if err != nil {
			return nil, err
		}

	}
	return aggregation, nil
}
