package ui

import (
	"github.com/daylioti/docker-commander/config"
	"github.com/daylioti/docker-commander/docker"
	"github.com/daylioti/docker-commander/ui/commands"
	"github.com/gizak/termui/v3"
)


const (
	MenuCommandsKey = iota
	MenuHelpKey
)

// UI main user interface struct.
type UI struct {
	ConfigUi *config.UIConfig

	Cnf          *config.Config
	DockerClient *docker.Docker

	// Terminal size, may change on window resize.
	TermWidth  int
	TermHeight int

	Commands *commands.Commands
	// Flag for switching control between menu and terminal area.
	SelectedMenu byte

	onRender bool
	//onRenderChan chan byte
}

// Init initialize all render components.
func (ui *UI) Init(cnf *config.Config, dockerClient *docker.Docker, configUi *config.UIConfig) {

	ui.TermWidth, ui.TermHeight = termui.TerminalDimensions()
	ui.Cnf = cnf
	ui.ConfigUi = configUi


	ui.Commands = &commands.Commands{
		ConfigUi:     ui.ConfigUi,
		TermWidth:    ui.TermWidth,
		TermHeight:   ui.TermHeight,
		DockerClient: dockerClient,
		RenderAll:    ui.Render,
		Cnf:          ui.Cnf,
	}
	ui.Commands.Init()

	ui.Render()
}

// Handle keyboard keys.
func (ui *UI) Handle(key string) {
	ui.Commands.Handle(key)
	//if !ui.switchMenu(key) {
	//	switch ui.SelectedMenu {
	//	case MenuCommandsKey:
	//		ui.Commands.Handle(key)
	//
	//	}
	//} else {
	//	ui.Render()
	//}
}

//func (ui *UI) switchMenu(key string) bool {
//	switch ui.SelectedMenu {
//	case MenuCommandsKey:
//		if key == "h" || key == "H" {
//			ui.SelectedMenu = MenuHelpKey
//			return true
//		}
//	case MenuHelpKey:
//		if key == "<Tab>" || key == "c" || key == "C" {
//			ui.SelectedMenu = MenuCommandsKey
//			return true
//		}
//	}
//	return false
//}

// Render main render function, that render all ui.
func (ui *UI) Render() {
	if !ui.onRender {
		ui.onRender = true
		termui.Clear()
		ui.Commands.Render()
		ui.onRender = false
	} else {
		for {
			if ui.onRender {
				ui.Render()
				return
			}
		}
	}
}
