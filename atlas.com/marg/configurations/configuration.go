package configurations

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Configurator struct {
	l *log.Logger
}

func NewConfigurator(l *log.Logger) *Configurator {
	return &Configurator{l}
}

type Configuration struct {
	RespawnInterval int `yaml:"respawnInterval"`
}

func (c *Configurator) GetConfiguration() (*Configuration, error) {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		c.l.Printf("yamlFile.Get err   #%v ", err)
		return nil, err
	}

	con := &Configuration{}
	err = yaml.Unmarshal(yamlFile, con)
	if err != nil {
		c.l.Fatalf("Unmarshal: %v", err)
		return nil, err
	}

	return con, nil
}
