package theme

import "github.com/gdamore/tcell"

const (
	LabelNormal  LabelVariant = "normal"
	LabelSuccess LabelVariant = "success"
	LabelWarn    LabelVariant = "warn"
	LabelError   LabelVariant = "error"
)

var (
	Base = tcell.StyleDefault

	AppHeader = AppHeaderTheme{
		Bar: Base.
			Background(tcell.ColorTeal).
			Foreground(tcell.ColorGray),
		Action: Base.
			Background(tcell.ColorTeal).
			Foreground(tcell.ColorRed),
	}

	Table = TableTheme{
		TH: LabelTheme{
			Normal:  Base.Bold(true),
			Success: Base.Bold(true),
			Warn:    Base.Bold(true),
			Error:   Base.Bold(true),
		},
		TD: LabelTheme{
			Normal:  Base,
			Success: Base,
			Warn:    Base,
			Error:   Base,
		},
		TDSelected: LabelTheme{
			Normal:  Base.Reverse(true),
			Success: Base.Reverse(true),
			Warn:    Base.Reverse(true),
			Error:   Base.Reverse(true),
		},
	}
)

type AppHeaderTheme struct {
	Bar    tcell.Style
	Action tcell.Style
}

type TableTheme struct {
	TH         LabelTheme
	TD         LabelTheme
	TDSelected LabelTheme
}

type LabelVariant string

type LabelTheme struct {
	Normal  tcell.Style
	Success tcell.Style
	Warn    tcell.Style
	Error   tcell.Style
}

func (t LabelTheme) Get(v LabelVariant) tcell.Style {
	switch v {
	case LabelNormal:
		return t.Normal
	case LabelSuccess:
		return t.Success
	case LabelWarn:
		return t.Warn
	case LabelError:
		return t.Error
	default:
		return t.Normal
	}
}
