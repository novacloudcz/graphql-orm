package model

import (
	"github.com/graphql-go/graphql/language/kinds"

	"github.com/graphql-go/graphql/language/ast"
)

func createObjectSortEnum() *ast.EnumDefinition {
	values := []*ast.EnumValueDefinition{
		&ast.EnumValueDefinition{
			Kind: kinds.EnumValueDefinition,
			Name: nameNode("ASC"),
		},
		&ast.EnumValueDefinition{
			Kind: kinds.EnumValueDefinition,
			Name: nameNode("DESC"),
		},
	}
	return &ast.EnumDefinition{
		Kind:   kinds.EnumDefinition,
		Name:   nameNode("ObjectSortType"),
		Values: values,
	}
}
func createObjectSortType(obj Object) *ast.InputObjectDefinition {
	fields := []*ast.InputValueDefinition{}

	for _, col := range obj.Columns() {
		if col.IsReadonlyType() {
			continue
		}

		field := ast.InputValueDefinition{
			Kind: kinds.InputValueDefinition,
			Name: nameNode(col.Name()),
			Type: namedType("ObjectSortType"),
		}
		fields = append(fields, &field)
	}

	for _, rel := range obj.Relationships() {
		field := ast.InputValueDefinition{
			Kind: kinds.InputValueDefinition,
			Name: nameNode(rel.Name()),
			Type: namedType(rel.Target().Name() + "SortType"),
		}
		fields = append(fields, &field)
	}

	return &ast.InputObjectDefinition{
		Kind:   kinds.InputObjectDefinition,
		Name:   nameNode(obj.Name() + "SortType"),
		Fields: fields,
	}
}
