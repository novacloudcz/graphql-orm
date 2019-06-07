package model

import (
	"github.com/graphql-go/graphql/language/kinds"

	"github.com/graphql-go/graphql/language/ast"
)

func namedType(name string) ast.Type {
	var t ast.Type
	t = &ast.Named{
		Kind: kinds.Named,
		Name: &ast.Name{Kind: kinds.Name, Value: name},
	}
	return t
}

func getNamedType(t ast.Type) ast.Type {
	if t.GetKind() == kinds.Named {
		return t
	}
	switch t.GetKind() {
	case kinds.List:
		return getNamedType(t.(*ast.List).Type)
	case kinds.NonNull:
		return getNamedType(t.(*ast.NonNull).Type)
	}
	panic("unable to get named type of " + t.String())
}

func nonNull(t ast.Type) ast.Type {
	if t.GetKind() == kinds.NonNull {
		return t
	}
	return &ast.NonNull{
		Kind: kinds.NonNull,
		Type: t,
	}
}

func nameNode(name string) *ast.Name {
	return &ast.Name{
		Kind:  kinds.Name,
		Value: name,
	}
}
