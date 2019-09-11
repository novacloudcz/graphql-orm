package model

import (
	"github.com/graphql-go/graphql/language/kinds"

	"github.com/graphql-go/graphql/language/ast"
)

// https://github.com/99designs/gqlgen/issues/681 for nested fields
// graphql.CollectFieldsCtx()

// EnrichModelObjects ...
func EnrichModelObjects(m *Model) error {
	id := fieldDefinition("id", "ID", true)
	createdAt := fieldDefinition("createdAt", "Time", true)
	updatedAt := fieldDefinition("updatedAt", "Time", false)
	createdBy := fieldDefinition("createdBy", "ID", false)
	updatedBy := fieldDefinition("updatedBy", "ID", false)

	for _, o := range m.Objects() {
		o.Def.Fields = append([]*ast.FieldDefinition{id}, o.Def.Fields...)
		for _, rel := range o.Relationships() {
			if rel.IsToOne() {
				o.Def.Fields = append(o.Def.Fields, fieldDefinition(rel.Name()+"Id", "ID", false))
			}
		}
		o.Def.Fields = append(o.Def.Fields, updatedAt, createdAt, updatedBy, createdBy)
	}
	return nil
}

// EnrichModel ...
func EnrichModel(m *Model) error {
	if m.HasFederatedTypes() {
		m.Doc.Definitions = append(m.Doc.Definitions, createFederationEntityUnion(m))
	}

	definitions := []ast.Node{}
	for _, o := range m.Objects() {
		for _, rel := range o.Relationships() {
			if rel.IsToMany() {
				o.Def.Fields = append(o.Def.Fields, fieldDefinitionWithType(rel.Name()+"Ids", nonNull(listType(nonNull(namedType("ID"))))))
			}
		}
		definitions = append(definitions, createObjectDefinition(o), updateObjectDefinition(o), createObjectSortType(o), createObjectFilterType(o))
		definitions = append(definitions, objectResultTypeDefinition(&o))
	}

	schemaHeaderNodes := []ast.Node{
		scalarDefinition("Time"),
		scalarDefinition("_Any"),
		schemaDefinition(m),
		queryDefinition(m),
		mutationDefinition(m),
	}
	m.Doc.Definitions = append(schemaHeaderNodes, m.Doc.Definitions...)
	m.Doc.Definitions = append(m.Doc.Definitions, definitions...)
	m.Doc.Definitions = append(m.Doc.Definitions, createFederationServiceObject())

	return nil
}

func BuildFederatedModel(m *Model) error {

	for _, e := range m.ObjectExtensions() {
		if e.IsFederatedType() {
			m.Doc.Definitions = append(m.Doc.Definitions, getObjectDefinitionFromFederationExtension(e.Object.Def))
			m.RemoveObjectExtension(&e)
		}
	}

	for _, obj := range m.Objects() {
		if obj.HasDirective("key") {
			obj.Def.Directives = filterDirective(obj.Def.Directives, "key")
		}
	}

	return nil
}

func scalarDefinition(name string) *ast.ScalarDefinition {
	return &ast.ScalarDefinition{
		Name: &ast.Name{
			Kind:  kinds.Name,
			Value: name,
		},
		Kind: "ScalarDefinition",
	}
}

func fieldDefinition(fieldName, fieldType string, isNonNull bool) *ast.FieldDefinition {
	t := namedType(fieldType)
	if isNonNull {
		t = nonNull(t)
	}
	return fieldDefinitionWithType(fieldName, t)
}
func fieldDefinitionWithType(fieldName string, t ast.Type) *ast.FieldDefinition {
	return &ast.FieldDefinition{
		Name: nameNode(fieldName),
		Kind: kinds.FieldDefinition,
		Type: t,
	}
}

func schemaDefinition(m *Model) *ast.SchemaDefinition {
	return &ast.SchemaDefinition{
		Kind: kinds.SchemaDefinition,
		OperationTypes: []*ast.OperationTypeDefinition{
			&ast.OperationTypeDefinition{
				Operation: "query",
				Kind:      kinds.OperationTypeDefinition,
				Type: &ast.Named{
					Kind: kinds.Named,
					Name: &ast.Name{
						Kind:  kinds.Name,
						Value: "Query",
					},
				},
			},
			&ast.OperationTypeDefinition{
				Operation: "mutation",
				Kind:      kinds.OperationTypeDefinition,
				Type: &ast.Named{
					Kind: kinds.Named,
					Name: &ast.Name{
						Kind:  kinds.Name,
						Value: "Mutation",
					},
				},
			},
		},
	}
}
