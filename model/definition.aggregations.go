package model

import (
	"fmt"

	"github.com/graphql-go/graphql/language/kinds"

	"github.com/graphql-go/graphql/language/ast"
)

func createObjectAggregationsType(obj Object) *ast.ObjectDefinition {
	name := obj.Name() + "Aggregations"

	fields := []*ast.FieldDefinition{}

	for _, col := range obj.Columns() {
		if obj.IsToManyColumn(col) || !col.IsAggregable() {
			continue
		}
		for _, fn := range col.Aggregations() {
			name := fmt.Sprintf("%s%s", col.Name(), fn.Name)
			fields = append(fields, columnDefinitionWithType(name, getNullableType(col.Def.Type)))
		}
	}

	return &ast.ObjectDefinition{
		Kind:   kinds.ObjectDefinition,
		Name:   nameNode(name),
		Fields: fields,
	}
}
