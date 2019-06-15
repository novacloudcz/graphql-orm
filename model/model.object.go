package model

import (
	"fmt"

	"github.com/jinzhu/inflection"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/iancoleman/strcase"
)

type Object struct {
	Def   *ast.ObjectDefinition
	Model *Model
}

func (o *Object) Name() string {
	return o.Def.Name.Value
}
func (o *Object) PluralName() string {
	return inflection.Plural(o.Name())
}
func (o *Object) LowerName() string {
	return strcase.ToLowerCamel(o.Def.Name.Value)
}
func (o *Object) TableName() string {
	return inflection.Plural(o.LowerName())
}
func (o *Object) Columns() []ObjectColumn {
	columns := []ObjectColumn{}
	for _, f := range o.Def.Fields {
		if o.isColumn(f) {
			columns = append(columns, ObjectColumn{f, o})
		}
	}
	return columns
}
func (o *Object) Relationships() []*ObjectRelationship {
	relationships := []*ObjectRelationship{}
	for _, f := range o.Def.Fields {
		if o.isRelationship(f) {
			relationships = append(relationships, &ObjectRelationship{f, o})
		}
	}
	return relationships
}

func (o *Object) Relationship(name string) *ObjectRelationship {
	for _, rel := range o.Relationships() {
		if rel.Name() == name {
			return rel
		}
	}
	panic(fmt.Sprintf("relationship %s->%s not found", o.Name(), name))
}
func (o *Object) HasRelationships() bool {
	return len(o.Relationships()) > 0
}

func (o *Object) isColumn(f *ast.FieldDefinition) bool {
	return !o.Model.HasObject(getNamedType(f.Type).(*ast.Named).Name.Value) && !o.isRelationship(f)
}
func (o *Object) isRelationship(f *ast.FieldDefinition) bool {
	for _, d := range f.Directives {
		if d != nil && d.Name.Value == "relationship" {
			return true
		}
	}
	return false
}
