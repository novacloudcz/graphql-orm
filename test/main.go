package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/novacloudcz/graphql-orm/events"
	"github.com/novacloudcz/graphql-orm/test/gen"
	"github.com/novacloudcz/graphql-orm/test/src"
	"github.com/rs/cors"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "graphql-orm"
	app.Usage = "This tool is for generating "
	app.Version = "0.0.0"

	app.Commands = []cli.Command{
		startCmd,
		migrateCmd,
		automigrateCmd,
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
	Usage: "run migration sequence (using gomigrate)",
	Action: func(ctx *cli.Context) error {
		fmt.Println("starting migration")
		if err := migrate(); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		fmt.Println("migration complete")
		return nil
	},
}

var automigrateCmd = cli.Command{
	Name:  "automigrate",
	Usage: "migrate schema database using basic gorm migration",
	Action: func(ctx *cli.Context) error {
		fmt.Println("starting automigration")
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
	return db.AutoMigrate()
}

func migrate() error {
	db := gen.NewDBFromEnvVars()
	defer db.Close()
	return db.Migrate(src.GetMigrations(db))
}

func startServer(enableCors bool, port string) error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	db := gen.NewDBFromEnvVars()
	defer db.Close()

	eventController, err := events.NewEventController()
	if err != nil {
		return err
	}

	mux := gen.GetHTTPServeMux(src.New(db, &eventController), db, src.GetMigrations(db))

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

	h := &http.Server{Addr: ":" + port, Handler: handler}

	go func() {
		log.Printf("connect to http://localhost:%s/graphql for GraphQL playground", port)
		log.Fatal(h.ListenAndServe())
	}()

	<-stop

	log.Println("\nShutting down the server...")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	err = h.Shutdown(ctx)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	log.Println("Server gracefully stopped")

	err = db.Close()
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	log.Println("Database connection closed")

	return nil
}
