package config

import (
	"strconv"
)

type UIConfig struct {
	UI struct {
		Commands map[string]string `yaml:"commands"`
	}
}

func (uc *UIConfig) GetCommandsHeight() int {
	var height int
	if h, exist := uc.UI.Commands["height"]; !exist {
		height = 5
	} else {
		height, _ = strconv.Atoi(h)
	}
	return height + 2
}
