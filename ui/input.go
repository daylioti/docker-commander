package ui

import (
	"github.com/atotto/clipboard"
	commanderWidgets "github.com/daylioti/docker-commander/ui/widgets"
	"github.com/gizak/termui/v3"
)

// InputFieldHeight - border sizes.
const (
	InputFieldHeight = 3
)

// Input UI struct.
type Input struct {
	ui           *UI
	Fields       []*commanderWidgets.TextBox
	ActiveField  int
	inputChannel *chan map[string]string
}

// Init initialize input render component.
func (in *Input) Init(ui *UI) {
	in.ui = ui
}

// Handle keyboard keys.
func (in *Input) Handle(key string) {
	switch key {
	case "<Enter>":
		in.GetInputValues()
		in.ui.Cmd.UpdateRenderElements(in.ui.Cmd.cnf)
		in.Fields = nil
		// Render all ui twice to fix bug with partly rendering old elements.
		in.ui.Render()
		in.ui.Render()
	case "<Tab>":
		if in.ActiveField+1 >= len(in.Fields) {
			in.ActiveField = 0
		} else {
			in.ActiveField++
		}
		in.Render()
	case "<Up>":
		if in.ActiveField != 0 {
			in.ActiveField--
			in.Render()
		}
	case "<Down>":
		if in.ActiveField+1 < len(in.Fields) {
			in.ActiveField++
			in.Render()
		}
	case "<Backspace>":
		in.Fields[in.ActiveField].Backspace()
		termui.Render(in.Fields[in.ActiveField])
	case "<Space>":
		in.Fields[in.ActiveField].InsertText(" ")
		termui.Render(in.Fields[in.ActiveField])
	case "<Left>":
		in.Fields[in.ActiveField].MoveCursorLeft()
		termui.Render(in.Fields[in.ActiveField])
	case "<Right>":
		in.Fields[in.ActiveField].MoveCursorRight()
		termui.Render(in.Fields[in.ActiveField])
	case "<C-v>":
		// @Todo implement clipboard with better way.
		// It requires additional tools xsel, xclip, wl-clipboard.
		clip := in.ReadFromClipboard()
		if clip != "" {
			in.Fields[in.ActiveField].InsertText(clip)
		}
		termui.Render(in.Fields[in.ActiveField])
	case "<Escape>":
		in.Fields = nil
		in.ui.Render()
	default:
		if in.allowedInput(key) {
			in.Fields[in.ActiveField].InsertText(key)
		}
		termui.Render(in.Fields[in.ActiveField])
	}
}

// Render function, that render input component.
func (in *Input) Render() {
	in.Fields[in.ActiveField].BorderStyle = termui.NewStyle(termui.ColorGreen)
	for i, field := range in.Fields {
		if i != in.ActiveField {
			field.BorderStyle = termui.NewStyle(termui.ColorWhite)
		}
		termui.Render(field)
	}
}

// ReadFromClipboard get string from clipboard.
func (in *Input) ReadFromClipboard() string {
	clip, err := clipboard.ReadAll()
	if err != nil {
		return ""
	}
	return clip
}

// allowedInput - filter allowed to paste in input field keyboard keys.
func (in *Input) allowedInput(key string) bool {
	return key != "<MouseLeft>" && key != "<MouseRelease>" && key != "<MouseRight>"
}

// GetInputValues - get input values, using chanel.
func (in *Input) GetInputValues() {
	values := make(map[string]string)
	for _, input := range in.Fields {
		values[input.ID] = input.GetText()
	}
	*in.inputChannel <- values
}

// NewInputs - create and render input fields.
func (in *Input) NewInputs(inputs map[string]string, cn *chan map[string]string) {
	var inputsCount int
	in.Fields = nil
	in.inputChannel = cn
	var box *commanderWidgets.TextBox
	for id, title := range inputs {
		box = commanderWidgets.NewTextBox()
		box.Title = title
		box.ID = id
		box.SetRect(int(in.ui.TermWidth/4), inputsCount*InputFieldHeight, in.ui.TermWidth-int(in.ui.TermWidth/4),
			inputsCount*InputFieldHeight+InputFieldHeight)
		box.ShowCursor = true
		in.Fields = append(in.Fields, box)
		inputsCount++
	}
	in.Fields[0].BorderStyle = termui.NewStyle(termui.ColorGreen)
	// Un-focus all other render elements.
	for _, list := range in.ui.Cmd.Lists {
		list.BorderStyle = termui.NewStyle(termui.ColorWhite)
	}
	termui.Clear()
	in.ui.Cmd.UnFocus()
	in.ui.Term.UnFocus()
}
