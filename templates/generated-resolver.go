package templates

var GeneratedResolver = `package gen

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
	DeleteAll{{.PluralName}} func (ctx context.Context, r *GeneratedMutationResolver) (bool, error) 
	Query{{$obj.Name}} func (ctx context.Context, r *GeneratedQueryResolver, id *string, q *string, filter *{{$obj.Name}}FilterType) (*{{$obj.Name}}, error)
	Query{{$obj.PluralName}} func (ctx context.Context, r *GeneratedQueryResolver, offset *int, limit *int, q *string, sort []{{$obj.Name}}SortType, filter *{{$obj.Name}}FilterType) (*{{$obj.Name}}ResultType, error)
	{{range $col := $obj.Columns}}{{if $col.IsReadonlyType}}
	{{$obj.Name}}{{$col.MethodName}} func (ctx context.Context,r *Generated{{$obj.Name}}Resolver, obj *{{$obj.Name}}) (res {{$col.GoType}}, err error)
	{{end}}{{end}}
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

{{range .Model.Objects}}
func (r *GeneratedResolver) {{.Name}}ResultType() {{.Name}}ResultTypeResolver {
	return &Generated{{.Name}}ResultTypeResolver{r}
}
{{if .HasAnyRelationships}}
func (r *GeneratedResolver) {{.Name}}() {{.Name}}Resolver {
	return &Generated{{.Name}}Resolver{r}
}
{{end}}
{{end}}

type GeneratedMutationResolver struct{ *GeneratedResolver }

{{range .Model.Objects}}
func (r *GeneratedMutationResolver) Create{{.Name}}(ctx context.Context, input map[string]interface{}) (item *{{.Name}}, err error) {
	return r.Handlers.Create{{.Name}}(ctx, r, input)
}
func Create{{.Name}}Handler(ctx context.Context, r *GeneratedMutationResolver, input map[string]interface{}) (item *{{.Name}}, err error) {
	principalID := getPrincipalIDFromContext(ctx)
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
	return r.Handlers.Update{{.Name}}(ctx, r, id, input)
}
func Update{{.Name}}Handler(ctx context.Context, r *GeneratedMutationResolver, id string, input map[string]interface{}) (item *{{.Name}}, err error) {
	principalID := getPrincipalIDFromContext(ctx)
	item = &{{.Name}}{}
	now := time.Now()
	tx := r.DB.db.Begin()

	event := events.NewEvent(events.EventMetadata{
		Type:        events.EventTypeUpdated,
		Entity:      "{{.Name}}",
		EntityID:    id,
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
		// data, _ := json.Marshal(event)
		// fmt.Println("?",string(data))
	}

	return 
}
func (r *GeneratedMutationResolver) Delete{{.Name}}(ctx context.Context, id string) (item *{{.Name}}, err error) {
	return r.Handlers.Delete{{.Name}}(ctx, r, id)
}
func Delete{{.Name}}Handler(ctx context.Context, r *GeneratedMutationResolver, id string) (item *{{.Name}}, err error) {
	principalID := getPrincipalIDFromContext(ctx)
	item = &{{.Name}}{}
	now := time.Now()
	tx := r.DB.db.Begin()

	err = resolvers.GetItem(ctx, tx, item, &id)
	if err != nil {
		return 
	}

	event := events.NewEvent(events.EventMetadata{
		Type:        events.EventTypeDeleted,
		Entity:      "{{.Name}}",
		EntityID:    id,
		Date:        now,
		PrincipalID: principalID,
	})

	err = tx.Delete(item, "{{.TableName}}.id = ?", id).Error
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return
	}

	err = r.EventController.SendEvent(ctx, &event)
	
	return 
}
func (r *GeneratedMutationResolver) DeleteAll{{.PluralName}}(ctx context.Context) (bool, error) {
	return r.Handlers.DeleteAll{{.PluralName}}(ctx, r)
}
func DeleteAll{{.PluralName}}Handler(ctx context.Context, r *GeneratedMutationResolver) (bool, error) {
	err := r.DB.db.Delete(&{{.Name}}{}).Error
	return err == nil, err
}
{{end}}

type GeneratedQueryResolver struct{ *GeneratedResolver }


func (r *GeneratedQueryResolver) _service(ctx context.Context) (*_Service, error) {
	sdl := SchemaSDL
	return &_Service{
		Sdl: &sdl,
	}, nil
}

{{if .Model.HasFederatedTypes}}
func (r *GeneratedQueryResolver) _entities(ctx context.Context, representations []interface{}) (res []_Entity, err error) {
	res = []_Entity{}
	for _, repr := range representations {
		anyValue, ok := repr.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("The _entities resolver received invalid representation type")
			break
		}
		typename, ok := anyValue["__typename"].(string)
		if !ok {
			err = fmt.Errorf("The _entities resolver received invalid representation type (missing __typename field)")
			break
		}
		id, identifierExists := anyValue["id"].(string)
		
		switch typename { {{range $obj := .Model.Objects}}{{if $obj.HasDirective "key"}}
		case "{{$obj.Name}}":
			if identifierExists {
				item, qerr := r.Handlers.Query{{$obj.Name}}(ctx, r, &id, nil, nil)
				err = qerr
				if err != nil {
					return
				}
				res = append(res, item)
			} else {
				res = append(res, &{{$obj.Name}}{})
			}
			break;{{end}}{{end}}
		default:
			err = fmt.Errorf("The _entities resolver tried to load an entity for type \"%s\", but no object type of that name was found in the schema", typename)
			return
		}
	}
	return res, err
}
{{end}}

{{range $object := .Model.Objects}}
func (r *GeneratedQueryResolver) {{$object.Name}}(ctx context.Context, id *string, q *string, filter *{{$object.Name}}FilterType) (*{{$object.Name}}, error) {
	return r.Handlers.Query{{$object.Name}}(ctx, r, id, q, filter)
}
func Query{{$object.Name}}Handler(ctx context.Context, r *GeneratedQueryResolver, id *string, q *string, filter *{{$object.Name}}FilterType) (*{{$object.Name}}, error) {
	query := {{$object.Name}}QueryFilter{q}
	offset := 0
	limit := 1
	rt := &{{$object.Name}}ResultType{
		EntityResultType: resolvers.EntityResultType{
			Offset: &offset,
			Limit:  &limit,
			Query:  &query,
			Filter: filter,
		},
	}
	qb := r.DB.Query()
	if id != nil {
		qb = qb.Where("{{$object.TableName}}.id = ?", *id)
	}

	var items []*{{$object.Name}}
	err := rt.GetItems(ctx, qb, "{{$object.TableName}}", &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("{{$object.Name}} not found")
	}
	return items[0], err
}
func (r *GeneratedQueryResolver) {{$object.PluralName}}(ctx context.Context, offset *int, limit *int, q *string, sort []{{$object.Name}}SortType, filter *{{$object.Name}}FilterType) (*{{$object.Name}}ResultType, error) {
	return r.Handlers.Query{{$object.PluralName}}(ctx, r, offset, limit, q , sort, filter)
}
func Query{{$object.PluralName}}Handler(ctx context.Context, r *GeneratedQueryResolver, offset *int, limit *int, q *string, sort []{{$object.Name}}SortType, filter *{{$object.Name}}FilterType) (*{{$object.Name}}ResultType, error) {
	_sort := []resolvers.EntitySort{}
	for _, s := range sort {
		_sort = append(_sort, s)
	}
	query := {{$object.Name}}QueryFilter{q}
	
	var selectionSet *ast.SelectionSet
	for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
		if f.Field.Name == "items" {
			selectionSet = &f.Field.SelectionSet
		}
	}
	
	return &{{$object.Name}}ResultType{
		EntityResultType: resolvers.EntityResultType{
			Offset: offset,
			Limit:  limit,
			Query:  &query,
			Sort: _sort,
			Filter: filter,
			SelectionSet: selectionSet,
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

{{if .HasAnyRelationships}}
type Generated{{$object.Name}}Resolver struct { *GeneratedResolver }

{{range $col := .Columns}}
{{if $col.IsReadonlyType}}
func (r *Generated{{$object.Name}}Resolver) {{$col.MethodName}}(ctx context.Context, obj *{{$object.Name}}) (res {{$col.GoType}}, err error) {
	return r.Handlers.{{$object.Name}}{{$col.MethodName}}(ctx, r, obj)
}
func {{$object.Name}}{{$col.MethodName}}Handler(ctx context.Context,r *Generated{{$object.Name}}Resolver, obj *{{$object.Name}}) (res {{$col.GoType}}, err error) {
	return nil, fmt.Errorf("Resolver for {{$object.Name}}.{{$col.MethodName}} not implemented")
}
{{end}}
{{end}}

{{range $index, $relationship := .Relationships}}
func (r *Generated{{$object.Name}}Resolver) {{$relationship.MethodName}}(ctx context.Context, obj *{{$object.Name}}) (res {{.ReturnType}}, err error) {
{{if $relationship.IsToMany}}
	items := []*{{.TargetType}}{}
	err = r.DB.Query().Model(obj).Related(&items, "{{$relationship.MethodName}}").Error
	res = items
{{else}}
	loaders := ctx.Value("loaders").(map[string]*dataloader.Loader)
	if obj.{{$relationship.MethodName}}ID != nil {
		item, _err := loaders["{{$relationship.Target.Name}}"].Load(ctx, dataloader.StringKey(*obj.{{$relationship.MethodName}}ID))()
		res, _ = item.({{.ReturnType}})
		err = _err
	}
{{end}}
	return 
}
{{if $relationship.IsToMany}}

func (r *Generated{{$object.Name}}Resolver) {{$relationship.MethodName}}Ids(ctx context.Context, obj *{{$object.Name}}) (ids []string, err error) {
	ids = []string{}

	items := []*{{$relationship.TargetType}}{}
	err = r.DB.Query().Model(obj).Select("{{$relationship.Target.TableName}}.id").Related(&items, "{{$relationship.MethodName}}").Error

	for _, item := range items {
		ids = append(ids, item.ID)
	}

	return
}
{{end}}

{{end}}

{{end}}

{{end}}
`
