package model

import (
	"github.com/graphql-go/graphql/language/kinds"

	"github.com/graphql-go/graphql/language/ast"
)

func createObjectDefinition(obj Object) *ast.InputObjectDefinition {
	fields := []*ast.InputValueDefinition{}
	for _, col := range obj.Columns() {
		if col.Name() == "createdAt" || col.Name() == "updatedAt" || col.Name() == "createdBy" || col.Name() == "updatedBy" {
			continue
		}
		t := col.Def.Type
		if col.Name() == "id" {
			t = getNamedType(t)
		}
		fields = append(fields, &ast.InputValueDefinition{
			Kind:        kinds.InputValueDefinition,
			Name:        col.Def.Name,
			Description: col.Def.Description,
			Type:        t,
		})
	}
	for _, rel := range obj.Relationships() {
		if rel.IsToMany() {
			fields = append(fields, &ast.InputValueDefinition{
				Kind:        kinds.InputValueDefinition,
				Name:        nameNode(rel.Name() + "Ids"),
				Description: rel.Def.Description,
				Type:        listType(nonNull(namedType("ID"))),
			})
		}
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
		if col.Name() == "id" || col.Name() == "createdAt" || col.Name() == "updatedAt" || col.Name() == "createdBy" || col.Name() == "updatedBy" {
			continue
		}
		fields = append(fields, &ast.InputValueDefinition{
			Kind:        kinds.InputValueDefinition,
			Name:        col.Def.Name,
			Description: col.Def.Description,
			Type:        getNamedType(col.Def.Type),
		})
	}
	for _, rel := range obj.Relationships() {
		if rel.IsToMany() {
			fields = append(fields, &ast.InputValueDefinition{
				Kind:        kinds.InputValueDefinition,
				Name:        nameNode(rel.Name() + "Ids"),
				Description: rel.Def.Description,
				Type:        listType(nonNull(namedType("ID"))),
			})
		}
	}
	return &ast.InputObjectDefinition{
		Kind:   kinds.InputObjectDefinition,
		Name:   nameNode(obj.Name() + "UpdateInput"),
		Fields: fields,
	}
}
