package server

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/pat"
	"github.com/richardjennings/invertedindex/inverted"
	"github.com/richardjennings/invertedindex/query"
	"net/http"
	"net/http/httptest"
)

type Server struct {
	engine  inverted.Engine
	httpApi *httpApi
}

type httpApi struct {
	srv    *http.Server
	engine *inverted.Engine
}

func NewServer() Server {
	s := Server{engine: inverted.New()}
	s.httpApi = s.newHttpApi()
	return s
}

func (s *Server) Serve() error {
	return s.httpApi.srv.ListenAndServe()
}

func (s *Server) ServeHTTPMock(resp *httptest.ResponseRecorder, req *http.Request) {
	s.httpApi.srv.Handler.ServeHTTP(resp, req)
}

func (s *Server) Close() error {
	return s.httpApi.srv.Close()
}

func (s *Server) newHttpApi() *httpApi {
	api := &httpApi{srv: &http.Server{Addr: "127.0.0.1:8080"}, engine: &s.engine}
	api.srv.Handler = api.createMux()
	return api
}

func (a *httpApi) createMux() *http.ServeMux {
	mux := http.NewServeMux()
	router := pat.New()

	// search api
	router.Get("/{name}/_search", a.search)
	router.Post("/{name}/_search", a.search)

	// index api

	// index document with id and single field (plain text body)
	router.Put("/{name}/{uri}/{field}", a.indexPut)

	// index document with id (uri)
	router.Put("/{name}/{uri}", a.indexPutBody)

	//create index
	router.Put("/{name}", a.indexCreate)

	//delete index
	router.Delete("/{name}", a.indexDelete)

	// get index info
	router.Get("/{name}", a.index)

	// list indexes
	router.Get("/", a.indexes)

	mux.Handle("/", router)

	return mux
}

func (a *httpApi) handleError(err error, w http.ResponseWriter) {
	switch err.Error() {
	case "index not found":
		w.WriteHeader(404)
	default:
		w.WriteHeader(500)
	}
	/*
		_, err2 := w.Write([]byte(err.Error()))
		if err2 != nil {
			panic(err2)
		}
	*/
}

func (a *httpApi) jsonResponse(o interface{}, w http.ResponseWriter) {
	w.Header()["Content-Type"] = []string{"application/json"}
	j, err := json.Marshal(o)
	if err != nil {
		// how to test ?
		a.handleError(err, w)
	}
	_, err = w.Write(j)
	if err != nil {
		// what are the error conditions to handle here ?
		return
	}
}

//
func (a *httpApi) search(w http.ResponseWriter, r *http.Request) {
	if qs := r.URL.Query().Get("q"); qs != "" {
		a.searchQueryString(w, r)
		return
	}
	a.searchBody(w, r)
}

// Perform a search using the Request Body Query DSL
func (a *httpApi) searchBody(w http.ResponseWriter, r *http.Request) {
	indexName := r.URL.Query().Get(":name")
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	defer r.Body.Close()
	if err != nil {
		a.handleError(err, w)
		return
	}
	q, err := query.ParseReq(buf.Bytes())
	if err != nil {
		a.handleError(err, w)
		return
	}
	res, err := a.engine.Search(indexName, q)
	if err != nil {
		a.handleError(err, w)
		return
	}
	a.jsonResponse(res, w)
}

// Perform a search using the Query String DSL
func (a *httpApi) searchQueryString(w http.ResponseWriter, r *http.Request) {
	indexName := r.URL.Query().Get(":name")
	c := r.URL.Query()
	q, err := query.ParseReqFromQuery(c)
	if err != nil {
		a.handleError(err, w)
		return
	}
	res, err := a.engine.Search(indexName, q)
	if err != nil {
		a.handleError(err, w)
		return
	}
	a.jsonResponse(res, w)
}

// list all indexes
func (a *httpApi) indexes(w http.ResponseWriter, r *http.Request) {
	a.jsonResponse(a.engine.IndexList(), w)
}

// get index stats
func (a *httpApi) index(w http.ResponseWriter, r *http.Request) {
	indexName := r.URL.Query().Get(":name")
	stats, err := a.engine.IndexStats(indexName)
	if err != nil {
		a.handleError(err, w)
		return
	}
	a.jsonResponse(stats, w)
}

// create an index
func (a *httpApi) indexCreate(w http.ResponseWriter, r *http.Request) {
	indexName := r.URL.Query().Get(":name")
	// allow configuration using post body
	cfg := make(map[string]map[string]map[string]string)
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		a.handleError(err, w)
		return
	}
	if buf.Len() > 0 {
		err = json.Unmarshal(buf.Bytes(), &cfg)
		if err != nil {
			a.handleError(err, w)
			return
		}
	}

	// if there is config
	conf, ok := cfg["mapping"]
	if ok {
		// create the index with the config
		_, err := a.engine.NewIndex(indexName, conf)
		if err != nil {
			a.handleError(err, w)
			return
		}
		stats, err := a.engine.IndexStats(indexName)
		if err != nil {
			a.handleError(err, w)
			return
		}
		a.jsonResponse(stats, w)
		return
	}

	// otherwise create index with no config
	_, err = a.engine.NewIndex(indexName, nil)
	if err != nil {
		a.handleError(err, w)
		return
	}
	stats, err := a.engine.IndexStats(indexName)
	if err != nil {
		a.handleError(err, w)
		return
	}
	a.jsonResponse(stats, w)
}

// delete an index
func (a *httpApi) indexDelete(w http.ResponseWriter, r *http.Request) {
	indexName := r.URL.Query().Get(":name")
	err := a.engine.DeleteIndex(indexName)
	if err != nil {
		a.handleError(err, w)
	}
	a.jsonResponse(true, w)
}

func (a *httpApi) indexPutBody(w http.ResponseWriter, r *http.Request) {
	indexName := r.URL.Query().Get(":name")
	uri := r.URL.Query().Get(":uri")
	// read body into map[string]string .. later need to work out how to
	// work with other types, ints, bools, arrays, ....
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		a.handleError(err, w)
		return
	}
	var content map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &content)
	if err != nil {
		a.handleError(err, w)
		return
	}
	err = a.engine.Index(indexName, uri, content)
	if err != nil {
		a.handleError(err, w)
		return
	}
	a.jsonResponse(true, w)
}

func (a *httpApi) indexPut(w http.ResponseWriter, r *http.Request) {
	indexName := r.URL.Query().Get(":name")
	uri := r.URL.Query().Get(":uri")
	field := r.URL.Query().Get(":field")

	content := map[string]interface{}{field: r.Body}

	// router ensures field not empty
	err := a.engine.Index(indexName, uri, content)
	if err != nil {
		a.handleError(err, w)
		return
	}
	a.jsonResponse(true, w)
}
