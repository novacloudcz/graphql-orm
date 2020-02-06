package main

import (
	"github.com/akrylysov/algnhsa"
	"github.com/novacloudcz/graphql-orm/events"
	"github.com/novacloudcz/graphql-orm/test/gen"
	"github.com/novacloudcz/graphql-orm/test/src"
)

func main() {
	db := gen.NewDBFromEnvVars()

	eventController, err := events.NewEventController()
	if err != nil {
		panic(err)
	}

	handler := gen.GetHTTPServeMux(src.New(db, &eventController), db, gen.GetMigrations(db))
	algnhsa.ListenAndServe(handler, &algnhsa.Options{
		UseProxyPath: true,
	})
}
