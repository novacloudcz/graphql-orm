package templates

var Resolver = `package gen

import (
	"context"
	
	"github.com/novacloudcz/graphql-orm/resolvers"
	uuid "github.com/satori/go.uuid"
)

type Resolver struct {
	DB *DB
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
{{range .Objects}}
func (r *Resolver) {{.Name}}ResultType() {{.Name}}ResultTypeResolver {
	return &{{.LowerName}}ResultTypeResolver{r}
}
{{end}}

type mutationResolver struct{ *Resolver }

{{range .Objects}}
func (r *mutationResolver) Create{{.Name}}(ctx context.Context, input  map[string]interface{}) (item *{{.Name}}, err error) {
	item = &{{.Name}}{ID:uuid.Must(uuid.NewV4()).String()}
	err = resolvers.CreateItem(ctx, r.DB.db, item, input)
	return 
}
func (r *mutationResolver) Update{{.Name}}(ctx context.Context, id string, input  map[string]interface{}) (item *{{.Name}}, err error) {
	item = &{{.Name}}{}
	err = resolvers.GetItem(ctx, r.DB.Query(), item, &id)
	if err != nil {
		return 
	}

	err = resolvers.UpdateItem(ctx, r.DB.Query(), item, input)

	return 
}
func (r *mutationResolver) Delete{{.Name}}(ctx context.Context, id string) (item *{{.Name}}, err error) {
	item = &{{.Name}}{}
	err = resolvers.GetItem(ctx, r.DB.Query(), item, &id)
	if err != nil {
		return 
	}

	err = resolvers.DeleteItem(ctx, r.DB.Query(), item, id)

	return 
}
{{end}}

type queryResolver struct{ *Resolver }

{{range .Objects}}
func (r *queryResolver) {{.Name}}(ctx context.Context, id *string, q *string) (*{{.Name}}, error) {
	t := {{.Name}}{}
	err := resolvers.GetItem(ctx, r.DB.Query(), &t, id)
	return &t, err
}
func (r *queryResolver) {{.Name}}s(ctx context.Context, offset *int, limit *int, q *string) (*{{.Name}}ResultType, error) {
	return &{{.Name}}ResultType{}, nil
}

type {{.LowerName}}ResultTypeResolver struct{ *Resolver }

func (r *{{.LowerName}}ResultTypeResolver) Items(ctx context.Context, obj *{{.Name}}ResultType) (items []{{.Name}}, err error) {
	err = resolvers.GetResultTypeItems(ctx, r.DB.db, &items)
	return
}

func (r *{{.LowerName}}ResultTypeResolver) Count(ctx context.Context, obj *{{.Name}}ResultType) (count int, err error) {
	return resolvers.GetResultTypeCount(ctx, r.DB.db, &{{.Name}}{})
}

{{end}}
`
