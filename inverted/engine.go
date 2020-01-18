package inverted

import (
	"errors"
	"github.com/richardjennings/invertedindex/index"
)

var (
	IndexNotFound      = "index not found"
	IndexAlreadyExists = "index already exists"
)

type Engine struct {
	Indexes map[string]*index.Index
}

type SearchRequest struct {
	Query *Query
	Agg   *Aggregation
}

func New() Engine {
	e := Engine{}
	e.Indexes = make(map[string]*index.Index)
	return e
}

func (e *Engine) GetIndex(indexName string) (*index.Index, error) {
	cidx, ok := e.Indexes[indexName]
	if !ok {
		return nil, errors.New(IndexNotFound)
	}

	return cidx, nil
}

func (e *Engine) IndexList() []string {
	indexes := e.Indexes
	list := make([]string, len(indexes))
	i := 0
	for name := range indexes {
		list[i] = name
		i++
	}
	return list
}

func (e *Engine) IndexStats(indexName string) (*index.Stats, error) {
	idx, err := e.GetIndex(indexName)
	if err != nil {
		return nil, err
	}
	stats := idx.Stats()
	return stats, nil
}

func (e *Engine) NewIndex(indexName string, cf map[string]map[string]string) (*index.Index, error) {
	_, exists := e.Indexes[indexName]
	if exists {
		return nil, errors.New(IndexAlreadyExists)
	}
	cidx, err := index.NewIndex(cf)
	if err != nil {
		return nil, err
	}
	e.Indexes[indexName] = cidx
	//stats, err := e.IndexStats(indexName)
	return cidx, err
}

func (e *Engine) DeleteIndex(indexName string) error {
	_, err := e.GetIndex(indexName)
	if err != nil {
		return err
	}
	delete(e.Indexes, indexName)
	return nil
}

func (e *Engine) Index(indexName string, uri string, content map[string]interface{}) error {
	idx, err := e.GetIndex(indexName)
	if err != nil {
		return err
	}
	return idx.Index(uri, content)
}

func (e *Engine) Search(idxName string, req *SearchRequest) (map[string][]int, error) {
	cidx, err := e.GetIndex(idxName)
	if err != nil {
		return nil, err
	}
	result := make(map[string][]int)
	if req.Query != nil {
		res, err := req.Query.Run(cidx)
		if err != nil {
			return nil, err
		}
		result["hits"] = res.Docs()
	}
	return result, nil
}
