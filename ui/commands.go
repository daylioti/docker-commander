package ui

import (
	"github.com/daylioti/docker-commander/config"
	"github.com/daylioti/docker-commander/docker"
	"github.com/gizak/termui"
)

type Commands struct {
	cnf          *config.Config
	lists        []*termui.List
	listCols     []*termui.Row
	terminal     *termui.Paragraph
	dockerClient *docker.Docker
}

func (cmd *Commands) Init(cnf *config.Config, dockerClient *docker.Docker) {
	cmd.cnf = cnf
	cmd.dockerClient = dockerClient

	cmd.dockerClient.Exec.SetTerminalUpdateFn(cmd.updateTerminal)

	cmd.updateStatus(cmd.cnf, cmd.Path(cmd.cnf))

	cmd.TerminalRender()

	cmd.Render(cmd.cnf)
	termui.Clear()
}

func (cmd *Commands) Handle(key string) {
	switch key {
	case "<Up>":
		cmd.changeSelected(0, -1, cmd.cnf)
		break
	case "<Left>":
		cmd.changeSelected(-1, 0, cmd.cnf)
		break
	case "<Right>":
		cmd.changeSelected(1, 0, cmd.cnf)
		break
	case "<Down>":
		cmd.changeSelected(0, 1, cmd.cnf)
		break
	case "<Enter>":
		termui.Body.Rows[2] = termui.NewRow(termui.NewCol(12, 0, cmd.terminal))
		selected := cmd.getSelected()
		if selected.Command != "" {
			cmd.TerminalRender()
			cmd.dockerClient.Exec.SetTerminalHeight(cmd.terminal.Height)
			cmd.dockerClient.Exec.CommandExecute(selected.Command, cmd.Path(cmd.cnf), selected.Container)
		}
		break
	case "<Resize>":
		cmd.updateRenderElements(cmd.cnf)
		termui.Body.Width = termui.TermWidth()
		termui.Body.Align()
		termui.Clear() // Delete this line to avoid the crash
		termui.Render(termui.Body)
		break
	}
}

func (cmd *Commands) SetDockerClient(client *docker.Docker) {
	cmd.dockerClient = client
}

func (cmd *Commands) updateTerminal() {
	if cmd.terminal != nil && cmd.terminal.Text != *cmd.dockerClient.Exec.RunningTerminal {
		cmd.terminal.Text = *cmd.dockerClient.Exec.RunningTerminal
		cmd.TerminalRender()
	}
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
	cmd.updateRenderElements(cnf)
	termui.Clear()
	termui.Body.Align()
	termui.Render(termui.Body)
}

func (cmd *Commands) TerminalRender() {
	if cmd.terminal == nil {
		cmd.terminal = termui.NewParagraph("")
	}
	var h int
	for i := 0; i < len(cmd.lists); i++ {
		if cmd.lists[i].Height > h {
			h = cmd.lists[i].Height
		}
	}
	cmd.terminal.Width = termui.Body.Width
	cmd.terminal.Y = h
	cmd.terminal.Height = termui.TermHeight() - h
	termui.Render(termui.Body)
}

func (cmd *Commands) preRender(cnf *config.Config) {
	var path []int
	path = cmd.Path(cnf)

	cmd.resetStatus(cnf)
	cmd.updateStatus(cnf, path)
}

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

func (cmd *Commands) resetStatus(cnf *config.Config) {
	cnf.Selected = false
	cnf.Status = false
	if len(cnf.Config) > 0 {
		for i := 0; i < len(cnf.Config); i++ {
			cmd.resetStatus(&cnf.Config[i])
		}
	}
}

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

func (cmd *Commands) updateRenderElements(c *config.Config) {
	var item string
	var width int
	var x int
	path := append(cmd.Path(cmd.cnf), 0)
	cmd.lists = nil
	cmd.listCols = nil
	for i, pathIndex := range path {
		if len(c.Config) <= pathIndex {
			break
		}
		width = 0
		cmd.lists = append(cmd.lists, termui.NewList())
		for p, cnf := range c.Config {
			if len(cnf.Name) > width {
				width = len(cnf.Name)
			}
			if cnf.Selected == true {
				item = StringColor(cnf.Name, "fg-white,bg-green")
			} else {
				if len(path) > i && pathIndex == p && i < len(path)-1 {
					item = StringColor(cnf.Name, "fg-green")
				} else {
					item = cnf.Name
				}
			}
			cmd.lists[i].Items = append(cmd.lists[i].Items, item)
		}

		cmd.lists[i].Height = len(cmd.lists[i].Items) + 2
		cmd.lists[i].Width = width + 3
		cmd.lists[i].X = x
		x += cmd.lists[i].Width

		c = &c.Config[pathIndex]
		cmd.listCols = append(cmd.listCols, termui.NewCol(cmd.getSpan(cmd.lists[i].Width), 0, cmd.lists[i]))
	}
	termui.Body.Rows[0] = termui.NewRow(cmd.listCols...)
}

func (cmd *Commands) getSpan(maxTextLength int) int {
	oneColumn := termui.TermWidth() / 12
	if oneColumn < maxTextLength && oneColumn != 0 {
		return int(maxTextLength/oneColumn) + 1
	}
	return 1
}
