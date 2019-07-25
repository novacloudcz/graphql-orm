package model

import (
	"github.com/graphql-go/graphql/language/kinds"

	"github.com/graphql-go/graphql/language/ast"
)

func createObjectDefinition(obj Object) *ast.InputObjectDefinition {
	fields := []*ast.InputValueDefinition{}
	for _, col := range obj.Columns() {
		if !col.IsCreatable() {
			continue
		}
		t := col.Def.Type
		if col.Name() == "id" {
			t = getNamedType(t)
		}
		if isListType(getNullableType(t)) {
			t = getNullableType(t)
		}
		fields = append(fields, &ast.InputValueDefinition{
			Kind:        kinds.InputValueDefinition,
			Name:        col.Def.Name,
			Description: col.Def.Description,
			Type:        t,
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
	for _, col := range obj.Columns() {
		if !col.IsUpdatable() {
			continue
		}
		fields = append(fields, &ast.InputValueDefinition{
			Kind:        kinds.InputValueDefinition,
			Name:        col.Def.Name,
			Description: col.Def.Description,
			Type:        getNullableType(col.Def.Type),
		})
	}
	return &ast.InputObjectDefinition{
		Kind:   kinds.InputObjectDefinition,
		Name:   nameNode(obj.Name() + "UpdateInput"),
		Fields: fields,
	}
}
