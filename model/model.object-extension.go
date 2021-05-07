package model

import (
	"github.com/graphql-go/graphql/language/ast"
)

// ObjectExtension ...
type ObjectExtension struct {
	Def    *ast.TypeExtensionDefinition
	Model  *Model
	Object *Object
}

// func (oe *ObjectExtension) GetObject() *Object {
// 	return &Object{
// 		Def:   oe.Def.Definition,
// 		Model: oe.Model,
// 		Extension: oe,
// 	}
// }

// IsFederatedType ...
func (oe *ObjectExtension) IsFederatedType() bool {
	return oe.Object.IsFederatedType()
}

// ExtendsLocalObject ...
func (oe *ObjectExtension) ExtendsLocalObject() bool {
	return oe.Model.HasObject(oe.Object.Name())
}

// IsExternal ...
func (oe *ObjectExtension) HasAnyNonExternalField() bool {
	for _, f := range oe.Object.Fields() {
		if !f.IsExternal() {
			return true
		}
	}
	return false
}
