package ui

import (
	"docker-commander/config"
	"docker-commander/docker"
	"github.com/gizak/termui"
)

type UI struct {
	cnf      *config.Config
	lists    []*termui.List
	terminal *termui.Par
}

func (ui *UI) Init(configPath string) {
	err := termui.Init()
	if err != nil {
		panic(err)
	}

	defer termui.Close()
	ui.cnf = &config.Config{}
	ui.cnf.Init(configPath)
	ui.cnf.Status = true
	ui.cnf.Config[0].Selected = true
	for i := 0; i < len(ui.cnf.Config); i++ {
		ui.cnf.Config[i].Status = true
	}
	ui.cnf.Config[0].Status = true
	if len(ui.cnf.Config[0].Config) > 0 {
		for i := 0; i < len(ui.cnf.Config[0].Config); i++ {
			ui.cnf.Config[0].Config[i].Status = true
		}
	}

	dockerExecute := new(docker.Docker)

	updateTerminal := func() {
		if ui.terminal != nil && ui.terminal.Text != *dockerExecute.RunningTerminal {
			ui.terminal.Text = *dockerExecute.RunningTerminal
			ui.TerminalRender()
		}
	}
	go func() {
		dockerExecute.Init(updateTerminal)
	}()

	ui.updateStatus(ui.cnf, ui.Path(ui.cnf))
	ui.Render(ui.cnf)

	termui.Body.Align()

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Handle("/sys/kbd/<up>", func(termui.Event) {
		ui.changeSelected(0, -1, ui.cnf)
		ui.Render(ui.cnf)
		dockerExecute.ChangeTerminal(ui.Path(ui.cnf))
	})
	termui.Handle("/sys/kbd/<left>", func(termui.Event) {
		ui.changeSelected(-1, 0, ui.cnf)
		ui.Render(ui.cnf)
		dockerExecute.ChangeTerminal(ui.Path(ui.cnf))
	})
	termui.Handle("/sys/kbd/<right>", func(termui.Event) {
		ui.changeSelected(1, 0, ui.cnf)
		ui.Render(ui.cnf)
		dockerExecute.ChangeTerminal(ui.Path(ui.cnf))
	})
	termui.Handle("/sys/kbd/<down>", func(termui.Event) {
		ui.changeSelected(0, 1, ui.cnf)
		ui.Render(ui.cnf)
		dockerExecute.ChangeTerminal(ui.Path(ui.cnf))
	})
	termui.Handle("/sys/kbd/<enter>", func(termui.Event) {
		selected := ui.getSelected()
		if selected.Command != "" {
			ui.TerminalRender()
			dockerExecute.SetTerminalHeight(ui.terminal.Height)
			dockerExecute.Exec(selected.Command, ui.Path(ui.cnf), selected.Container)

		}
	})

	termui.Loop()
}

func (ui *UI) changeSelected(x int, y int, cnf *config.Config) {
	path := ui.Path(ui.cnf)
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
}

func (ui *UI) getSelected() config.Config {
	var c *config.Config
	c = ui.cnf
	for _, path := range ui.Path(ui.cnf) {
      c = &c.Config[path]
	}
	return *c
}

func (ui *UI) Render(cnf *config.Config) {
	ui.preRender(cnf)
	ui.updateRenderElements(cnf)
	termui.Clear()
	for i := 0; i < len(ui.lists); i++ {
		termui.Render(ui.lists[i])
	}
}

func (ui *UI) TerminalRender() {
	if ui.terminal == nil {
		ui.terminal = termui.NewPar("")
	}
	var h int
	for i := 0; i < len(ui.lists); i++ {
		if ui.lists[i].Height > h {
			h = ui.lists[i].Height
		}
	}
	ui.terminal.Width = termui.Body.Width
	ui.terminal.Y = h
	ui.terminal.Height = termui.TermHeight() - h
	termui.Render(ui.terminal)
}

func (ui *UI) preRender(cnf *config.Config) {
	var path []int
	path = ui.Path(cnf)

	ui.resetStatus(cnf)
	ui.updateStatus(cnf, path)
}

func (ui *UI) Path(cnf *config.Config) []int {
	var path []int
	var p []int
	ui.getSelectedPath(&path, cnf)
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] > 0 {
			p = append(p, int(path[i])-1)
		}
	}
	return p
}

func (ui *UI) updateStatus(cnf *config.Config, path []int) {
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
		ui.updateStatus(cnf, nextPath)
	}
}

func (ui *UI) resetStatus(cnf *config.Config) {
	cnf.Selected = false
	cnf.Status = false
	if len(cnf.Config) > 0 {
		for i := 0; i < len(cnf.Config); i++ {
			ui.resetStatus(&cnf.Config[i])
		}
	}
}

func (ui *UI) getSelectedPath(path *[]int, cnf *config.Config) bool {
	if cnf.Selected == true {
		return true
	}
	for i := 0; i < len(cnf.Config); i++ {
		if ui.getSelectedPath(path, &cnf.Config[i]) == true {
			*path = append(*path, i+1)
			return true
		}
	}
	return false
}

func (ui *UI) updateRenderElements(c *config.Config) {
	var item string
	var width int
	var x int
	ui.lists = nil
    for i, path := range append(ui.Path(ui.cnf), 0) {
		if len(c.Config) <= path {
			break
		}
    	width = 0
		ui.lists = append(ui.lists, termui.NewList())
        for _, cnf := range c.Config {
            if len(cnf.Name) > width {
            	width = len(cnf.Name)
			}
        	if cnf.Selected == true {
              item =  StringColor(cnf.Name, "fg-white,bg-green")
			} else {
				item = cnf.Name
			}
			ui.lists[i].Items = append(ui.lists[i].Items, item)
		}

		ui.lists[i].Height = len(ui.lists[i].Items) + 2
		ui.lists[i].Width = width + 3
		ui.lists[i].X = x
		x += ui.lists[i].Width

		c = &c.Config[path]
	}
}


func StringColor(text string, color string) string {
	return "[" + text + "](" + color + ")"
}
