package ui

import (
	"github.com/daylioti/docker-commander/config"
	"github.com/daylioti/docker-commander/docker"
	"github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
	"math"
)

type Commands struct {
	ui           *UI
	cnf          *config.Config
	lists        []MenuListsItems
	terminal     *widgets.List
	dockerClient *docker.Docker
}

type MenuListsItems struct {
	list *widgets.List
	ratio float64
}


func (cmd *Commands) Init(cnf *config.Config, dockerClient *docker.Docker, ui *UI) {
	cmd.cnf = cnf
	cmd.dockerClient = dockerClient
	cmd.ui = ui

	cmd.dockerClient.Exec.SetTerminalUpdateFn(cmd.updateTerminal)
	cmd.updateStatus(cmd.cnf, cmd.Path(cmd.cnf))

	cmd.UpdateRenderElements(cnf)
}

func (cmd *Commands) Handle(key string) {
	switch key {
	case "<Up>", "K", "k":
		cmd.changeSelected(0, -1, cmd.cnf)
		break
	case "<Left>", "H", "h":
		cmd.changeSelected(-1, 0, cmd.cnf)
		break
	case "<Right>", "L", "l":
		cmd.changeSelected(1, 0, cmd.cnf)
		break
	case "<Down>", "J", "j":
		cmd.changeSelected(0, 1, cmd.cnf)
		break
	case "<Enter>":
		selected := cmd.getSelected()
		if selected.Exec.Cmd != "" {
			id := GetIdFromPath(cmd.Path(cmd.cnf))
			term := &docker.TerminalRun{
				ContainerName: selected.Exec.Connect.ContainerName,
				ContainerId: selected.Exec.Connect.ContainerId,
				FromImage: selected.Exec.Connect.FromImage,
			    Command: selected.Exec.Cmd,
			    Id: id,
			}
			cmd.ui.Term.Execute(term)
		}
		break
	}
}

func (cmd *Commands) SetDockerClient(client *docker.Docker) {
	cmd.dockerClient = client
}

func (cmd *Commands) updateTerminal() {
	//if cmd.terminal != nil && len(cmd.terminal.Rows) != len(*cmd.dockerClient.Exec.RunningTerminal) {
	//	cmd.terminal.Rows = *cmd.dockerClient.Exec.RunningTerminal
	//	//cmd.TerminalRender()
	//}
}

func (cmd *Commands) changeSelected(x int, y int, cnf *config.Config) {
	if cmd.ui.SelectedRow != 0 {
		return
	}
	path := cmd.Path(cmd.cnf)
	var terminalRender bool
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
		terminalRender = cmd.renderTerminal(&c.Config[0])
	} else if x == -1 && cp != nil && cp.Name != "" {
		c.Selected = false
		cp.Selected = true
		terminalRender = cmd.renderTerminal(cp)
	} else if y == 1 && cd != nil {
		c.Selected = false
		cd.Selected = true
		terminalRender = cmd.renderTerminal(cd)
	} else if y == -1 && cu != nil {
		c.Selected = false
		cu.Selected = true
		terminalRender = cmd.renderTerminal(cu)
	}
	cmd.updateStatus(cnf, cmd.Path(cnf))
	cmd.UpdateRenderElements(cmd.cnf)
	if terminalRender == true {
		cmd.renderTerminal(cnf)
	}

	cmd.ui.Render()
}

func (cmd *Commands) renderTerminal(cnf *config.Config) bool {
	if cnf.Exec.Cmd != "" {
		return true
	}
	return false
}

func (cmd *Commands) getSelected() config.Config {
	var c *config.Config
	c = cmd.cnf
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

// Clean all Selected and Status props in commands structure.
func (cmd *Commands) resetStatus(cnf *config.Config) {
	cnf.Selected = false
	cnf.Status = false
	if len(cnf.Config) > 0 {
		for i := 0; i < len(cnf.Config); i++ {
			cmd.resetStatus(&cnf.Config[i])
		}
	}
}

// Get selected item path via int array
func (cmd *Commands) getSelectedPath(path *[]int, cnf *config.Config) bool {
	if cnf.Selected == true {
		return true
	}
	for i := 0; i < len(cnf.Config); i++ {
		if cmd.getSelectedPath(path, &cnf.Config[i]) == true {
			*path = append(*path, i+1)
			return true
		}
	}
	return false
}

func (cmd *Commands) UpdateRenderElements(c *config.Config) {
	var width int
	//var height int
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
			if cnf.Selected == true {
				if cmd.ui.SelectedRow == 0 {
					menuList.list.BorderStyle = termui.NewStyle(termui.ColorGreen)
				}
				menuList.list.SelectedRow = len(menuList.list.Rows)
				menuList.list.SelectedRowStyle = termui.NewStyle(termui.ColorWhite, termui.ColorGreen)
			} else if cnf.Selected == false && len(path) > i && pathIndex == p && i < len(path)-1 {
				menuList.list.SelectedRow = len(menuList.list.Rows)
				menuList.list.SelectedRowStyle = termui.NewStyle(termui.ColorGreen)
			}
			menuList.list.Rows = append(menuList.list.Rows, cnf.Name)
			if len(cnf.Name) > width {
				width = len(cnf.Name)
			}
			menuList.ratio = cmd.getMenuColumnRatio(width)
		}
		c = &c.Config[pathIndex]
	}
}

func (cmd *Commands) GetLists() []MenuListsItems  {
  return cmd.lists
}

func (cmd *Commands) getMenuColumnRatio(maxTextLength int) float64 {
	width, _ := termui.TerminalDimensions()
	return math.Floor(float64(maxTextLength/(width/100))) / 100
}
