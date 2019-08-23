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
	"github.com/rs/cors"
	"{{.Config.Package}}/gen"
	"{{.Config.Package}}/src"
)

func main() {
	app := cli.NewApp()
	app.Name = "graphql-orm"
	app.Usage = "This tool is for generating "
	app.Version = "0.0.0"

	app.Commands = []cli.Command{
		startCmd,
		migrateCmd,
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

var startCmd = cli.Command{
	Name:  "start",
	Usage: "start api server",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "cors",
			Usage: "Enable cors",
		},
		cli.StringFlag{
			Name:   "p,port",
			Usage:  "Port to listen to",
			Value:  "80",
			EnvVar: "PORT",
		},
	},
	Action: func(ctx *cli.Context) error {
		cors := ctx.Bool("cors")
		port := ctx.String("port")
		if err := startServer(cors, port); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		return nil
	},
}

var migrateCmd = cli.Command{
	Name:  "migrate",
	Usage: "migrate schema database",
	Action: func(ctx *cli.Context) error {
		fmt.Println("starting migration")
		if err := automigrate(); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		fmt.Println("migration complete")
		return nil
	},
}

func automigrate() error {
	db := gen.NewDBFromEnvVars()
	defer db.Close()
	return db.AutoMigrate().Error
}

func startServer(enableCors bool, port string) error {

	db := gen.NewDBFromEnvVars()
	defer db.Close()
	db.AutoMigrate()

	eventController, err := events.NewEventController()
	if err != nil {
		return err
	}

	mux := gen.GetHTTPServeMux(src.New(db, &eventController), db)

	mux.HandleFunc("/healthcheck", func(res http.ResponseWriter, req *http.Request) {
		if err := db.Ping(); err != nil {
			res.WriteHeader(400)
			res.Write([]byte("ERROR"))
			return
		}
		res.WriteHeader(200)
		res.Write([]byte("OK"))
	})

	var handler http.Handler
	if enableCors {
		handler = cors.AllowAll().Handler(mux)
	} else {
		handler = mux
	}

	log.Printf("connect to http://localhost:%s/graphql for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
	return nil
}
`
