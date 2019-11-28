test: test-generate test-build-lambda test-run-sqlite3 test-run-postgres test-run-mysql
test-generate:
	GO111MODULE=on go run main.go init test
test-start:
	cd test && make run && cd ..
test-start:
	cd test && make run && cd ..
# test-build:
# 	cd test && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app *.go && cd ..
test-run-sqlite3:
	cd test && DATABASE_URL=sqlite3://test.db make test && cd ..
test-run-postgres:
	cd test && DATABASE_URL=postgres://postgres@localhost:5432/test?sslmode=disable make test && cd ..
test-run-mysql:
	cd test && DATABASE_URL=mysql://root:@localhost/test make test && cd ..
test-build-lambda:
	cd test && make build-lambda-function && cd ..
test-cleanup:
	docker-compose -f test/docker-compose.yml down