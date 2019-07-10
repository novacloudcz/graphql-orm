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

	if f.Or != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, or := range f.Or {
			err := or.ApplyWithAlias(ctx, alias, &cs, &vs, &js)
			if err != nil {
				return err
			}
		}
		*wheres = append(*wheres, "("+strings.Join(cs, " OR ")+")")
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}
	if f.And != nil {
		cs := []string{}
		vs := []interface{}{}
		js := []string{}
		for _, and := range f.And {
			err := and.ApplyWithAlias(ctx, alias, &cs, &vs, &js)
			if err != nil {
				return err
			}
		}
		*wheres = append(*wheres, strings.Join(cs, " AND "))
		*values = append(*values, vs...)
		*joins = append(*joins, js...)
	}
	
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

{{range $col := $object.Columns}}
{{range $fm := $col.FilterMapping}} {{$varName := (printf "f.%s%s" $col.MethodName $fm.SuffixCamel)}}
	if {{$varName}} != nil {
		conditions = append(conditions, aliasPrefix + "{{$col.Name}} {{$fm.Operator}}")
		values = append(values, {{$fm.WrapValueVariable $varName}})
	}{{end}}{{end}}

	return
}

// AndWith convenience method for combining two or more filters with AND statement
func (f *{{$object.Name}}FilterType) AndWith(f2 ...*{{$object.Name}}FilterType) *{{$object.Name}}FilterType {
	_f2 := f2[:0]
	for _, x := range f2 {
		if x != nil {
			_f2 = append(_f2, x)
		}
	}
	if len(_f2) == 0 {
		return f
	}
	return &{{$object.Name}}FilterType{
		And: append(_f2,f),
	}
}

// OrWith convenience method for combining two or more filters with OR statement
func (f *{{$object.Name}}FilterType) OrWith(f2 ...*{{$object.Name}}FilterType) *{{$object.Name}}FilterType {
	_f2 := f2[:0]
	for _, x := range f2 {
		if x != nil {
			_f2 = append(_f2, x)
		}
	}
	if len(_f2) == 0 {
		return f
	}
	return &{{$object.Name}}FilterType{
		Or: append(_f2,f),
	}
}

{{end}}
`
