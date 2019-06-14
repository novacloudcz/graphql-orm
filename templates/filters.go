package templates

var Filters = `package gen

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
)

{{range $object := .Model.Objects}}

func (f *{{$object.Name}}FilterType) Apply(db *gorm.DB) (*gorm.DB, error) {
{{range $column := $object.Columns}}
{{range $fm := $column.FilterMapping}} {{$varName := (printf "f.%s%s" $column.MethodName $fm.SuffixCamel)}}
	if {{$varName}} != nil {
		db = db.Where("{{$column.Name}} {{$fm.Operator}}",{{$fm.WrapValueVariable $varName}})
	}{{end}}{{end}}

	return db, nil
}

{{end}}
`
