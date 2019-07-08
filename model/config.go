package model

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
)

type Config struct {
	Package    string `json:"package"`
	Connection *struct {
		MaxIdleConnections *uint `json:"maxIdleConnections"`
		MaxOpenConnections *uint `json:"maxOpenConnections"`
	} `json:"connection,omitempty"`
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

func (c *Config) MaxIdleConnections() uint {
	if c.Connection != nil && (*c.Connection).MaxIdleConnections != nil {
		return *(*c.Connection).MaxIdleConnections
	}
	return 0
}

func (c *Config) MaxOpenConnections() uint {
	if c.Connection != nil && (*c.Connection).MaxOpenConnections != nil {
		return *(*c.Connection).MaxOpenConnections
	}
	return 10
}
