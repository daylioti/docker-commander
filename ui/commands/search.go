package commands

import (
	"github.com/daylioti/docker-commander/ui/render_lock"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"strconv"
	"strings"
)

// Search UI struct
type Search struct {
	Commands      *Commands
	Text          *widgets.Paragraph
	Input          string
	searchIndexes []int
	selectedIndex int
}

// Init initialize search render component.
func (s *Search) Init() {
  s.Text = widgets.NewParagraph()
  s.Text.Border = false
  width, _ := termui.TerminalDimensions()
  s.Text.PaddingBottom = -1
  s.Text.SetRect(0, s.Commands.TermHeight, width, s.Commands.TermHeight +2)
}

// Handle keyboard keys.
func (s *Search) Handle(key string) {
	switch key {
	case "<Backspace>":
		if len(s.Input) > 0 {
			s.Input = s.Input[:len(s.Input)-1]
		}
		if len(s.Input) == 0 {
			s.Reset()
			return
		}
	case "<C-n>":
		s.Next()
	default:
		if !strings.Contains(key, "<") {
			s.searchIndexes = nil
			s.selectedIndex = 0
			s.Input += key
		}
	}
  	s.Search()
}

// Render render search box.
func (s *Search) Render() {
	var pager string
	if len(s.searchIndexes) > 0 {
		pager = strconv.Itoa(s.selectedIndex+1) + "/" + strconv.Itoa(len(s.searchIndexes))
	} else {
		pager = "0/0"
	}
	width, _ := termui.TerminalDimensions()
	spaces := strings.Repeat(" ", width-len(s.Input)-len(pager)-10)
	s.Text.Text = "Search: " + s.Input + spaces + pager
	render_lock.RenderLock(s.Text)
}

// Reset reset search box to defaults.
func (s *Search) Reset() {
	s.Text.Text = ""
	s.Input = ""
	s.searchIndexes = make([]int, 0)
	s.selectedIndex = 0
	render_lock.RenderLock(s.Text)
}

// Search execute search.
func (s *Search) Search() {
	if s.Commands.SelectedArea == KeySelectedCommands {
		s.searchIndexes = s.Commands.Menu.Search(s.Input)
		s.Commands.Menu.setCommandsSelectedIndex(s.getSearchIndex())
	}
	s.Render()
}

// getSearchIndex return config index.
func (s *Search) getSearchIndex() int {
	index := 0
	if s.selectedIndex+1 > len(s.searchIndexes) {
		s.selectedIndex = 0
	} else if len(s.searchIndexes) != 0 {
		index = s.searchIndexes[s.selectedIndex]
	}
	return index
}

// Next display next search item.
func (s *Search) Next() {
	if s.selectedIndex+1 > len(s.searchIndexes) {
		s.selectedIndex = 0
	} else {
		s.selectedIndex++
	}
	s.Commands.Menu.setCommandsSelectedIndex(s.getSearchIndex())
	s.Render()
}
