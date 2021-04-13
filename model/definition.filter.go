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
		if obj.IsToManyColumn(col) || !col.IsFilterable() {
			continue
		}
		fields = append(fields, filterInputValues(&col, col.InputType())...)
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

func filterInputValues(column *ObjectField, t ast.Type) []*ast.InputValueDefinition {
	values := []*ast.InputValueDefinition{}
	for _, val := range column.FilterMapping() {
		name := fmt.Sprintf("%s%s", column.Name(), val.Suffix)
		values = append(values, filterInputValue(name, val.InputType))
		for _, agg := range column.Aggregations() {
			name := fmt.Sprintf("%s%s%s", column.Name(), agg.Name, val.Suffix)
			values = append(values, filterInputValue(name, val.InputType))
		}
	}
	nullName := fmt.Sprintf("%s_null", column.Name())
	values = append(values, filterInputValue(nullName, namedType("Boolean")))
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
