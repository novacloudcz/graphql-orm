package templates

var Resolver = `package gen

import (
	"context"
	"time"
	
	"github.com/novacloudcz/graphql-orm/resolvers"
	uuid "github.com/satori/go.uuid"
)

func getPrincipalID(ctx context.Context) string {
	v, _ := ctx.Value(KeyPrincipalID).(string)
	return v
}

type Resolver struct {
	DB *DB
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
{{range .Model.Objects}}
func (r *Resolver) {{.Name}}ResultType() {{.Name}}ResultTypeResolver {
	return &{{.LowerName}}ResultTypeResolver{r}
}
{{if .HasRelationships}}
func (r *Resolver) {{.Name}}() {{.Name}}Resolver {
	return &{{.LowerName}}Resolver{r}
}
{{end}}
{{end}}

type mutationResolver struct{ *Resolver }

{{range .Model.Objects}}
func (r *mutationResolver) Create{{.Name}}(ctx context.Context, input map[string]interface{}) (item *{{.Name}}, err error) {
	principalID := getPrincipalID(ctx)
	item = &{{.Name}}{ID: uuid.Must(uuid.NewV4()).String(), CreatedBy: principalID}
	tx := r.DB.db.Begin()

{{range $col := .Columns}}{{if $col.IsCreatable}}
	if val, ok := input["{{$col.Name}}"].({{$col.GoTypeWithPointer false}}); ok && ({{if $col.IsOptional}}item.{{$col.MethodName}} == nil || *{{end}}item.{{$col.MethodName}} != val) {
		item.{{$col.MethodName}} = {{if $col.IsOptional}}&{{end}}val
	}
{{end}}
{{end}}

{{range $rel := .Relationships}}
{{if $rel.IsToMany}}
	if ids,ok:=input["{{$rel.Name}}Ids"].([]interface{}); ok {
		items := []{{$rel.TargetType}}{}
		tx.Find(&items, "id IN (?)", ids)
		association := tx.Model(&item).Association("{{$rel.MethodName}}")
		association.Replace(items)
	}
{{end}}
{{end}}

	err = tx.Create(item).Error
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit().Error
	return 
}
func (r *mutationResolver) Update{{.Name}}(ctx context.Context, id string, input map[string]interface{}) (item *{{.Name}}, err error) {
	item = &{{.Name}}{}
	tx := r.DB.db.Begin()
	
	err = resolvers.GetItem(ctx, tx, item, &id)
	if err != nil {
		return 
	}

	principalID := getPrincipalID(ctx)
	item.UpdatedBy = &principalID

{{range $col := .Columns}}{{if $col.IsUpdatable}}
	if val, ok := input["{{$col.Name}}"].({{$col.GoTypeWithPointer false}}); ok && ({{if $col.IsOptional}}item.{{$col.MethodName}} == nil || *{{end}}item.{{$col.MethodName}} != val) {
		item.{{$col.MethodName}} = {{if $col.IsOptional}}&{{end}}val
	}
{{end}}
{{end}}

{{range $rel := .Relationships}}
{{if $rel.IsToMany}}
	if ids,ok:=input["{{$rel.Name}}Ids"].([]interface{}); ok {
		items := []{{$rel.TargetType}}{}
		tx.Find(&items, "id IN (?)", ids)
		association := tx.Model(&item).Association("{{$rel.MethodName}}")
		association.Replace(items)
	}
{{end}}
{{end}}

	err = tx.Save(item).Error
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit().Error
	return 
}
func (r *mutationResolver) Delete{{.Name}}(ctx context.Context, id string) (item *{{.Name}}, err error) {
	item = &{{.Name}}{}
	err = resolvers.GetItem(ctx, r.DB.Query(), item, &id)
	if err != nil {
		return 
	}

	err = r.DB.Query().Delete(item, "id = ?", id).Error

	return 
}
{{end}}

type queryResolver struct{ *Resolver }

{{range $object := .Model.Objects}}
func (r *queryResolver) {{$object.Name}}(ctx context.Context, id *string, q *string) (*{{$object.Name}}, error) {
	t := {{$object.Name}}{}
	err := resolvers.GetItem(ctx, r.DB.Query(), &t, id)
	return &t, err
}
func (r *queryResolver) {{$object.PluralName}}(ctx context.Context, offset *int, limit *int, q *string,sort []{{$object.Name}}SortType,filter *{{$object.Name}}FilterType) (*{{$object.Name}}ResultType, error) {
	_sort := []resolvers.EntitySort{}
	for _, s := range sort {
		_sort = append(_sort, s)
	}
	query := {{$object.Name}}QueryFilter{q}
	return &{{$object.Name}}ResultType{
		EntityResultType: resolvers.EntityResultType{
			Offset: offset,
			Limit:  limit,
			Query:  &query,
			Sort: _sort,
			Filter: filter,
		},
	}, nil
}

type {{$object.LowerName}}ResultTypeResolver struct{ *Resolver }

func (r *{{$object.LowerName}}ResultTypeResolver) Items(ctx context.Context, obj *{{$object.Name}}ResultType) (items []*{{$object.Name}}, err error) {
	err = obj.GetItems(ctx, r.DB.db, "{{$object.TableName}}", &items)
	return
}

func (r *{{$object.LowerName}}ResultTypeResolver) Count(ctx context.Context, obj *{{$object.Name}}ResultType) (count int, err error) {
	return obj.GetCount(ctx, r.DB.db, &{{$object.Name}}{})
}

{{if .HasRelationships}}
type {{$object.LowerName}}Resolver struct { *Resolver }

{{range $index, $relationship := .Relationships}}
func (r *{{$object.LowerName}}Resolver) {{$relationship.MethodName}}(ctx context.Context, obj *{{$object.Name}}) (res {{.ReturnType}}, err error) {
{{if $relationship.IsToMany}}
	items := []*{{.TargetType}}{}
	err = r.DB.Query().Model(obj).Related(&items, "{{$relationship.MethodName}}").Error
	res = items
{{else}}
	item := {{.TargetType}}{}
	_res := r.DB.Query().Model(obj).Related(&item, "{{$relationship.MethodName}}")
	if _res.RecordNotFound() {
		return
	} else {
		err = _res.Error
	}
	res = &item
{{end}}
	return 
}
{{end}}
{{end}}

{{end}}
`
