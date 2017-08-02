package theme

type Themeable interface {
	SetTheme(Theme)
	Theme() Theme
}

type ThemeableWidget struct {
	theme Theme
}

func (w *ThemeableWidget) SetTheme(theme Theme) {
	w.theme = theme
}

func (w *ThemeableWidget) Theme() Theme {
	return w.theme
}
