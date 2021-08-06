package templates

// Makefile ...
var Makefile = `generate:
	GO111MODULE=on go run github.com/novacloudcz/graphql-orm

reinit:
	GO111MODULE=on go run github.com/novacloudcz/graphql-orm init

migrate:
	DATABASE_URL=sqlite3://test.db PORT=8080 go run *.go migrate

automigrate:
	DATABASE_URL=sqlite3://test.db PORT=8080 go run *.go automigrate

run:
	DATABASE_URL=sqlite3://test.db PORT=8080 go run *.go start --cors

voyager:
	docker run --rm -v ` + "`" + `pwd` + "`" + `/gen/schema.graphql:/app/schema.graphql -p 8080:80 graphql/voyager

build-lambda-function:
	GO111MODULE=on GOOS=linux CGO_ENABLED=0 go build -o main lambda/main.go && zip lambda.zip main && rm main

test-sqlite:
	GO111MODULE=on go build -o app *.go && DATABASE_URL=sqlite3://test.db ./app migrate && (ENABLE_DELETE_ALL_RESOLVERS=true DATABASE_URL=sqlite3://test.db PORT=8080 ./app start& export app_pid=$$! && make test-godog || test_result=$$? && kill $$app_pid && exit $$test_result)
test:
	GO111MODULE=on go build -o app *.go && ./app migrate && (ENABLE_DELETE_ALL_RESOLVERS=true PORT=8080 ./app start& export app_pid=$$! && make test-godog || test_result=$$? && kill $$app_pid && exit $$test_result)
// TODO: add detection of host ip (eg. host.docker.internal) for other OS
test-godog:
	docker run --rm --network="host" -v "${PWD}/features:/godog/features" -e GRAPHQL_URL=http://$$(if [[ $${OSTYPE} == darwin* ]]; then echo host.docker.internal;else echo localhost;fi):8080/graphql jakubknejzlik/godog-graphql
`
