module github.com/novacloudcz/graphql-orm/test

go 1.12

require (
	github.com/99designs/gqlgen v0.10.2
	github.com/akrylysov/algnhsa v0.11.0 // indirect
	github.com/cloudevents/sdk-go v0.10.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/graph-gophers/dataloader v5.0.0+incompatible
	github.com/iancoleman/strcase v0.0.0-20190422225806-e506e3ef7365
	github.com/jakubknejzlik/cloudevents-aws-transport v0.1.4
	github.com/jinzhu/gorm v1.9.10
	github.com/mitchellh/mapstructure v0.0.0-20180203102830-a4e142e9c047
	github.com/novacloudcz/graphql-orm v0.0.0
	github.com/vektah/gqlparser v1.2.0
	gopkg.in/gormigrate.v1 v1.6.0
)

replace github.com/novacloudcz/graphql-orm v0.0.0 => ../
