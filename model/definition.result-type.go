package model

import (
	"github.com/graphql-go/graphql/language/kinds"

	"github.com/graphql-go/graphql/language/ast"
)

func objectResultTypeDefinition(o *Object) *ast.ObjectDefinition {
	return &ast.ObjectDefinition{
		Kind: kinds.ObjectDefinition,
		Name: nameNode(o.Name() + "ResultType"),
		Fields: []*ast.FieldDefinition{
			{
				Kind: kinds.FieldDefinition,
				Name: nameNode("items"),
				Type: nonNull(listType(nonNull(namedType(o.Name())))),
			},
			{
				Kind: kinds.FieldDefinition,
				Name: nameNode("aggregations"),
				Type: nonNull(namedType(o.Name() + "Aggregations")),
			},
			{
				Kind: kinds.FieldDefinition,
				Name: nameNode("count"),
				Type: nonNull(namedType("Int")),
			},
		},
	}
}
