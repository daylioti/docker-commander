package config

import "gopkg.in/yaml.v2"

// ExecConfig docker exec configs (what to execute).
type ExecConfig struct {
	Connect    ExecConnect   `yaml:"connect"`
	Env        []string      `yaml:"env"`     // Environment variables.
	WorkingDir string        `yaml:"workdir"` // Working directory.
	Cmd        string        `yaml:"cmd"`     // Execution commands and args
	Input      yaml.MapSlice `yaml:"input"`
}
