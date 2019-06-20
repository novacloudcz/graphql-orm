# graphql-orm

Golang GraphQL API generator using [https://gqlgen.com](gqlgen) and [https://gorm.io](gorm)

# Installation

_NOTE: Make sure you have Go installed on your system._

1. Create new project repository
1. run `go run github.com/novacloudcz/graphql-orm init`
1. follow initialization instruction (creating makefile is suggested)
1. open create `model.graphql` and create your custom model schema
1. each time you change model, run `make generate` or `go run github.com/novacloudcz/graphql-orm` to recreate generated source codes

All generated stuff is stored in `./gen/` folder.

# What's this library for?

While following microservices design patterns we ended up with many "model services". gqlgen is perfect tool, but implementing resolvers in every service is getting more and more cumbersome. Using this tool we only have to update `model.graphql` and all resolvers get generated automatically.
