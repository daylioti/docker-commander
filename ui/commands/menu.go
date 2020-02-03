package commands

import (
	"github.com/daylioti/docker-commander/config"
	"github.com/daylioti/docker-commander/docker"
	"github.com/daylioti/docker-commander/ui/render_lock"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"strconv"
)

// Commands UI struct.
type Menu struct {
	DockerClient *docker.Docker
	Commands     *Commands
	Lists        []*widgets.List
}

// Init initialize commands render component.
func (m *Menu) Init() {
	m.updateStatus(m.Commands.Cnf, m.Path(m.Commands.Cnf))
	m.UpdateRenderElements(m.Commands.Cnf)
}

// Handle keyboard keys.
func (m *Menu) Handle(key string) {
	switch key {
	case "<Up>", "K", "k":
		m.changeSelected(0, -1, m.Commands.Cnf)
	case "<Left>", "H", "h":
		m.changeSelected(-1, 0, m.Commands.Cnf)
	case "<Right>", "L", "l":
		m.changeSelected(1, 0, m.Commands.Cnf)
	case "<Down>", "J", "j":
		m.changeSelected(0, 1, m.Commands.Cnf)
	case "<Enter>":
		selected := m.getSelected()
		if selected.Exec.Cmd != "" {
			m.ExecuteSelectedCommand(selected)
		}
	}
}

// Render function, that render commands component.
func (m *Menu) Render() {
	for listIndex := 0; listIndex < len(m.Lists); listIndex++ {
		render_lock.RenderLock(m.Lists[listIndex])
	}
}

// Focus commands lists, set borders.
func (m *Menu) Focus() {
	m.UpdateRenderElements(m.Commands.Cnf)
	m.Render()
}

// UnFocus commands lists, remove borders.
func (m *Menu) UnFocus() {
	for _, list := range m.Lists {
		list.BorderStyle = termui.NewStyle(termui.ColorWhite)
	}
	m.Render()
}

// ExecuteSelectedCommand start execution process, open input if needed.
func (m *Menu) ExecuteSelectedCommand(cnf config.Config) {
	if len(cnf.Exec.Input) > 0 {
		// Wait for input fields.
		m.Commands.Input.NewInputs(cnf.Exec.Input, cnf)
	} else {
		placeholders := m.Commands.Cnf.GetPlaceholders(m.Path(m.Commands.Cnf), make(map[string]string), m.Commands.Cnf)
		m.Commands.Cnf.ReplacePlaceholders(placeholders, &cnf)
		m.commandExecProcess(cnf)
	}
}

// commandExecProcess execute command in docker.
func (m *Menu) commandExecProcess(cnf config.Config) {
	cnf.Exec.Input = nil
	id := m.Commands.Terminal.GetIDFromPath(m.Path(m.Commands.Cnf))
	term := m.Commands.Terminal.NewTerminal(cnf, id)
	m.Commands.Terminal.Execute(term)
}

// SetDockerClient
func (m *Menu) SetDockerClient(client *docker.Docker) {
	m.DockerClient = client
}

// getNearestConfigs return nearest from selected configs.
func (m *Menu) getNearestConfigs() (*config.Config, *config.Config, *config.Config, *config.Config) {
	cnf := m.Commands.Cnf
	path := m.Path(m.Commands.Cnf)
	var c, cp, cu, cd *config.Config
	cnf.Selected = false
	c = cnf
	for i := 0; i <= len(path); i++ {
		if i == len(path)-1 {
			cp = c
			if len(cp.Config)-1 >= path[i]+1 {
				cd = &cp.Config[path[i]+1]
			}
			if len(cp.Config) > 0 && path[i]-1 >= 0 {
				cu = &cp.Config[path[i]-1]
			}
			c = &cp.Config[path[i]]
			break
		} else if len(path) >= i-1 && len(c.Config) >= path[i] {
			c = &c.Config[path[i]]
		}
	}
	return c, cp, cu, cd
}

// changeCommandsSelected change selected in cnf struct, depends on current selected item.
// return bool - requires re-render or not.
func (m *Menu) changeCommandsSelected(x, y int, c, cp, cu, cd *config.Config) bool {
	switch {
	case x == 1 && c.Config != nil:
		c.Selected = false
		c.Config[0].Selected = true
		return true
	case x == -1 && cp != nil && cp.Name != "":
		c.Selected = false
		cp.Selected = true
		return true
	case y == 1 && cd != nil:
		c.Selected = false
		cd.Selected = true
		if cd.Config != nil || c.Config != nil {
			return true
		}
	case y == -1 && cu != nil:
		c.Selected = false
		cu.Selected = true
		if cu.Config != nil || c.Config != nil {
			return true
		}
	}
	return false
}

// changeSelected change selected, depends on current selected item.
func (m *Menu) changeSelected(x int, y int, cnf *config.Config) {
	c, cp, cu, cd := m.getNearestConfigs()
	if c.Name == "" {
		return
	}
	clear := m.changeCommandsSelected(x, y, c, cp, cu, cd)
	m.updateStatus(cnf, m.Path(cnf))
	m.UpdateRenderElements(m.Commands.Cnf)
	if !clear {
		// No new elements and elements to remove.
		// Just render command lists
		m.Render()
	} else {
		// Re-render all.
		termui.Clear()
		m.Commands.RenderAll()
	}
}

// getSelected
func (m *Menu) getSelected() config.Config {
	c := m.Commands.Cnf
	for _, path := range m.Path(m.Commands.Cnf) {
		c = &c.Config[path]
	}
	return *c
}

// Path get array with path to selected item.
func (m *Menu) Path(cnf *config.Config) []int {
	var path []int
	var p []int
	m.getSelectedPath(&path, cnf)
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] > 0 {
			p = append(p, path[i]-1)
		}
	}
	return p
}

// Update Status (display flag) prop in commands structure, depends on current selected item.
func (m *Menu) updateStatus(cnf *config.Config, path []int) {
	if len(path) == 0 {
		return
	}
	cnf = &cnf.Config[path[0]]

	if len(cnf.Config) > 0 {
		for i := 0; i < len(cnf.Config); i++ {
			cnf.Config[i].Status = true
		}
	}
	if len(path) == 1 {
		cnf.Selected = true
	} else {
		m.updateStatus(cnf, path[1:])
	}
}

// Get selected item path via int array
func (m *Menu) getSelectedPath(path *[]int, cnf *config.Config) bool {
	if cnf.Selected {
		return true
	}
	for i := 0; i < len(cnf.Config); i++ {
		if m.getSelectedPath(path, &cnf.Config[i]) {
			*path = append(*path, i+1)
			return true
		}
	}
	return false
}

// UpdateRenderElements update lists from config.
func (m *Menu) UpdateRenderElements(c *config.Config) {
	var width, height int
	if h, exist := m.Commands.ConfigUi.UI.Commands["height"]; !exist {
		height = 5
	} else {
		height, _ = strconv.Atoi(h)
	}

	var menuList *widgets.List
	borderSize := 2
	widthPrev := 0
	path := append(m.Path(m.Commands.Cnf), 0)
	m.Lists = nil
	for i, pathIndex := range path {
		if len(c.Config) <= pathIndex {
			break
		}
		width = 0
		m.Lists = append(m.Lists, widgets.NewList())
		menuList = m.Lists[i]
		menuList.SelectedRow = 0
		for p, cnf := range c.Config {
			if cnf.Selected {
				menuList.SelectedRow = len(menuList.Rows)
				menuList.BorderStyle = termui.NewStyle(termui.ColorGreen)
				menuList.SelectedRowStyle = termui.NewStyle(termui.ColorBlack, termui.ColorGreen)
			} else if !cnf.Selected && len(path) > i && pathIndex == p && i < len(path)-1 {
				menuList.SelectedRow = len(menuList.Rows)
				menuList.SelectedRowStyle = termui.NewStyle(termui.ColorGreen)
			}
			menuList.Rows = append(menuList.Rows, cnf.Name)
			if len(cnf.Name) > width {
				width = len(cnf.Name)
			}
		}
		menuList.Border = true
		width += borderSize
		menuList.SetRect(widthPrev, 0, widthPrev+width, height+borderSize)
		widthPrev += width
		c = &c.Config[pathIndex]
	}
}
