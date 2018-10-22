package ui

import (
	"github.com/gizak/termui"
	"provisioner/config"
	"provisioner/docker"
)

type UI struct {
	cnf      *config.Config
	lists    []*termui.List
	terminal *termui.Par
}

func (ui *UI) Init() {
	err := termui.Init()
	if err != nil {
		panic(err)
	}

	defer termui.Close()
	ui.cnf = &config.Config{}
	ui.cnf.Init("./config.yml")
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
	updateTerminal := func(term *string) {
		ui.terminal.Text = *term
		ui.TerminalRender()
	}
	go func() {
		dockerExecute.Init(updateTerminal)
	}()

	for i := 0; i < 8; i++ {
		ui.lists = append(ui.lists, termui.NewList())
	}
	ui.terminal = termui.NewPar("")
	ui.terminal.Width = termui.Body.Width
	ui.terminal.Height = 60

	ui.NormilizeStatus(ui.cnf, ui.Path(ui.cnf))
	ui.Render(ui.cnf)

	termui.Body.Align()

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Handle("/sys/kbd/<up>", func(termui.Event) {
		ui.changeSelected(0, -1, ui.cnf)
		ui.Render(ui.cnf)
	})
	termui.Handle("/sys/kbd/<left>", func(termui.Event) {
		ui.changeSelected(-1, 0, ui.cnf)
		ui.Render(ui.cnf)
	})
	termui.Handle("/sys/kbd/<right>", func(termui.Event) {
		ui.changeSelected(1, 0, ui.cnf)
		ui.Render(ui.cnf)
	})
	termui.Handle("/sys/kbd/<down>", func(termui.Event) {
		ui.changeSelected(0, 1, ui.cnf)
		ui.Render(ui.cnf)
	})
	termui.Handle("/sys/kbd/<enter>", func(termui.Event) {
		selected := ui.getSelected()
		if selected.Command != "" {
			dockerExecute.Exec(selected.Command, ui.Path(ui.cnf), selected.Container)
			ui.TerminalRender()
		}
	})

	termui.Loop()
}

func (ui *UI) changeSelected(x int, y int, cnf *config.Config) {
	path := ui.Path(ui.cnf)
	var c *config.Config  // Current selected
	var cp *config.Config // Parrent of selected
	var cu *config.Config // Above selected
	var cd *config.Config // Below selected.
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
	var path []int
	path = ui.Path(ui.cnf)
	c = ui.cnf
	for i := 0; i < len(path); i++ {
		c = &c.Config[path[i]]
	}
	return *c
}

func (ui *UI) Render(cnf *config.Config) {
	ui.preRender(cnf)
	ui.GetTermuiLists(cnf)
	termui.Clear()
	for i := 0; i < len(ui.lists); i++ {
		termui.Render(ui.lists[i])
	}
	ui.TerminalRender()
}

func (ui *UI) TerminalRender() {
	var h int
	for i := 0; i < len(ui.lists); i++ {
		if ui.lists[i].Height > h {
			h = ui.lists[i].Height
		}
	}
	ui.terminal.Width = termui.Body.Width
	ui.terminal.Y = h
	termui.Render(ui.terminal)
}

func (ui *UI) preRender(cnf *config.Config) {
	var path []int
	path = ui.Path(cnf)

	ui.resetStatus(cnf)
	ui.NormilizeStatus(cnf, path)
}

func (ui *UI) Path(c *config.Config) []int {
	var path []int
	var p []int
	ui.getSelectedPath(&path, c)
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] > 0 {
			p = append(p, int(path[i])-1)
		}
	}
	return p
}

func (ui *UI) NormilizeStatus(c *config.Config, path []int) {
	if len(path) == 0 {
		return
	}
	var cnf *config.Config
	var p int
	var nextPath []int
	p = path[0]
	nextPath = path[1:]
	cnf = &c.Config[p]

	if len(c.Config) > 0 {
		for i := 0; i < len(c.Config); i++ {
			c.Config[i].Status = true
		}
	}
	if len(cnf.Config) > 0 {
		for i := 0; i < len(cnf.Config); i++ {
			cnf.Config[i].Status = true
		}
	}
	if len(nextPath) == 0 {
		cnf.Selected = true
	} else {
		ui.NormilizeStatus(&c.Config[p], nextPath)
	}
}

func (ui *UI) resetStatus(cnf *config.Config) {
	var c *config.Config
	c = cnf
	cnf.Selected = false
	cnf.Status = false
	if len(c.Config) > 0 {
		for i := 0; i < len(c.Config); i++ {
			ui.resetStatus(&c.Config[i])
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

func (ui *UI) GetTermuiLists(c *config.Config) {
	var path []int
	var items [8][]string
	path = ui.Path(ui.cnf)
	for i := 0; i < len(path); i++ {
		for j := 0; j < len(c.Config); j++ {
			if c.Config[j].Status == true {
				if c.Config[j].Selected == true {
					items[i] = append(items[i], StringColor(c.Config[j].Name, "fg-white,bg-green"))
				} else {
					items[i] = append(items[i], c.Config[j].Name)
				}
			}
		}
		c = &c.Config[path[i]]
		if i == len(path)-1 {
			// Next menu
			for oo := 0; oo < len(c.Config); oo++ {
				items[i+1] = append(items[i+1], c.Config[oo].Name)
			}
		}
	}
	ui.lists = nil
	for i := 0; i < 8; i++ {
		if i < len(items) && len(items[i]) > 0 {
			if len(ui.lists)-1 < i {
				// Create new item.
				ui.lists = append(ui.lists, termui.NewList())
			}
			ui.lists[i].Items = items[i]
			ui.lists[i].Height = len(items[i]) + 2
		} else {
			ui.lists = append(ui.lists, termui.NewList())
		}
		ui.lists[i].Width = termui.Body.Width / 8
		ui.lists[i].X = i * termui.Body.Width / 8
	}
}

func StringColor(text string, color string) string {
	return "[" + text + "](" + color + ")"
}
