test: test-generate test-build test-run
test-generate:
	go run main.go init test
test-build:
	cd test && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app *.go && cd ..
test-run:
	docker-compose -f test/docker-compose.yml up --build test
test-cleanup:
	docker-compose -f test/docker-compose.yml down