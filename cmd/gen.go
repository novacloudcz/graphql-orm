package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/novacloudcz/goclitools"

	"github.com/novacloudcz/graphql-orm/model"
	"github.com/novacloudcz/graphql-orm/templates"
	"github.com/urfave/cli"
)

var genCmd = cli.Command{
	Name:  "generate",
	Usage: "generate contents",
	Action: func(ctx *cli.Context) error {
		if err := generate("model.graphql", "."); err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

func generate(filename, p string) error {
	filename = path.Join(p, filename)
	fmt.Println("Generating contents from", filename, "...")
	modelSource, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	m, err := model.Parse(string(modelSource))
	if err != nil {
		return err
	}

	c, err := model.LoadConfigFromPath(p)
	if err != nil {
		return err
	}

	genPath := path.Join(p, "gen")
	if _, err := os.Stat(genPath); os.IsNotExist(err) {
		os.Mkdir(genPath, 0777)
	}

	err = model.EnrichModelObjects(&m)
	if err != nil {
		return err
	}

	err = generateFiles(p, &m, &c)
	if err != nil {
		return err
	}

	err = model.EnrichModel(&m)
	if err != nil {
		return err
	}

	schema, err := model.PrintSchema(m)
	if err != nil {
		return err
	}

	schema = "# This schema is generated, please don't update it manually\n\n" + schema

	if err := ioutil.WriteFile(path.Join(p, "gen/schema.graphql"), []byte(schema), 0644); err != nil {
		return err
	}

	fmt.Printf("Running gqlgen generator in %s ...\n", path.Join(p, "gen"))
	if err := goclitools.RunInteractiveInDir("go run github.com/99designs/gqlgen", path.Join(p, "gen")); err != nil {
		return err
	}

	// for _, obj := range plainModel.Objects() {
	// 	s1 := fmt.Sprintf("type %s struct {", obj.Name())
	// 	s2 := fmt.Sprintf("type %s struct {\n\t%sExtensions", obj.Name(), obj.Name())
	// 	if err := replaceStringInFile("gen/models_gen.go", s1, s2); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func generateFiles(p string, m *model.Model, c *model.Config) error {
	data := templates.TemplateData{Model: m, Config: c}
	if err := templates.WriteTemplate(templates.Database, path.Join(p, "gen/database.go"), data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.GQLGen, path.Join(p, "gen/gqlgen.yml"), data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.Model, path.Join(p, "gen/models.go"), data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.Filters, path.Join(p, "gen/filters.go"), data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.QueryFilters, path.Join(p, "gen/query-filters.go"), data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.Keys, path.Join(p, "gen/keys.go"), data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.Loaders, path.Join(p, "gen/loaders.go"), data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.HTTPHandler, path.Join(p, "gen/http-handler.go"), data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.GeneratedResolver, path.Join(p, "gen/resolver.go"), data); err != nil {
		return err
	}

	return nil
}
