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

// ReplacePlaceholders replace placeholders for config.
func (cfg *Config) ReplacePlaceholders(placeholders map[string]string, c *Config) {
	for k, v := range placeholders {
		cfg.ReplacePlaceholder(k, v, c)
	}
}

// GetPlaceholders all placeholders for selected config.
func (cfg *Config) GetPlaceholders(path []int, placeholders map[string]string, c *Config) map[string]string {
	if len(path) < 1 {
		return cfg.mergePlaceholders(c, placeholders)
	}
	return cfg.GetPlaceholders(path[1:],  cfg.mergePlaceholders(c, placeholders), &c.Config[path[0]])
}

// mergePlaceholders
func (cfg *Config) mergePlaceholders(c *Config, placeholders map[string]string) map[string]string {
	for key, value := range c.Placeholders {
		for k, v := range placeholders {
			cfg.Replace(&value, k, v)
		}
		placeholders[key] = value
	}
	return placeholders
}

// Replace replace placeholder strings.
func (cfg *Config) Replace(str *string, placeholder string, value string) {
	*str = strings.Replace(*str, "@"+placeholder, value, -1)
	*str = strings.Replace(*str, "["+placeholder+"]", value, -1)
}

// ReplacePlaceholder replace placeholders in all available fields.
func (cfg *Config) ReplacePlaceholder(placeholder string, value string, c *Config) {
	cfg.Replace(&c.Exec.WorkingDir, placeholder, value)
	cfg.Replace(&c.Exec.Connect.FromImage, placeholder, value)
	cfg.Replace(&c.Exec.Connect.ContainerID, placeholder, value)
	cfg.Replace(&c.Exec.Cmd, placeholder, value)
	for i := 0; i < len(c.Exec.Env); i++ {
		cfg.Replace(&c.Exec.Env[i], placeholder, value)
	}
	replacedPlaceholders := make(map[string]string)
	for k, v := range c.Placeholders {
		cfg.Replace(&k, placeholder, value)
		cfg.Replace(&v, placeholder, value)
		replacedPlaceholders[k] = v
	}
	c.Placeholders = replacedPlaceholders
	for k, v := range c.Exec.Input {
		c.Exec.Input[k].Key = strings.Replace(fmt.Sprintf("%v", v.Key), "@"+placeholder, value, -1)
		c.Exec.Input[k].Key = strings.Replace(fmt.Sprintf("%v", v.Key), "["+placeholder+"]", value, -1)
		c.Exec.Input[k].Value = strings.Replace(fmt.Sprintf("%v", v.Value), "@"+placeholder, value, -1)
		c.Exec.Input[k].Value = strings.Replace(fmt.Sprintf("%v", v.Value), "["+placeholder+"]", value, -1)
	}
}
