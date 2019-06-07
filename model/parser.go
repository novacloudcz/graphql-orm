package model

import (
	"github.com/graphql-go/graphql/language/parser"
)

// Parse
func Parse(m string) (Model, error) {
	var model Model
	astDoc, err := parser.Parse(parser.ParseParams{
		Source: m,
		Options: parser.ParseOptions{
			NoLocation: true,
		},
	})
	if err != nil {
		return model, err
	}

	model = Model{astDoc}
	return model, nil
}
