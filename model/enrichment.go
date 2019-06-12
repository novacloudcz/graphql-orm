package model

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/kinds"

	"github.com/graphql-go/graphql/language/ast"
)

// https://github.com/99designs/gqlgen/issues/681 for nested fields
// graphql.CollectFieldsCtx()

// EnrichModel ...
func EnrichModel(m *Model) error {

	id := fieldDefinition("id", "ID", true)
	createdAt := fieldDefinition("createdAt", "Time", true)
	updatedAt := fieldDefinition("updatedAt", "Time", true)
	deletedAt := fieldDefinition("deletedAt", "Time", false)

	definitions := []ast.Node{}
	for _, o := range m.Objects() {
		definitions = append(definitions, createObjectDefinition(o), updateObjectDefinition(o))
		o.Def.Fields = append(append([]*ast.FieldDefinition{id}, o.Def.Fields...), updatedAt, createdAt, deletedAt)
		definitions = append(definitions, objectResultTypeDefinition(&o))
	}

	schemaHeaderNodes := []ast.Node{
		scalarDefinition("Time"),
		relationshipDirectiveDefinition(),
		schemaDefinition(m),
		queryDefinition(m),
		mutationDefinition(m),
	}
	m.Doc.Definitions = append(schemaHeaderNodes, m.Doc.Definitions...)
	m.Doc.Definitions = append(m.Doc.Definitions, definitions...)

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
	return &ast.FieldDefinition{
		Name: nameNode(fieldName),
		Kind: kinds.FieldDefinition,
		Type: t,
	}
}

func relationshipDirectiveDefinition() *ast.DirectiveDefinition {
	return &ast.DirectiveDefinition{
		Kind: kinds.DirectiveDefinition,
		Name: nameNode("relationship"),
		Arguments: []*ast.InputValueDefinition{
			&ast.InputValueDefinition{
				Kind: kinds.InputValueDefinition,
				Name: nameNode("inverse"),
				Type: nonNull(namedType("String")),
			},
		},
		Locations: []*ast.Name{
			nameNode(graphql.DirectiveLocationFieldDefinition),
		},
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
