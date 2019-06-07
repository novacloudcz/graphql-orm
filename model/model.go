package model

import (
	"strings"

	"github.com/graphql-go/graphql/language/ast"
)

type Object struct {
	Def *ast.ObjectDefinition
}

func (o *Object) Name() string {
	return o.Def.Name.Value
}
func (o *Object) LowerName() string {
	return strings.ToLower(o.Def.Name.Value)
}
func isColumn(f *ast.FieldDefinition) bool {
	v, ok := getNamedType(f.Type).(*ast.Named)
	if ok {
		switch v.Name.Value {
		case "String":
			fallthrough
		case "Time":
			fallthrough
		case "ID":
			fallthrough
		case "Float":
			fallthrough
		case "Int":
			fallthrough
		case "Boolean":
			return true
		}
	}
	return false
}
func isRelationship(f *ast.FieldDefinition) bool {
	for _, d := range f.Directives {
		if d != nil && d.Name.Value == "relationship" {
			return true
		}
	}
	return false
}
func (o *Object) Columns() []*ast.FieldDefinition {
	fields := []*ast.FieldDefinition{}
	for _, f := range o.Def.Fields {
		if isColumn(f) {
			fields = append(fields, f)
		}
	}
	return fields
}
func (o *Object) Relationships() []*ast.FieldDefinition {
	fields := []*ast.FieldDefinition{}
	for _, f := range o.Def.Fields {
		if isRelationship(f) {
			fields = append(fields, f)
		}
	}
	return fields
}

type Model struct {
	Doc *ast.Document
	// Objects []Object
}

func (m *Model) Objects() []Object {
	objs := []Object{}
	for _, def := range m.Doc.Definitions {
		op, ok := def.(*ast.ObjectDefinition)
		if ok {
			objs = append(objs, Object{op})
		}
	}
	return objs
}
