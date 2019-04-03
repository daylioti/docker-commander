package ui

import (
	commanderWidgets "github.com/daylioti/docker-commander/ui/widgets"
	"github.com/gizak/termui/v3"
)

const (
	InputFieldHeight = 3
)

type Input struct {
	ui           *UI
	Fields       []*commanderWidgets.TextBox
	ActiveField  int
	inputChannel *chan map[string]string
}

func (in *Input) Init(ui *UI) {
	in.ui = ui
}

func (in *Input) Handle(key string) {
	switch key {
	case "<Enter>":
		in.GetInputValues()
		in.ui.Cmd.UpdateRenderElements(in.ui.Cmd.cnf)
		in.Fields = nil
		in.ui.ClearRender = true
		in.ui.Render()
	case "<Tab>":
		if in.ActiveField+1 <= len(in.Fields) {
			in.ActiveField = 0
		} else {
			in.ActiveField++
		}
		for i := 0; i < len(in.Fields); i++ {
			in.Fields[i].BorderStyle = termui.NewStyle(termui.ColorWhite)
		}
		in.Fields[in.ActiveField].BorderStyle = termui.NewStyle(termui.ColorGreen)
		in.ui.Render()
	default:
		in.Fields[in.ActiveField].InsertText(key)
	}
	for i := 0; i < len(in.Fields); i++ {
		termui.Render(in.Fields[i])
	}
}

func (in *Input) GetInputValues() {
	values := make(map[string]string)
	for _, input := range in.Fields {
		values[input.ID] = input.GetText()
	}
	*in.inputChannel <- values
}

func (in *Input) NewInputs(inputs map[string]string, cn *chan map[string]string) {
	var i int
	in.Fields = nil
	in.inputChannel = cn
	var box *commanderWidgets.TextBox
	for id, title := range inputs {
		box = commanderWidgets.NewTextBox()
		box.Title = title
		box.ID = id
		box.SetRect(int(in.ui.TermWidth/4), i*InputFieldHeight, in.ui.TermWidth-int(in.ui.TermWidth/4), i*InputFieldHeight+InputFieldHeight)
		in.Fields = append(in.Fields, box)
		i++
	}
	// Un-focus all other render elements.
	for _, list := range in.ui.Cmd.Lists {
		list.BorderStyle = termui.NewStyle(termui.ColorWhite)
	}
	in.ui.Term.TabPane.BorderStyle = termui.NewStyle(termui.ColorWhite)
	in.ui.Term.DisplayTerminal.BorderStyle = termui.NewStyle(termui.ColorWhite)
	in.ui.ClearRender = true
	in.ui.Render()
}
