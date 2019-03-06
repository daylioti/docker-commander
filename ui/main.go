package ui

import (
	"github.com/daylioti/docker-commander/config"
	"github.com/daylioti/docker-commander/docker"
	"github.com/gizak/termui"
)

type UI struct {
	Cmd         *Commands
	Term        *TerminalUi
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

	ui.Term = &TerminalUi{}
	ui.Term.Init(ui, dockerClient)

	ui.Render()
}

func (ui *UI) Handle(key string) {
	switch key {
	case "<Tab>":
		if ui.SelectedRow >= len(ui.Grid.Items) {
			ui.SelectedRow++
		} else {
			ui.SelectedRow = 0
		}
	default:
		if ui.SelectedRow == 0 {
			ui.Cmd.Handle(key)
		}
	}
}

func (ui *UI) Render() {
    var cols []interface{}
	termWidth, termHeight := termui.TerminalDimensions()
	ui.Grid = nil
	ui.Grid = termui.NewGrid()

	ui.Grid.SetRect(0, 0, termWidth, termHeight)

	for _, list := range ui.Cmd.GetLists() {
		cols = append(cols, termui.NewCol(list.ratio, list.list))
	}
	ui.Grid.Set(
		termui.NewRow(0.2, cols...),
		//termui.NewRow(0.8, cmd.terminal),
	)
	termui.Render(ui.Grid)
}
