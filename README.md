# graphql-orm

[![Build Status](https://travis-ci.org/novacloudcz/graphql-orm.svg?branch=master)](https://travis-ci.org/novacloudcz/graphql-orm)
[![Go Report Card](https://goreportcard.com/badge/github.com/novacloudcz/graphql-orm)](https://goreportcard.com/report/github.com/novacloudcz/graphql-orm)

Golang GraphQL API generator using [gqlgen](https://gqlgen.com) and [gorm](https://gorm.io)

## Why

While following microservices design patterns we ended up with many "model services". gqlgen is perfect tool, but implementing resolvers in every service is getting more and more cumbersome. Using this tool we only have to update `model.graphql` and all resolvers get generated automatically.

## Installation

Before you start, please make sure you have `goimports` installed:

```sh
go get golang.org/x/tools/cmd/goimports
```

_NOTE: Make sure you have Go installed on your system._

1. Create new project repository
1. run `go mod init [MODULE]` to initialize your project with go modules
1. run `go run github.com/novacloudcz/graphql-orm init`
1. all necessary files should be created and You can run `make run` to start service with dummy model
1. to update model open crated `model.graphql` and create your custom model schema
1. each time you change model, run `make generate` or `go run github.com/novacloudcz/graphql-orm` to recreate generated source codes

_NOTE: graphql-orm requires Go modules for installation. If you are running in \$GOPATH, make sure you are running init command with GO111MODULE=on_

## Running locally

For running locally you can use:

```sh
make run
```

or without makefile:

```sh
DATABASE_URL=sqlite3://test.db PORT=8080 go run *.go
```

## Environment variables

- `DATABASE_URL` - connection string for database in format `db://user:password@host:port/tablename` (eg. `mysql://root:pass@localhost:3306/test`; required)
- `EXPOSE_MIGRATION_ENDPOINT` - expose `/migration` endpoint which triggers database migration (migrates to latest database schema; default: false)
- `TABLE_NAME_PREFIX` - set global prefix for all table names (default: "")
- `EVENT_TRANSPORT_URL` - destination url for sending mutation events (array supported in format `EVENT_TRANSPORT_URL_[INDEX]`) see [Events transport](#installation)
- `EVENT_TRANSPORT_SOURCE` - custom value for CloudEvent source attribute (default: `http://{hostname}/graphql`)

### Sqlite connection

In case you want to connect with sqlite, you can use local file storage:
`sqlite3://path/to/file.db`

Or use in-memory storage:
`sqlite3://:memory:`

## Example

You can find example project at [graphql-orm-example repo](https://github.com/novacloudcz/graphql-orm-example)

## Schema preview in Voyager

[GraphQL Voyager](https://apis.guru/graphql-voyager/) is very nice tool for previewing your GraphQL Schema, you can run it locally by:

```sh
make voyager
```

or without makefile:

```sh
docker run --rm -v `pwd`/gen/schema.graphql:/app/schema.graphql -p 8080:80 graphql/voyager
```

...after voyager starts up go to https://localhost:8080

All generated stuff is stored in `./gen/` folder.

## Building docker image

If you generated Dockerfile initialization it's as easy as running:

```sh
docker build -t {IMAGE_NAME} .
```

If you want to create your own docker image, you can check the example repository for generated Dockerfile: https://github.com/novacloudcz/graphql-orm-example/blob/master/Dockerfile

## Events transport

For event driven architecture it's necessary that the service is able to send events about changes in state.
Services built using this library automatically send event for every mutation using CloudEvents (entity created/updated/deleted and changed column and their values). Supported targets are:

- HTTP/HTTPS
- AWS Services using [cloudevents-aws-transport](github.com/jakubknejzlik/cloudevents-aws-transport) (SNS/SQS/EventBridge)

For more information about event structure see: https://github.com/novacloudcz/graphql-orm/blob/master/events/model.go

## Migrations and automigrations

Since version `0.4.0` the migrations using gormigrate are introduced and it's possible to write custom migrations with rollbacks.
The automigration (with foreign keys) is still available, but gormigrate migrations are used by default. You use following commands:

- `make migrate` - runs gormigrate migrations
- `make automigrate` - runs gorm basic automigration

The same applies for HTTP endpoints (when `EXPOSE_MIGRATION_ENDPOINT=true`):

- `POST /migrate` - runs gormigrate migrations
- `post /automigrate` - runs gorm basic automigration

To add new migration, edit `src/migrations` file and its GetMigrations method. For more information see [gormigrate Readme](https://github.com/go-gormigrate/gormigrate)
