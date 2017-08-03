package theme

import "github.com/gdamore/tcell"

const (
	LabelNormal  LabelVariant = "normal"
	LabelSuccess LabelVariant = "success"
	LabelWarn    LabelVariant = "warn"
	LabelError   LabelVariant = "error"
)

type Themeable interface {
	SetTheme(Theme)
}

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
		Foreground(tcell.ColorWhiteSmoke)
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
				Success: apbase.Foreground(tcell.ColorSpringGreen),
				Warn:    apbase.Foreground(tcell.ColorGold),
				Error:   apbase.Foreground(tcell.ColorCrimson),
			},
			TDSelected: LabelTheme{
				Normal:  apbase.Reverse(true),
				Success: apbase.Foreground(tcell.ColorSpringGreen).Reverse(true),
				Warn:    apbase.Foreground(tcell.ColorGold).Reverse(true),
				Error:   apbase.Foreground(tcell.ColorCrimson).Reverse(true),
			},
		},
		Deflist: DeflistTheme{
			Term: LabelTheme{
				Normal:  apbase.Bold(true),
				Success: apbase.Foreground(tcell.ColorSpringGreen).Bold(true),
				Warn:    apbase.Foreground(tcell.ColorGold).Bold(true),
				Error:   apbase.Foreground(tcell.ColorCrimson).Bold(true),
			},
			Definition: LabelTheme{
				Normal:  apbase,
				Success: apbase.Foreground(tcell.ColorSpringGreen),
				Warn:    apbase.Foreground(tcell.ColorGold),
				Error:   apbase.Foreground(tcell.ColorCrimson),
			},
		},
	}

	ipbase        = Base.Foreground(tcell.ColorDarkSlateGrey)
	ThemeInactive = Theme{
		Base: ipbase,
		Title: ipbase.
			Background(tcell.ColorDimGrey).
			Foreground(tcell.ColorAntiqueWhite),
		Label: LabelTheme{
			Normal:  ipbase,
			Success: ipbase.Foreground(tcell.ColorDarkOliveGreen),
			Warn:    ipbase.Foreground(tcell.ColorSandyBrown),
			Error:   ipbase.Foreground(tcell.ColorLightCoral),
		},
		Table: TableTheme{
			TH: LabelTheme{
				Normal:  ipbase.Bold(true),
				Success: ipbase.Foreground(tcell.ColorDarkOliveGreen).Bold(true),
				Warn:    ipbase.Foreground(tcell.ColorSandyBrown).Bold(true),
				Error:   ipbase.Foreground(tcell.ColorLightCoral).Bold(true),
			},
			TD: LabelTheme{
				Normal:  ipbase,
				Success: ipbase.Foreground(tcell.ColorDarkOliveGreen),
				Warn:    ipbase.Foreground(tcell.ColorSandyBrown),
				Error:   ipbase.Foreground(tcell.ColorLightCoral),
			},
			TDSelected: LabelTheme{
				Normal:  ipbase.Reverse(true),
				Success: ipbase.Foreground(tcell.ColorDarkOliveGreen).Reverse(true),
				Warn:    ipbase.Foreground(tcell.ColorSandyBrown).Reverse(true),
				Error:   ipbase.Foreground(tcell.ColorLightCoral).Reverse(true),
			},
		},
		Deflist: DeflistTheme{
			Term: LabelTheme{
				Normal:  ipbase.Bold(true),
				Success: ipbase.Foreground(tcell.ColorDarkOliveGreen).Bold(true),
				Warn:    ipbase.Foreground(tcell.ColorSandyBrown).Bold(true),
				Error:   ipbase.Foreground(tcell.ColorLightCoral).Bold(true),
			},
			Definition: LabelTheme{
				Normal:  ipbase,
				Success: ipbase.Foreground(tcell.ColorDarkOliveGreen),
				Warn:    ipbase.Foreground(tcell.ColorSandyBrown),
				Error:   ipbase.Foreground(tcell.ColorLightCoral),
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
