package templates

var GeneratedResolver = `package gen

import (
	"context"
	"time"
	
	"github.com/graph-gophers/dataloader"
	"github.com/99designs/gqlgen/graphql"
	"github.com/gofrs/uuid"
	"github.com/novacloudcz/graphql-orm/events"
	"github.com/vektah/gqlparser/ast"
)

type resolutionHandlers struct {
	{{range $obj := .Model.ObjectEntities}}
	{{if not $obj.IsExtended}}
	Create{{$obj.Name}} func (ctx context.Context, r *GeneratedMutationResolver, input map[string]interface{}) (item *{{$obj.Name}}, err error)
	Update{{$obj.Name}} func(ctx context.Context, r *GeneratedMutationResolver, id string, input map[string]interface{}) (item *{{$obj.Name}}, err error)
	Delete{{$obj.Name}} func(ctx context.Context, r *GeneratedMutationResolver, id string) (item *{{$obj.Name}}, err error)
	DeleteAll{{$obj.PluralName}} func (ctx context.Context, r *GeneratedMutationResolver) (bool, error) 
	Query{{$obj.Name}} func (ctx context.Context, r *GeneratedQueryResolver, id *string, q *string, filter *{{$obj.Name}}FilterType) (*{{$obj.Name}}, error)
	Query{{$obj.PluralName}} func (ctx context.Context, r *GeneratedQueryResolver, offset *int, limit *int, q *string, sort []{{$obj.Name}}SortType, filter *{{$obj.Name}}FilterType) (*{{$obj.Name}}ResultType, error)
	{{end}}
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
		{{range $obj := .Model.ObjectEntities}}
		{{if not $obj.IsExtended}}
		Create{{$obj.Name}}: Create{{$obj.Name}}Handler,
		Update{{$obj.Name}}: Update{{$obj.Name}}Handler,
		Delete{{$obj.Name}}: Delete{{$obj.Name}}Handler,
		DeleteAll{{$obj.PluralName}}: DeleteAll{{$obj.PluralName}}Handler,
		Query{{$obj.Name}}: Query{{$obj.Name}}Handler,
		Query{{$obj.PluralName}}: Query{{$obj.PluralName}}Handler,
		{{end}}
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

{{range $obj := .Model.ObjectEntities}}
{{if not $obj.IsExtended}}
func (r *GeneratedResolver) {{$obj.Name}}ResultType() {{$obj.Name}}ResultTypeResolver {
	return &Generated{{$obj.Name}}ResultTypeResolver{r}
}
{{end}}
{{if or $obj.HasAnyRelationships $obj.HasReadonlyColumns}}
func (r *GeneratedResolver) {{$obj.Name}}() {{$obj.Name}}Resolver {
	return &Generated{{$obj.Name}}Resolver{r}
}
{{end}}
{{end}}

type GeneratedMutationResolver struct{ *GeneratedResolver }

{{range $obj := .Model.ObjectEntities}}
{{if not $obj.IsExtended}}
func (r *GeneratedMutationResolver) Create{{$obj.Name}}(ctx context.Context, input map[string]interface{}) (item *{{$obj.Name}}, err error) {
	return r.Handlers.Create{{$obj.Name}}(ctx, r, input)
}
func Create{{$obj.Name}}Handler(ctx context.Context, r *GeneratedMutationResolver, input map[string]interface{}) (item *{{$obj.Name}}, err error) {
	principalID := GetPrincipalIDFromContext(ctx)
	now := time.Now()
	item = &{{$obj.Name}}{ID: uuid.Must(uuid.NewV4()).String(), CreatedAt: now, CreatedBy: principalID}
	tx := r.DB.db.Begin()

	event := events.NewEvent(events.EventMetadata{
		Type:        events.EventTypeCreated,
		Entity:      "{{$obj.Name}}",
		EntityID:    item.ID,
		Date:        now,
		PrincipalID: principalID,
	})

	var changes {{$obj.Name}}Changes
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
{{if $rel.IsToMany}}{{if not $rel.Target.IsExtended}}
	if ids,ok:=input["{{$rel.Name}}Ids"].([]interface{}); ok {
		items := []{{$rel.TargetType}}{}
		tx.Find(&items, "id IN (?)", ids)
		association := tx.Model(&item).Association("{{$rel.MethodName}}")
		association.Replace(items)
	}
{{end}}{{end}}
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
func (r *GeneratedMutationResolver) Update{{$obj.Name}}(ctx context.Context, id string, input map[string]interface{}) (item *{{$obj.Name}}, err error) {
	return r.Handlers.Update{{$obj.Name}}(ctx, r, id, input)
}
func Update{{$obj.Name}}Handler(ctx context.Context, r *GeneratedMutationResolver, id string, input map[string]interface{}) (item *{{$obj.Name}}, err error) {
	principalID := GetPrincipalIDFromContext(ctx)
	item = &{{$obj.Name}}{}
	now := time.Now()
	tx := r.DB.db.Begin()

	event := events.NewEvent(events.EventMetadata{
		Type:        events.EventTypeUpdated,
		Entity:      "{{$obj.Name}}",
		EntityID:    id,
		Date:        now,
		PrincipalID: principalID,
	})

	var changes {{$obj.Name}}Changes
	err = ApplyChanges(input, &changes)
	if err != nil {
		return 
	}

	err = GetItem(ctx, tx, item, &id)
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
{{if $rel.IsToMany}}{{if not $rel.Target.IsExtended}}
	if ids,ok:=input["{{$rel.Name}}Ids"].([]interface{}); ok {
		items := []{{$rel.TargetType}}{}
		tx.Find(&items, "id IN (?)", ids)
		association := tx.Model(&item).Association("{{$rel.MethodName}}")
		association.Replace(items)
	}
{{end}}{{end}}
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
func (r *GeneratedMutationResolver) Delete{{$obj.Name}}(ctx context.Context, id string) (item *{{$obj.Name}}, err error) {
	return r.Handlers.Delete{{$obj.Name}}(ctx, r, id)
}
func Delete{{$obj.Name}}Handler(ctx context.Context, r *GeneratedMutationResolver, id string) (item *{{$obj.Name}}, err error) {
	principalID := GetPrincipalIDFromContext(ctx)
	item = &{{$obj.Name}}{}
	now := time.Now()
	tx := r.DB.db.Begin()

	err = GetItem(ctx, tx, item, &id)
	if err != nil {
		return 
	}

	event := events.NewEvent(events.EventMetadata{
		Type:        events.EventTypeDeleted,
		Entity:      "{{$obj.Name}}",
		EntityID:    id,
		Date:        now,
		PrincipalID: principalID,
	})

	err = tx.Delete(item, TableName("{{$obj.TableName}}") + ".id = ?", id).Error
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
func (r *GeneratedMutationResolver) DeleteAll{{$obj.PluralName}}(ctx context.Context) (bool, error) {
	return r.Handlers.DeleteAll{{$obj.PluralName}}(ctx, r)
}
func DeleteAll{{$obj.PluralName}}Handler(ctx context.Context, r *GeneratedMutationResolver) (bool, error) {
	err := r.DB.db.Delete(&{{$obj.Name}}{}).Error
	return err == nil, err
}
{{end}}
{{end}}


type GeneratedQueryResolver struct{ *GeneratedResolver }


func (r *GeneratedQueryResolver) _service(ctx context.Context) (*_Service, error) {
	sdl := SchemaSDL
	return &_Service{
		Sdl: &sdl,
	}, nil
}

{{if .Model.HasFederatedTypes}}
func getExecutionContext(ctx context.Context) executionContext {
	e := ctx.Value(KeyExecutableSchema).(*executableSchema)
	return executionContext{graphql.GetRequestContext(ctx), e}
}

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
		
		switch typename { {{range $obj := .Model.ObjectEntities}}{{if $obj.HasDirective "key"}}
		case "{{$obj.Name}}":
			{{if $obj.IsExtended}}
				item := &{{$obj.Name}}{}
				{{range $col := $obj.Columns}}{{if $col.HasDirective "external"}}
				if v,ok:=anyValue["{{$col.Name}}"]; ok {
					_v,_ := v.({{$col.GoType}})
					item.{{$col.MethodName}} = _v
				}{{end}}{{end}}
			{{else}}
				ec := getExecutionContext(ctx)
				f, _err := ec.unmarshalInput{{$obj.Name}}FilterType(ctx, anyValue)
				err = _err
				if err != nil {
					return
				}
				item, qerr := r.Handlers.Query{{$obj.Name}}(ctx, r, nil, nil, &f)
				err = qerr
				if err != nil {
					return
				}
			{{end}}
			res = append(res, item)
			break;{{end}}{{end}}
		default:
			err = fmt.Errorf("The _entities resolver tried to load an entity for type \"%s\", but no object type of that name was found in the schema", typename)
			return
		}
	}
	return res, err
}
{{end}}

{{range $obj := .Model.ObjectEntities}}
{{if not $obj.IsExtended}}
func (r *GeneratedQueryResolver) {{$obj.Name}}(ctx context.Context, id *string, q *string, filter *{{$obj.Name}}FilterType) (*{{$obj.Name}}, error) {
	return r.Handlers.Query{{$obj.Name}}(ctx, r, id, q, filter)
}
func Query{{$obj.Name}}Handler(ctx context.Context, r *GeneratedQueryResolver, id *string, q *string, filter *{{$obj.Name}}FilterType) (*{{$obj.Name}}, error) {
	query := {{$obj.Name}}QueryFilter{q}
	offset := 0
	limit := 1
	rt := &{{$obj.Name}}ResultType{
		EntityResultType: EntityResultType{
			Offset: &offset,
			Limit:  &limit,
			Query:  &query,
			Filter: filter,
		},
	}
	qb := r.DB.Query()
	if id != nil {
		qb = qb.Where(TableName("{{$obj.TableName}}") + ".id = ?", *id)
	}

	var items []*{{$obj.Name}}
	err := rt.GetItems(ctx, qb, GetItemsOptions{Alias:TableName("{{$obj.TableName}}")}, &items)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("{{$obj.Name}} not found")
	}
	return items[0], err
}
func (r *GeneratedQueryResolver) {{$obj.PluralName}}(ctx context.Context, offset *int, limit *int, q *string, sort []{{$obj.Name}}SortType, filter *{{$obj.Name}}FilterType) (*{{$obj.Name}}ResultType, error) {
	return r.Handlers.Query{{$obj.PluralName}}(ctx, r, offset, limit, q , sort, filter)
}
func Query{{$obj.PluralName}}Handler(ctx context.Context, r *GeneratedQueryResolver, offset *int, limit *int, q *string, sort []{{$obj.Name}}SortType, filter *{{$obj.Name}}FilterType) (*{{$obj.Name}}ResultType, error) {
	_sort := []EntitySort{}
	for _, s := range sort {
		_sort = append(_sort, s)
	}
	query := {{$obj.Name}}QueryFilter{q}
	
	var selectionSet *ast.SelectionSet
	for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
		if f.Field.Name == "items" {
			selectionSet = &f.Field.SelectionSet
		}
	}
	
	return &{{$obj.Name}}ResultType{
		EntityResultType: EntityResultType{
			Offset: offset,
			Limit:  limit,
			Query:  &query,
			Sort: _sort,
			Filter: filter,
			SelectionSet: selectionSet,
		},
	}, nil
}

type Generated{{$obj.Name}}ResultTypeResolver struct{ *GeneratedResolver }

func (r *Generated{{$obj.Name}}ResultTypeResolver) Items(ctx context.Context, obj *{{$obj.Name}}ResultType) (items []*{{$obj.Name}}, err error) {
	err = obj.GetItems(ctx, r.DB.db, GetItemsOptions{Alias:TableName("{{$obj.TableName}}")}, &items)
	return
}

func (r *Generated{{$obj.Name}}ResultTypeResolver) Count(ctx context.Context, obj *{{$obj.Name}}ResultType) (count int, err error) {
	return obj.GetCount(ctx, r.DB.db, &{{$obj.Name}}{})
}
{{end}}

{{if or $obj.HasAnyRelationships $obj.HasReadonlyColumns}}
type Generated{{$obj.Name}}Resolver struct { *GeneratedResolver }

{{range $col := $obj.Columns}}
{{if $col.IsReadonlyType}}
func (r *Generated{{$obj.Name}}Resolver) {{$col.MethodName}}(ctx context.Context, obj *{{$obj.Name}}) (res {{$col.GoType}}, err error) {
	return r.Handlers.{{$obj.Name}}{{$col.MethodName}}(ctx, r, obj)
}
func {{$obj.Name}}{{$col.MethodName}}Handler(ctx context.Context,r *Generated{{$obj.Name}}Resolver, obj *{{$obj.Name}}) (res {{$col.GoType}}, err error) {
	return nil, fmt.Errorf("Resolver handler for {{$obj.Name}}{{$col.MethodName}} not implemented")
}
{{end}}
{{end}}

{{range $index, $rel := .Relationships}}
func (r *Generated{{$obj.Name}}Resolver) {{$rel.MethodName}}(ctx context.Context, obj *{{$obj.Name}}) (res {{$rel.ReturnType}}, err error) {
	return r.Handlers.{{$obj.Name}}{{$rel.MethodName}}(ctx, r, obj)
}
func {{$obj.Name}}{{$rel.MethodName}}Handler(ctx context.Context,r *Generated{{$obj.Name}}Resolver, obj *{{$obj.Name}}) (res {{$rel.ReturnType}}, err error) {
	{{if $rel.Target.IsExtended}}
		err = fmt.Errorf("not implemented")
	{{else}}
		{{if $rel.IsToMany}}
			items := []*{{$rel.TargetType}}{}
			err = r.DB.Query().Model(obj).Related(&items, "{{$rel.MethodName}}").Error
			res = items
		{{else}}
			loaders := ctx.Value(KeyLoaders).(map[string]*dataloader.Loader)
			if obj.{{$rel.MethodName}}ID != nil {
				item, _err := loaders["{{$rel.Target.Name}}"].Load(ctx, dataloader.StringKey(*obj.{{$rel.MethodName}}ID))()
				res, _ = item.({{$rel.ReturnType}})
				err = _err
			}
		{{end}}
	{{end}}
	return 
}
{{if $rel.IsToMany}}
func (r *Generated{{$obj.Name}}Resolver) {{$rel.MethodName}}Ids(ctx context.Context, obj *{{$obj.Name}}) (ids []string, err error) {
	ids = []string{}

	items := []*{{$rel.TargetType}}{}
	err = r.DB.Query().Model(obj).Select(TableName("{{$rel.Target.TableName}}") + ".id").Related(&items, "{{$rel.MethodName}}").Error

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
