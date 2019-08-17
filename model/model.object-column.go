package model

import (
	"fmt"

	"github.com/99designs/gqlgen/codegen/templates"
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
	name := o.Name()
	return templates.ToGo(name)
}

func (o *ObjectColumn) TargetType() string {
	nt := getNamedType(o.Def.Type).(*ast.Named)
	return nt.Name.Value
}
func (o *ObjectColumn) IsCreatable() bool {
	return !(o.Name() == "createdAt" || o.Name() == "updatedAt" || o.Name() == "createdBy" || o.Name() == "updatedBy") && !o.IsReadonlyType()
}
func (o *ObjectColumn) IsUpdatable() bool {
	return !(o.Name() == "id" || o.Name() == "createdAt" || o.Name() == "updatedAt" || o.Name() == "createdBy" || o.Name() == "updatedBy") && !o.IsReadonlyType()
}
func (o *ObjectColumn) IsReadonlyType() bool {
	return !(o.IsScalarType() || o.IsEnumType()) || o.Obj.Model.HasObject(o.TargetType())
}
func (o *ObjectColumn) IsWritableType() bool {
	return !o.IsReadonlyType()
}
func (o *ObjectColumn) IsScalarType() bool {
	return o.Obj.Model.HasScalar(o.TargetType())
}
func (o *ObjectColumn) IsEnumType() bool {
	return o.Obj.Model.HasEnum(o.TargetType())
}
func (o *ObjectColumn) IsOptional() bool {
	return !isNonNullType(o.Def.Type)
}
func (o *ObjectColumn) IsList() bool {
	return isListType(o.Def.Type)
}
func (o *ObjectColumn) IsSearchable() bool {
	t := getNamedType(o.Def.Type).(*ast.Named)
	return t.Name.Value == "String"
}
func (o *ObjectColumn) Directive(name string) *ast.Directive {
	for _, d := range o.Def.Directives {
		if d.Name.Value == name {
			return d
		}
	}
	return nil
}
func (o *ObjectColumn) HasDirective(name string) bool {
	return o.Directive(name) != nil
}

func (o *ObjectColumn) GoType() string {
	return o.GoTypeWithPointer(true)
}
func (o *ObjectColumn) GoTypeWithPointer(showPointer bool) string {
	t := o.Def.Type
	st := ""

	if o.IsOptional() && showPointer {
		st += "*"
	} else {
		t = getNullableType(t)
	}

	if isListType(t) {
		st += "[]*"
	}

	v, ok := getNamedType(o.Def.Type).(*ast.Named)
	if ok {
		_t, known := goTypeMap[v.Name.Value]
		if known {
			st += _t
		} else {
			st += v.Name.Value
		}
	}

	return st
}

func (o *ObjectColumn) ModelTags() string {
	_gorm := fmt.Sprintf("column:%s", o.Name())
	if o.Name() == "id" {
		_gorm += ";primary_key"
	}

	for _, d := range o.Def.Directives {
		if d.Name.Value == "column" {
			for _, arg := range d.Arguments {
				if arg.Name.Value == "type" {
					_gorm += fmt.Sprintf(";type:%v", arg.Value.GetValue())
				}
			}
		}
	}

	return fmt.Sprintf(`json:"%s" gorm:"%s"`, o.Name(), _gorm)
}

type FilterMappingItem struct {
	Suffix      string
	Operator    string
	InputType   ast.Type
	ValueFormat string
}

func (f *FilterMappingItem) SuffixCamel() string {
	return strcase.ToCamel(f.Suffix)
}
func (f *FilterMappingItem) WrapValueVariable(v string) string {
	return fmt.Sprintf(f.ValueFormat, v)
}

func (o *ObjectColumn) FilterMapping() []FilterMappingItem {
	t := getNamedType(o.Def.Type)
	mapping := []FilterMappingItem{
		FilterMappingItem{"", "= ?", t, "%s"},
		FilterMappingItem{"_ne", "!= ?", t, "%s"},
		FilterMappingItem{"_gt", "> ?", t, "%s"},
		FilterMappingItem{"_lt", "< ?", t, "%s"},
		FilterMappingItem{"_gte", ">= ?", t, "%s"},
		FilterMappingItem{"_lte", "<= ?", t, "%s"},
		FilterMappingItem{"_in", "IN (?)", listType(nonNull(t)), "%s"},
	}
	_t := t.(*ast.Named)
	if _t.Name.Value == "String" {
		mapping = append(mapping,
			FilterMappingItem{"_like", "LIKE ?", t, "strings.Replace(strings.Replace(*%s,\"?\",\"_\",-1),\"*\",\"%%\",-1)"},
			FilterMappingItem{"_prefix", "LIKE ?", t, "fmt.Sprintf(\"%%s%%%%\",*%s)"},
			FilterMappingItem{"_suffix", "LIKE ?", t, "fmt.Sprintf(\"%%%%%%s\",*%s)"},
		)
	}
	return mapping
}
