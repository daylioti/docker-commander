package ui

import (
	"github.com/daylioti/docker-commander/config"
	"github.com/daylioti/docker-commander/docker"
	"github.com/gizak/termui"
)

type UI struct {
	Cmd         *Commands
	Grid        *termui.Grid
	SelectedRow int
}

func (ui *UI) Init(cnf *config.Config, dockerClient *docker.Docker) {
	ui.Grid = termui.NewGrid()
	termWidth, termHeight := termui.TerminalDimensions()
	ui.Grid.SetRect(0, 0, termWidth, termHeight)
	ui.SelectedRow = 0
	ui.Cmd = &Commands{}
	ui.Cmd.Init(cnf, dockerClient, ui)
}

func StringColor(text string, color string) string {
	return "[" + text + "](" + color + ")"
}
