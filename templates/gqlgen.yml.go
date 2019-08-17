package templates

var GQLGen = `# Generated with graphql-orm
{{$config:=.Config}}
schema:
  - schema.graphql
exec:
  filename: generated.go
  package: gen
model:
  filename: models_gen.go
  package: gen
resolver:
  filename: resolver.go
  type: Resolver
  package: gen

models:
  {{range $obj := .Model.Objects}}
  {{$obj.Name}}:
    model: {{$config.Package}}/gen.{{$obj.Name}}
    fields:{{range $col := $obj.Columns}}{{if $col.IsReadonlyType}}
      {{$col.Name}}:
        resolver: true{{end}}{{end}}{{range $rel := $obj.Relationships}}
      {{$rel.Name}}:
        resolver: true{{end}}
  {{if not $obj.IsExtended}}
  {{$obj.Name}}ResultType:
    model: {{$config.Package}}/gen.{{$obj.Name}}ResultType
    fields:
      count:
        resolver: true
      items:
        resolver: true
  {{$obj.Name}}CreateInput:
    model: "map[string]interface{}"
  {{$obj.Name}}UpdateInput:
    model: "map[string]interface{}"
  {{end}}
  {{end}}
  _Any:
    model: {{$config.Package}}/gen._Any
`
