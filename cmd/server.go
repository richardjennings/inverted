package main

import (
	"github.com/richardjennings/invertedindex/server"
)

func main() {
	s := server.NewServer()
	err := s.Serve()
	if err != nil {
		panic(err)
	}
}
