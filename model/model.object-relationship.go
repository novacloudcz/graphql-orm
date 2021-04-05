package model

import (
	"fmt"
	"strings"

	"github.com/jinzhu/inflection"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/iancoleman/strcase"
)

type ObjectRelationship struct {
	Def *ast.FieldDefinition
	Obj *Object
}

func (o *ObjectRelationship) Name() string {
	return o.Def.Name.Value
}
func (o *ObjectRelationship) MethodName() string {
	return strcase.ToCamel(o.Def.Name.Value)
}
func (o *ObjectRelationship) ValueForRelationshipDirectiveAttribute(name string) (val interface{}, ok bool) {
	for _, d := range o.Def.Directives {
		if d.Name.Value == "relationship" {
			for _, arg := range d.Arguments {
				if arg.Name.Value == name {
					val = arg.Value.GetValue()
					ok = true
					return
				}
			}
		}
	}
	return
}
func (o *ObjectRelationship) StringForRelationshipDirectiveAttribute(name string) (val string, ok bool) {
	value, ok := o.ValueForRelationshipDirectiveAttribute(name)
	if !ok {
		return
	}
	val, ok = value.(string)
	if !ok {
		panic(fmt.Sprintf("invalid %s value for %s->%s relationship", name, o.Obj.Name(), o.Name()))
	}
	return
}
func (o *ObjectRelationship) BoolForRelationshipDirectiveAttribute(name string) (val bool, ok bool) {
	value, ok := o.ValueForRelationshipDirectiveAttribute(name)
	if !ok {
		return
	}
	val, ok = value.(bool)
	if !ok {
		panic(fmt.Sprintf("invalid %s value for %s->%s relationship", name, o.Obj.Name(), o.Name()))
	}
	return
}
func (o *ObjectRelationship) InverseRelationshipName() string {
	val, ok := o.StringForRelationshipDirectiveAttribute("inverse")
	if !ok {
		panic(fmt.Sprintf("missing inverse value for %s->%s relationship", o.Obj.Name(), o.Name()))
	}
	return val
}

func (o *ObjectRelationship) Preload() bool {
	val, _ := o.BoolForRelationshipDirectiveAttribute("preload")
	return val
}

func (o *ObjectRelationship) Target() *Object {
	target := o.Obj.Model.Object(o.TargetType())
	return &target
}
func (o *ObjectRelationship) InverseRelationship() *ObjectRelationship {
	return o.Target().Relationship(o.InverseRelationshipName())
}

func (o *ObjectRelationship) IsToMany() bool {
	t := getNullableType(o.Def.Type)
	return isListType(t)
}
func (o *ObjectRelationship) IsToOne() bool {
	return !o.IsToMany()
}

func (o *ObjectRelationship) IsManyToMany() bool {
	return o.IsToMany() && o.InverseRelationship().IsToMany()
}
func (o *ObjectRelationship) IsManyToOne() bool {
	return o.IsToMany() && !o.InverseRelationship().IsToMany()
}
func (o *ObjectRelationship) IsOneToMany() bool {
	return !o.IsToMany() && o.InverseRelationship().IsToMany()
}
func (o *ObjectRelationship) IsSelfReferencing() bool {
	inv := o.InverseRelationship()
	return o.Obj.Name() == inv.Obj.Name() && o.Name() == inv.Name()
}
func (o *ObjectRelationship) IsMainRelationshipForManyToMany() bool {
	main := o.MainRelationshipForManyToMany()
	return o.Obj.Name() == main.Obj.Name() && o.Name() == main.Name()
}
func (o *ObjectRelationship) IsNonNull() bool {
	return isNonNullType(o.Def.Type)
}

func (o *ObjectRelationship) ReturnType() string {
	nt := getNamedType(o.Def.Type).(*ast.Named)
	if o.IsToMany() {
		return fmt.Sprintf("[]*%s", nt.Name.Value)
	}
	return fmt.Sprintf("*%s", nt.Name.Value)
}
func (o *ObjectRelationship) TargetType() string {
	nt := getNamedType(o.Def.Type).(*ast.Named)
	return nt.Name.Value
}
func (o *ObjectRelationship) GoType() string {
	return o.ReturnType()
}
func (o *ObjectRelationship) ChangesName() string {
	name := o.MethodName()
	if o.IsToMany() {
		name += "IDs"
	} else {
		name += "ID"
	}
	return name
}
func (o *ObjectRelationship) ChangesType() string {
	if o.IsToMany() {
		return "[]*string"
	} else {
		return "*string"
	}
}
func (o *ObjectRelationship) ModelTags() string {
	tags := fmt.Sprintf(`json:"%s"`, o.Name())
	invrel := o.InverseRelationship()
	if o.IsManyToOne() {
		tags += fmt.Sprintf(" gorm:\"foreignkey:%sID\"", invrel.MethodName())
	} else if o.IsManyToMany() {
		rel := o.MainRelationshipForManyToMany()
		if o.IsSelfReferencing() {
			tags += fmt.Sprintf(" gorm:\"many2many:%s;jointable_foreignkey:%s_id;association_jointable_foreignkey:%s_id\"", rel.ManyToManyJoinTable(), inflection.Singular(strings.ToLower(o.Obj.Name())), inflection.Singular(o.InverseRelationshipName()))
		} else {
			tags += fmt.Sprintf(" gorm:\"many2many:%s;jointable_foreignkey:%s_id;association_jointable_foreignkey:%s_id\"", rel.ManyToManyJoinTable(), inflection.Singular(o.InverseRelationshipName()), inflection.Singular(o.Name()))
		}
	}
	return tags
}
func (o *ObjectRelationship) ManyToManyJoinTable() string {
	m := o.MainRelationshipForManyToMany()
	return m.Obj.LowerName() + "_" + m.Name()
}
func (o *ObjectRelationship) ManyToManyObjectName() string {
	m := o.MainRelationshipForManyToMany()
	return m.Obj.Name() + "_" + m.Name()
}
func (o *ObjectRelationship) MainRelationshipForManyToMany() *ObjectRelationship {
	inversed := o.InverseRelationship()
	if inversed.Name() > o.Name() {
		return inversed
	}
	return o
}
func (o *ObjectRelationship) JoinString() string {
	join := ""
	if o.IsManyToMany() {
		joinTable := o.ManyToManyJoinTable()
		join += fmt.Sprintf("\"LEFT JOIN \"+dialect.Quote(TableName(\"%[1]s\"))+\" \"+dialect.Quote(_alias+\"_jointable\")+\" ON \"+dialect.Quote(alias)+\".id = \"+dialect.Quote(_alias+\"_jointable\")+\".\"+dialect.Quote(\"%[3]s_id\")+\" LEFT JOIN \"+dialect.Quote(TableName(\"%[2]s\"))+\" \"+dialect.Quote(_alias)+\" ON \"+dialect.Quote(_alias+\"_jointable\")+\".\"+dialect.Quote(\"%[4]s_id\")+\" = \"+dialect.Quote(_alias)+\".id\"", joinTable, o.Target().TableName(), inflection.Singular(o.InverseRelationshipName()), inflection.Singular(o.Name()))
	} else if o.IsToOne() {
		join += fmt.Sprintf("\"LEFT JOIN \"+dialect.Quote(TableName(\"%[1]s\"))+\" \"+dialect.Quote(_alias)+\" ON \"+dialect.Quote(_alias)+\".id = \"+dialect.Quote(alias)+\".\"+dialect.Quote(\"%[2]sId\")", o.Target().TableName(), o.Name())
	} else if o.IsToMany() {
		join += fmt.Sprintf("\"LEFT JOIN \"+dialect.Quote(TableName(\"%[1]s\"))+\" \"+dialect.Quote(_alias)+\" ON \"+dialect.Quote(_alias)+\".\"+dialect.Quote(\"%[3]sId\")+\" = \"+dialect.Quote(alias)+\".id\"", o.Target().TableName(), o.Name(), o.InverseRelationshipName())
	}
	return join
}

func (o *ObjectRelationship) ForeignKeyDestinationColumn() string {
	if o.IsToOne() {
		return "id"
	}
	if o.IsManyToMany() {
		return inflection.Singular(o.InverseRelationshipName()) + "_id"
	}
	return ""
}
func (o *ObjectRelationship) OnDelete(def string) string {
	str, exists := o.StringForRelationshipDirectiveAttribute("onDelete")
	if !exists {
		return def
	}
	return str
}
func (o *ObjectRelationship) OnUpdate(def string) string {
	str, exists := o.StringForRelationshipDirectiveAttribute("onUpdate")
	if !exists {
		return def
	}
	return str
}
