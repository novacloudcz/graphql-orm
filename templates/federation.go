package templates

var Federation = `package gen

import (
	"encoding/json"
	"io"

	"github.com/99designs/gqlgen/graphql"
)

func Marshal_Any(v interface{}) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		err := json.NewEncoder(w).Encode(v)
		if err != nil {
			panic(err)
		}
	})
}

func Unmarshal_Any(v interface{}) (interface{}, error) {
	return v, nil
}

`
