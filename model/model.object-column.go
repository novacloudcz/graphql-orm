package model

import (
	"fmt"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/iancoleman/strcase"
)

var goTypeMap = map[string]string{
	"String":  "string",
	"Time":    "time.Time",
	"ID":      "string",
	"Float":   "float64",
	"Int":     "int64",
	"Boolean": "bool",
}

type ObjectColumn struct {
	Def *ast.FieldDefinition
	Obj *Object
}

func (o *ObjectColumn) Name() string {
	return o.Def.Name.Value
}
func (o *ObjectColumn) MethodName() string {
	return strcase.ToCamel(o.Def.Name.Value)
}

func (o *ObjectColumn) TargetType() string {
	nt := getNamedType(o.Def.Type).(*ast.Named)
	return nt.Name.Value
}
func (o *ObjectColumn) IsOptional() bool {
	return o.Def.Type.GetKind() != "NonNull"
}
func (o *ObjectColumn) GoType() string {
	t := ""

	if o.IsOptional() {
		t += "*"
	}

	v, ok := getNamedType(o.Def.Type).(*ast.Named)
	if ok {
		_t, known := goTypeMap[v.Name.Value]
		if known {
			t += _t
		} else {
			t += v.Name.Value
		}
	}
	return t
}

func (o *ObjectColumn) ModelTags() string {
	return fmt.Sprintf(`json:"%s"`, o.Name())
}
