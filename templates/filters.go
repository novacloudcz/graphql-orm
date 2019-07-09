package templates

var Filters = `package gen

import (
	"context"
	"fmt"
	"strings"

)

{{range $object := .Model.Objects}}

func (f *{{$object.Name}}FilterType) Apply(ctx context.Context, wheres *[]string, values *[]interface{}, joins *[]string) error {
	return f.ApplyWithAlias(ctx, "{{$object.TableName}}", wheres, values, joins)
}
func (f *{{$object.Name}}FilterType) ApplyWithAlias(ctx context.Context, alias string, wheres *[]string, values *[]interface{}, joins *[]string) error {
	if f == nil {
		return nil
	}
	aliasPrefix := alias + "."
	
	_where, _values := f.WhereContent(aliasPrefix)
	*wheres = append(*wheres, _where...)
	*values = append(*values, _values...)
	
{{range $rel := $object.Relationships}}
{{$varName := (printf "f.%s" $rel.MethodName)}}
	if {{$varName}} != nil {
		_alias := alias + "_{{$rel.Name}}"
		*joins = append(*joins, {{$rel.JoinString}})
		err := {{$varName}}.ApplyWithAlias(ctx, _alias, wheres, values, joins)
		if err != nil {
			return err
		}
	}{{end}}

	return nil
}

func (f *{{$object.Name}}FilterType) WhereContent(aliasPrefix string) (conditions []string, values []interface{}) {
	conditions = []string{}
	values = []interface{}{}

	if f.Or != nil {
		cs := []string{}
		vs := []interface{}{}
		for _, or := range f.Or {
			_cond, _values := or.WhereContent(aliasPrefix)
			cs = append(cs, _cond...)
			vs = append(vs, _values...)
		}
		conditions = append(conditions, "("+strings.Join(cs, " OR ")+")")
		values = append(values, vs...)
	}
	if f.And != nil {
		cs := []string{}
		vs := []interface{}{}
		for _, and := range f.And {
			_cond, _values := and.WhereContent(aliasPrefix)
			cs = append(cs, _cond...)
			vs = append(vs, _values...)
		}
		conditions = append(conditions, strings.Join(cs, " AND "))
		values = append(values, vs...)
	}

{{range $col := $object.Columns}}
{{range $fm := $col.FilterMapping}} {{$varName := (printf "f.%s%s" $col.MethodName $fm.SuffixCamel)}}
	if {{$varName}} != nil {
		conditions = append(conditions, aliasPrefix + "{{$col.Name}} {{$fm.Operator}}")
		values = append(values, {{$fm.WrapValueVariable $varName}})
	}{{end}}{{end}}

	return
}

{{end}}
`
