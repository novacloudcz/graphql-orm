package templates

var Filters = `package gen

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
)

{{range $object := .Model.Objects}}

func (f *{{$object.Name}}FilterType) Apply(db *gorm.DB) (*gorm.DB, error) {
	return f.ApplyWithAlias(db, "{{$object.TableName}}")
}
func (f *{{$object.Name}}FilterType) ApplyWithAlias(db *gorm.DB, alias string) (*gorm.DB, error) {
	if f == nil {
		return db, nil
	}
	aliasPrefix := alias
	if aliasPrefix != "" {
		aliasPrefix += "."
	}

conditions, values := f.WhereContent(aliasPrefix)
db = db.Where(conditions, values...)

{{range $rel := $object.Relationships}}
{{$varName := (printf "f.%s" $rel.MethodName)}}
	if {{$varName}} != nil {
		_alias := alias + "_{{$rel.Name}}"
		db = db.Joins({{$rel.JoinString}})
		_db, err := {{$varName}}.ApplyWithAlias(db, _alias)
		if err != nil {
			return db, err
		}
		db = _db
	}{{end}}

	return db, nil
}

func (f *{{$object.Name}}FilterType) WhereContent(aliasPrefix string) (conditions string, values []interface{}) {
	_conditions := []string{}
	values = []interface{}{}

	if f.Or != nil {
		cs := []string{}
		vs := []interface{}{}
		for _, or := range f.Or {
			_cond, _values := or.WhereContent(aliasPrefix)
			cs = append(cs, _cond)
			vs = append(vs, _values...)
		}
		_conditions = append(_conditions, "("+strings.Join(cs, " OR ")+")")
		values = append(values, vs...)
	}
	if f.And != nil {
		cs := []string{}
		vs := []interface{}{}
		for _, or := range f.Or {
			_cond, _values := or.WhereContent(aliasPrefix)
			cs = append(cs, _cond)
			vs = append(vs, _values...)
		}
		_conditions = append(_conditions, strings.Join(cs, " AND "))
		values = append(values, vs...)
	}

{{range $col := $object.Columns}}
{{range $fm := $col.FilterMapping}} {{$varName := (printf "f.%s%s" $col.MethodName $fm.SuffixCamel)}}
	if {{$varName}} != nil {
		_conditions = append(_conditions, aliasPrefix + "{{$col.Name}} {{$fm.Operator}}")
		values = append(values, {{$fm.WrapValueVariable $varName}})
	}{{end}}{{end}}

	conditions = strings.Join(_conditions, " AND ")
	return
}

{{end}}
`
