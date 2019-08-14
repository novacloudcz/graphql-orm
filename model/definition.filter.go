package model

import (
	"fmt"

	"github.com/graphql-go/graphql/language/kinds"

	"github.com/graphql-go/graphql/language/ast"
)

func createObjectFilterType(obj Object) *ast.InputObjectDefinition {
	name := obj.Name() + "FilterType"

	fields := []*ast.InputValueDefinition{
		filterInputValue("AND", listType(nonNull(namedType(name)))),
		filterInputValue("OR", listType(nonNull(namedType(name)))),
	}
	// fields = append(fields, filterInputValues("id", namedType("ID"))...)

	for _, col := range obj.Columns() {
		if obj.IsToManyColumn(col) || !col.IsScalarType() {
			continue
		}
		fields = append(fields, filterInputValues(&col, col.Def.Type)...)
	}
	for _, rel := range obj.Relationships() {
		fields = append(fields, filterInputValue(rel.Name(), namedType(rel.Target().Name()+"FilterType")))
	}

	return &ast.InputObjectDefinition{
		Kind:   kinds.InputObjectDefinition,
		Name:   nameNode(name),
		Fields: fields,
	}
}

func filterInputValues(column *ObjectColumn, t ast.Type) []*ast.InputValueDefinition {
	values := []*ast.InputValueDefinition{}
	for _, val := range column.FilterMapping() {
		values = append(values, filterInputValue(fmt.Sprintf("%s%s", column.Name(), val.Suffix), val.InputType))
	}
	return values
}

func filterInputValue(name string, t ast.Type) *ast.InputValueDefinition {
	return &ast.InputValueDefinition{
		Kind: kinds.InputValueDefinition,
		Name: nameNode(name),
		Description: &ast.StringValue{
			Kind:  kinds.StringValue,
			Value: "Entity identifier. If not specified, the generated UUID is used.",
		},
		Type: t,
	}
}
