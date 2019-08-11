package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
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
		p := ctx.Args().First()
		if p == "" {
			p = "./"
		}

		fmt.Printf("Initializing project in %s ...\n", p)

		if !fileExists(path.Join(p, "graphql-orm.yml")) {
			if err := createConfigFile(p); err != nil {
				return cli.NewExitError(err, 1)
			}
		}

		if !fileExists(path.Join(p, "model.graphql")) {
			if err := createDummyModelFile(p); err != nil {
				return cli.NewExitError(err, 1)
			}
		}

		if err := createMainFile(p); err != nil {
			return cli.NewExitError(err, 1)
		}

		if !fileExists(path.Join(p, "resolver.go")) {
			if err := createResolverFile(p); err != nil {
				return cli.NewExitError(err, 1)
			}
		}

		if err := createMakeFile(p); err != nil {
			return cli.NewExitError(err, 1)
		}

		if err := createDockerFile(p); err != nil {
			return cli.NewExitError(err, 1)
		}

		if !fileExists(path.Join(p, "go.mod")) {
			if err := initModules(p); err != nil {
				return cli.NewExitError(err, 1)
			}
		}

		if err := runGenerate(p); err != nil {
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

func createConfigFile(p string) error {
	defaultPackagep := ""
	if os.Getenv("GOp") != "" {
		cw, _ := os.Getwd()
		defaultPackagep, _ = filepath.Rel(os.Getenv("GOp")+"/src", cw)
	}
	packagep := goclitools.Prompt(fmt.Sprintf("Package p (default %s)", defaultPackagep))
	if packagep != "" {
		defaultPackagep = packagep
	}
	c := model.Config{Package: defaultPackagep}

	content, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(p, "graphql-orm.yml"), content, 0644)
	return err
}
func createMainFile(p string) error {
	c, err := model.LoadConfigFromPath(p)
	if err != nil {
		return err
	}
	return templates.WriteTemplate(templates.Main, path.Join(p, "main.go"), templates.TemplateData{Config: &c})
}
func createDummyModelFile(p string) error {
	data := templates.TemplateData{Model: nil, Config: nil}
	return templates.WriteTemplate(templates.DummyModel, path.Join(p, "model.graphql"), data)
}
func createMakeFile(p string) error {
	data := templates.TemplateData{Model: nil, Config: nil}
	return templates.WriteTemplate(templates.Makefile, path.Join(p, "makefile"), data)
}
func createDockerFile(p string) error {
	c, err := model.LoadConfigFromPath(p)
	if err != nil {
		return err
	}
	data := templates.TemplateData{Model: nil, Config: &c}
	return templates.WriteTemplate(templates.Dockerfile, path.Join(p, "Dockerfile"), data)
}

func initModules(p string) error {
	c, err := model.LoadConfigFromPath(p)
	if err != nil {
		return err
	}
	return goclitools.RunInteractiveInDir(fmt.Sprintf("go mod init %s", c.Package), p)
}

func createResolverFile(p string) error {
	c, err := model.LoadConfigFromPath(p)
	if err != nil {
		return err
	}
	data := templates.TemplateData{Model: nil, Config: &c}
	return templates.WriteTemplate(templates.Resolver, path.Join(p, "resolver.go"), data)
}

func runGenerate(p string) error {
	return goclitools.RunInteractiveInDir("go run github.com/novacloudcz/graphql-orm", p)
}
