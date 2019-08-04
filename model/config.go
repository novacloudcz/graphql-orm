package model

import (
	"io/ioutil"
	"time"

	"github.com/ghodss/yaml"
)

type Config struct {
	Package    string `json:"package"`
	Connection *struct {
		MaxIdleConnections *uint   `json:"maxIdleConnections"`
		ConnMaxLifetime    *string `json:"connMaxLifetime"`
		MaxOpenConnections *uint   `json:"maxOpenConnections"`
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
	return 5
}

func (c *Config) MaxOpenConnections() uint {
	if c.Connection != nil && (*c.Connection).MaxOpenConnections != nil {
		return *(*c.Connection).MaxOpenConnections
	}
	return 10
}
func (c *Config) ConnMaxLifetime() float64 {
	if c.Connection != nil && (*c.Connection).ConnMaxLifetime != nil {
		val := *(*c.Connection).ConnMaxLifetime
		dur, err := time.ParseDuration(val)
		if err != nil {
			panic("failed to parse config connMaxLifetime duration, error: " + err.Error())
		}
		return dur.Seconds()
	}
	return time.Minute.Seconds()
}
