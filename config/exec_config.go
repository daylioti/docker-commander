package config

type ExecConfig struct {
	Connect    ExecConnect       `yaml:"connect"`
	Env        []string          `yaml:"env"`     // Environment variables.
	WorkingDir string            `yaml:"workdir"` // Working directory.
	Cmd        string            `yaml:"cmd"`     // Execution commands and args
	Input      map[string]string `yaml:"input"`
}
