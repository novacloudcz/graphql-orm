package model

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
)

type Config struct {
	Package string `json:"package"`
}

func LoadConfig() (c Config, err error) {
	configSource, err := ioutil.ReadFile("graphql-orm.yml")
	if err != nil {
		return
	}
	err = yaml.Unmarshal(configSource, &c)
	if err != nil {
		return
	}
	return
}
