package templates

var ResolverCore = `package gen

import (
	"context"
	"time"
	
	"github.com/graph-gophers/dataloader"
	"github.com/99designs/gqlgen/graphql"
	"github.com/gofrs/uuid"
	"github.com/novacloudcz/graphql-orm/events"
	"github.com/novacloudcz/graphql-orm/resolvers"
	"github.com/vektah/gqlparser/ast"
)

type resolutionHandlers struct {
	{{range $obj := .Model.Objects}}
		Create{{$obj.Name}} func (ctx context.Context, r *GeneratedMutationResolver, input map[string]interface{}) (item *{{$obj.Name}}, err error)
		Update{{$obj.Name}} func(ctx context.Context, r *GeneratedMutationResolver, id string, input map[string]interface{}) (item *{{$obj.Name}}, err error)
		Delete{{$obj.Name}} func(ctx context.Context, r *GeneratedMutationResolver, id string) (item *{{$obj.Name}}, err error)
		DeleteAll{{$obj.PluralName}} func (ctx context.Context, r *GeneratedMutationResolver) (bool, error) 
		Query{{$obj.Name}} func (ctx context.Context, r *GeneratedQueryResolver, id *string, q *string, filter *{{$obj.Name}}FilterType) (*{{$obj.Name}}, error)
		Query{{$obj.PluralName}} func (ctx context.Context, r *GeneratedQueryResolver, offset *int, limit *int, q *string, sort []{{$obj.Name}}SortType, filter *{{$obj.Name}}FilterType) (*{{$obj.Name}}ResultType, error)
		{{range $col := $obj.Columns}}{{if $col.IsReadonlyType}}
			{{$obj.Name}}{{$col.MethodName}} func (ctx context.Context,r *Generated{{$obj.Name}}Resolver, obj *{{$obj.Name}}) (res {{$col.GoType}}, err error)
		{{end}}{{end}}
		{{range $rel := $obj.Relationships}}
			{{$obj.Name}}{{$rel.MethodName}} func (ctx context.Context,r *Generated{{$obj.Name}}Resolver, obj *{{$obj.Name}}) (res {{$rel.ReturnType}}, err error)
		{{end}}
	{{end}}
	{{range $ext := .Model.ObjectExtensions}}
		{{$obj := $ext.Object}}
		{{range $col := $obj.Columns}}{{if $col.IsReadonlyType}}
			{{$obj.Name}}{{$col.MethodName}} func (ctx context.Context,r *Generated{{$obj.Name}}Resolver, obj *{{$obj.Name}}) (res {{$col.GoType}}, err error)
		{{end}}{{end}}
		{{range $rel := $obj.Relationships}}
			{{$obj.Name}}{{$rel.MethodName}} func (ctx context.Context,r *Generated{{$obj.Name}}Resolver, obj *{{$obj.Name}}) (res {{$rel.ReturnType}}, err error)
		{{end}}
	{{end}}
}

func NewResolver(db *DB, ec *events.EventController) *GeneratedResolver {
	handlers := resolutionHandlers{
		{{range $obj := .Model.Objects}}
			Create{{$obj.Name}}: Create{{$obj.Name}}Handler,
			Update{{$obj.Name}}: Update{{$obj.Name}}Handler,
			Delete{{$obj.Name}}: Delete{{$obj.Name}}Handler,
			DeleteAll{{$obj.PluralName}}: DeleteAll{{$obj.PluralName}}Handler,
			Query{{$obj.Name}}: Query{{$obj.Name}}Handler,
			Query{{$obj.PluralName}}: Query{{$obj.PluralName}}Handler,
			{{range $col := $obj.Columns}}{{if $col.IsReadonlyType}}
				{{$obj.Name}}{{$col.MethodName}}: {{$obj.Name}}{{$col.MethodName}}Handler,
			{{end}}{{end}}
			{{range $rel := $obj.Relationships}}
				{{$obj.Name}}{{$rel.MethodName}}: {{$obj.Name}}{{$rel.MethodName}}Handler,
			{{end}}
		{{end}}
		{{range $ext := .Model.ObjectExtensions}}
			{{$obj := $ext.Object}}
			{{range $col := $obj.Columns}}{{if $col.IsReadonlyType}}
				{{$obj.Name}}{{$col.MethodName}}: {{$obj.Name}}{{$col.MethodName}}Handler,
			{{end}}{{end}}
			{{range $rel := $obj.Relationships}}
				{{$obj.Name}}{{$rel.MethodName}}: {{$obj.Name}}{{$rel.MethodName}}Handler,
			{{end}}
		{{end}}
	}
	return &GeneratedResolver{Handlers: handlers, DB: db, EventController: ec}
}

type GeneratedResolver struct {
	Handlers resolutionHandlers
	DB *DB
	EventController *events.EventController
}

func (r *GeneratedResolver) Mutation() MutationResolver {
	return &GeneratedMutationResolver{r}
}
func (r *GeneratedResolver) Query() QueryResolver {
	return &GeneratedQueryResolver{r}
}

{{range $obj := .Model.Objects}}
	func (r *GeneratedResolver) {{$obj.Name}}ResultType() {{$obj.Name}}ResultTypeResolver {
		return &Generated{{$obj.Name}}ResultTypeResolver{r}
	}
	{{if or $obj.HasAnyRelationships $obj.HasReadonlyColumns}}
		func (r *GeneratedResolver) {{$obj.Name}}() {{$obj.Name}}Resolver {
			return &Generated{{$obj.Name}}Resolver{r}
		}
	{{end}}
{{end}}
{{range $ext := .Model.ObjectExtensions}}
	{{$obj := $ext.Object}}
	{{if or $obj.HasAnyRelationships $obj.HasReadonlyColumns}}
		func (r *GeneratedResolver) {{$obj.Name}}() {{$obj.Name}}Resolver {
			return &Generated{{$obj.Name}}Resolver{r}
		}
	{{end}}
{{end}}
`
