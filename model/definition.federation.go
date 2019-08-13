package model

import (
	"github.com/graphql-go/graphql/language/kinds"

	"github.com/graphql-go/graphql/language/ast"
)

func createFederationServiceObject() *ast.ObjectDefinition {
	return &ast.ObjectDefinition{
		Kind: kinds.ObjectDefinition,
		Name: nameNode("_Service"),
		Fields: []*ast.FieldDefinition{
			&ast.FieldDefinition{
				Kind: kinds.FieldDefinition,
				Name: nameNode("sdl"),
				Type: namedType("String"),
			},
		},
	}
}

func createFederationServiceQueryField() *ast.FieldDefinition {
	return &ast.FieldDefinition{
		Kind: kinds.FieldDefinition,
		Name: nameNode("_service"),
		Type: namedType("_Service"),
	}
}

func createFederationEntityUnion(m *Model) *ast.UnionDefinition {
	types := []*ast.Named{}

	for _, o := range m.Objects() {
		t := namedType(o.Name())
		types = append(types, t.(*ast.Named))
	}

	return &ast.UnionDefinition{
		Kind:  kinds.UnionDefinition,
		Name:  nameNode("_Entity"),
		Types: types,
	}
}
func createFederationEntitiesQueryField() *ast.FieldDefinition {
	return &ast.FieldDefinition{
		Kind: kinds.FieldDefinition,
		Name: nameNode("_entities"),
		Type: nonNull(listType(namedType("_Entity"))),
		Arguments: []*ast.InputValueDefinition{
			&ast.InputValueDefinition{
				Kind: kinds.InputValueDefinition,
				Name: nameNode("representations"),
				Type: nonNull(listType(nonNull(namedType("_Any")))),
			},
		},
	}
}
