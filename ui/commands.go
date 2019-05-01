package ui

import (
	"github.com/daylioti/docker-commander/config"
	"github.com/daylioti/docker-commander/docker"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"strconv"
)

// Commands UI struct.
type Commands struct {
	ui           *UI
	cnf          *config.Config
	dockerClient *docker.Docker
	Lists        []*widgets.List
}

// Init initialize commands render component.
func (cmd *Commands) Init(cnf *config.Config, dockerClient *docker.Docker, ui *UI) {
	cmd.cnf = cnf
	cmd.dockerClient = dockerClient
	cmd.ui = ui

	cmd.updateStatus(cmd.cnf, cmd.Path(cmd.cnf))

	cmd.UpdateRenderElements(cnf)
}

// Handle keyboard keys.
func (cmd *Commands) Handle(key string) {
	switch key {
	case "<Up>", "K", "k":
		cmd.changeSelected(0, -1, cmd.cnf)
	case "<Left>", "H", "h":
		cmd.changeSelected(-1, 0, cmd.cnf)
	case "<Right>", "L", "l":
		cmd.changeSelected(1, 0, cmd.cnf)
	case "<Down>", "J", "j":
		cmd.changeSelected(0, 1, cmd.cnf)
	case "<Enter>":
		selected := cmd.getSelected()
		if selected.Exec.Cmd != "" {
			cmd.ExecuteSelectedCommand(selected)
		}
	}
}

// Render function, that render commands component.
func (cmd *Commands) Render() {
	for listIndex := 0; listIndex < len(cmd.Lists); listIndex++ {
		termui.Render(cmd.Lists[listIndex])
	}
}

// Focus commands lists, set borders.
func (cmd *Commands) Focus() {
	cmd.UpdateRenderElements(cmd.cnf)
	cmd.Render()
}

// UnFocus commands lists, remove borders.
func (cmd *Commands) UnFocus() {
	for _, list := range cmd.Lists {
		list.BorderStyle = termui.NewStyle(termui.ColorWhite)
	}
	cmd.Render()
}

// ExecuteSelectedCommand start execution process, open input if needed.
func (cmd *Commands) ExecuteSelectedCommand(cnf config.Config) {
	if len(cnf.Exec.Input) > 0 {
		// Wait for input fields.
		cn := make(chan map[string]string)
		cmd.ui.Input.NewInputs(cnf.Exec.Input, &cn)
		go func() {
			for k, v := range <-cn {
				cnf.ReplacePlaceholder(k, v, &cnf)
			}
			cmd.commandExecProcess(cnf)
		}()
	} else {
		cmd.commandExecProcess(cnf)
	}
}

// commandExecProcess execute command in docker.
func (cmd *Commands) commandExecProcess(cnf config.Config) {
	cnf.Exec.Input = nil
	id := cmd.ui.Term.GetIDFromPath(cmd.Path(cmd.cnf))
	term := cmd.ui.Term.NewTerminal(cnf, id)
	cmd.ui.Term.Execute(term)
}

// SetDockerClient
func (cmd *Commands) SetDockerClient(client *docker.Docker) {
	cmd.dockerClient = client
}

// getNearestConfigs return nearest from selected configs.
func (cmd *Commands) getNearestConfigs() (*config.Config, *config.Config, *config.Config, *config.Config) {
	cnf := cmd.cnf
	path := cmd.Path(cmd.cnf)
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
func (cmd *Commands) changeCommandsSelected(x, y int, c, cp, cu, cd *config.Config) bool {
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
func (cmd *Commands) changeSelected(x int, y int, cnf *config.Config) {
	c, cp, cu, cd := cmd.getNearestConfigs()
	if c.Name == "" {
		return
	}

	clear := cmd.changeCommandsSelected(x, y, c, cp, cu, cd)

	cmd.updateStatus(cnf, cmd.Path(cnf))
	cmd.UpdateRenderElements(cmd.cnf)
	if !clear {
		// No new elements and elements to remove.
		// Just render command lists
		cmd.Render()
	} else {
		// Re-render all.
		termui.Clear()
		cmd.ui.Render()
	}
}

// getSelected
func (cmd *Commands) getSelected() config.Config {
	c := cmd.cnf
	for _, path := range cmd.Path(cmd.cnf) {
		c = &c.Config[path]
	}
	return *c
}

// Path get array with path to selected item.
func (cmd *Commands) Path(cnf *config.Config) []int {
	var path []int
	var p []int
	cmd.getSelectedPath(&path, cnf)
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] > 0 {
			p = append(p, path[i]-1)
		}
	}
	return p
}

// Update Status (display flag) prop in commands structure, depends on current selected item.
func (cmd *Commands) updateStatus(cnf *config.Config, path []int) {
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
		cmd.updateStatus(cnf, path[1:])
	}
}

// Get selected item path via int array
func (cmd *Commands) getSelectedPath(path *[]int, cnf *config.Config) bool {
	if cnf.Selected {
		return true
	}
	for i := 0; i < len(cnf.Config); i++ {
		if cmd.getSelectedPath(path, &cnf.Config[i]) {
			*path = append(*path, i+1)
			return true
		}
	}
	return false
}

// UpdateRenderElements update lists from config.
func (cmd *Commands) UpdateRenderElements(c *config.Config) {
	var width, height int
	if h, exist := cmd.ui.configUi.UI.Commands["height"]; !exist {
		height = 5
	} else {
		height, _ = strconv.Atoi(h)
	}

	var menuList *widgets.List
	borderSize := 2
	widthPrev := 0
	path := append(cmd.Path(cmd.cnf), 0)
	cmd.Lists = nil
	for i, pathIndex := range path {
		if len(c.Config) <= pathIndex {
			break
		}
		width = 0
		cmd.Lists = append(cmd.Lists, widgets.NewList())
		menuList = cmd.Lists[i]
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
