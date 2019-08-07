SHELL := /bin/bash

test: test-generate test-build test-run
test-generate:
	cd test && go run ../main.go && cd ..
test-build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o test/app test/*.go
test-run:
	docker-compose -f test/docker-compose.yml up --build test
test-cleanup:
	docker-compose -f test/docker-compose.yml down