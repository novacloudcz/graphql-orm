test: test-generate test-build test-run
test-generate:
	cd test && go run ../main.go init && cd ..
test-build:
	cd test && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app *.go && cd ..
test-run:
	docker-compose -f test/docker-compose.yml up --build test
test-cleanup:
	docker-compose -f test/docker-compose.yml down