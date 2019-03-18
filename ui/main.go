package ui

import (
	"github.com/daylioti/docker-commander/config"
	"github.com/daylioti/docker-commander/docker"
	"github.com/gizak/termui/v3"
)

type UI struct {
	Cmd                 *Commands
	Term                *TerminalUI
	Grid                *termui.Grid
	SelectedRowTerminal bool
	widthDimension      int
}

func (ui *UI) Init(cnf *config.Config, dockerClient *docker.Docker) {
	ui.Grid = termui.NewGrid()
	termWidth, termHeight := termui.TerminalDimensions()
	ui.Grid.SetRect(0, 0, termWidth, termHeight)
	ui.SelectedRowTerminal = false
	ui.widthDimension, _ = termui.TerminalDimensions()

	ui.Cmd = &Commands{}
	ui.Cmd.Init(cnf, dockerClient, ui)

	ui.Term = &TerminalUI{}
	ui.Term.Init(ui, dockerClient)

	ui.Render()
}

func (ui *UI) Handle(key string) {
	switch key {
	case "<Tab>":
		ui.SelectedRowTerminal = !ui.SelectedRowTerminal
		if ui.SelectedRowTerminal {
			ui.Term.TabPane.BorderStyle = termui.NewStyle(termui.ColorGreen)
			ui.Term.DisplayTerminal.BorderStyle = termui.NewStyle(termui.ColorGreen)
		} else {
			ui.Term.TabPane.BorderStyle = termui.NewStyle(termui.ColorWhite)
			ui.Term.DisplayTerminal.BorderStyle = termui.NewStyle(termui.ColorWhite)
		}
		ui.Render()
	default:
		if !ui.SelectedRowTerminal {
			ui.Cmd.Handle(key)
		} else {
			ui.Term.Handle(key)
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
	index := ui.Term.client.Exec.GetActiveTerminalIndex()
	if index >= 0 {
		ui.Term.DisplayTerminal = ui.Term.client.Exec.Terminals[index].List
	}
	ui.Term.TabPane.TabNames = ui.Term.GetTabPaneItems()

	ratio := float64(termHeight) * 0.8 * 4 / 100 / 100
	ui.Grid.Set(
		termui.NewRow(0.2, cols...),
		termui.NewRow(0.8,
			termui.NewRow(ratio, ui.Term.TabPane),
			termui.NewRow(1-ratio, ui.Term.DisplayTerminal),
		),
	)
	termui.Render(ui.Grid)
}
