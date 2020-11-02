package templates

var Filters = `package gen

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

{{range $obj := .Model.ObjectEntities}}
{{if not $obj.IsExtended}}
func (f *{{$obj.Name}}FilterType) IsEmpty(ctx context.Context, dialect *gorm.Statement) bool {
	wheres := []string{}
	havings := []string{}
	whereValues := []interface{}{}
	havingValues := []interface{}{}
	joins := []string{}
	err := f.ApplyWithAlias(ctx, dialect, "companies", &wheres, &whereValues, &havings, &havingValues, &joins)
	if err != nil {
		panic(err)
	}
	return len(wheres) == 0 && len(havings) == 0
}
func (f *{{$obj.Name}}FilterType) Apply(ctx context.Context, dialect *gorm.Statement, wheres *[]string, whereValues *[]interface{}, havings *[]string, havingValues *[]interface{}, joins *[]string) error {
	return f.ApplyWithAlias(ctx, dialect, TableName("{{$obj.TableName}}"), wheres, whereValues, havings, havingValues, joins)
}
func (f *{{$obj.Name}}FilterType) ApplyWithAlias(ctx context.Context, dialect *gorm.Statement, alias string, wheres *[]string, whereValues *[]interface{}, havings *[]string, havingValues *[]interface{}, joins *[]string) error {
	if f == nil {
		return nil
	}
	aliasPrefix := dialect.Quote(alias) + "."
	
	_where, _whereValues := f.WhereContent(dialect, aliasPrefix)
	_having, _havingValues := f.HavingContent(dialect, aliasPrefix)
	*wheres = append(*wheres, _where...)
	*havings = append(*havings, _having...)
	*whereValues = append(*whereValues, _whereValues...)
	*havingValues = append(*havingValues, _havingValues...)


	if f.Or != nil {
		ws := []string{}
		hs := []string{}
		wvs := []interface{}{}
		hvs := []interface{}{}
		js := []string{}
		for _, or := range f.Or {
			_ws := []string{}
			_hs := []string{}
			err := or.ApplyWithAlias(ctx, dialect, alias, &_ws, &wvs, &_hs, &hvs, &js)
			if err != nil {
				return err
			}
			if len(_ws) > 0 {
				ws = append(ws, strings.Join(_ws, " AND "))
			}
			if len(_hs) > 0 {
				hs = append(hs, strings.Join(_hs, " AND "))
			}
		}
		if len(ws) > 0 {
			*wheres = append(*wheres, "("+strings.Join(ws, " OR ")+")")
		}
		if len(hs) > 0 {
			*havings = append(*havings, "("+strings.Join(hs, " OR ")+")")
		}
		*whereValues = append(*whereValues, wvs...)
		*havingValues = append(*havingValues, hvs...)
		*joins = append(*joins, js...)
	}
	if f.And != nil {
		ws := []string{}
		hs := []string{}
		wvs := []interface{}{}
		hvs := []interface{}{}
		js := []string{}
		for _, and := range f.And {
			err := and.ApplyWithAlias(ctx, dialect, alias, &ws, &wvs, &hs, &hvs, &js)
			if err != nil {
				return err
			}
		}
		if len(ws) > 0 {
			*wheres = append(*wheres, strings.Join(ws, " AND "))
		}
		if len(hs) > 0 {
			*havings = append(*havings, strings.Join(hs, " AND "))
		}
		*whereValues = append(*whereValues, wvs...)
		*havingValues = append(*havingValues, hvs...)
		*joins = append(*joins, js...)
	}
	
	{{range $rel := $obj.Relationships}}
	{{if not $rel.Target.IsExtended}}
	{{$varName := (printf "f.%s" $rel.MethodName)}}
	if {{$varName}} != nil {
		_alias := alias + "_{{$rel.Name}}"
		*joins = append(*joins, {{$rel.JoinString}})
		err := {{$varName}}.ApplyWithAlias(ctx, dialect, _alias, wheres, whereValues, havings, havingValues, joins)
		if err != nil {
			return err
		}
	}{{end}}{{end}}

	return nil
}

func (f *{{$obj.Name}}FilterType) WhereContent(dialect *gorm.Statement, aliasPrefix string) (conditions []string, values []interface{}) {
	conditions = []string{}
	values = []interface{}{}

	{{range $col := $obj.Columns}} {{if $col.IsFilterable}}
		{{range $fm := $col.FilterMapping}} {{$varName := (printf "f.%s%s" $col.MethodName $fm.SuffixCamel)}}
			if {{$varName}} != nil {
				conditions = append(conditions, aliasPrefix + dialect.Quote("{{$col.Name}}")+" {{$fm.Operator}}")
				values = append(values, {{$fm.WrapValueVariable $varName}})
			}
		{{end}}
		if f.{{$col.MethodName}}Null != nil {
			if *f.{{$col.MethodName}}Null {
				conditions = append(conditions, aliasPrefix+dialect.Quote("{{$col.Name}}")+" IS NULL")
			} else {
				conditions = append(conditions, aliasPrefix+dialect.Quote("{{$col.Name}}")+" IS NOT NULL")
			}
		}
	{{end}} {{end}}

	return
}
func (f *{{$obj.Name}}FilterType) HavingContent(dialect *gorm.Statement, aliasPrefix string) (conditions []string, values []interface{}) {
	conditions = []string{}
	values = []interface{}{}

	{{range $col := $obj.Columns}} {{if $col.IsFilterable}}
		{{range $fm := $col.FilterMapping}} {{range $fn := $col.Aggregations}} {{$varName := (printf "f.%s%s%s" $col.MethodName $fn.Name $fm.SuffixCamel)}}
			if {{$varName}} != nil {
				conditions = append(conditions, "{{$fn.Name}}("+aliasPrefix + dialect.Quote("{{$col.Name}}")+") {{$fm.Operator}}")
				values = append(values, {{$fm.WrapValueVariable $varName}})
			}
		{{end}} {{end}}
	{{end}} {{end}}

	return
}

// AndWith convenience method for combining two or more filters with AND statement
func (f *{{$obj.Name}}FilterType) AndWith(f2 ...*{{$obj.Name}}FilterType) *{{$obj.Name}}FilterType {
	_f2 := f2[:0]
	for _, x := range f2 {
		if x != nil {
			_f2 = append(_f2, x)
		}
	}
	if len(_f2) == 0 {
		return f
	}
	return &{{$obj.Name}}FilterType{
		And: append(_f2,f),
	}
}

// OrWith convenience method for combining two or more filters with OR statement
func (f *{{$obj.Name}}FilterType) OrWith(f2 ...*{{$obj.Name}}FilterType) *{{$obj.Name}}FilterType {
	_f2 := f2[:0]
	for _, x := range f2 {
		if x != nil {
			_f2 = append(_f2, x)
		}
	}
	if len(_f2) == 0 {
		return f
	}
	return &{{$obj.Name}}FilterType{
		Or: append(_f2,f),
	}
}

{{end}}
{{end}}
`
