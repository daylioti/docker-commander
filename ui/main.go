package ui

import (
	"github.com/daylioti/docker-commander/config"
	"github.com/daylioti/docker-commander/docker"
	"github.com/gizak/termui/v3"
	"strconv"
)

// Main user interface struct.
type UI struct {
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

	// Render optimization variables.
	menuDisplayHash string
	ClearRender     bool
}

func (ui *UI) Init(cnf *config.Config, dockerClient *docker.Docker) {
	ui.TermWidth, ui.TermHeight = termui.TerminalDimensions()

	ui.SelectedRowTerminal = false

	ui.Term = &TerminalUI{}
	ui.Term.Init(ui, dockerClient)

	ui.Cmd = &Commands{}
	ui.Cmd.Init(cnf, dockerClient, ui)

	ui.Input = &Input{}
	ui.Input.Init(ui)

	ui.Render()
}

func (ui *UI) Handle(key string) {
	if len(ui.Input.Fields) > 0 {
		ui.Input.Handle(key)
		return
	}
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

// Main render function, that render all ui.
func (ui *UI) Render() {
	if ui.ClearRender {
		// This part needs for fix blinking border on blocks.
		// can be disabled by ClearRender for one time f.e after receive output from command.
		var hashMenu string
		for listIndex := 0; listIndex < len(ui.Cmd.Lists); listIndex++ {
			for rowIndex := 0; rowIndex < len(ui.Cmd.Lists[listIndex].Rows); rowIndex++ {
				hashMenu += ui.Cmd.Lists[listIndex].Rows[rowIndex]
			}
		}
		if hashMenu+strconv.Itoa(len(ui.Input.Fields)) != ui.menuDisplayHash {
			termui.Clear()
		}
		ui.menuDisplayHash = hashMenu
	}
	ui.ClearRender = true

	for listIndex := 0; listIndex < len(ui.Cmd.Lists); listIndex++ {
		termui.Render(ui.Cmd.Lists[listIndex])
	}
	termui.Render(ui.Term.TabPane, ui.Term.DisplayTerminal)
	if len(ui.Input.Fields) > 0 {
		for fieldIndex := 0; fieldIndex < len(ui.Input.Fields); fieldIndex++ {
			termui.Render(ui.Input.Fields[fieldIndex])
		}
	}
}
