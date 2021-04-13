package model

import (
	"github.com/graphql-go/graphql/language/kinds"

	"github.com/graphql-go/graphql/language/ast"
)

func objectResultTypeDefinition(o *Object) *ast.ObjectDefinition {
	fields := []*ast.FieldDefinition{
		{
			Kind: kinds.FieldDefinition,
			Name: nameNode("items"),
			Type: nonNull(&ast.List{
				Kind: kinds.List,
				Type: nonNull(namedType(o.Name())),
			}),
		},
		{
			Kind: kinds.FieldDefinition,
			Name: nameNode("count"),
			Type: nonNull(namedType("Int")),
		},
	}

	if o.HasAggregableColumn() {
		fields = append(fields, &ast.FieldDefinition{
			Kind: kinds.FieldDefinition,
			Name: nameNode("aggregations"),
			Type: nonNull(namedType(o.Name() + "ResultAggregations")),
		})
	}

	return &ast.ObjectDefinition{
		Kind:   kinds.ObjectDefinition,
		Name:   nameNode(o.Name() + "ResultType"),
		Fields: fields,
	}
}

func objectResultTypeAggregationsDefinition(o *Object) *ast.ObjectDefinition {
	fields := []*ast.FieldDefinition{}

	for _, column := range o.Columns() {
		if column.IsAggregable() {
			for _, aggregation := range column.Aggregations() {
				fields = append(fields, &ast.FieldDefinition{
					Kind: kinds.FieldDefinition,
					Name: nameNode(column.Name() + aggregation.Name),
					Type: aggregation.Type,
				})
			}
		}
	}

	return &ast.ObjectDefinition{
		Kind:   kinds.ObjectDefinition,
		Name:   nameNode(o.Name() + "ResultAggregations"),
		Fields: fields,
	}
}
