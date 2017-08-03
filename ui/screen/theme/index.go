package theme

import (
	"sort"

	"github.com/boz/kubetop/ui/elements"
	"github.com/boz/kubetop/ui/theme"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

func NewIndex(ctx elements.Context, req elements.Request) (elements.Screen, error) {
	ctx = ctx.New("theme/index")

	layout := elements.NewVPanes(ctx.Env(), true)

	layout.Append(newIndexBar("app header"))
	layout.Append(newFieldBar("bar", theme.AppHeader.Bar))
	layout.Append(newFieldBar("action", theme.AppHeader.Action))

	layout.Append(newIndexBar("popup"))
	layout.Append(newLabelList("popup", theme.LabelTheme(theme.Popup)))

	panels := []struct {
		name  string
		theme theme.Theme
	}{
		{"active", theme.ThemeActive},
		{"inactive", theme.ThemeInactive},
	}

	for _, ptheme := range panels {
		layout.Append(newIndexBar(ptheme.name + " panel"))
		layout.Append(newFieldBar("base", ptheme.theme.Base))
		layout.Append(newFieldBar("title", ptheme.theme.Title))
		layout.Append(newIndexBar("table"))
		layout.Append(newLabelList("TH", ptheme.theme.Table.TH))
		layout.Append(newLabelList("TD", ptheme.theme.Table.TD))
		layout.Append(newLabelList("TDSelected", ptheme.theme.Table.TDSelected))
	}

	layout.Append(newIndexBar("colorwheel"))

	layout.Append(colorwheel())

	return elements.NewScreen(ctx, req, "Theme", layout), nil
}

func newIndexBar(name string) views.Widget {
	bar := views.NewTextBar()
	bar.SetCenter(name, theme.Base)
	return bar
}

func newLabelList(name string, lt theme.LabelTheme) views.Widget {
	layout := views.NewBoxLayout(views.Horizontal)

	label := views.NewText()
	label.SetAlignment(views.HAlignLeft)
	label.SetText(name)
	label.SetStyle(theme.Base)
	layout.AddWidget(label, 0)

	for _, v := range theme.LabelVariants {
		label := views.NewText()
		label.SetText(string(v))
		label.SetStyle(lt.Get(v))
		label.SetAlignment(views.HAlignCenter)
		layout.AddWidget(label, 1.0)
	}
	return layout
}

func newFieldBar(name string, style tcell.Style) views.Widget {
	bar := views.NewText()
	bar.SetText(name)
	bar.SetStyle(style)
	bar.SetAlignment(views.HAlignLeft)
	return bar
}

func colorwheel() views.Widget {
	model := newcwmodel()
	view := views.NewCellView()
	view.SetModel(model)
	return view
}

type cwmodel struct {
	width  int
	height int
	pwidth int
	names  []string

	message string
}

func newcwmodel() cwmodel {
	m := cwmodel{
		message: "  a b c d e - ! ? - ",
		names:   make([]string, 0, len(tcell.ColorNames)),
	}

	for k, _ := range tcell.ColorNames {
		if l := len(k); l > m.pwidth {
			m.pwidth = l
		}
		m.names = append(m.names, k)
	}

	sort.Strings(m.names)

	m.pwidth += 2
	m.height = len(m.names)
	m.width = m.pwidth + len(m.message)*3

	return m
}

func (m cwmodel) GetCell(x int, y int) (rune, tcell.Style, []rune, int) {
	name := m.names[y]

	if x < len(name) {
		return rune(name[x]), theme.Base, nil, 1
	}

	if x < m.pwidth {
		return ' ', theme.Base, nil, 1
	}

	color := tcell.ColorNames[name]

	delta := x - m.pwidth
	style := theme.Base.Foreground(color)

	switch {
	case delta >= len(m.message)*3:
		return 0, theme.Base, nil, 1
	case delta >= len(m.message)*2:
		style = theme.Base.Background(color)
	case delta >= len(m.message)*1:
		style = theme.Base.Foreground(color).Bold(true)
	}

	delta %= len(m.message)

	return rune(m.message[delta]), style, nil, 1
}

func (m cwmodel) GetBounds() (int, int) {
	return m.width, m.height
}

func (m cwmodel) SetCursor(int, int) {
}

func (m cwmodel) GetCursor() (int, int, bool, bool) {
	return 0, 0, false, false
}

func (m cwmodel) MoveCursor(offx int, offy int) {
}
