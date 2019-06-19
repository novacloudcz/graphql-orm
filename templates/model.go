package templates

var Model = `package gen

import (
	"time"
	"github.com/novacloudcz/graphql-orm/resolvers"
)


{{range $object := .Model.Objects}}

type {{.Name}}ResultType struct {
	resolvers.EntityResultType
}

type {{.Name}} struct {
{{range $col := $object.Columns}}
	{{$col.MethodName}} {{$col.GoType}} ` + "`" + `{{$col.ModelTags}}` + "`" + `{{end}}

{{range $rel := $object.Relationships}}
{{$rel.MethodName}} {{$rel.GoType}} ` + "`" + `{{$rel.ModelTags}}` + "`" + `
{{end}}
}

{{end}}
`
