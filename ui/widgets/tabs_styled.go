package widgets

import (
	. "github.com/gizak/termui/v3"
	"image"
)

type TabsPaneStyled struct {
	Block
	TabNames         []*TabItem
	ActiveTabIndex   int
	ActiveTabStyle   Style
	InactiveTabStyle Style
}

type TabItem struct {
	Style Style
	Name  string
}

func NewTabPaneStyled() *TabsPaneStyled {
	return &TabsPaneStyled{
		Block:            *NewBlock(),
		ActiveTabStyle:   Theme.Tab.Active,
		InactiveTabStyle: Theme.Tab.Inactive,
	}
}

func (tp *TabsPaneStyled) FocusLeft() {
	if tp.ActiveTabIndex > 0 {
		tp.ActiveTabIndex--
	}
}

func (tp *TabsPaneStyled) FocusRight() {
	if tp.ActiveTabIndex < len(tp.TabNames)-1 {
		tp.ActiveTabIndex++
	}
}

func (tp *TabsPaneStyled) Draw(buf *Buffer) {
	tp.Block.Draw(buf)

	xCoordinate := tp.Inner.Min.X
	for i, name := range tp.TabNames {
		ColorPair := name.Style

		buf.SetString(
			TrimString(name.Name, tp.Inner.Max.X-xCoordinate),
			ColorPair,
			image.Pt(xCoordinate, tp.Inner.Min.Y),
		)

		xCoordinate += 1 + len(name.Name)

		if i < len(tp.TabNames)-1 && xCoordinate < tp.Inner.Max.X {
			buf.SetCell(
				NewCell(VERTICAL_LINE, NewStyle(ColorWhite)),
				image.Pt(xCoordinate, tp.Inner.Min.Y),
			)
		}

		xCoordinate += 2
	}
}
