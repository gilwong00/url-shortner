package main

import (
	"github.com/gilwong00/url-shortner/pkg/server"
)

func main() {
	s := server.NewServer()
	s.StartServer()
}
