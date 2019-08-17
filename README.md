# graphql-orm

[![Build Status](https://travis-ci.org/novacloudcz/graphql-orm.svg?branch=master)](https://travis-ci.org/novacloudcz/graphql-orm)

Golang GraphQL API generator using [gqlgen](https://gqlgen.com) and [gorm](https://gorm.io)

# Installation

Before you start, please make sure you have `goimports` installed:

```
go get golang.org/x/tools/cmd/goimports
```

_NOTE: Make sure you have Go installed on your system._

1. Create new project repository
1. run `go run github.com/novacloudcz/graphql-orm init`
1. follow initialization instruction (creating makefile is suggested)
1. open create `model.graphql` and create your custom model schema
1. each time you change model, run `make generate` or `go run github.com/novacloudcz/graphql-orm` to recreate generated source codes

## Running locally

For running locally you can use:

```
make run
```

or without makefile:

```
DATABASE_URL=sqlite3://test.db PORT=8080 go run *.go
```

# Example

You can find example project at [graphql-orm-example repo](https://github.com/novacloudcz/graphql-orm-example)

# Schema preview in Voyager

[GraphQL Voyager](https://apis.guru/graphql-voyager/) is very nice tool for previewing your GraphQL Schema, you can run it locally by:

```
make voyager
```

or without makefile:

```
docker run --rm -v `pwd`/gen/schema.graphql:/app/schema.graphql -p 8080:80 graphql/voyager
```

...after voyager starts up go to https://localhost:8080

All generated stuff is stored in `./gen/` folder.

# Building docker image

If you generated Dockerfile initialization it's as easy as running:

```
docker build -t {IMAGE_NAME} .
```

If you want to create your own docker image, you can check the example repository for generated Dockerfile: https://github.com/novacloudcz/graphql-orm-example/blob/master/Dockerfile

# What's this library for?

While following microservices design patterns we ended up with many "model services". gqlgen is perfect tool, but implementing resolvers in every service is getting more and more cumbersome. Using this tool we only have to update `model.graphql` and all resolvers get generated automatically.
