package ui

import (
	"github.com/daylioti/docker-commander/config"
	"github.com/daylioti/docker-commander/docker"
	commanderWidgets "github.com/daylioti/docker-commander/ui/widgets"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"strconv"
)

type TerminalUI struct {
	ui              *UI
	client          *docker.Docker
	TabPane         *commanderWidgets.TabsPaneStyled
	DisplayTerminal *widgets.List
}

func (t *TerminalUI) Init(ui *UI, client *docker.Docker) {
	t.ui = ui
	t.client = client
	t.DisplayTerminal = widgets.NewList()
	t.TabPane = commanderWidgets.NewTabPaneStyled()
	t.TabPane.Border = true
	t.client.Exec.SetTerminalUpdateFn(t.TerminalUpdate)
}

func (t *TerminalUI) TerminalUpdate(term *docker.TerminalRun, finished bool) {
	if finished {
		term.Running = false
		term.TabItem.Style = termui.NewStyle(termui.ColorRed)
		//t.UpdateRunningStatus()
		t.ui.Render()
	}
	if term.Active {
		t.ui.Render()
	}
}

func (t *TerminalUI) Handle(key string) {
	switch key {
	case "<Up>", "K", "k":
		if len(t.DisplayTerminal.Rows) > 0 {
			t.DisplayTerminal.ScrollUp()
			t.ui.Render()
		}
	case "<Down>", "J", "j":
		if len(t.DisplayTerminal.Rows) > 0 {
			t.DisplayTerminal.ScrollDown()
			t.ui.Render()
		}
	case "<Left>", "H", "h":
		t.TabPane.FocusLeft()
		index := t.client.Exec.GetActiveTerminalIndex()
		if index > 0 {
			t.SwitchTerminal(t.client.Exec.Terminals[index-1].ID)
		}
		t.ui.Render()
	case "<Right>", "L", "l":
		t.TabPane.FocusRight()
		index := t.client.Exec.GetActiveTerminalIndex()
		if index >= 0 && index < len(t.client.Exec.Terminals)-1 {
			t.SwitchTerminal(t.client.Exec.Terminals[index+1].ID)
		}
		t.ui.Render()
	}
}

func (t *TerminalUI) GetIDFromPath(path []int) string {
	id := "0"
	for _, i := range path {
		id += strconv.Itoa(i)
	}
	return id
}

func (t *TerminalUI) SwitchTerminal(id string) {
	for _, term := range t.client.Exec.Terminals {
		if term.ID == id {
			term.Active = true
			t.DisplayTerminal = term.List
		} else {
			term.Active = false
		}
	}
	t.UpdateRunningStatus()
	t.ui.Render()
}

func (t *TerminalUI) UpdateRunningStatus() {
	var terminal *docker.TerminalRun
	t.TabPane.TabNames = nil
	for i := 0; i < len(t.client.Exec.Terminals); i++ {
		terminal = t.client.Exec.Terminals[i]
		if terminal.Active {
			terminal.TabItem.Style = termui.NewStyle(termui.ColorWhite, termui.ColorGreen)
			t.DisplayTerminal = terminal.List
		} else if terminal.Running {
			terminal.TabItem.Style = termui.NewStyle(termui.ColorGreen)
		} else {
			terminal.TabItem.Style = termui.NewStyle(termui.ColorRed)
		}
		t.TabPane.TabNames = append(t.TabPane.TabNames, terminal.TabItem)
	}
}

func (t *TerminalUI) Execute(term *docker.TerminalRun) {
	t.client.Exec.Terminals = append(t.client.Exec.Terminals, term)
	t.client.Exec.CommandRun(term)
	t.ui.ClearRender = true
	t.SwitchTerminal(term.ID)
}

func (t *TerminalUI) NewTerminal(config config.Config, id string) *docker.TerminalRun {
	return &docker.TerminalRun{
		TabItem: &commanderWidgets.TabItem{
			Name:  config.Name,
			Style: termui.NewStyle(termui.ColorGreen),
		},
		List:        widgets.NewList(),
		Active:      true,
		Running:     false,
		ContainerID: t.client.Exec.GetContainerID(config),
		Name:        config.Name,
		ID:          id,
		WorkDir:     config.Exec.WorkingDir,
		Command:     config.Exec.Cmd,
	}
}
