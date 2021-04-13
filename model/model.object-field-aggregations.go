package model

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
)

// ObjectFieldAggregation specifies which aggregation functions are supported for given field
type ObjectFieldAggregation struct {
	Field string
	Name  string
	Type  ast.Type
}

// FieldName ...
func (a *ObjectFieldAggregation) FieldName() string {
	return a.Field + a.Name
}

// SQLColumn ...
func (a *ObjectFieldAggregation) SQLColumn() string {
	return fmt.Sprintf("%s(%s) as %s_%s", strings.ToUpper(a.Name), a.Field, a.Field, strings.ToLower(a.Name))
}

// IsAggregable ...
func (o *ObjectField) IsAggregable() bool {
	return o.IsString() || o.IsNumeric()
}

// Aggregations ...
func (o *ObjectField) Aggregations() []ObjectFieldAggregation {
	res := []ObjectFieldAggregation{
		{Field: o.Name(), Name: "Min", Type: o.Def.Type},
		{Field: o.Name(), Name: "Max", Type: o.Def.Type},
	}
	if o.IsNumeric() {
		res = append(res,
			ObjectFieldAggregation{Field: o.Name(), Name: "Avg", Type: namedType("Float")},
			ObjectFieldAggregation{Field: o.Name(), Name: "Sum", Type: o.Def.Type},
		)
	}
	return res
}
