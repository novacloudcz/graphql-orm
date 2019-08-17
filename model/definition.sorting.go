package model

import (
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/graphql-go/graphql/language/kinds"

	"github.com/graphql-go/graphql/language/ast"
)

func createObjectSortType(obj Object) *ast.EnumDefinition {
	values := []*ast.EnumValueDefinition{}

	for _, col := range obj.Columns() {
		if col.IsReadonlyType() {
			continue
		}
		colName := strings.ToUpper(strcase.ToSnake(col.Name()))
		asc := ast.EnumValueDefinition{
			Kind: kinds.EnumValueDefinition,
			Name: nameNode(colName + "_ASC"),
		}
		desc := ast.EnumValueDefinition{
			Kind: kinds.EnumValueDefinition,
			Name: nameNode(colName + "_DESC"),
		}
		values = append(values, &asc, &desc)
	}

	return &ast.EnumDefinition{
		Kind:   kinds.EnumDefinition,
		Name:   nameNode(obj.Name() + "SortType"),
		Values: values,
	}
}
