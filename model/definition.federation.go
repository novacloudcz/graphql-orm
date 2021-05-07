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
		Type: nonNull(namedType("_Service")),
	}
}

func createFederationEntityUnion(m *Model) *ast.UnionDefinition {
	types := []*ast.Named{}

	for _, o := range m.Objects() {
		if o.IsFederatedType() {
			t := namedType(o.Name())
			types = append(types, t.(*ast.Named))
		}
	}
	for _, e := range m.ObjectExtensions() {
		if e.IsFederatedType() {
			t := namedType(e.Object.Name())
			types = append(types, t.(*ast.Named))
		}
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

func getObjectDefinitionFromFederationExtension(def *ast.ObjectDefinition) *ast.ObjectDefinition {
	federationDirectives := []string{"requires", "provides", "key", "extends", "external"}
	for _, dir := range federationDirectives {
		def.Directives = filterDirective(def.Directives, dir)
	}
	for _, field := range def.Fields {
		for _, dir := range federationDirectives {
			field.Directives = filterDirective(field.Directives, dir)
		}
	}
	return def
}
