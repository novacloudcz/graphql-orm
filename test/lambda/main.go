package main

import (
	"github.com/akrylysov/algnhsa"
	"github.com/novacloudcz/graphql-orm/events"
	"github.com/novacloudcz/graphql-orm/test/gen"
	"github.com/novacloudcz/graphql-orm/test/resolver"
)

func main() {
	db := gen.NewDBFromEnvVars()

	eventController, err := events.NewEventController()
	if err != nil {
		panic(err)
	}

	handler := gen.GetHTTPServeMux(resolver.New(db, &eventController), db)
	algnhsa.ListenAndServe(handler, nil)
}
