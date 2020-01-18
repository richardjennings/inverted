package test

import (
	"bytes"
	"github.com/richardjennings/invertedindex/server"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func Test_Acceptance(t *testing.T) {
	api := server.NewServer()
	defer func(t *testing.T, server server.Server) {
		err := server.Close()
		assert.Nil(t, err)
	}(t, api)
	go func() {
		_ = api.Serve()
	}()
	var netClient = &http.Client{
		Timeout: time.Second * 1,
	}

	for _, tcase := range []struct {
		name       string
		method     string
		url        string
		body       io.Reader
		want       string
		wantStatus int
	}{
		// Empty list of Indexes
		{"/", http.MethodGet, "http://127.0.0.1:8080/", nil, "[]", 200},

		// Create a Test Index
		{
			"create testindex",
			http.MethodPut,
			"http://127.0.0.1:8080/testindex",
			bytes.NewBufferString(`
				{
					"mapping":{
						"name":{
							"type":"keyword"
						}, 
						"brand":{
							"type":"keyword"
						},
						"category": {
							"type":"keyword"		
						},
						"title":{
							"type":"text"
						}, 
						"description":{
							"type":"text"
						}
					}
				}`),
			`{"DocumentCount":0,"Fields":{"brand":{"TermCount":0},"category":{"TermCount":0},"description":{"TermCount":0},"name":{"TermCount":0},"title":{"TermCount":0}}}`,
			200,
		},

		// Index a document
		{
			"create new document with uri 1 in testindex",
			http.MethodPut,
			"http://127.0.0.1:8080/testindex/1",
			bytes.NewBufferString(`
				{"brand":"apple", "category": "wearable", "title": "apple watch 4", "description": "smart watch with heart rate monitor"}
			`),
			`true`,
			200,
		},
		// Index another document
		{
			"create new document with uri 2 in testindex",
			http.MethodPut,
			"http://127.0.0.1:8080/testindex/2",
			bytes.NewBufferString(`
				{"brand":"apple", "category": "tablet", "title": "ipad pro", "description": "touch screen tablet"}
			`),
			`true`,
			200,
		},
		// match branch = apple
		{
			"match brand = apple testindex",
			http.MethodPost,
			"http://127.0.0.1:8080/testindex/_search",
			bytes.NewBufferString(`
				{"query":{"term":{"brand":"apple"}}}
			`),
			`{"hits":[0,1]}`,
			200,
		},
		// must match branch = apple and title match ipad
		{
			"match brand = apple testindex",
			http.MethodPost,
			"http://127.0.0.1:8080/testindex/_search",
			bytes.NewBufferString(`
				{
					"query": {
						"bool": {
							"must": [
								{"term":{"brand":"apple"}}, 
								{"match":{"title":"ipad"}}
							]
						}
					}
				}
			`),
			`{"hits":[1]}`,
			200,
		},
	} {
		req, err := http.NewRequest(
			tcase.method,
			tcase.url,
			tcase.body, //
		)
		assert.Nil(t, err)
		resp, err := netClient.Do(req)
		assert.Nil(t, err)
		b, err := ioutil.ReadAll(resp.Body)
		assert.Nil(t, err)
		assert.Equal(t, tcase.wantStatus, resp.StatusCode, tcase.name)
		assert.Equal(t, tcase.want, string(b[:]), tcase.name)
	}

}
