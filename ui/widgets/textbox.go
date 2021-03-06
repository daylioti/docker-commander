package widgets

import (
	"github.com/gizak/termui/v3"
	"image"
)

// TextBox main text box struct.
type TextBox struct {
	termui.Block
	WrapText    bool
	TextStyle   termui.Style
	CursorStyle termui.Style
	ShowCursor  bool
	ID          string
	text        [][]termui.Cell
	cursorPoint image.Point
}

// TextBoxTheme text box theme.
var TextBoxTheme = TextBoxThemeType{
	Text:   termui.NewStyle(termui.ColorWhite),
	Cursor: termui.NewStyle(termui.ColorWhite, termui.ColorClear, termui.ModifierReverse),
}

// TextBoxThemeType theme type.
type TextBoxThemeType struct {
	Text   termui.Style
	Cursor termui.Style
}

// NewTextBox create new text box.
func NewTextBox() *TextBox {
	return &TextBox{
		Block:       *termui.NewBlock(),
		WrapText:    false,
		TextStyle:   TextBoxTheme.Text,
		CursorStyle: TextBoxTheme.Cursor,

		text:        [][]termui.Cell{[]termui.Cell{}},
		cursorPoint: image.Pt(1, 1),
	}
}

// Draw implements the Drawable interface.
func (tb *TextBox) Draw(buf *termui.Buffer) {
	tb.Block.Draw(buf)

	yCoordinate := 0
	for _, line := range tb.text {
		if tb.WrapText {
			line = termui.WrapCells(line, uint(tb.Inner.Dx()))
		}
		lines := termui.SplitCells(line, '\n')
		for _, line := range lines {
			for _, cx := range termui.BuildCellWithXArray(line) {
				x, cell := cx.X, cx.Cell
				buf.SetCell(cell, image.Pt(x, yCoordinate).Add(tb.Inner.Min))
			}
			yCoordinate++
		}
		if yCoordinate > tb.Inner.Max.Y {
			break
		}
	}

	if tb.ShowCursor {
		point := tb.cursorPoint.Add(tb.Inner.Min).Sub(image.Pt(1, 1))
		cell := buf.GetCell(point)
		cell.Style = tb.CursorStyle
		buf.SetCell(cell, point)
	}
}

// Backspace remove previous char.
func (tb *TextBox) Backspace() {
	if tb.cursorPoint == image.Pt(1, 1) {
		return
	}
	if tb.cursorPoint.X == 1 {
		index := tb.cursorPoint.Y - 1
		tb.cursorPoint.X = len(tb.text[index-1]) + 1
		tb.text = append(
			tb.text[:index-1],
			append(
				[][]termui.Cell{append(tb.text[index-1], tb.text[index]...)},
				tb.text[index+1:len(tb.text)]...,
			)...,
		)
		tb.cursorPoint.Y--
	} else {
		index := tb.cursorPoint.Y - 1
		tb.text[index] = append(
			tb.text[index][:tb.cursorPoint.X-2],
			tb.text[index][tb.cursorPoint.X-1:]...,
		)
		tb.cursorPoint.X--
	}
}

// InsertText inserts the given text at the cursor position.
func (tb *TextBox) InsertText(input string) {
	cells := termui.ParseStyles(input, tb.TextStyle)
	lines := termui.SplitCells(cells, '\n')
	index := tb.cursorPoint.Y - 1
	cellsAfterCursor := tb.text[index][tb.cursorPoint.X-1:]
	tb.text[index] = append(tb.text[index][:tb.cursorPoint.X-1], lines[0]...)
	for i, line := range lines[1:] {
		index := tb.cursorPoint.Y + i
		tb.text = append(tb.text[:index], append([][]termui.Cell{line}, tb.text[index:]...)...)
	}
	tb.cursorPoint.Y += len(lines) - 1
	index = tb.cursorPoint.Y - 1
	tb.text[index] = append(tb.text[index], cellsAfterCursor...)
	if len(lines) > 1 {
		tb.cursorPoint.X = len(lines[len(lines)-1]) + 1
	} else {
		tb.cursorPoint.X += len(lines[0])
	}
}

// ClearText clears the text and resets the cursor position.
func (tb *TextBox) ClearText() {
	tb.text = [][]termui.Cell{[]termui.Cell{}}
	tb.cursorPoint = image.Pt(1, 1)
}

// SetText sets the text to the given text.
func (tb *TextBox) SetText(input string) {
	tb.ClearText()
	tb.InsertText(input)
}

// GetText get text from box.
func (tb *TextBox) GetText() string {
	var text string
	for _, r := range tb.text {
		for _, t := range r {
			text += string(t.Rune)
		}
	}
	return text
}

// MoveCursorLeft move cursor to left.
func (tb *TextBox) MoveCursorLeft() {
	tb.MoveCursor(tb.cursorPoint.X-1, tb.cursorPoint.Y)
}

// MoveCursorRight move cursor to right.
func (tb *TextBox) MoveCursorRight() {
	tb.MoveCursor(tb.cursorPoint.X+1, tb.cursorPoint.Y)
}

// MoveCursorUp move cursor to up.
func (tb *TextBox) MoveCursorUp() {
	tb.MoveCursor(tb.cursorPoint.X, tb.cursorPoint.Y-1)
}

// MoveCursorDown move cursor to down.
func (tb *TextBox) MoveCursorDown() {
	tb.MoveCursor(tb.cursorPoint.X, tb.cursorPoint.Y+1)
}

// MoveCursor move cursor to coordinates.
func (tb *TextBox) MoveCursor(x, y int) {
	tb.cursorPoint.Y = termui.MinInt(termui.MaxInt(1, y), len(tb.text))
	tb.cursorPoint.X = termui.MinInt(termui.MaxInt(1, x), len(tb.text[tb.cursorPoint.Y-1])+1)
}
