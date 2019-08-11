package templates

var Main = `package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/99designs/gqlgen/handler"
	"github.com/novacloudcz/graphql-orm/events"
	jwtgo "github.com/dgrijalva/jwt-go"
	// "github.com/rs/cors"
	"{{.Config.Package}}/gen"
)

const (
	defaultPort = "80"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db := gen.NewDBFromEnvVars()
	defer db.Close()
	db.AutoMigrate()

	eventController, err := events.NewEventController()
	if err != nil {
		panic(err)
	}

	mux := gen.GetHTTPServeMux(NewResolver(db, &eventController), db)

	mux.HandleFunc("/healthcheck", func(res http.ResponseWriter, req *http.Request) {
		if err := db.Ping(); err != nil {
			res.WriteHeader(400)
			res.Write([]byte("ERROR"))
			return
		}
		res.WriteHeader(200)
		res.Write([]byte("OK"))
	})

	handler := mux
	// use this line to allow cors for all origins/methods/headers (for development)
	// handler := cors.AllowAll().Handler(mux)

	log.Printf("connect to http://localhost:%s/graphql for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
`
