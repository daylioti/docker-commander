package commands

import (
	"github.com/daylioti/docker-commander/config"
	"github.com/daylioti/docker-commander/docker"
	"github.com/daylioti/docker-commander/ui/render_lock"
)

const (
	KeySelectedCommands = iota
	KeySelectedTerminal
)

type Commands struct {
	ConfigUi   *config.UIConfig
	TermWidth  int
	TermHeight int

	DockerClient *docker.Docker
	Cnf          *config.Config
	Input        *Input
	RenderAll    func()
	SelectedArea byte

	Menu     *Menu
	Terminal *Terminal
}

func (cmd *Commands) Init() {
	if cmd.Menu == nil {
		cmd.Menu = &Menu{
			DockerClient: cmd.DockerClient,
			Commands:     cmd,
		}
	}
	cmd.Menu.Init()
    if cmd.Terminal == nil {
		cmd.Terminal = &Terminal{
			Commands: cmd,
		}
	}
	cmd.Terminal.Init()
    if cmd.Input == nil {
		cmd.Input = &Input{
			Commands: cmd,
		}
	}
}

// Handle keyboard keys.
func (cmd *Commands) Handle(key string) {
	if len(cmd.Input.Fields) > 0 {
		cmd.Input.Handle(key)
		return
	}
	switch key {
	case "<Tab>":
		cmd.HandleSwitchMenu()
	default:
		cmd.HandleKeys(key)
	}
}

// HandleSwitchMenu - switch menu
func (cmd *Commands) HandleSwitchMenu() {
	switch cmd.SelectedArea {

	case KeySelectedCommands:
		cmd.Terminal.Focus()
		cmd.Menu.UnFocus()
		cmd.SelectedArea = KeySelectedTerminal
	case KeySelectedTerminal:
		cmd.Terminal.UnFocus()
		cmd.Menu.Focus()
		cmd.SelectedArea = KeySelectedCommands
	}
	cmd.RenderAll()
}

// HandleKeys - handle other keys.
func (cmd *Commands) HandleKeys(key string) {
	switch cmd.SelectedArea {
	case KeySelectedCommands:
		cmd.Menu.Handle(key)
	case KeySelectedTerminal:
		cmd.Terminal.Handle(key)
	}
	cmd.RenderAll()
}

// Render main render function, that render commands ui.
func (cmd *Commands) Render() {
	cmd.Menu.Render()
	cmd.Terminal.Render()
	if len(cmd.Input.Fields) > 0 {
		for _, field := range cmd.Input.Fields {
			render_lock.RenderLock(field)
		}
	}
}
