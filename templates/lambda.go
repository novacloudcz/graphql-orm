package templates

var Lambda = `package main

import (
	"github.com/akrylysov/algnhsa"
	"gitlab.com/novacloud/reportingdokapsy/cms/orm/gen"
)

func main() {
	handler := gen.GetHTTPHandler()
	algnhsa.ListenAndServe(handler, nil)
}
`
