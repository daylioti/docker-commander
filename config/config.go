package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Config main config structure for commands lists.
type Config struct {
	Name         string            `yaml:"name"` // Display name
	Selected     bool              // Selected config or not
	Status       bool              // Display or not
	Config       []Config          `yaml:"config"` // Sub-configs (recursive)
	Exec         ExecConfig        `yaml:"exec"`   // Docker exec config.
	Placeholders map[string]string `yaml:"placeholders"`
}

// CnfInit unmarshal yml by structures.
func CnfInit(path string, configs ...interface{}) {
	var err error
	var data []byte
	if data, err = ioutil.ReadFile(path); err != nil {
		_, parseErr := url.Parse(path)
		if parseErr == nil {
			// Get from url
			client := &http.Client{Timeout: time.Second}
			if r, responseErr := client.Get(path); responseErr == nil {
				data, err = ioutil.ReadAll(r.Body)
				if err != nil {
					fmt.Printf("Can't open config from %v", path)
					os.Exit(0)
				}
			} else {
				fmt.Printf("Can't open config from %v", path)
				os.Exit(0)
			}
		} else {
			fmt.Printf("Can't open config from %v", path)
			os.Exit(0)
		}
	}
	for _, cfg := range configs {
		data = []byte(os.ExpandEnv(string(data)))
		if err = yaml.Unmarshal(data, cfg); err != nil {
			panic(err)
		}
	}
}

// Init set default selected items, replace placeholders.
func (cfg *Config) Init() {
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

// ChildConfigsPlaceholders replace placeholders in children menu items.
func (cfg *Config) ChildConfigsPlaceholders(placeholders map[string]string, c *Config) map[string]string {
	for i := 0; i < len(c.Config); i++ {
		for key, value := range c.Placeholders {
			placeholders[key] = value
		}
		for key, value := range c.Config[i].Placeholders {
			placeholders[key] = value
		}
		for key, value := range placeholders  {
			cfg.ReplacePlaceholder(key, value, &c.Config[i])
		}
		cfg.ChildConfigsPlaceholders(placeholders, &c.Config[i])
	}
	for key, value := range c.Placeholders {
		cfg.ReplacePlaceholder(key, value, c)
	}
	return placeholders
}

// ReplacePlaceholder replace placeholders in all available fields.
func (cfg *Config) ReplacePlaceholder(placeholder string, value string, c *Config) {
	c.Exec.WorkingDir = strings.Replace(c.Exec.WorkingDir, "@"+placeholder, value, -1)
	c.Exec.Connect.FromImage = strings.Replace(c.Exec.Connect.FromImage, "@"+placeholder, value, -1)
	c.Exec.Connect.ContainerID = strings.Replace(c.Exec.Connect.ContainerID, "@"+placeholder, value, -1)
	c.Exec.Cmd = strings.Replace(c.Exec.Cmd, "@"+placeholder, value, -1)
	for i := 0; i < len(c.Exec.Env); i++ {
		c.Exec.Env[i] = strings.Replace(c.Exec.Env[i], "@"+placeholder, value, -1)
	}
	for k, v := range c.Placeholders {
		c.Placeholders[k] = strings.Replace(v, "@"+placeholder, value, -1)
	}
	for k, v := range c.Exec.Input {
		c.Exec.Input[k].Key = strings.Replace(fmt.Sprintf("%v", v.Key), "@"+placeholder, value, -1)
		c.Exec.Input[k].Value = strings.Replace(fmt.Sprintf("%v", v.Value), "@"+placeholder, value, -1)
	}
}
