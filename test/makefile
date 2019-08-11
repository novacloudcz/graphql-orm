generate:
	go run github.com/novacloudcz/graphql-orm

reinit:
	go run github.com/novacloudcz/graphql-orm init

run:
	DATABASE_URL=sqlite3://test.db PORT=8080 go run *.go

voyager:
	docker run --rm -v `pwd`/gen/schema.graphql:/app/schema.graphql -p 8080:80 graphql/voyager
