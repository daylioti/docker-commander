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
	lists        []*widgets.List
	listCols     []interface{}
	terminal     *widgets.List
	dockerClient *docker.Docker
}

func (cmd *Commands) Init(cnf *config.Config, dockerClient *docker.Docker, ui *UI) {
	cmd.cnf = cnf
	cmd.dockerClient = dockerClient
	cmd.ui = ui

	cmd.dockerClient.Exec.SetTerminalUpdateFn(cmd.updateTerminal)
	cmd.updateStatus(cmd.cnf, cmd.Path(cmd.cnf))

	cmd.TerminalRender()
	cmd.Render(cmd.cnf)
	termui.Clear()

}

func (cmd *Commands) Handle(key string) {
	switch key {
	case "<Tab>":
		if cmd.ui.SelectedRow >= len(cmd.ui.Grid.Items) {
			cmd.ui.SelectedRow++
		} else {
			cmd.ui.SelectedRow = 0
		}
		cmd.Render(cmd.cnf)
		cmd.TerminalRender()
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
		cmd.ui.Grid.Items[1].Entry = termui.NewCol(1.0, cmd.terminal)
		selected := cmd.getSelected()
		if selected.Command != "" {
			cmd.TerminalRender()
			//cmd.dockerClient.Exec.SetTerminalHeight(cmd.terminal.Dy())
			cmd.dockerClient.Exec.CommandExecute(selected.Command, cmd.Path(cmd.cnf), selected.Container)
		}
		//termui.Body.Rows[2] = termui.NewRow(termui.NewCol(12, 0, cmd.terminal))
		//selected := cmd.getSelected()
		//if selected.Command != "" {
		//	cmd.TerminalRender()
		//	cmd.dockerClient.Exec.SetTerminalHeight(cmd.terminal.Height)
		//	cmd.dockerClient.Exec.CommandExecute(selected.Command, cmd.Path(cmd.cnf), selected.Container)
		//}
		break
	}
}

func (cmd *Commands) SetDockerClient(client *docker.Docker) {
	cmd.dockerClient = client
}

func (cmd *Commands) updateTerminal() {
	if cmd.terminal != nil && len(cmd.terminal.Rows) != len(*cmd.dockerClient.Exec.RunningTerminal) {
		cmd.terminal.Rows = *cmd.dockerClient.Exec.RunningTerminal
		cmd.TerminalRender()
	}
}

func (cmd *Commands) changeSelected(x int, y int, cnf *config.Config) {
	if cmd.ui.SelectedRow != 0 {
		return
	}
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
	cmd.Render(cmd.cnf)
	cmd.dockerClient.Exec.ChangeTerminal(cmd.Path(cmd.cnf))
}

func (cmd *Commands) getSelected() config.Config {
	var c *config.Config
	c = cmd.cnf
	for _, path := range cmd.Path(cmd.cnf) {
		c = &c.Config[path]
	}
	return *c
}

func (cmd *Commands) Render(cnf *config.Config) {
	cmd.preRender(cnf)
	cmd.UpdateRenderElements(cnf)
	termui.Clear()
	termui.Render(cmd.ui.Grid)
}

func (cmd *Commands) TerminalRender() {
	if cmd.terminal == nil {
		cmd.terminal = widgets.NewList()
		cmd.terminal.WrapText = false
	}
	if cmd.ui.SelectedRow == 1 {
		cmd.terminal.BorderStyle = termui.NewStyle(termui.ColorGreen)
	}

	//var h int
	//for i := 0; i < len(cmd.lists); i++ {
	//if cmd.lists[i].Max > h {
	//	h = cmd.lists[i].Max
	//}
	//if cmd.lists[i].Height > h {
	//	h = cmd.lists[i].Height
	//}
	//}
	//cmd.terminal.Width = termui.Body.Width
	//cmd.terminal.Y = h
	//cmd.terminal.Height = termui.TermHeight() - h
	cmd.Render(cmd.cnf)
}

func (cmd *Commands) preRender(cnf *config.Config) {
	//cmd.resetStatus(cnf)
	cmd.updateStatus(cnf, cmd.Path(cnf))
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
	var nextPath []int
	cnf = &cnf.Config[path[0]]
	nextPath = path[1:]

	if len(cnf.Config) > 0 {
		for i := 0; i < len(cnf.Config); i++ {
			cnf.Config[i].Status = true
		}
	}
	if len(nextPath) == 0 {
		cnf.Selected = true
	} else {
		cmd.updateStatus(cnf, nextPath)
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
	var height int
	path := append(cmd.Path(cmd.cnf), 0)
	cmd.lists = nil
	cmd.listCols = nil
	for i, pathIndex := range path {
		if len(c.Config) <= pathIndex {
			break
		}
		width = 0
		cmd.lists = append(cmd.lists, widgets.NewList())
		cmd.lists[i].SelectedRow = 0
		for p, cnf := range c.Config {
			if cnf.Selected == true {
				if cmd.ui.SelectedRow == 0 {
					cmd.lists[i].BorderStyle = termui.NewStyle(termui.ColorGreen)
				}
				cmd.lists[i].SelectedRow = uint(len(cmd.lists[i].Rows))
				cmd.lists[i].SelectedRowStyle = termui.NewStyle(termui.ColorWhite, termui.ColorGreen)
			} else if cnf.Selected == false && len(path) > i && pathIndex == p && i < len(path)-1 {
				cmd.lists[i].SelectedRow = uint(len(cmd.lists[i].Rows))
				cmd.lists[i].SelectedRowStyle = termui.NewStyle(termui.ColorGreen)
			}
			cmd.lists[i].Rows = append(cmd.lists[i].Rows, cnf.Name)
			if len(cnf.Name) > width {
				width = len(cnf.Name)
			}
			if len(cmd.lists[i].Rows) > height {
				height = len(cmd.lists[i].Rows)
			}
		}
		termui.TerminalDimensions()
		c = &c.Config[pathIndex]
		cmd.listCols = append(cmd.listCols, termui.NewCol(cmd.getMenuColumnRatio(width), cmd.lists[i]))
	}
	cmd.ui.Grid = termui.NewGrid()
	termWidth, termHeight := termui.TerminalDimensions()
	cmd.ui.Grid.SetRect(0, 0, termWidth, termHeight)
	cmd.ui.Grid.Set(
		termui.NewRow(0.2, cmd.listCols...),
		termui.NewRow(0.8, cmd.terminal),
	)
}

func (cmd *Commands) getMenuColumnRatio(maxTextLength int) float64 {
	width, _ := termui.TerminalDimensions()
	return math.Floor(float64(maxTextLength/(width/100))) / 100
}
