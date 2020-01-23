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
		def, ok := def.(*ast.ObjectDefinition)
		if ok {
			objs = append(objs, Object{Def: def, Model: m})
		}
	}
	return objs
}
func (m *Model) ObjectEntities() []Object {
	objs := []Object{}
	for _, def := range m.Doc.Definitions {
		def, ok := def.(*ast.ObjectDefinition)
		if ok {
			obj := Object{Def: def, Model: m}
			if obj.HasDirective("entity") {
				objs = append(objs, obj)
			}
		}
	}
	return objs
}

func (m *Model) EmbeddedObjects() []Object {
	objs := []Object{}
	objsMap := map[string]bool{}
	for _, obj := range m.ObjectEntities() {
		for _, col := range obj.Columns() {
			if col.IsEmbeddedColumn() {
				obj := col.TargetObject()
				if _, exists := objsMap[obj.Name()]; !exists {
					objs = append(objs, *obj)
					objsMap[obj.Name()] = true
				}
			}
		}
	}
	return objs
}

func (m *Model) ObjectExtensions() []ObjectExtension {
	objs := []ObjectExtension{}
	for _, def := range m.Doc.Definitions {
		def, ok := def.(*ast.TypeExtensionDefinition)
		if ok {
			obj := &Object{Def: def.Definition, Model: m}
			objs = append(objs, ObjectExtension{Def: def, Model: m, Object: obj})
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
func (m *Model) ObjectExtension(name string) ObjectExtension {
	for _, e := range m.ObjectExtensions() {
		if e.Object.Name() == name {
			return e
		}
	}
	panic(fmt.Sprintf("Extension for object with name %s not found in model", name))
}

func (m *Model) HasObject(name string) bool {
	if name == "Query" || name == "Mutation" || name == "Subscription" {
		return true
	}
	for _, o := range m.Objects() {
		if o.Name() == name {
			return true
		}
	}
	return false
}

func (m *Model) HasObjectExtension(name string) bool {
	for _, e := range m.ObjectExtensions() {
		if e.Object.Name() == name {
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

func (m *Model) RemoveObjectExtension(oe *ObjectExtension) {
	newDefinitions := []ast.Node{}
	for _, d := range m.Doc.Definitions {
		if d != oe.Def {
			newDefinitions = append(newDefinitions, d)
		}
	}
	m.Doc.Definitions = newDefinitions
}
