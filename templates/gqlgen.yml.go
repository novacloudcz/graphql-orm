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
  {{range .Model.Objects}}
  {{.Name}}ResultType:
    model: {{$config.Package}}/gen.{{.Name}}ResultType
    fields:
      count:
        resolver: true
      items:
        resolver: true
  {{.Name}}:
    model: {{$config.Package}}/gen.{{.Name}}
    fields:{{range .Relationships}}
      {{.Name}}:
        resolver: true{{end}}
  {{.Name}}CreateInput:
    model: "map[string]interface{}"
  {{.Name}}UpdateInput:
    model: "map[string]interface{}"
  {{end}}
`
