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

func (t *TerminalUI) TerminalUpdate(id string, finished bool) {
	term := t.client.Exec.GetTerminal(id)
	if term != nil {
		if finished {
			term.Running = false
			term.TabItem.Style = termui.NewStyle(termui.ColorRed)
		}
		if term.Active {
			t.ui.Render()
		}

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
	case "<Right>", "L", "l":
		t.TabPane.FocusRight()
		index := t.client.Exec.GetActiveTerminalIndex()
		if index >= 0 && index < len(t.client.Exec.Terminals)-1 {
			t.SwitchTerminal(t.client.Exec.Terminals[index+1].ID)
		}
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
		} else {
			term.Active = false
		}
	}
	t.UpdateRunningStatus()
	t.ui.Render()
}

func (t *TerminalUI) UpdateRunningStatus() {
	for i := 0; i < len(t.client.Exec.Terminals); i++ {
		if t.client.Exec.Terminals[i].Active {
			t.client.Exec.Terminals[i].TabItem.Style = termui.NewStyle(termui.ColorWhite, termui.ColorGreen)
		} else if t.client.Exec.Terminals[i].Running {
			t.client.Exec.Terminals[i].TabItem.Style = termui.NewStyle(termui.ColorGreen)
		} else {
			t.client.Exec.Terminals[i].TabItem.Style = termui.NewStyle(termui.ColorRed)
		}
	}
}

func (t *TerminalUI) GetTabPaneItems() []commanderWidgets.TabItem {
	var items []commanderWidgets.TabItem
	for i := 0; i < len(t.client.Exec.Terminals); i++ {
		items = append(items, t.client.Exec.Terminals[i].TabItem)
	}
	return items
}

func (t *TerminalUI) Execute(term *docker.TerminalRun) {
	t.client.Exec.Terminals = append(t.client.Exec.Terminals, term)
	t.client.Exec.CommandRun(term)
	t.ui.Render()
}

func (t *TerminalUI) NewTerminal(config config.Config, id string) *docker.TerminalRun {
	return &docker.TerminalRun{
		TabItem: commanderWidgets.TabItem{
			Name:  config.Name,
			Style: termui.NewStyle(termui.ColorGreen),
		},
		List:        widgets.NewList(),
		Active:      true,
		Running:     true,
		ContainerID: t.client.Exec.GetContainerID(config),
		Name:        config.Name,
		ID:          id,
		WorkDir:     config.Exec.WorkingDir,
		Command:     config.Exec.Cmd,
	}
}
