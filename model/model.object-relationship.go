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
func (o *ObjectRelationship) InverseRelationshipName() string {
	for _, d := range o.Def.Directives {
		if d.Name.Value == "relationship" {
			for _, arg := range d.Arguments {
				if arg.Name.Value == "inverse" {
					v, ok := arg.Value.GetValue().(string)
					if !ok {
						panic(fmt.Sprintf("invalid inverse value for %s->%s relationship", o.Obj.Name(), o.Name()))
					}
					return v
				}
			}
		}
	}
	panic(fmt.Sprintf("missing relationship directive/inverse argument for %s->%s relationship", o.Obj.Name(), o.Name()))
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
		} else if o.IsMainRelationshipForManyToMany() {
			tags += fmt.Sprintf(" gorm:\"many2many:%s;jointable_foreignkey:%s_id;association_jointable_foreignkey:%s_id\"", rel.ManyToManyJoinTable(), inflection.Singular(o.Name()), inflection.Singular(o.InverseRelationshipName()))
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
		join += fmt.Sprintf("\"LEFT JOIN \"+dialect.Quote(\"%[1]s\")+\" \"+dialect.Quote(_alias)+\"_jointable ON \"+dialect.Quote(alias)+\".id = \"+dialect.Quote(_alias+\"_jointable\")+\".\"+dialect.Quote(\"%[3]s_id\")+\" LEFT JOIN \"+dialect.Quote(\"%[2]s\")+\" \"+dialect.Quote(_alias)+\" ON \"+dialect.Quote(_alias+\"_jointable\")+\".\"+dialect.Quote(\"%[4]s_id\")+\" = \"+dialect.Quote(_alias)+\".id\"", joinTable, o.Target().TableName(), inflection.Singular(o.InverseRelationshipName()), inflection.Singular(o.Name()))
	} else if o.IsToOne() {
		join += fmt.Sprintf("\"LEFT JOIN \"+dialect.Quote(\"%[1]s\")+\" \"+dialect.Quote(_alias)+\" ON \"+dialect.Quote(_alias)+\".id = \"+alias+\".\"+dialect.Quote(\"%[2]sId\")", o.Target().TableName(), o.Name())
	} else if o.IsToMany() {
		join += fmt.Sprintf("\"LEFT JOIN \"+dialect.Quote(\"%[1]s\")+\" \"+dialect.Quote(_alias)+\" ON \"+dialect.Quote(_alias)+\".\"+dialect.Quote(\"%[3]sId\")+\" = \"+dialect.Quote(alias)+\".id\"", o.Target().TableName(), o.Name(), o.InverseRelationshipName())
	}
	return join
}
