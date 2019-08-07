SHELL := /bin/bash

test: generate-test build-test run-test
generate-test:
	cd test && go run ../main.go && cd ..
build-test:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o test/app test/*.go
run-test:
	docker-compose -f test/docker-compose.yml up --build test
cleanup-test:
	docker-compose -f test/docker-compose.yml down