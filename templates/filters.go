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
{{range $col := $object.Columns}}
{{range $fm := $col.FilterMapping}} {{$varName := (printf "f.%s%s" $col.MethodName $fm.SuffixCamel)}}
	if {{$varName}} != nil {
		db = db.Where(aliasPrefix + "{{$col.Name}} {{$fm.Operator}}",{{$fm.WrapValueVariable $varName}})
	}{{end}}{{end}}

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

{{end}}
`
