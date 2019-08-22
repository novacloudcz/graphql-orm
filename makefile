test: test-generate test-build-lambda test-run
test-generate:
	GO111MODULE=on go run main.go init test
test-start:
	cd test && make run && cd ..
# test-build:
# 	cd test && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app *.go && cd ..
test-run:
	cd test && make test && cd ..
	# docker-compose -f test/docker-compose.yml up --build test
test-build-lambda:
	cd test && make build-lambda-function && cd ..
test-cleanup:
	docker-compose -f test/docker-compose.yml down