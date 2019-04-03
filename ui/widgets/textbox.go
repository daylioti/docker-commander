package widgets

import (
	"image"

	. "github.com/gizak/termui/v3"
)

type TextBox struct {
	Block
	WrapText    bool
	TextStyle   Style
	CursorStyle Style
	ShowCursor  bool
	ID          string
	text        [][]Cell
	cursorPoint image.Point
}

var TextBoxTheme = TextBoxThemeType{
	Text:   NewStyle(ColorWhite),
	Cursor: NewStyle(ColorWhite, ColorClear, ModifierReverse),
}

type TextBoxThemeType struct {
	Text   Style
	Cursor Style
}

func NewTextBox() *TextBox {
	return &TextBox{
		Block:       *NewBlock(),
		WrapText:    false,
		TextStyle:   TextBoxTheme.Text,
		CursorStyle: TextBoxTheme.Cursor,

		text:        [][]Cell{[]Cell{}},
		cursorPoint: image.Pt(1, 1),
	}
}

func (tb *TextBox) Draw(buf *Buffer) {
	tb.Block.Draw(buf)

	yCoordinate := 0
	for _, line := range tb.text {
		if tb.WrapText {
			line = WrapCells(line, uint(tb.Inner.Dx()))
		}
		lines := SplitCells(line, '\n')
		for _, line := range lines {
			for _, cx := range BuildCellWithXArray(line) {
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
				[][]Cell{append(tb.text[index-1], tb.text[index]...)},
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
	cells := ParseStyles(input, tb.TextStyle)
	lines := SplitCells(cells, '\n')
	index := tb.cursorPoint.Y - 1
	cellsAfterCursor := tb.text[index][tb.cursorPoint.X-1:]
	tb.text[index] = append(tb.text[index][:tb.cursorPoint.X-1], lines[0]...)
	for i, line := range lines[1:] {
		index := tb.cursorPoint.Y + i
		tb.text = append(tb.text[:index], append([][]Cell{line}, tb.text[index:]...)...)
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
	tb.text = [][]Cell{[]Cell{}}
	tb.cursorPoint = image.Pt(1, 1)
}

// SetText sets the text to the given text.
func (tb *TextBox) SetText(input string) {
	tb.ClearText()
	tb.InsertText(input)
}

func (tb *TextBox) GetText() string {
	var text string
	for _, r := range tb.text {
		for _, t := range r {
			text += string(t.Rune)
		}
	}
	return text
}

func (tb *TextBox) MoveCursorLeft() {
	tb.MoveCursor(tb.cursorPoint.X-1, tb.cursorPoint.Y)
}

func (tb *TextBox) MoveCursorRight() {
	tb.MoveCursor(tb.cursorPoint.X+1, tb.cursorPoint.Y)
}

func (tb *TextBox) MoveCursorUp() {
	tb.MoveCursor(tb.cursorPoint.X, tb.cursorPoint.Y-1)
}

func (tb *TextBox) MoveCursorDown() {
	tb.MoveCursor(tb.cursorPoint.X, tb.cursorPoint.Y+1)
}

func (tb *TextBox) MoveCursor(x, y int) {
	tb.cursorPoint.Y = MinInt(MaxInt(1, y), len(tb.text))
	tb.cursorPoint.X = MinInt(MaxInt(1, x), len(tb.text[tb.cursorPoint.Y-1])+1)
}
