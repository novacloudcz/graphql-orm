package model

import (
	"github.com/graphql-go/graphql/language/kinds"

	"github.com/graphql-go/graphql/language/ast"
)

func createFederationServiceObject() *ast.ObjectDefinition {
	return &ast.ObjectDefinition{
		Kind: kinds.ObjectDefinition,
		Name: nameNode("_Service"),
		Fields: []*ast.FieldDefinition{
			&ast.FieldDefinition{
				Kind: kinds.FieldDefinition,
				Name: nameNode("sdl"),
				Type: namedType("String"),
			},
		},
	}
}

func createFederationServiceQueryField() *ast.FieldDefinition {
	return &ast.FieldDefinition{
		Kind: kinds.FieldDefinition,
		Name: nameNode("_service"),
		Type: nonNull(namedType("_Service")),
	}
}

func createFederationEntityUnion(m *Model) *ast.UnionDefinition {
	types := []*ast.Named{}

	for _, o := range m.Objects() {
		if o.HasDirective("key") {
			t := namedType(o.Name())
			types = append(types, t.(*ast.Named))
		}
	}

	return &ast.UnionDefinition{
		Kind:  kinds.UnionDefinition,
		Name:  nameNode("_Entity"),
		Types: types,
	}
}
func createFederationEntitiesQueryField() *ast.FieldDefinition {
	return &ast.FieldDefinition{
		Kind: kinds.FieldDefinition,
		Name: nameNode("_entities"),
		Type: nonNull(listType(namedType("_Entity"))),
		Arguments: []*ast.InputValueDefinition{
			&ast.InputValueDefinition{
				Kind: kinds.InputValueDefinition,
				Name: nameNode("representations"),
				Type: nonNull(listType(nonNull(namedType("_Any")))),
			},
		},
	}
}

func getObjectDefinitionFromFederationExtension(def *ast.TypeExtensionDefinition) *ast.ObjectDefinition {
	federationDirectives := []string{"requires", "provides", "key", "extends", "external"}
	objDef := def.Definition
	for _, dir := range federationDirectives {
		objDef.Directives = filterDirective(objDef.Directives, dir)
	}
	for _, field := range objDef.Fields {
		for _, dir := range federationDirectives {
			field.Directives = filterDirective(field.Directives, dir)
		}
	}
	return objDef
}

// func getObjectResolverReferenceField(o *Object) *ast.FieldDefinition {
// 	d := o.Directive("key")
// 	var fields *ast.Argument

// 	for _, arg := range d.Arguments {
// 		if arg.Name.Value == "fields" {
// 			fields = arg
// 		}
// 	}

// 	fieldsString := fields.Value.GetValue().(string)

// 	args := []*ast.InputValueDefinition{}
// 	for _, field := range strings.Split(fieldsString, ",") {
// 		args = append(args, &ast.InputValueDefinition{
// 			Kind: kinds.InputValueDefinition,
// 			Name: nameNode(field),
// 			Type: o.Column(field).Def.Type,
// 		})
// 	}
// 	return &ast.FieldDefinition{
// 		Kind:      kinds.FieldDefinition,
// 		Name:      nameNode("resolveReference"),
// 		Type:      namedType(o.Name()),
// 		Arguments: args,
// 	}
// }
