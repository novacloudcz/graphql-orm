package cmd

import (
	"fmt"

	"github.com/urfave/cli"
)

var initCmd = cli.Command{
	Name:  "init",
	Usage: "initialize new project",
	Action: func(ctx *cli.Context) error {
		fmt.Println("initializing")
		return nil
	},
}
