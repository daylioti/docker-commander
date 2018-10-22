package config

import (
	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type Config struct {
	Name         string `yaml:"name"`            // Display name
	Selected     bool                            // Selected config or not
	Status       bool                            // Display or not
	Config       []Config `yaml:"config"`        // Sub-configs (recursive)
	Command      string   `yaml:"command"`       // bash command
	Container    string   `yaml:"container"`     // The name of the image from which the container is made
	ChildConfigs []Config `yaml:"child_configs"` // Able to insert insert configs structure into child configs.
}

func (cfg *Config) Init(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	_ = yaml.Unmarshal(data, cfg)
	cfg.ChildConfigsInsert(cfg)
}

func (cfg *Config) ChildConfigsInsert(c *Config) {
	var i int
	var name string
	var replace Config
	if c.Config != nil {
		for i = 0; i < len(c.Config); i++ {
			cfg.ChildConfigsInsert(&c.Config[i])
		}
	}
	if c.ChildConfigs != nil {
		for i = 0; i < len(c.Config); i++ {
			for j:=0;j<len(c.ChildConfigs);j++ {
				//for _, cc := range c.ChildConfigs {
					name = c.Config[i].Name
					replace = cfg.ChildConfigsPlaceholders(name, c.ChildConfigs[j])
					//cfg.ChildConfigsPlaceholders(name, c.ChildConfigs[j])
					c.Config[i].Config = append(c.Config[i].Config, replace)
				//}
			}
		}
		//c.ChildConfigs = nil
	}
}

func (cfg *Config) ChildConfigsPlaceholders(name string, c Config) Config {
	config := Config{}
	_ = copier.Copy(config, c)
	if c.Config != nil {
		for i := 0; i < len(c.Config); i++ {
			config.Config = append(config.Config, cfg.ChildConfigsPlaceholders(name, c.Config[i]))
			//config.Config[i] = cfg.ChildConfigsPlaceholders(name, c.Config[i])
		}
	}
	config.Name = strings.Replace(c.Name, "@parent", name, 1)
	config.Command = strings.Replace(c.Command, "@parent", name, 1)
	config.Container = c.Container
	return config
}
