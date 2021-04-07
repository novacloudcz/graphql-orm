package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/novacloudcz/graphql-orm/model"
	"github.com/novacloudcz/graphql-orm/templates"
	"github.com/novacloudcz/graphql-orm/tools"
	"github.com/urfave/cli"
)

var genCmd = cli.Command{
	Name:  "generate",
	Usage: "generate contents",
	Action: func(ctx *cli.Context) error {
		if err := generate("model*.graphql", "."); err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

func generate(filePattern, p string) error {
	filePattern = path.Join(p, filePattern)
	matches, err := filepath.Glob(filePattern)
	if err != nil {
		return err
	}

	fmt.Println("Generating contents from", matches, "...")
	modelSource := ""
	for _, file := range matches {
		fmt.Println("Appending content from model file", file)
		source, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		modelSource += string(source) + "\n"
	}

	m, err := model.Parse(modelSource)
	if err != nil {
		return err
	}

	c, err := model.LoadConfigFromPath(p)
	if err != nil {
		return err
	}

	genPath := path.Join(p, "gen")
	ensureDir(genPath)

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

	schemaSDL, err := model.PrintSchema(m)
	if err != nil {
		return err
	}

	err = model.BuildFederatedModel(&m)
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

	var re = regexp.MustCompile(`(?sm)schema\W{[^}]+}`)
	schemaSDL = re.ReplaceAllString(schemaSDL, ``)
	var re2 = regexp.MustCompile(`(?sm)type _Service {[^}]+}`)
	schemaSDL = re2.ReplaceAllString(schemaSDL, ``)
	schemaSDL = strings.Replace(schemaSDL, "\n  _service: _Service!", "", 1)
	schemaSDL = strings.Replace(schemaSDL, "\n  _entities(representations: [_Any!]!): [_Entity]!", "", 1)
	schemaSDL = strings.Replace(schemaSDL, "\nscalar _Any", "", 1)
	var re3 = regexp.MustCompile(`(?sm)[\n]{3,}`)
	schemaSDL = re3.ReplaceAllString(schemaSDL, "\n\n")
	schemaSDL = strings.Trim(schemaSDL, "\n")
	constants := map[string]interface{}{
		"SchemaSDL": schemaSDL,
	}
	if err := templates.WriteTemplateRaw(templates.Constants, path.Join(p, "gen/constants.go"), constants); err != nil {
		return err
	}

	fmt.Printf("Running gqlgen generator in %s ...\n", path.Join(p, "gen"))
	if err := tools.RunInteractiveInDir("go run github.com/99designs/gqlgen", path.Join(p, "gen")); err != nil {
		return err
	}

	return nil
}

func ensureDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0777)
		if err != nil {
			panic(err)
		}
	}
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
	if err := templates.WriteTemplate(templates.Sorting, path.Join(p, "gen/sorting.go"), data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.QueryFilters, path.Join(p, "gen/query-filters.go"), data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.Loaders, path.Join(p, "gen/loaders.go"), data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.HTTPHandler, path.Join(p, "gen/http-handler.go"), data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.ResolverCore, path.Join(p, "gen/resolver.go"), data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.ResolverQueries, path.Join(p, "gen/resolver-queries.go"), data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.ResolverMutations, path.Join(p, "gen/resolver-mutations.go"), data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.ResolverExtensions, path.Join(p, "gen/resolver-extensions.go"), data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.ResolverFederation, path.Join(p, "gen/resolver-federation.go"), data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.Federation, path.Join(p, "gen/federation.go"), data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.ResultType, path.Join(p, "gen/result-type.go"), data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.ResolverSrcGen, path.Join(p, "src/resolver_gen.go"), data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.Migrations, path.Join(p, "gen/migrations.go"), data); err != nil {
		return err
	}
	if err := templates.WriteTemplate(templates.EventsController, path.Join(p, "gen/events-controller.go"), data); err != nil {
		return err
	}
	if !fileExists(path.Join(p, "src/migrations.go")) {
		if err := templates.WriteTemplate(templates.MigrationsSrc, path.Join(p, "src/migrations.go"), data); err != nil {
			return err
		}
	}

	return nil
}
