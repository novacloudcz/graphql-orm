package templates

var ResolverFederation = `package gen

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
		
		switch typename { {{range $obj := .Model.Objects}}{{if $obj.IsFederatedType}}
		case "{{$obj.Name}}":
			ec := getExecutionContext(ctx)
			f, _err := ec.unmarshalInput{{$obj.Name}}FilterType(ctx, anyValue)
			err = _err
			if err != nil {
				return
			}

			if f.IsEmpty(ctx, r.DB.Query().Dialect()) {
				res = append(res, nil)
				continue
			}
			
			item, qerr := r.{{$obj.Name}}(ctx, nil, nil, &f)
			if qerr != nil {
				if _, isNotFound := qerr.(*NotFoundError); !isNotFound {
					err = qerr
					return
				}
				res = append(res, nil)
			} else {
				res = append(res, item)
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

{{range $ext := .Model.ObjectExtensions}}
{{$obj := $ext.Object}}

	{{if or $obj.HasAnyRelationships $obj.HasReadonlyColumns}}

		{{range $index, $rel := $obj.Relationships}}
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
		{{end}}

	{{end}}

{{end}}
`
