package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type Config struct {
	Name         string   `yaml:"name"` // Display name
	Selected     bool     // Selected config or not
	Status       bool     // Display or not
	Config       []Config `yaml:"config"`        // Sub-configs (recursive)
	Exec         ExecConfig `yaml:"exec"`  // Docker exec config.
	Placeholders map[string]string `yaml:"placeholders"`
}

type ExecConfig struct {
	Connect ExecConnect `yaml:"connect"`
	Env []string      `yaml:"env"`        // Environment variables.
	WorkingDir string `yaml:"workdir"`    // Working directory.
	Cmd string      `yaml:"cmd"`        // Execution commands and args
}

type ExecConnect struct {
	FromImage string  `yaml:"container_image"` // The name of the image from which the container is made.
	ContainerName string `yaml:"container_name"`    // Container Name
	ContainerId string `yaml:"container_id"`// Container id
}


func (cfg *Config) Init(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		panic(err)
	}
	cfg.ChildConfigsPlaceholders(make(map[string]string), cfg)

	// Set default config data.
	cfg.Status = true
	cfg.Config[0].Selected = true
	for i := 0; i < len(cfg.Config); i++ {
		cfg.Config[i].Status = true
	}
	cfg.Config[0].Status = true
	if len(cfg.Config[0].Config) > 0 {
		for i := 0; i < len(cfg.Config[0].Config); i++ {
			cfg.Config[0].Config[i].Status = true
		}
	}
}

func (cfg *Config) ChildConfigsPlaceholders(placeholders map[string]string, c *Config) map[string]string {
	for i:=0;i<len(c.Config);i++ {
		for k, v := range c.Placeholders {
			placeholders[k] = v
		}
		for key, value := range placeholders {
          cfg.replacePlaceholder(key, value, &c.Config[i])
		}
		cfg.ChildConfigsPlaceholders(placeholders, &c.Config[i])
	}
	return placeholders
}

func (cfg *Config) replacePlaceholder(placeholder string, value string, c *Config) {
	c.Exec.WorkingDir = strings.Replace(c.Exec.WorkingDir, "@"+placeholder, value, 1)
	c.Exec.Connect.FromImage = strings.Replace(c.Exec.Connect.FromImage, "@"+placeholder, value, 1)
	c.Exec.Connect.ContainerId = strings.Replace(c.Exec.Connect.ContainerId, "@"+placeholder, value, 1)
	c.Exec.Cmd = strings.Replace(c.Exec.Cmd, "@"+placeholder, value, 1)
	for i:=0;i<len(c.Exec.Env);i++ {
		c.Exec.Env[i] = strings.Replace(c.Exec.Env[i], "@"+placeholder, value, 1)
	}
	for k, v := range c.Placeholders {
      c.Placeholders[k] =  strings.Replace(v, "@"+placeholder, value, 1)
	}
}
