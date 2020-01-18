package main

import (
	"fmt"
	"github.com/richardjennings/invertedindex/index"
	"github.com/richardjennings/invertedindex/inverted"
	"log"
)

func main() {
	var idx *index.Index
	var err error

	// create new inverted
	e := inverted.New()

	// index definition
	idxDef := index.Schema{
		"firstname": {
			"type": index.Keyword,
		},
		"lastname": {
			"type": index.Keyword,
		},
		"technology": {
			"type": index.Keyword,
		},
	}

	// create a new composite index called programmers with name and technology keyword type fields
	if idx, err = e.NewIndex("programmers", idxDef);
		err != nil {
		log.Fatal(err)
	}

	// add some data to the index
	for _, v := range [][4]string{
		{"1", "Dennis", "Ritchie", "C"},
		{"2", "Linus", "Torvalds", "Linux"},
		{"3", "Ken", "Thompson", "Unix"},
		{"4", "Donald", "Knuth", "Algorithm Analysis"},
		{"5", "Bjarne", "Stroustrup", "c++"},
	} {
		row := map[string]interface{}{"firstname": v[1], "lastname": v[2], "technology": v[3]}
		if err = idx.Index(v[0], row); err != nil {
			log.Fatal(err)
		}
	}

	r, err := e.Search("programmers", &inverted.SearchRequest{Query: &inverted.Query{Leaf: &inverted.TermQuery{Term: "Ken", Field: "firstname"}}},)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range r["hits"] {
		uri, _ := idx.Doc(d)
		fmt.Println("matched: ", uri)
	}


}
