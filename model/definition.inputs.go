package model

import (
	"github.com/graphql-go/graphql/language/kinds"

	"github.com/graphql-go/graphql/language/ast"
)

func createObjectDefinition(obj Object) *ast.InputObjectDefinition {
	fields := []*ast.InputValueDefinition{
		&ast.InputValueDefinition{
			Kind: kinds.InputValueDefinition,
			Name: nameNode("id"),
			Description: &ast.StringValue{
				Kind:  kinds.StringValue,
				Value: "Entity identifier. If not specified, the generated UUID is used.",
			},
			Type: namedType("String"),
		},
	}
	for _, col := range obj.Columns() {
		fields = append(fields, &ast.InputValueDefinition{
			Kind:        kinds.InputValueDefinition,
			Name:        col.Def.Name,
			Description: col.Def.Description,
			Type:        col.Def.Type,
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
		} else {
			fields = append(fields, &ast.InputValueDefinition{
				Kind:        kinds.InputValueDefinition,
				Name:        nameNode(rel.Name() + "Id"),
				Description: rel.Def.Description,
				Type:        namedType("ID"),
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
		} else {
			fields = append(fields, &ast.InputValueDefinition{
				Kind:        kinds.InputValueDefinition,
				Name:        nameNode(rel.Name() + "Id"),
				Description: rel.Def.Description,
				Type:        namedType("ID"),
			})
		}
	}
	return &ast.InputObjectDefinition{
		Kind:   kinds.InputObjectDefinition,
		Name:   nameNode(obj.Name() + "UpdateInput"),
		Fields: fields,
	}
}
