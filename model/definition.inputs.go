package model

import (
	"github.com/graphql-go/graphql/language/kinds"

	"github.com/graphql-go/graphql/language/ast"
)

func createObjectDefinition(obj Object) *ast.InputObjectDefinition {
	fields := []*ast.InputValueDefinition{}
	for _, f := range obj.Columns() {
		fields = append(fields, &ast.InputValueDefinition{
			Kind:        kinds.InputValueDefinition,
			Name:        f.Name,
			Description: f.Description,
			Type:        f.Type,
		})
	}
	for _, f := range obj.Relationships() {
		fields = append(fields, &ast.InputValueDefinition{
			Kind:        kinds.InputValueDefinition,
			Name:        nameNode(f.Name.Value + "_id"),
			Description: f.Description,
			Type:        namedType("ID"),
		})
	}
	return &ast.InputObjectDefinition{
		Kind:   kinds.InputObjectDefinition,
		Name:   nameNode(obj.Name() + "CreateInput"),
		Fields: fields,
	}
}

func updateObjectDefinition(obj Object) *ast.InputObjectDefinition {
	fields := []*ast.InputValueDefinition{}
	for _, f := range obj.Columns() {
		fields = append(fields, &ast.InputValueDefinition{
			Kind:        kinds.InputValueDefinition,
			Name:        f.Name,
			Description: f.Description,
			Type:        getNamedType(f.Type),
		})
	}
	for _, f := range obj.Relationships() {
		fields = append(fields, &ast.InputValueDefinition{
			Kind:        kinds.InputValueDefinition,
			Name:        nameNode(f.Name.Value + "_id"),
			Description: f.Description,
			Type:        namedType("ID"),
		})
	}
	return &ast.InputObjectDefinition{
		Kind:   kinds.InputObjectDefinition,
		Name:   nameNode(obj.Name() + "UpdateInput"),
		Fields: fields,
	}
}
