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

	fmt.Println("create index test")
	fmt.Println(request("PUT", "/test",`{"mapping":{"content":{"type": "text"}}}`, s))

	fmt.Println("add some data to the index")
	fmt.Println(request("PUT", "/test/1/content",`a b c`, s))

	fmt.Println("list indexes")
	fmt.Println(request("GET", "/", "", s))

	fmt.Println("match query on test index")
	fmt.Println(request("GET", "/test/_search", `{"query":{"match":{"content": "a"}}}`, s))
}

func request(method string, url string, body string, s server.Server) string {
	req := httptest.NewRequest(method, url, bytes.NewBufferString(body))
	resp := httptest.NewRecorder()
	s.ServeHTTPMock(resp, req)
	b, _ := ioutil.ReadAll(resp.Body)
	return string(b[:])
}