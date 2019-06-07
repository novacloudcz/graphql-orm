package cmd

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/inloop/goclitools"

	"github.com/novacloudcz/graphql-orm/model"
	"github.com/novacloudcz/graphql-orm/templates"
	"github.com/urfave/cli"
)

var genCmd = cli.Command{
	Name:  "generate",
	Usage: "generate contents",
	Action: func(ctx *cli.Context) error {
		if err := generate("model.graphql"); err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

func generate(filename string) error {
	modelSource, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	m, err := model.Parse(string(modelSource))
	if err != nil {
		return err
	}

	if _, err := os.Stat("./gen"); os.IsNotExist(err) {
		os.Mkdir("./gen", 0777)
	}

	err = generateFiles(m)
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

	ioutil.WriteFile("gen/schema.graphql", []byte(schema), 0644)

	return goclitools.RunInteractiveInDir("go run github.com/99designs/gqlgen", "./gen")
}

func generateFiles(m model.Model) error {
	if err := writeTemplate(templates.Database, "gen/database.go", &m); err != nil {
		return err
	}
	if err := writeTemplate(templates.Resolver, "gen/resolver.go", &m); err != nil {
		return err
	}
	if err := writeTemplate(templates.GQLGen, "gen/gqlgen.yml", &m); err != nil {
		return err
	}

	// for _, obj := range m.Objects() {
	// 	if err := writeTemplate(templates.Model, fmt.Sprintf("gen/models_%s.go", strings.ToLower(obj.Name())), &obj); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func writeTemplate(t, filename string, data interface{}) error {
	temp, err := template.New("filename").Parse(t)
	if err != nil {
		return err
	}
	var content bytes.Buffer
	writer := io.Writer(&content)

	err = temp.Execute(writer, data)
	if err != nil {
		return err
	}
	ioutil.WriteFile(filename, content.Bytes(), 0777)
	return nil
}
