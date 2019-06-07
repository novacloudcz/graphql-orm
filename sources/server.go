package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	"github.com/graphql-services/memberships"
	"github.com/graphql-services/memberships/database"
)

const (
	defaultPort = "80"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	urlString := os.Getenv("DATABASE_URL")
	if urlString == "" {
		panic(fmt.Errorf("database url must be provided"))
	}

	db := database.NewDBWithString(urlString)
	defer db.Close()
	db.AutoMigrate(&memberships.Member{}, &memberships.Membership{})

	gqlHandler := handler.GraphQL(memberships.NewExecutableSchema(memberships.Config{Resolvers: &memberships.Resolver{DB: db}}))
	http.Handle("/", handler.Playground("GraphQL playground", "/graphql"))
	http.HandleFunc("/graphql", func(res http.ResponseWriter, req *http.Request) {
		ctx := context.WithValue(req.Context(), memberships.DBContextKey, db)
		req = req.WithContext(ctx)
		gqlHandler(res, req)
	})

	http.HandleFunc("/healthcheck", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write([]byte("OK"))
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
