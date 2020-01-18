package index

import (
	"errors"
	"fmt"
)

// The Document Index
type Index struct {
	DocumentIndex map[string]int
	Documents     []Document
	Idxs          map[string]Idx
}

type Stats struct {
	DocumentCount int
	Fields        map[string]IdxStats
}

type IdxStats struct {
	TermCount int
}

type Document struct {
	URI string
}

type Result interface {
	Docs() []int
}

type Schema map[string]map[string]string

// Field Index Interface
type Idx interface {
	Stats() IdxStats
	Index(docId int, content interface{}) error
}

// Query Interfaces
type Match interface {
	MatchQuery(query string) (TermFreqResult, error)
}
type Phrase interface {
	PhraseQuery(query string) (PostingResult, error)
}
type Term interface {
	TermQuery(query string) (KeywordResult, error)
}
type Terms interface {
	TermsQuery(query []string) (KeywordResult, error)
}

// Aggregation Interfaces
type TermsAggregation interface {
	TermsAgg() (KeywordResult, error)
}

func NewIndex(cf map[string]map[string]string) (*Index, error) {
	cidx := Index{}
	cidx.Idxs = make(map[string]Idx)
	cidx.DocumentIndex = make(map[string]int)

	if len(cf) > 0 {
		// create the mapping specified
		for field, v := range cf {
			typ, hasType := v["type"]
			if hasType {
				switch typ {
				case "text":
					// currently error cannot occur because field duplication case
					// is prevented by the map key in this function
					_, _ = cidx.newFieldIndex(field, NewTextIndex())
				case "keyword":
					_, _ = cidx.newFieldIndex(field, NewKeywordIndex())
				default:
					return nil, errors.New("unknown field type")
				}
			} else {
				return nil, errors.New("missing type")
			}
		}
	}
	return &cidx, nil
}

func (ci *Index) Stats() *Stats {
	stats := Stats{}
	stats.DocumentCount = len(ci.Documents)
	stats.Fields = make(map[string]IdxStats)
	for name, idx := range ci.Idxs {
		stats.Fields[name] = idx.Stats()
	}
	return &stats
}

func (ci *Index) newFieldIndex(field string, idx Idx) (Idx, error) {
	_, ok := ci.Idxs[field]
	if ok {
		return nil, errors.New("field index already exists")
	}
	ci.Idxs[field] = idx
	return ci.Idxs[field], nil
}

func (ci *Index) GetFieldIdx(field string) (Idx, error) {
	idx, ok := ci.Idxs[field]
	if ok {
		return idx, nil
	}
	return nil, errors.New("field not found")
}

func (ci *Index) Index(uri string, content map[string]interface{}) error {
	_, ok := ci.DocumentIndex[uri]
	if ok {
		return errors.New("document uri already exists")
	}

	// add document to documents list
	ci.Documents = append(ci.Documents, Document{URI: uri})
	docId := len(ci.Documents) - 1
	ci.DocumentIndex[uri] = docId

	for field, txt := range content {
		idx, ok := ci.Idxs[field]
		if !ok {
			return errors.New("field not found")
		}
		err := idx.Index(docId, txt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ci Index) Doc(id int) (string, error) {
	if id < len(ci.Documents) {
		return ci.Documents[id].URI, nil
	}
	return "", fmt.Errorf("document id %d not found", id)
}
