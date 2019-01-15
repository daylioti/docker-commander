package ui

import (
	"docker-commander/docker"
	"github.com/gizak/termui"
)

type UI struct {
	Cmd      *Commands
}

func (ui *UI) Init(configPath string, dockerClient *docker.Docker) {


	termui.Body.AddRows(
		termui.NewRow(),
		termui.NewRow(),
		termui.NewRow(),
	)


	ui.Cmd = &Commands{}
	ui.Cmd.Init(configPath, dockerClient)
}

func StringColor(text string, color string) string {
	return "[" + text + "](" + color + ")"
}
