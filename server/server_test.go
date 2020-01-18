package server

import (
	"bytes"
	"github.com/richardjennings/invertedindex/test"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestServer_NewServer_HttpApi(t *testing.T) {
	s := NewServer()
	handler := s.httpApi.srv.Handler

	tcases := []struct {
		name       string
		method     string
		uri        string
		body       io.Reader
		wantStatus int
		wantBody   string
	}{
		{
			"list indexes",
			"GET",
			"/",
			nil,
			200,
			`[]`,
		},
		{
			"create index",
			"PUT",
			"/test",
			nil,
			200,
			`{"DocumentCount":0,"Fields":{}}`,
		},
		{
			"create duplicate index with config",
			"PUT", "/test",
			bytes.NewBufferString(`{"mapping":{"content":{"type": "text"}}}`),
			500,
			``,
		},
		{
			"create index with config",
			"PUT",
			"/testcfg",
			bytes.NewBufferString(`{"mapping":{"content":{"type": "text"}}}`),
			200,
			`{"DocumentCount":0,"Fields":{"content":{"TermCount":0}}}`,
		},
		{
			"index text",
			"PUT",
			"/testcfg/docid/content",
			bytes.NewBufferString(`a b c`),
			200,
			`true`,
		},
		{
			"single term query",
			"GET",
			"/testcfg/_search?q=content:a",
			nil,
			200,
			`{"hits":[0]}`,
		},
		{
			"single term query index not exists",
			"GET", "/testsasdasdg/_search?q=content:a",
			nil,
			404,
			``,
		},
		{
			"phrase query",
			"GET",
			"/testcfg/_search?q=content:%22a%20b%22",
			nil,
			200,
			`{"hits":[0]}`,
		},
		{
			"invalid query",
			"GET",
			"/testcfg/_search?q=content:%22a%20b",
			nil,
			500,
			``,
		},
		{
			"get index that does not exist",
			"GET",
			"/notexist",
			nil,
			404,
			"",
		},
		{
			"delete index does not exist",
			"DELETE",
			"/notexist",
			nil,
			404,
			`true`,
		},
		{
			"create duplicate index",
			"PUT",
			"/test",
			nil,
			500,
			``,
		},
		{
			"delete index",
			"DELETE",
			"/test",
			nil,
			200,
			`true`,
		},
		{
			"query index does not exist",
			"GET",
			"/tsdsf/content/a%20b",
			nil,
			404,
			``,
		},
		{
			"index text non existent index",
			"PUT", "/test2/docid/content",
			bytes.NewBufferString(``),
			404,
			``,
		},
		{
			"create index multiple fields",
			"PUT",
			"/mf",
			bytes.NewBufferString(`{"mapping":{"a":{"type": "text"}, "b":{"type": "text"}}}`),
			200,
			`{"DocumentCount":0,"Fields":{"a":{"TermCount":0},"b":{"TermCount":0}}}`,
		},
		{
			"create index invalid json",
			"PUT",
			"/mf",
			bytes.NewBufferString(`{"mapping":{"a":{"type": "text",`),
			500,
			``,
		},
		{
			"index content from body",
			"PUT",
			"/mf/1",
			bytes.NewBufferString(`{"a":"hello","b":"world"}`),
			200,
			"true",
		},
		{
			"index content from body invalid json",
			"PUT",
			"/mf/2",
			bytes.NewBufferString(`{"a":"hello:}`),
			500,
			"",
		},
		{"index content from body invalid index",
			"PUT",
			"/mfsssss/2",
			bytes.NewBufferString(`{"a":"hello"}`),
			404,
			"",
		},
		{
			"index stats",
			"GET",
			"/mf",
			nil,
			200,
			`{"DocumentCount":1,"Fields":{"a":{"TermCount":1},"b":{"TermCount":1}}}`,
		},
		{
			"index create body read error",
			"PUT",
			"/invalidbody",
			test.NewErrReadCloser(),
			500,
			``,
		},
		{
			"index content from error read closer",
			"PUT",
			"/mf/error",
			test.NewErrReadCloser(),
			500,
			"",
		},
		{
			"match query post body",
			"GET",
			"/mf/_search",
			bytes.NewBufferString(`{"query":{"match":{"a": "hello"}}}`),
			200,
			`{"hits":[0]}`,
		},
		{
			"match query post body nested",
			"GET",
			"/mf/_search",
			bytes.NewBufferString(`{"query":{"match":{"a": {"query": "hello"}}}}`),
			200,
			`{"hits":[0]}`,
		},
		{
			"query post body invalid json",
			"GET",
			"/mf/_search",
			bytes.NewBufferString(`{"query":{"mat}`),
			500,
			``,
		},
	}

	for _, tcase := range tcases {
		req := httptest.NewRequest(tcase.method, tcase.uri, tcase.body)
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		have := resp.Result()
		b, _ := ioutil.ReadAll(resp.Body)
		body := string(b[:])
		assert.Equal(t, tcase.wantStatus, have.StatusCode, tcase.name)
		assert.Equal(t, tcase.wantBody, body, tcase.name)
	}
}
