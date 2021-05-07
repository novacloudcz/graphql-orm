package model

import (
	"fmt"
	"strings"

	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/iancoleman/strcase"
)

// ObjectField ...
type ObjectField struct {
	Def *ast.FieldDefinition
	Obj *Object
}

// Name ...
func (o *ObjectField) Name() string {
	return o.Def.Name.Value
}

// MethodName ...
func (o *ObjectField) MethodName() string {
	name := o.Name()
	return templates.ToGo(name)
}

// TargetType ...
func (o *ObjectField) TargetType() string {
	nt := getNamedType(o.Def.Type).(*ast.Named)
	return nt.Name.Value
}

// IsColumn ...
func (o *ObjectField) IsColumn() bool {
	return o.HasDirective("column")
}

// IsExternal ...
func (o *ObjectField) IsExternal() bool {
	return o.HasDirective("external")
}

// IsIdentifier ...
func (o *ObjectField) IsIdentifier() bool {
	return o.Name() == "id"
}

// IsRelationshipIdentifier ...
func (o *ObjectField) IsRelationshipIdentifier() bool {
	return strings.HasSuffix(o.Name(), "Id") || strings.HasSuffix(o.Name(), "Ids")
}

// IsRelationship ...
func (o *ObjectField) IsRelationship() bool {
	return o.HasDirective("relationship")
}

// IsCreatable ...
func (o *ObjectField) IsCreatable() bool {
	return !(o.Name() == "createdAt" || o.Name() == "updatedAt" || o.Name() == "createdBy" || o.Name() == "updatedBy") && !o.IsReadonlyType()
}

// IsUpdatable ...
func (o *ObjectField) IsUpdatable() bool {
	return !(o.IsIdentifier() || o.Name() == "createdAt" || o.Name() == "updatedAt" || o.Name() == "createdBy" || o.Name() == "updatedBy") && !o.IsReadonlyType()
}

// IsReadonlyType ..
func (o *ObjectField) IsReadonlyType() bool {
	if o.IsEmbeddedColumn() {
		return false
	}
	return !(o.IsScalarType() || o.IsEnumType()) || o.Obj.Model.HasObject(o.TargetType())
}

// IsWritableType ...
func (o *ObjectField) IsWritableType() bool {
	return !o.IsReadonlyType()
}

// IsFilterable ...
func (o *ObjectField) IsFilterable() bool {
	return !o.IsReadonlyType() && !o.IsEmbedded()
}

// IsScalarType ...
func (o *ObjectField) IsScalarType() bool {
	return o.Obj.Model.HasScalar(o.TargetType())
}

// IsEnumType ...
func (o *ObjectField) IsEnumType() bool {
	return o.Obj.Model.HasEnum(o.TargetType())
}

// IsOptional ...
func (o *ObjectField) IsOptional() bool {
	return !isNonNullType(o.Def.Type)
}

// IsList ...
func (o *ObjectField) IsList() bool {
	return isListType(o.Def.Type)
}

// IsEmbedded ...
func (o *ObjectField) IsEmbedded() bool {
	return !o.IsColumn() && !o.IsRelationship() || o.IsEmbeddedColumn()
}

// IsEmbeddedColumn ...
func (o *ObjectField) IsEmbeddedColumn() bool {
	return (o.IsColumn() && o.ColumnType() == "embedded")
}

// HasTargetObject ...
func (o *ObjectField) HasTargetObject() bool {
	return o.Obj.Model.HasObject(o.TargetType())
}

// TargetObject ...
func (o *ObjectField) TargetObject() *Object {
	obj := o.Obj.Model.Object(o.TargetType())
	return &obj
}

// HasTargetObjectExtension ...
func (o *ObjectField) HasTargetObjectExtension() bool {
	return o.Obj.Model.HasObjectExtension(o.TargetType())
}

// TargetObjectExtension ...
func (o *ObjectField) TargetObjectExtension() *ObjectExtension {
	e := o.Obj.Model.ObjectExtension(o.TargetType())
	return &e
}

// IsSortable ...
func (o *ObjectField) IsSortable() bool {
	return !o.IsReadonlyType() && o.IsScalarType()
}

// IsSearchable ...
func (o *ObjectField) IsSearchable() bool {
	return o.IsString() || o.IsNumeric()
}

// IsNumeric ...
func (o *ObjectField) IsNumeric() bool {
	t := getNamedType(o.Def.Type).(*ast.Named)
	return t.Name.Value == "Int" || t.Name.Value == "Float"
}

// IsString ...
func (o *ObjectField) IsString() bool {
	t := getNamedType(o.Def.Type).(*ast.Named)
	return t.Name.Value == "String"
}

// NeedsQueryResolver ...
func (o *ObjectField) NeedsQueryResolver() bool {
	return o.IsEmbedded()
}

// HasTargetTypeWithIDField ...
func (o *ObjectField) HasTargetTypeWithIDField() bool {
	if o.HasTargetObject() && o.TargetObject().HasField("id") {
		return true
	}
	if o.HasTargetObjectExtension() && o.TargetObjectExtension().Object.HasField("id") {
		return true
	}
	return false
}

// GoType ...
func (o *ObjectField) GoType() string {
	return o.GoTypeWithPointer(true, false)
}

// GoTypeWithPointer ...
func (o *ObjectField) GoTypeWithPointer(showPointer, ignoreEmbedded bool) string {
	t := o.Def.Type
	st := ""

	if o.IsOptional() && showPointer {
		st += "*"
	} else {
		t = getNullableType(t)
	}

	if o.IsEmbeddedColumn() && !ignoreEmbedded {
		return st + "string"
	}

	if isListType(t) {
		st = "[]*"
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

// GoResultType ...
func (o *ObjectField) GoResultType() string {
	return o.GoTypeWithPointer(true, true)
}

// InputType ...
func (o *ObjectField) InputType() ast.Type {
	t := o.Def.Type
	if o.IsIdentifier() {
		t = nonNull(getNamedType(t))
	}
	isList := o.IsList()
	isOptional := o.IsOptional()

	if o.IsEmbeddedColumn() {
		_t := getNamedType(t).(*ast.Named)
		t = namedType(_t.Name.Value + "Input")

		if isList {
			t = listType(t)
		}
		if !isOptional {
			t = nonNull(t)
		}
	}
	if o.IsRelationshipIdentifier() {
		t = getNullableType(t)
	}

	return t
}

// InputTypeName ...
func (o *ObjectField) InputTypeName() string {
	t := o.InputType()
	return astTypeToGoType(t)
}

// FilterMappingItem ...
type FilterMappingItem struct {
	Suffix      string
	Operator    string
	InputType   ast.Type
	ValueFormat string
}

// SuffixCamel ...
func (f *FilterMappingItem) SuffixCamel() string {
	return strcase.ToCamel(f.Suffix)
}

// WrapValueVariable ...
func (f *FilterMappingItem) WrapValueVariable(v string) string {
	return fmt.Sprintf(f.ValueFormat, v)
}

// FilterMapping ...
func (o *ObjectField) FilterMapping() []FilterMappingItem {
	t := getNamedType(o.Def.Type)
	mapping := []FilterMappingItem{
		{"", "= ?", t, "%s"},
		{"_ne", "!= ?", t, "%s"},
		{"_gt", "> ?", t, "%s"},
		{"_lt", "< ?", t, "%s"},
		{"_gte", ">= ?", t, "%s"},
		{"_lte", "<= ?", t, "%s"},
		{"_in", "IN (?)", listType(nonNull(t)), "%s"},
		{"_not_in", "NOT IN (?)", listType(nonNull(t)), "%s"},
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
