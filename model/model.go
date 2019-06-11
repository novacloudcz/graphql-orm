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
			objs = append(objs, Object{op, m})
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
