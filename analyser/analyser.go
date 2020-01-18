package analyser

import (
	"bytes"
	"errors"
	"io"
)

type Analyser interface {
	Analyse(content interface{}) ([]string, error)
}

type FullTextAnalyser struct {
	Tokenizer Tokenizer
}

func (a *FullTextAnalyser) Analyse(content interface{}) ([]string, error) {
	switch content.(type) {
	case string:
		return a.Tokenizer.Tokenize(content.(string)), nil
	case io.ReadCloser:
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(content.(io.ReadCloser))
		if err != nil {
			return nil, err
		}
		return a.Tokenizer.Tokenize(buf.String()), nil
	default:
		return nil, errors.New("string or io.ReadCloser type required")
	}
}

type KeywordAnalyser struct{}

func (a *KeywordAnalyser) Analyse(content interface{}) ([]string, error) {
	switch content.(type) {
	case string:
		return []string{content.(string)}, nil
	case []string:
		return content.([]string), nil
	default:
		return nil, errors.New("expecting string or []string")
	}
}
