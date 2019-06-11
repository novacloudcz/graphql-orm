package templates

var Model = `package gen

import (
	"time"
)


{{range $object := .Objects}}

type {{.Name}} struct {
	ID        string     ` + "`" + `json:"id" gorm:"primary_key"` + "`" + `
{{range $col := $object.Columns}}
	{{$col.MethodName}} {{$col.GoType}} ` + "`" + `{{$col.ModelTags}}` + "`" + `{{end}}

{{range $rel := $object.Relationships}}
{{$rel.MethodName}} {{$rel.GoType}} ` + "`" + `{{$rel.ModelTags}}` + "`" + `
{{if $rel.IsToOne}}{{$rel.MethodName}}ID string {{end}}
{{end}}

	UpdatedAt time.Time  ` + "`" + `json:"updatedAt"` + "`" + `
	CreatedAt time.Time  ` + "`" + `json:"createdAt"` + "`" + `
	DeletedAt *time.Time ` + "`" + `json:"deletedAt"` + "`" + `
}

{{end}}
`
