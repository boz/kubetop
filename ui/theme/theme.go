package theme

import "github.com/gdamore/tcell"

const (
	LabelNormal  LabelVariant = "normal"
	LabelSuccess LabelVariant = "success"
	LabelWarn    LabelVariant = "warn"
	LabelError   LabelVariant = "error"
)

var (
	LabelVariants = []LabelVariant{
		LabelNormal,
		LabelSuccess,
		LabelWarn,
		LabelError,
	}

	Base = tcell.StyleDefault

	AppHeader = AppHeaderTheme{
		Bar: Base.
			Background(tcell.ColorTeal).
			Foreground(tcell.ColorGray),
		Action: Base.
			Background(tcell.ColorTeal).
			Foreground(tcell.ColorRed),
	}

	Popup = PopupTheme{
		Normal:  Base.Foreground(tcell.ColorSeashell).Bold(true),
		Success: Base.Foreground(tcell.ColorSpringGreen).Bold(true),
		Warn:    Base.Foreground(tcell.ColorGold).Bold(true),
		Error:   Base.Foreground(tcell.ColorCrimson).Bold(true),
	}

	apbase = Base.
		Foreground(tcell.ColorWhite)
	ThemeActive = Theme{
		Base: apbase,
		Title: Base.
			Background(tcell.ColorAquaMarine).
			Foreground(tcell.ColorBlack).Bold(true),
		Label: LabelTheme{
			Normal:  apbase.Foreground(tcell.ColorSeashell).Bold(true),
			Success: apbase.Foreground(tcell.ColorSpringGreen).Bold(true),
			Warn:    apbase.Foreground(tcell.ColorGold).Bold(true),
			Error:   apbase.Foreground(tcell.ColorCrimson).Bold(true),
		},
		Table: TableTheme{
			TH: LabelTheme{
				Normal:  apbase.Bold(true),
				Success: apbase.Foreground(tcell.ColorSpringGreen).Bold(true),
				Warn:    apbase.Foreground(tcell.ColorGold).Bold(true),
				Error:   apbase.Foreground(tcell.ColorCrimson).Bold(true),
			},
			TD: LabelTheme{
				Normal:  apbase,
				Success: apbase.Foreground(tcell.ColorSpringGreen).Bold(true),
				Warn:    apbase.Foreground(tcell.ColorGold).Bold(true),
				Error:   apbase.Foreground(tcell.ColorCrimson).Bold(true),
			},
			TDSelected: LabelTheme{
				Normal:  apbase.Reverse(true),
				Success: apbase.Foreground(tcell.ColorSpringGreen).Bold(true).Reverse(true),
				Warn:    apbase.Foreground(tcell.ColorGold).Bold(true).Reverse(true),
				Error:   apbase.Foreground(tcell.ColorCrimson).Bold(true).Reverse(true),
			},
		},
		Deflist: DeflistTheme{
			Term: LabelTheme{
				Normal:  apbase.Bold(true),
				Success: apbase.Bold(true),
				Warn:    apbase.Bold(true),
				Error:   apbase.Bold(true),
			},
			Definition: LabelTheme{
				Normal:  apbase,
				Success: apbase,
				Warn:    apbase,
				Error:   apbase,
			},
		},
	}

	ipbase        = Base
	ThemeInactive = Theme{
		Base: ipbase,
		Title: ipbase.
			Background(tcell.ColorDarkCyan).
			Foreground(tcell.ColorAntiqueWhite),
		Label: LabelTheme{
			Normal:  ipbase,
			Success: ipbase,
			Warn:    ipbase,
			Error:   ipbase,
		},
		Table: TableTheme{
			TH: LabelTheme{
				Normal:  ipbase.Bold(true),
				Success: ipbase.Bold(true),
				Warn:    ipbase.Bold(true),
				Error:   ipbase.Bold(true),
			},
			TD: LabelTheme{
				Normal:  ipbase,
				Success: ipbase,
				Warn:    ipbase,
				Error:   ipbase,
			},
			TDSelected: LabelTheme{
				Normal:  ipbase.Reverse(true),
				Success: ipbase.Reverse(true),
				Warn:    ipbase.Reverse(true),
				Error:   ipbase.Reverse(true),
			},
		},
		Deflist: DeflistTheme{
			Term: LabelTheme{
				Normal:  ipbase.Bold(true),
				Success: ipbase.Bold(true),
				Warn:    ipbase.Bold(true),
				Error:   ipbase.Bold(true),
			},
			Definition: LabelTheme{
				Normal:  ipbase,
				Success: ipbase,
				Warn:    ipbase,
				Error:   ipbase,
			},
		},
	}
)

type AppHeaderTheme struct {
	Bar    tcell.Style
	Action tcell.Style
}

type Theme struct {
	Base    tcell.Style
	Title   tcell.Style
	Label   LabelTheme
	Table   TableTheme
	Deflist DeflistTheme
}

type TableTheme struct {
	TH         LabelTheme
	TD         LabelTheme
	TDSelected LabelTheme
}

type DeflistTheme struct {
	Term       LabelTheme
	Definition LabelTheme
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

type PopupTheme LabelTheme
