package widgets

import (
	"github.com/gizak/termui/v3"
	"image"
)

type TabsPaneStyled struct {
	termui.Block
	TabNames         []*TabItem
	ActiveTabIndex   int
	ActiveTabStyle   termui.Style
	InactiveTabStyle termui.Style
}

type TabItem struct {
	Style termui.Style
	Name  string
}

func NewTabPaneStyled() *TabsPaneStyled {
	return &TabsPaneStyled{
		Block:            *termui.NewBlock(),
		ActiveTabStyle:   termui.Theme.Tab.Active,
		InactiveTabStyle: termui.Theme.Tab.Inactive,
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

func (tp *TabsPaneStyled) Draw(buf *termui.Buffer) {
	tp.Block.Draw(buf)

	xCoordinate := tp.Inner.Min.X
	for i, name := range tp.TabNames {
		ColorPair := name.Style

		buf.SetString(
			termui.TrimString(name.Name, tp.Inner.Max.X-xCoordinate),
			ColorPair,
			image.Pt(xCoordinate, tp.Inner.Min.Y),
		)

		xCoordinate += 1 + len(name.Name)

		if i < len(tp.TabNames)-1 && xCoordinate < tp.Inner.Max.X {
			buf.SetCell(
				termui.NewCell(termui.VERTICAL_LINE, termui.NewStyle(termui.ColorWhite)),
				image.Pt(xCoordinate, tp.Inner.Min.Y),
			)
		}

		xCoordinate += 2
	}
}
