package commands

import (
	"github.com/daylioti/docker-commander/config"
	"github.com/daylioti/docker-commander/docker"
	"github.com/daylioti/docker-commander/ui/render_lock"
	commanderWidgets "github.com/daylioti/docker-commander/ui/widgets"
	"github.com/gizak/termui/v3"
	"strconv"
)

// TerminalUI UI struct.
type Terminal struct {
	//ui              *UI
	TabPane         *commanderWidgets.TabsPaneStyled
	DisplayTerminal *commanderWidgets.TerminalList
	Commands        *Commands
}

// Init initialize terminal render component.
func (t *Terminal) Init() {
	t.DisplayTerminal = t.InitDisplayTerminal()

	t.TabPane = commanderWidgets.NewTabPaneStyled()
	t.TabPane.SetRect(0, t.Commands.ConfigUi.GetCommandsHeight(),
		t.Commands.TermWidth, t.Commands.ConfigUi.GetCommandsHeight()+3)

	t.Commands.DockerClient.Exec.SetTerminalUpdateFn(t.TerminalUpdate)
}

// InitDisplayTerminal initialize command output list.
func (t *Terminal) InitDisplayTerminal() *commanderWidgets.TerminalList {
	list := commanderWidgets.NewTerminalList()
	list.SetRect(0, t.Commands.ConfigUi.GetCommandsHeight()+3,
		t.Commands.TermWidth, t.Commands.TermHeight)
	return list
}

// Render function, that render terminal component.
func (t *Terminal) Render() {
	t.TabPaneRender()
	t.DisplayTerminalRender()
}

// TabPaneRender render only tab.
func (t *Terminal) TabPaneRender() {
	render_lock.RenderLock(t.TabPane)
}

// DisplayTerminalRender render only terminal.
func (t *Terminal) DisplayTerminalRender() {
	render_lock.RenderLock(t.DisplayTerminal)
}

// TerminalUpdate calls on receive updates from docker process.
func (t *Terminal) TerminalUpdate(term *docker.TerminalRun, finished bool) {
	term.List.SelectedRow = len(term.List.Rows)
	if finished {
		term.Running = false
		term.TabItem.Style = termui.NewStyle(termui.ColorRed)
		t.Commands.RenderAll()
	}
	if term.Active {
		t.DisplayTerminalRender()
	}
}

// Handle keyboard keys.
func (t *Terminal) Handle(key string) {
	switch key {
	case "<Up>", "K", "k", "<MouseWheelUp>":
		if len(t.DisplayTerminal.Rows) > 0 {
			t.DisplayTerminal.ScrollUp()
			t.DisplayTerminalRender()
		}
	case "<Down>", "J", "j", "<MouseWheelDown>":
		if len(t.DisplayTerminal.Rows) > 0 {
			t.DisplayTerminal.ScrollDown()
			t.DisplayTerminalRender()
		}
	case "<PageUp>":
		if len(t.DisplayTerminal.Rows) > 0 {
			t.DisplayTerminal.ScrollPageUp()
			t.DisplayTerminalRender()
		}
	case "<PageDown>":
		if len(t.DisplayTerminal.Rows) > 0 {
			t.DisplayTerminal.ScrollPageDown()
			t.DisplayTerminalRender()
		}
	case "<Home>":
		if len(t.DisplayTerminal.Rows) > 0 {
			t.DisplayTerminal.SelectedRow = 0
			t.DisplayTerminalRender()
		}
	case "<End>":
		if len(t.DisplayTerminal.Rows) > 0 {
			t.DisplayTerminal.SelectedRow = len(t.DisplayTerminal.Rows) - 1
			t.DisplayTerminalRender()
		}
	case "<Left>", "H", "h":
		t.TabPane.FocusLeft()
		index := t.Commands.DockerClient.Exec.GetActiveTerminalIndex()
		if index > 0 {
			t.SwitchTerminal(t.Commands.DockerClient.Exec.Terminals[index-1])
		}
	case "<Right>", "L", "l":
		t.TabPane.FocusRight()
		index := t.Commands.DockerClient.Exec.GetActiveTerminalIndex()
		if index >= 0 && index < len(t.Commands.DockerClient.Exec.Terminals)-1 {
			t.SwitchTerminal(t.Commands.DockerClient.Exec.Terminals[index+1])
		}
	case "<C-r>":
		index := t.Commands.DockerClient.Exec.GetActiveTerminalIndex()
		t.DisplayTerminal = t.InitDisplayTerminal()
		t.Commands.DockerClient.Exec.Terminals[index] = t.Commands.DockerClient.Exec.Terminals[len(t.Commands.DockerClient.Exec.Terminals)-1]
		t.Commands.DockerClient.Exec.Terminals = t.Commands.DockerClient.Exec.Terminals[:len(t.Commands.DockerClient.Exec.Terminals)-1]
		if len(t.Commands.DockerClient.Exec.Terminals) > 0 {
			t.Commands.DockerClient.Exec.Terminals[0].Active = true
		}
		t.UpdateRunningStatus()
		t.Commands.RenderAll()
	}
}

// GetIDFromPath get id from selected commands list item.
func (t *Terminal) GetIDFromPath(path []int) string {
	id := "0"
	for _, i := range path {
		id += strconv.Itoa(i)
	}
	return id
}

// SwitchTerminal set selected terminal with id.
func (t *Terminal) SwitchTerminal(term *docker.TerminalRun) {
	t.unActivateTerminals()
	term.Active = true
	t.DisplayTerminal = term.List
	t.UpdateRunningStatus()
	t.Commands.Terminal.DisplayTerminal.BorderStyle = termui.NewStyle(termui.ColorGreen)
	t.Commands.RenderAll()
}

// Focus commands lists, set borders.
func (t *Terminal) Focus() {
	t.TabPane.BorderStyle = termui.NewStyle(termui.ColorGreen)
	t.DisplayTerminal.BorderStyle = termui.NewStyle(termui.ColorGreen)
	t.UpdateRunningStatus()
	t.Commands.RenderAll()
}

// UnFocus commands lists, remove borders.
func (t *Terminal) UnFocus() {
	t.UpdateRunningStatus()
	t.TabPane.BorderStyle = termui.NewStyle(termui.ColorWhite)
	t.DisplayTerminal.BorderStyle = termui.NewStyle(termui.ColorWhite)
	t.Commands.Render()
}

// SetDisplayTerminal change display terminal.
func (t *Terminal) SetDisplayTerminal(term *commanderWidgets.TerminalList) {
	t.DisplayTerminal = term
}

// UpdateRunningStatus change tabs colors depends on selected tab, running or not process in docker.
func (t *Terminal) UpdateRunningStatus() {
	var terminal *docker.TerminalRun
	t.TabPane.TabNames = nil
	for i := 0; i < len(t.Commands.DockerClient.Exec.Terminals); i++ {
		terminal = t.Commands.DockerClient.Exec.Terminals[i]
		if terminal.Active {
			if terminal.Running {
				terminal.TabItem.Style = termui.NewStyle(termui.ColorBlack, termui.ColorGreen)
			} else {
				terminal.TabItem.Style = termui.NewStyle(termui.ColorBlack, termui.ColorRed)
			}
			t.SetDisplayTerminal(terminal.List)
		} else if terminal.Running {
			terminal.TabItem.Style = termui.NewStyle(termui.ColorGreen)
		} else {
			terminal.TabItem.Style = termui.NewStyle(termui.ColorRed)
		}
		t.TabPane.TabNames = append(t.TabPane.TabNames, terminal.TabItem)
	}
}

// Execute start new process in docker.
func (t *Terminal) Execute(term *docker.TerminalRun) {
	t.Commands.DockerClient.Exec.Terminals = append(t.Commands.DockerClient.Exec.Terminals, term)
	t.Commands.DockerClient.Exec.CommandRun(term)
	t.SwitchTerminal(term)
}

// removeFinishedTerminals remove first terminal object if finished and length of names bigger that tab width.
func (t *Terminal) removeFinishedTerminals() {
	var tabItemsLength int
	tabBorder := 3
	for _, term := range t.Commands.DockerClient.Exec.Terminals {
		tabItemsLength += len(term.TabItem.Name) + tabBorder*2
	}
	if tabItemsLength-tabBorder*2 >= t.Commands.TermWidth-2 {
		for i, term := range t.Commands.DockerClient.Exec.Terminals {
			if !term.Running {
				t.Commands.DockerClient.Exec.Terminals = append(t.Commands.DockerClient.Exec.Terminals[:i],
					t.Commands.DockerClient.Exec.Terminals[i+1:]...)
				t.Commands.Render()
				return
			}
		}
	}
}

// unActivateTerminals set active to false on all terminals.
func (t *Terminal) unActivateTerminals() {
	for _, term := range t.Commands.DockerClient.Exec.Terminals {
		term.Active = false
	}
}

// NewTerminal return new terminal object,
func (t *Terminal) NewTerminal(config config.Config, id string) *docker.TerminalRun {
	list := commanderWidgets.NewTerminalList()
	list.SelectedRowStyle = termui.NewStyle(termui.ColorBlack, termui.ColorGreen)
	list.SetRect(0, t.Commands.ConfigUi.GetCommandsHeight()+3,
		t.Commands.TermWidth, t.Commands.TermHeight)
	t.removeFinishedTerminals()
	t.unActivateTerminals()
	return &docker.TerminalRun{
		TabItem: &commanderWidgets.TabItem{
			Name:  config.Name,
			Style: termui.NewStyle(termui.ColorGreen),
		},
		List:        list,
		Active:      true,
		Running:     true,
		ContainerID: t.Commands.DockerClient.Exec.GetContainerID(config),
		Name:        config.Name,
		ID:          id,
		WorkDir:     config.Exec.WorkingDir,
		Command:     config.Exec.Cmd,
	}
}
