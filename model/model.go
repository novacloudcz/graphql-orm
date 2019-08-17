package model

import (
	"fmt"

	"github.com/graphql-go/graphql/language/ast"
)

type Model struct {
	Doc *ast.Document
	// Objects []Object
}

func (m *Model) Objects() []Object {
	objs := []Object{}
	for _, def := range m.Doc.Definitions {
		op, ok := def.(*ast.ObjectDefinition)
		if ok {
			objs = append(objs, Object{Def: op, Model: m, IsExtended: false})
		}
	}
	for _, def := range m.Doc.Definitions {
		op, ok := def.(*ast.TypeExtensionDefinition)
		if ok {
			objs = append(objs, Object{Def: op.Definition, Model: m, IsExtended: true})
		}
	}
	return objs
}

func (m *Model) Object(name string) Object {
	for _, o := range m.Objects() {
		if o.Name() == name {
			return o
		}
	}
	panic(fmt.Sprintf("Object with name %s not found in model", name))
}

func (m *Model) HasObject(name string) bool {
	for _, o := range m.Objects() {
		if o.Name() == name {
			return true
		}
	}
	return false
}

var defaultScalars map[string]bool = map[string]bool{
	"Int":     true,
	"Float":   true,
	"String":  true,
	"Boolean": true,
	"ID":      true,
	"Any":     true,
	"Time":    true,
}

func (m *Model) HasScalar(name string) bool {
	if _, ok := defaultScalars[name]; ok {
		return true
	}
	for _, def := range m.Doc.Definitions {
		scalar, ok := def.(*ast.ScalarDefinition)
		if ok && scalar.Name.Value == name {
			return true
		}
	}
	return false
}

func (m *Model) HasEnum(name string) bool {
	if _, ok := defaultScalars[name]; ok {
		return true
	}
	for _, def := range m.Doc.Definitions {
		e, ok := def.(*ast.EnumDefinition)
		if ok && e.Name.Value == name {
			return true
		}
	}
	return false
}
