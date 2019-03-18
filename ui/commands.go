package ui

import (
	"github.com/daylioti/docker-commander/config"
	"github.com/daylioti/docker-commander/docker"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Commands struct {
	ui           *UI
	cnf          *config.Config
	lists        []MenuListsItems
	dockerClient *docker.Docker
}

type MenuListsItems struct {
	list  *widgets.List
	ratio float64
}

func (cmd *Commands) Init(cnf *config.Config, dockerClient *docker.Docker, ui *UI) {
	cmd.cnf = cnf
	cmd.dockerClient = dockerClient
	cmd.ui = ui

	cmd.updateStatus(cmd.cnf, cmd.Path(cmd.cnf))

	cmd.UpdateRenderElements(cnf)
}

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
			id := cmd.ui.Term.GetIDFromPath(cmd.Path(cmd.cnf))
			term := cmd.ui.Term.NewTerminal(selected, id)
			cmd.ui.Term.Execute(term)
		}
	}
}

func (cmd *Commands) SetDockerClient(client *docker.Docker) {
	cmd.dockerClient = client
}

func (cmd *Commands) changeSelected(x int, y int, cnf *config.Config) {
	path := cmd.Path(cmd.cnf)
	var c *config.Config
	var cp *config.Config
	var cu *config.Config
	var cd *config.Config
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
	if c.Name == "" {
		return
	}
	if x == 1 && c.Config != nil && len(c.Config) >= 0 {
		c.Selected = false
		c.Config[0].Selected = true
	} else if x == -1 && cp != nil && cp.Name != "" {
		c.Selected = false
		cp.Selected = true
	} else if y == 1 && cd != nil {
		c.Selected = false
		cd.Selected = true
	} else if y == -1 && cu != nil {
		c.Selected = false
		cu.Selected = true
	}
	cmd.updateStatus(cnf, cmd.Path(cnf))
	cmd.UpdateRenderElements(cmd.cnf)

	cmd.ui.Render()
}

func (cmd *Commands) getSelected() config.Config {
	c := cmd.cnf
	for _, path := range cmd.Path(cmd.cnf) {
		c = &c.Config[path]
	}
	return *c
}

// Get array with path to selected item.
func (cmd *Commands) Path(cnf *config.Config) []int {
	var path []int
	var p []int
	cmd.getSelectedPath(&path, cnf)
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] > 0 {
			p = append(p, int(path[i])-1)
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

func (cmd *Commands) UpdateRenderElements(c *config.Config) {
	var width int
	var menuList *MenuListsItems
	path := append(cmd.Path(cmd.cnf), 0)
	cmd.lists = nil
	for i, pathIndex := range path {
		if len(c.Config) <= pathIndex {
			break
		}
		width = 0
		cmd.lists = append(cmd.lists, MenuListsItems{
			list: widgets.NewList(),
		})
		menuList = &cmd.lists[i]
		cmd.lists[i].list.SelectedRow = 0
		for p, cnf := range c.Config {
			if cnf.Selected {
				menuList.list.SelectedRow = len(menuList.list.Rows)
				menuList.list.BorderStyle = termui.NewStyle(termui.ColorGreen)
				menuList.list.SelectedRowStyle = termui.NewStyle(termui.ColorWhite, termui.ColorGreen)
			} else if !cnf.Selected && len(path) > i && pathIndex == p && i < len(path)-1 {
				menuList.list.SelectedRow = len(menuList.list.Rows)
				menuList.list.SelectedRowStyle = termui.NewStyle(termui.ColorGreen)
			}
			menuList.list.Rows = append(menuList.list.Rows, cnf.Name)
			if len(cnf.Name) > width {
				width = len(cnf.Name)
			}
		}
		menuList.ratio = cmd.getMenuColumnRatio(cmd.ui.widthDimension, width)
		c = &c.Config[pathIndex]
	}
}

func (cmd *Commands) GetLists() []MenuListsItems {
	return cmd.lists
}

func (cmd *Commands) getMenuColumnRatio(widthDimension int, maxTextLength int) float64 {
	return float64(maxTextLength+2/(widthDimension/100)) / 100
}
