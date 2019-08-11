package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/novacloudcz/graphql-orm/templates"

	"github.com/novacloudcz/goclitools"

	"gopkg.in/yaml.v2"

	"github.com/novacloudcz/graphql-orm/model"
	"github.com/urfave/cli"
)

var initCmd = cli.Command{
	Name:  "init",
	Usage: "initialize new project",
	Action: func(ctx *cli.Context) error {
		fmt.Println("Initializing project...")

		if !fileExists("graphql-orm.yml") {
			if err := createConfigFile(); err != nil {
				return cli.NewExitError(err, 1)
			}
		}

		if !fileExists("model.graphql") {
			if err := createDummyModelFile(); err != nil {
				return cli.NewExitError(err, 1)
			}
		}

		if err := createMainFile(); err != nil {
			return cli.NewExitError(err, 1)
		}

		if !fileExists("resolver.go") {
			if err := createResolverFile(); err != nil {
				return cli.NewExitError(err, 1)
			}
		}

		if err := createMakeFile(); err != nil {
			return cli.NewExitError(err, 1)
		}

		if err := createDockerFile(); err != nil {
			return cli.NewExitError(err, 1)
		}

		if err := runGenerate(); err != nil {
			return cli.NewExitError(err, 1)
		}

		return nil
	},
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		return true
	}
	return false
}

func createConfigFile() error {
	defaultPackagePath := ""
	if os.Getenv("GOPATH") != "" {
		cw, _ := os.Getwd()
		defaultPackagePath, _ = filepath.Rel(os.Getenv("GOPATH")+"/src", cw)
	}
	packagePath := goclitools.Prompt(fmt.Sprintf("Package path (default %s)", defaultPackagePath))
	if packagePath != "" {
		defaultPackagePath = packagePath
	}
	c := model.Config{Package: defaultPackagePath}

	content, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("graphql-orm.yml", content, 0644)
	return err
}
func createMainFile() error {
	c, err := model.LoadConfig()
	if err != nil {
		return err
	}
	return templates.WriteTemplate(templates.Main, "main.go", templates.TemplateData{Config: &c})
}
func createDummyModelFile() error {
	data := templates.TemplateData{Model: nil, Config: nil}
	return templates.WriteTemplate(templates.DummyModel, "model.graphql", data)
}
func createMakeFile() error {
	data := templates.TemplateData{Model: nil, Config: nil}
	return templates.WriteTemplate(templates.Makefile, "makefile", data)
}
func createDockerFile() error {
	c, err := model.LoadConfig()
	if err != nil {
		return err
	}
	data := templates.TemplateData{Model: nil, Config: &c}
	return templates.WriteTemplate(templates.Dockerfile, "Dockerfile", data)
}

func createResolverFile() error {
	c, err := model.LoadConfig()
	if err != nil {
		return err
	}
	data := templates.TemplateData{Model: nil, Config: &c}
	return templates.WriteTemplate(templates.Resolver, "resolver.go", data)
}

func runGenerate() error {
	return goclitools.RunInteractive("go run github.com/novacloudcz/graphql-orm")
}
