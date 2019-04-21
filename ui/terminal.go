package ui

import (
	"github.com/daylioti/docker-commander/config"
	"github.com/daylioti/docker-commander/docker"
	commanderWidgets "github.com/daylioti/docker-commander/ui/widgets"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"strconv"
)

// TerminalUI UI struct.
type TerminalUI struct {
	ui              *UI
	client          *docker.Docker
	TabPane         *commanderWidgets.TabsPaneStyled
	DisplayTerminal *widgets.List
}

// Init initialize terminal render component.
func (t *TerminalUI) Init(ui *UI, client *docker.Docker) {
	t.ui = ui
	t.client = client
	t.DisplayTerminal = t.InitDisplayTerminal()

	t.TabPane = commanderWidgets.NewTabPaneStyled()
	t.TabPane.SetRect(0, t.ui.configUi.GetCommandsHeight(), t.ui.TermWidth, t.ui.configUi.GetCommandsHeight()+3)
	t.TabPane.Border = true

	t.client.Exec.SetTerminalUpdateFn(t.TerminalUpdate)
}

func (t *TerminalUI) InitDisplayTerminal() *widgets.List {
	list := widgets.NewList()
	list.SetRect(0, t.ui.configUi.GetCommandsHeight()+3, t.ui.TermWidth, t.ui.TermHeight)
	return list
}

// Render function, that render terminal component.
func (t *TerminalUI) Render() {
	t.TabPaneRender()
	t.DisplayTerminalRender()
}

// TabPaneRender render only tab.
func (t *TerminalUI) TabPaneRender() {
	termui.Render(t.TabPane)
}

// DisplayTerminalRender render only terminal.
func (t *TerminalUI) DisplayTerminalRender() {
	termui.Render(t.DisplayTerminal)
}

// TerminalUpdate calls on receive updates from docker process.
func (t *TerminalUI) TerminalUpdate(term *docker.TerminalRun, finished bool) {
	if finished {
		term.Running = false
		term.TabItem.Style = termui.NewStyle(termui.ColorRed)
		t.ui.Render()
	}
	term.List.SelectedRow = len(term.List.Rows)
	if term.Active {
		t.DisplayTerminalRender()
	}
}

// Handle keyboard keys.
func (t *TerminalUI) Handle(key string) {
	switch key {
	case "<Up>", "K", "k":
		if len(t.DisplayTerminal.Rows) > 0 {
			t.DisplayTerminal.ScrollUp()
			t.DisplayTerminalRender()
		}
	case "<Down>", "J", "j":
		if len(t.DisplayTerminal.Rows) > 0 {
			t.DisplayTerminal.ScrollDown()
			t.DisplayTerminalRender()
		}
	case "<Left>", "H", "h":
		t.TabPane.FocusLeft()
		index := t.client.Exec.GetActiveTerminalIndex()
		if index > 0 {
			t.SwitchTerminal(t.client.Exec.Terminals[index-1])
		}
	case "<Right>", "L", "l":
		t.TabPane.FocusRight()
		index := t.client.Exec.GetActiveTerminalIndex()
		if index >= 0 && index < len(t.client.Exec.Terminals)-1 {
			t.SwitchTerminal(t.client.Exec.Terminals[index+1])
		}
	case "<C-r>":
		index := t.client.Exec.GetActiveTerminalIndex()
		t.DisplayTerminal = t.InitDisplayTerminal()
		t.client.Exec.Terminals[index] = t.client.Exec.Terminals[len(t.client.Exec.Terminals)-1]
		t.client.Exec.Terminals = t.client.Exec.Terminals[:len(t.client.Exec.Terminals)-1]
		if len(t.client.Exec.Terminals) > 0 {
			t.client.Exec.Terminals[0].Active = true
		}
		t.UpdateRunningStatus()
		t.ui.Render()
	}
}

// GetIDFromPath get id from selected commands list item.
func (t *TerminalUI) GetIDFromPath(path []int) string {
	id := "0"
	for _, i := range path {
		id += strconv.Itoa(i)
	}
	return id
}

// SwitchTerminal set selected terminal with id.
func (t *TerminalUI) SwitchTerminal(term *docker.TerminalRun) {
	t.unActivateTerminals()
	term.Active = true
	t.DisplayTerminal = term.List
	t.UpdateRunningStatus()
	t.ui.Term.DisplayTerminal.BorderStyle = termui.NewStyle(termui.ColorGreen)
	t.ui.Render()
}

// Focus commands lists, set borders.
func (t *TerminalUI) Focus() {
	t.TabPane.BorderStyle = termui.NewStyle(termui.ColorGreen)
	t.DisplayTerminal.BorderStyle = termui.NewStyle(termui.ColorGreen)
	t.UpdateRunningStatus()
	t.ui.Render()
}

// UnFocus commands lists, remove borders.
func (t *TerminalUI) UnFocus() {
	t.UpdateRunningStatus()
	t.TabPane.BorderStyle = termui.NewStyle(termui.ColorWhite)
	t.DisplayTerminal.BorderStyle = termui.NewStyle(termui.ColorWhite)
	t.ui.Render()
}

// SetDisplayTerminal change display terminal.
func (t *TerminalUI) SetDisplayTerminal(term *widgets.List) {
	t.DisplayTerminal = term
}

// UpdateRunningStatus change tabs colors depends on selected tab, running or not process in docker.
func (t *TerminalUI) UpdateRunningStatus() {
	var terminal *docker.TerminalRun
	t.TabPane.TabNames = nil
	for i := 0; i < len(t.client.Exec.Terminals); i++ {
		terminal = t.client.Exec.Terminals[i]
		if terminal.Active {
			if terminal.Running {
				terminal.TabItem.Style = termui.NewStyle(termui.ColorWhite, termui.ColorGreen)
			} else {
				terminal.TabItem.Style = termui.NewStyle(termui.ColorWhite, termui.ColorRed)
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
func (t *TerminalUI) Execute(term *docker.TerminalRun) {
	t.client.Exec.Terminals = append(t.client.Exec.Terminals, term)
	t.client.Exec.CommandRun(term)
	t.SwitchTerminal(term)
}

// removeFinishedTerminals remove first terminal object if finished and length of names bigger that tab width.
func (t *TerminalUI) removeFinishedTerminals() {
	var tabItemsLength int
	tabBorder := 3
	for _, term := range t.client.Exec.Terminals {
		tabItemsLength += len(term.TabItem.Name) + tabBorder*2
	}
	if tabItemsLength-tabBorder*2 >= t.ui.TermWidth-2 {

		for i, term := range t.client.Exec.Terminals {
			if !term.Running {
				t.client.Exec.Terminals[i] = t.client.Exec.Terminals[len(t.client.Exec.Terminals)-1]
				t.client.Exec.Terminals = t.client.Exec.Terminals[:len(t.client.Exec.Terminals)-1]
				t.ui.Render()
				return
			}
		}
	}
}

// unActivateTerminals set active to false on all terminals.
func (t *TerminalUI) unActivateTerminals() {
	for _, term := range t.client.Exec.Terminals {
		term.Active = false
	}
}

// NewTerminal return new terminal object,
func (t *TerminalUI) NewTerminal(config config.Config, id string) *docker.TerminalRun {
	list := widgets.NewList()
	list.SelectedRowStyle = termui.NewStyle(termui.ColorWhite, termui.ColorGreen)
	list.SetRect(0, t.ui.configUi.GetCommandsHeight()+3, t.ui.TermWidth, t.ui.TermHeight)
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
		ContainerID: t.client.Exec.GetContainerID(config),
		Name:        config.Name,
		ID:          id,
		WorkDir:     config.Exec.WorkingDir,
		Command:     config.Exec.Cmd,
	}
}
