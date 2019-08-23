package cmd

import (
	"os"

	"github.com/urfave/cli"
)

// Execute ...
func Execute() {
	app := cli.NewApp()
	app.Name = "graphql-orm"
	app.Usage = "This tool is for generating GraphQL API using gqlgen and gorm"
	app.Version = "0.0.1"

	app.Action = genCmd.Action
	app.Usage = genCmd.Usage
	app.Flags = genCmd.Flags

	app.Commands = []cli.Command{
		initCmd,
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
