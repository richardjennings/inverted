package test

import (
	"errors"
	"io"
	"io/ioutil"
)

// reader that returns an error
type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func NewErrReadCloser() io.ReadCloser {
	var r errReader
	return ioutil.NopCloser(r)
}
