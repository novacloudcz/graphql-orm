package templates

var ResolverMutations = `package gen

import (
	"context"
	"time"
	
	"github.com/graph-gophers/dataloader"
	"github.com/99designs/gqlgen/graphql"
	"github.com/gofrs/uuid"
	"github.com/novacloudcz/graphql-orm/events"
	"github.com/vektah/gqlparser/ast"
)

type GeneratedMutationResolver struct{ *GeneratedResolver }

{{range $obj := .Model.ObjectEntities}}
	func (r *GeneratedMutationResolver) Create{{$obj.Name}}(ctx context.Context, input map[string]interface{}) (item *{{$obj.Name}}, err error) {
		return r.Handlers.Create{{$obj.Name}}(ctx, r.GeneratedResolver, input)
	}
	func Create{{$obj.Name}}Handler(ctx context.Context, r *GeneratedResolver, input map[string]interface{}) (item *{{$obj.Name}}, err error) {
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
		{{end}}{{end}}
		
		err = tx.Create(item).Error
		if err != nil {
			tx.Rollback()
			return
		}
		
		{{range $rel := $obj.Relationships}}
			{{if $rel.IsToMany}}{{if not $rel.Target.IsExtended}}
				if ids,exists:=input["{{$rel.Name}}Ids"]; exists {
					items := []{{$rel.TargetType}}{}
					tx.Find(&items, "id IN (?)", ids)
					association := tx.Model(&item).Association("{{$rel.MethodName}}")
					association.Replace(items)
				}
			{{end}}{{end}}
		{{end}}

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
		return r.Handlers.Update{{$obj.Name}}(ctx, r.GeneratedResolver, id, input)
	}
	func Update{{$obj.Name}}Handler(ctx context.Context, r *GeneratedResolver, id string, input map[string]interface{}) (item *{{$obj.Name}}, err error) {
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
		
		err = tx.Save(item).Error
		if err != nil {
			tx.Rollback()
			return
		}

		{{range $rel := $obj.Relationships}}
		{{if $rel.IsToMany}}{{if not $rel.Target.IsExtended}}
			if ids,exists:=input["{{$rel.Name}}Ids"]; exists {
				items := []{{$rel.TargetType}}{}
				tx.Find(&items, "id IN (?)", ids)
				association := tx.Model(&item).Association("{{$rel.MethodName}}")
				association.Replace(items)
			}
		{{end}}{{end}}
		{{end}}

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
		return r.Handlers.Delete{{$obj.Name}}(ctx, r.GeneratedResolver, id)
	}
	func Delete{{$obj.Name}}Handler(ctx context.Context, r *GeneratedResolver, id string) (item *{{$obj.Name}}, err error) {
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
		return r.Handlers.DeleteAll{{$obj.PluralName}}(ctx, r.GeneratedResolver)
	}
	func DeleteAll{{$obj.PluralName}}Handler(ctx context.Context, r *GeneratedResolver) (bool, error) {
		err := r.DB.db.Delete(&{{$obj.Name}}{}).Error
		return err == nil, err
	}
{{end}}
`
