package templates

var Lambda = `package main

import (
	"github.com/akrylysov/algnhsa"
	"{{.Config.Package}}/gen"
)

func main() {
	handler := gen.GetHTTPHandler()
	algnhsa.ListenAndServe(handler, nil)
}
`
