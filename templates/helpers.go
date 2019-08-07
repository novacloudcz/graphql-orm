package templates

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"text/template"

	"github.com/novacloudcz/goclitools"
	"github.com/novacloudcz/graphql-orm/model"
)

type TemplateData struct {
	Model  *model.Model
	Config *model.Config
}

func WriteTemplate(t, filename string, data TemplateData) error {
	temp, err := template.New(filename).Parse(t)
	if err != nil {
		return err
	}
	var content bytes.Buffer
	writer := io.Writer(&content)

	err = temp.Execute(writer, &data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, content.Bytes(), 0777)
	if err != nil {
		return err
	}
	if path.Ext(filename) == ".go" {
		return goclitools.RunInteractive(fmt.Sprintf("goimports -w %s", filename))
	}
	return nil
}
