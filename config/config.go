package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
)

func Load(path string) (*Config, error) {
	conf := &Config{}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read config: %s ", err)
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		return nil, fmt.Errorf("could not deserialize config: %s ", err)
	}
	return conf, nil
}

func (c *Config) Validate() error {
	validate := validator.New()
	err := validate.Struct(c)
	if err != nil {
		return err
	}
	return nil
}
