package main

import (
	"bytes"
	"fmt"
	"github.com/richardjennings/invertedindex/server"
	"io/ioutil"
	"net/http/httptest"
)

func main() {
	s := server.NewServer()
	fmt.Println("create an index for emails")
	fmt.Println(
		request(
		"PUT",
		"/emails",
		`{
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
			}`,
			s,
		),
	)
	fmt.Println("add some emails ...")
	for i, v := range [][4]string{
		{"a@b.com", "b@a.com", "testing 123", "full-text search capabilities"},
		{"b@a.com", "a@a.com", "testing 123", "yeah i know right!"},
		{"b@a.com", "a@b.com", "what are you working on?", "anything new"},
		{"a@b.com", "b@a.com", "what are you working on?", "a full-text search inverted!"},
	} {
		request(
			"PUT",
			fmt.Sprintf("/emails/%d", i+1),
			fmt.Sprintf(`{"from":"%s","to":"%s","subject":"%s", "body": "%s"}`,v[0],v[1],v[2],v[3]),
			s,
		)
	}

	fmt.Println("get email index stats:")
	fmt.Println(request("GET", "/emails", "", s))
	//time.Sleep(time.Duration(2*time.Second))

	// search for emails from b@a.com with "maybe think" in the body
	fmt.Println(
		request(
			"GET",
			"/emails/_search",
			`
			{
				"query":{
					"bool": {
						"must": {
							"match_phrase":{
								"body": "a full-text search inverted"
							},
							"term":{
								"from": "a@b.com"
							}
						}
					}
				}
			}`,
			s,
		),
	)


}

func request(method string, url string, body string, s server.Server) string {
	req := httptest.NewRequest(method, url, bytes.NewBufferString(body))
	resp := httptest.NewRecorder()
	s.ServeHTTPMock(resp, req)
	b, _ := ioutil.ReadAll(resp.Body)
	return string(b[:])
}