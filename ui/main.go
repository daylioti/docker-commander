package ui

import (
	"github.com/daylioti/docker-commander/config"
	"github.com/daylioti/docker-commander/docker"
	"github.com/daylioti/docker-commander/ui/helpers"
	"github.com/gizak/termui/v3"
)

// UI main user interface struct.
type UI struct {
	configUi *config.UIConfig

	// Terminal size, may change on window resize.
	TermWidth  int
	TermHeight int

	// Menu struct.
	Cmd *Commands

	// Terminal area struct.
	Term *TerminalUI

	// Inputs struct
	Input *Input

	// Flag for switching control between menu and terminal area.
	SelectedRowTerminal bool
}

// Init initialize all render components.
func (ui *UI) Init(cnf *config.Config, dockerClient *docker.Docker, configUi *config.UIConfig) {
	ui.TermWidth, ui.TermHeight = termui.TerminalDimensions()

	ui.SelectedRowTerminal = false

	ui.configUi = configUi

	ui.Term = &TerminalUI{}
	ui.Term.Init(ui, dockerClient)

	ui.Cmd = &Commands{}
	ui.Cmd.Init(cnf, dockerClient, ui)

	ui.Input = &Input{}
	ui.Input.Init(ui)

	ui.Render()

	termui.StyleParserColorMap = helpers.GetAllTermColors()
}

// Handle keyboard keys.
func (ui *UI) Handle(key string) {
	if len(ui.Input.Fields) > 0 {
		ui.Input.Handle(key)
		return
	}
	switch key {
	case "<Tab>":
		ui.SelectedRowTerminal = !ui.SelectedRowTerminal
		if ui.SelectedRowTerminal {
			ui.Term.Focus()
			ui.Cmd.UnFocus()
		} else {
			ui.Term.UnFocus()
			ui.Cmd.Focus()
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

// Render main render function, that render all ui.
func (ui *UI) Render() {
	termui.Clear()
	ui.Cmd.Render()
	ui.Term.Render()
	if len(ui.Input.Fields) > 0 {
		for _, field := range ui.Input.Fields {
			termui.Render(field)
		}
	}
}
