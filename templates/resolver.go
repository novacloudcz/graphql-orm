package templates

var Resolver = `package gen

import (
	"context"
	"time"
	
	"github.com/novacloudcz/graphql-orm/resolvers"
	uuid "github.com/satori/go.uuid"
)

func getPrincipalID(ctx context.Context) *string {
	v, _ := ctx.Value(KeyPrincipalID).(*string)
	return v
}

type GeneratedResolver struct {
	DB *DB
	EventController *events.EventController
}

func (r *GeneratedResolver) Mutation() MutationResolver {
	return &GeneratedMutationResolver{r}
}
func (r *GeneratedResolver) Query() QueryResolver {
	return &GeneratedQueryResolver{r}
}
{{range .Model.Objects}}
func (r *GeneratedResolver) {{.Name}}ResultType() {{.Name}}ResultTypeResolver {
	return &Generated{{.Name}}ResultTypeResolver{r}
}
{{if .HasRelationships}}
func (r *GeneratedResolver) {{.Name}}() {{.Name}}Resolver {
	return &Generated{{.Name}}Resolver{r}
}
{{end}}
{{end}}

type GeneratedMutationResolver struct{ *GeneratedResolver }

{{range .Model.Objects}}
func (r *GeneratedMutationResolver) Create{{.Name}}(ctx context.Context, input map[string]interface{}) (item *{{.Name}}, err error) {
	principalID := getPrincipalID(ctx)
	now := time.Now()
	item = &{{.Name}}{ID: uuid.Must(uuid.NewV4()).String(), CreatedAt: now, CreatedBy: principalID}
	tx := r.DB.db.Begin()

	event := events.NewEvent(events.EventMetadata{
		Type:        events.EventTypeCreated,
		Entity:      "{{.Name}}",
		EntityID:    item.ID,
		Date:        now,
		PrincipalID: principalID,
	})

	var changes {{.Name}}Changes
	err = ApplyChanges(input, &changes)
	if err != nil {
		return 
	}

{{range $col := .Columns}}{{if $col.IsCreatable}}
	if _, ok := input["{{$col.Name}}"]; ok && (item.{{$col.MethodName}} != changes.{{$col.MethodName}}){{if $col.IsOptional}} && (item.{{$col.MethodName}} == nil || changes.{{$col.MethodName}} == nil || *item.{{$col.MethodName}} != *changes.{{$col.MethodName}}){{end}} {
		item.{{$col.MethodName}} = changes.{{$col.MethodName}}
		event.AddNewValue("{{$col.Name}}", changes.{{$col.MethodName}})
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
	if err != nil {
		tx.Rollback()
		return
	}

	if len(event.Changes) > 0 {
		err = r.EventController.SendEvent(ctx, &event)
	}

	return 
}
func (r *GeneratedMutationResolver) Update{{.Name}}(ctx context.Context, id string, input map[string]interface{}) (item *{{.Name}}, err error) {
	principalID := getPrincipalID(ctx)
	item = &{{.Name}}{}
	now := time.Now()
	tx := r.DB.db.Begin()

	event := events.NewEvent(events.EventMetadata{
		Type:        events.EventTypeCreated,
		Entity:      "{{.Name}}",
		EntityID:    item.ID,
		Date:        now,
		PrincipalID: principalID,
	})

	var changes {{.Name}}Changes
	err = ApplyChanges(input, &changes)
	if err != nil {
		return 
	}

	err = resolvers.GetItem(ctx, tx, item, &id)
	if err != nil {
		return 
	}

	item.UpdatedBy = principalID

{{range $col := .Columns}}{{if $col.IsUpdatable}}
	if _, ok := input["{{$col.Name}}"]; ok && (item.{{$col.MethodName}} != changes.{{$col.MethodName}}){{if $col.IsOptional}} && (item.{{$col.MethodName}} == nil || changes.{{$col.MethodName}} == nil || *item.{{$col.MethodName}} != *changes.{{$col.MethodName}}){{end}} {
		event.AddOldValue("{{$col.Name}}", item.{{$col.MethodName}})
		event.AddNewValue("{{$col.Name}}", changes.{{$col.MethodName}})
		item.{{$col.MethodName}} = changes.{{$col.MethodName}}
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
	if err != nil {
		tx.Rollback()
		return
	}

	if len(event.Changes) > 0 {
		err = r.EventController.SendEvent(ctx, &event)
		data, _ := json.Marshal(event)
		fmt.Println("??", string(data))
	}

	return 
}
func (r *GeneratedMutationResolver) Delete{{.Name}}(ctx context.Context, id string) (item *{{.Name}}, err error) {
	item = &{{.Name}}{}
	err = resolvers.GetItem(ctx, r.DB.Query(), item, &id)
	if err != nil {
		return 
	}

	err = r.DB.Query().Delete(item, "id = ?", id).Error

	return 
}
{{end}}

type GeneratedQueryResolver struct{ *GeneratedResolver }

{{range $object := .Model.Objects}}
func (r *GeneratedQueryResolver) {{$object.Name}}(ctx context.Context, id *string, q *string) (*{{$object.Name}}, error) {
	t := {{$object.Name}}{}
	err := resolvers.GetItem(ctx, r.DB.Query(), &t, id)
	return &t, err
}
func (r *GeneratedQueryResolver) {{$object.PluralName}}(ctx context.Context, offset *int, limit *int, q *string,sort []{{$object.Name}}SortType,filter *{{$object.Name}}FilterType) (*{{$object.Name}}ResultType, error) {
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

type Generated{{$object.Name}}ResultTypeResolver struct{ *GeneratedResolver }

func (r *Generated{{$object.Name}}ResultTypeResolver) Items(ctx context.Context, obj *{{$object.Name}}ResultType) (items []*{{$object.Name}}, err error) {
	err = obj.GetItems(ctx, r.DB.db, "{{$object.TableName}}", &items)
	return
}

func (r *Generated{{$object.Name}}ResultTypeResolver) Count(ctx context.Context, obj *{{$object.Name}}ResultType) (count int, err error) {
	return obj.GetCount(ctx, r.DB.db, &{{$object.Name}}{})
}

{{if .HasRelationships}}
type Generated{{$object.Name}}Resolver struct { *GeneratedResolver }

{{range $index, $relationship := .Relationships}}
func (r *Generated{{$object.Name}}Resolver) {{$relationship.MethodName}}(ctx context.Context, obj *{{$object.Name}}) (res {{.ReturnType}}, err error) {
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
