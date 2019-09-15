package commands

import (
	"github.com/atotto/clipboard"
	"github.com/daylioti/docker-commander/ui/render_lock"
	"github.com/gizak/termui/v3"
)

import (
	"github.com/daylioti/docker-commander/config"
	commanderWidgets "github.com/daylioti/docker-commander/ui/widgets"
)

// InputFieldHeight - border sizes.
const (
	InputFieldHeight = 3
)

// Input UI struct.
type Input struct {
	Fields      []*commanderWidgets.TextBox
	Commands    *Commands
	ActiveField int
	cnf         config.Config
}

// Handle keyboard keys.
func (in *Input) Handle(key string) {
	switch key {
	case "<Enter>":
		values := in.GetInputValues()
		for k, v := range values {
			in.cnf.ReplacePlaceholder(k, v, &in.cnf)
		}
		in.Commands.Menu.commandExecProcess(in.cnf)
		in.Commands.Menu.UpdateRenderElements(in.Commands.Cnf)
		in.Fields = nil
		in.Commands.RenderAll()

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
		render_lock.RenderLock(in.Fields[in.ActiveField])
	case "<Space>":
		in.Fields[in.ActiveField].InsertText(" ")
		render_lock.RenderLock(in.Fields[in.ActiveField])
	case "<Left>":
		in.Fields[in.ActiveField].MoveCursorLeft()
		render_lock.RenderLock(in.Fields[in.ActiveField])
	case "<Right>":
		in.Fields[in.ActiveField].MoveCursorRight()
		render_lock.RenderLock(in.Fields[in.ActiveField])
	case "<C-v>":
		// @Todo implement clipboard with better way.
		// It requires additional tools xsel, xclip, wl-clipboard.
		clip := in.ReadFromClipboard()
		if clip != "" {
			in.Fields[in.ActiveField].InsertText(clip)
		}
		render_lock.RenderLock(in.Fields[in.ActiveField])
	case "<Escape>":
		in.Fields = nil
		in.Commands.RenderAll()
	default:
		if in.allowedInput(key) {
			in.Fields[in.ActiveField].InsertText(key)
		}
		render_lock.RenderLock(in.Fields[in.ActiveField])
	}
}

// Render function, that render input component.
func (in *Input) Render() {
	in.Fields[in.ActiveField].BorderStyle = termui.NewStyle(termui.ColorGreen)
	for i, field := range in.Fields {
		if i != in.ActiveField {
			field.BorderStyle = termui.NewStyle(termui.ColorWhite)
		}
		render_lock.RenderLock(field)
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
func (in *Input) GetInputValues() map[string]string {
	values := make(map[string]string)
	for _, input := range in.Fields {
		values[input.ID] = input.GetText()
	}
	return values
}

// NewInputs - create and render input fields.
func (in *Input) NewInputs(inputs map[string]string, cnf config.Config) {
	var inputsCount int
	in.Fields = nil
	in.cnf = cnf
	var box *commanderWidgets.TextBox
	for id, title := range inputs {
		box = commanderWidgets.NewTextBox()
		box.Title = title
		box.ID = id
		box.SetRect(int(in.Commands.TermWidth/4), inputsCount*InputFieldHeight, in.Commands.TermWidth-int(in.Commands.TermWidth/4),
			inputsCount*InputFieldHeight+InputFieldHeight)
		box.ShowCursor = true
		in.Fields = append(in.Fields, box)
		inputsCount++
	}
	in.Fields[0].BorderStyle = termui.NewStyle(termui.ColorGreen)
	// Un-focus all other render elements.
	for _, list := range in.Commands.Menu.Lists {
		list.BorderStyle = termui.NewStyle(termui.ColorWhite)
	}
	termui.Clear()
	in.Commands.Menu.UnFocus()
	in.Commands.Terminal.UnFocus()
}
