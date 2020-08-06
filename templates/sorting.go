package templates

var Sorting = `package gen

import (
	"context"
	
	"github.com/jinzhu/gorm"
)

{{range $obj := .Model.ObjectEntities}}
{{if not $obj.IsExtended}}
func (s {{$obj.Name}}SortType) Apply(ctx context.Context, dialect gorm.Dialect, sorts *[]SortInfo, joins *[]string) error {
	return s.ApplyWithAlias(ctx, dialect, TableName("{{$obj.TableName}}"), sorts, joins)
}
func (s {{$obj.Name}}SortType) ApplyWithAlias(ctx context.Context, dialect gorm.Dialect, alias string, sorts *[]SortInfo, joins *[]string) error {
	aliasPrefix := dialect.Quote(alias) + "."
	
	{{range $col := $obj.Columns}} {{if $col.IsSortable}} 
	if s.{{$col.MethodName}} != nil {
		sort := SortInfo{Field: aliasPrefix+dialect.Quote("{{$col.Name}}"), Direction: s.{{$col.MethodName}}.String()}
		*sorts = append(*sorts, sort)
	}
	{{range $agg := $col.Aggregations}}
	if s.{{$col.MethodName}}{{$agg.Name}} != nil {
		sort := SortInfo{Field: "{{$agg.Name}}("+aliasPrefix+dialect.Quote("{{$col.Name}}")+")", Direction: s.{{$col.MethodName}}{{$agg.Name}}.String(),IsAggregation: true}
		*sorts = append(*sorts, sort)
	}
	{{end}}
	{{end}}{{end}}
	
	{{range $rel := $obj.Relationships}}
	{{if not $rel.Target.IsExtended}}
	{{$varName := (printf "s.%s" $rel.MethodName)}}
	if {{$varName}} != nil {
		_alias := alias + "_{{$rel.Name}}"
		*joins = append(*joins, {{$rel.JoinString}})
		err := {{$varName}}.ApplyWithAlias(ctx, dialect, _alias, sorts, joins)
		if err != nil {
			return err
		}
	}{{end}}{{end}}

	return nil
}
{{end}}
{{end}}
`
