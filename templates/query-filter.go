package templates

var QueryFilters = `package gen

import (
	"context"
	"strings"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/ast"

)

{{range $object := .Model.Objects}}

type {{$object.Name}}QueryFilter struct {
	Query *string
}

func (qf *{{$object.Name}}QueryFilter) Apply(ctx context.Context, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if qf.Query == nil {
		return nil
	}
	fields := []*ast.Field{}
	for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
		fields = append(fields, f.Field)
	}

	ors := []string{}

	queryParts := strings.Split(*qf.Query, " ")
	for _, part := range queryParts {
		if err := qf.applyQueryWithFields(fields, part, "{{$object.TableName}}", &ors, values, joins); err != nil {
			return err
		}
		*wheres = append(*wheres, "("+strings.Join(ors, " OR ")+")")
	}
	return nil
}

func (qf *{{$object.Name}}QueryFilter) applyQueryWithFields(fields []*ast.Field, query, alias string, ors *[]string, values *[]interface{}, joins *[]string) error {
	if len(fields) == 0 {
		return nil
	}
	aliasPrefix := alias + "."
	
	fieldsMap := map[string]*ast.Field{}
	for _, f := range fields {
		fieldsMap[f.Name] = f
	}

	{{range $col := $object.Columns}}{{if $col.IsSearchable}}
	if _, ok := fieldsMap["{{$col.Name}}"]; ok {
		*ors = append(*ors, fmt.Sprintf("%[1]s{{$col.Name}} LIKE ? OR %[1]s{{$col.Name}} LIKE ?", aliasPrefix))
		*values = append(*values, fmt.Sprintf("%s%%", query), fmt.Sprintf("%% %s%%", query))
	}
	{{end}}
	{{end}}

	{{range $rel := $object.Relationships}}
	if f, ok := fieldsMap["employees"]; ok {
		_fields := []*ast.Field{}
		_alias := alias + "_{{$rel.Name}}"
		*joins = append(*joins,{{$rel.JoinString}})
		
		for _, s := range f.SelectionSet {
			if f, ok := s.(*ast.Field); ok {
				_fields = append(_fields, f)
			}
		}
		q := {{$rel.Target.Name}}QueryFilter{qf.Query}
		err := q.applyQueryWithFields(_fields, query, _alias, ors, values, joins)
		if err != nil {
			return err
		}
	}
	{{end}}
	
	return nil
}

{{end}}
`
