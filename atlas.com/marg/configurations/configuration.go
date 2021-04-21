package configurations

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Configurator struct {
	l logrus.FieldLogger
}

func NewConfigurator(l logrus.FieldLogger) *Configurator {
	return &Configurator{l}
}

type Configuration struct {
	RespawnInterval int `yaml:"respawnInterval"`
}

func (c *Configurator) GetConfiguration() (*Configuration, error) {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		c.l.WithError(err).Errorf("Unable to read yaml file")
		return nil, err
	}

	con := &Configuration{}
	err = yaml.Unmarshal(yamlFile, con)
	if err != nil {
		c.l.WithError(err).Errorf("Unable to unmarshal configuration file.")
		return nil, err
	}

	return con, nil
}
