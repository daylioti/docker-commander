package ui

import (
	"github.com/daylioti/docker-commander/config"
	"github.com/daylioti/docker-commander/docker"
	"github.com/daylioti/docker-commander/ui/commands"
	"github.com/daylioti/docker-commander/ui/helpers"
	"github.com/gizak/termui/v3"
)


// UI main user interface struct.
type UI struct {
	ConfigUi *config.UIConfig

	Cnf          *config.Config
	DockerClient *docker.Docker

	Commands *commands.Commands
	// Flag for switching control between menu and terminal area.
	SelectedMenu byte

	onRender bool
}

// Init initialize all render components.
func (ui *UI) Init(cnf *config.Config, dockerClient *docker.Docker, configUi *config.UIConfig) {

	if ui.Cnf == nil {
		ui.Cnf = cnf
	}
	if ui.ConfigUi == nil {
		ui.ConfigUi = configUi
	}
    if ui.Commands == nil {
		ui.Commands = &commands.Commands{
			ConfigUi:     ui.ConfigUi,
			DockerClient: dockerClient,
			RenderAll:    ui.Render,
			Cnf:          ui.Cnf,
		}
	}
	ui.Commands.TermWidth, ui.Commands.TermHeight = termui.TerminalDimensions()
	ui.Commands.Init()
	termui.StyleParserColorMap = helpers.GetAllTermColors()
	ui.Render()
}

// Handle keyboard keys.
func (ui *UI) Handle(key string) {
	ui.Commands.Handle(key)
}

// Render main render function, that render all ui.
func (ui *UI) Render() {
	termui.Clear()
	ui.Commands.Render()
}
