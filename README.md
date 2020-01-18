# Inverted Full Text Search Engine

## About
Inverted is a Full Text Search engine functioning both as a stand-alone HTTP based server and a library to be used within a Go application.

### Query DSL
Inverted provides a JSON based query syntax to provide sophisticated query capabilities. Query syntax targets compatibility 
with Elasticsearch to aid interoperability with existing Elasticsearch usage. 

### Indexes
Inverted supports composite indexes with field of differing types and configuration.

Example:
```
PUT /emails
{
  "mapping": {
    "from": {
      "type": "keyword"
    },
    "to": {
      "type": "keyword"
    },
    "subject": {
      "type": "text"
    },
    "body": {
      "type": "text"
    }
  }
}
```

### Text Fields
Text fields provide full-text search capabilities. Bodies of text are broken down into a sequence of individual tokens
by an Analyser. The choice of Analyser is configurable and defaults to splitting text on natural word boundaries removing
punctuation.

### Text Queries
Text fields support querying by Match, Multi Match or Match Phrase. Match queries count the number of times a term appears in each body
of text, returning results that are by default ordered by term frequency, Match Phrase queries look for the occurrence of
a sequence of tokens in bodies of text and returns results that are by default ordered by the closest match or the highest
number of exact matches. 

### Keyword Fields
Keyword fields represent exact values and are most useful for filtering and aggregations.

### Keyword Queries
Keyword fields support querying by Term or Terms. A Term query will only match the exact value searched for. A Terms query
matches exactly one or more of the supplied Term values.

### Queries
Queries can be constructed using logical containers.

Example:

```
GET /emails/_search
{
    "query":{
        "bool": {
            "must": {
                "match_phrase":{
                    "body": "a full-text search engine!"
                },
                "term":{
                    "from": "a@b.com"
                }
            }
        }
    }
}
```

### Package API
```go
    e := inverted.New()
    e.NewIndex(
        "my_index",
        index.Schema{
            "category": {"type": index.Keyword},
            "brand": {"type": index.Keyword},
            "name": {"type": index.Text},
            "description": { "type": index.Text},
        },
    )
    e.Index("my_index", "1", map[string]interface{}{"category":"laptops", "name": "latitude 7240", "description": "a laptop"})
    query := &inverted.SearchRequest{Query: &inverted.Query{Leaf: &inverted.TermQuery{Term: "laptops", Field: "category"}}}
    result, _ := e.Search("my_index", query)
    fmt.Println(result)
```






