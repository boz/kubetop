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
			Background(tcell.Color23).
			Foreground(tcell.ColorWhite),
		Action: Base.
			Background(tcell.Color23).
			Foreground(tcell.ColorRed),
	}

	Popup = PopupTheme{
		Normal:  Base.Foreground(tcell.Color23).Bold(true),
		Success: Base.Foreground(tcell.Color42).Bold(true),
		Warn:    Base.Foreground(tcell.Color220).Bold(true),
		Error:   Base.Foreground(tcell.Color166).Bold(true),
	}

	apbase      = Base
	ThemeActive = Theme{
		Base:  apbase,
		Title: Base.Background(tcell.Color21).Bold(true),
		Label: LabelTheme{
			Normal:  apbase.Foreground(tcell.Color255).Bold(true),
			Success: apbase.Foreground(tcell.Color42).Bold(true),
			Warn:    apbase.Foreground(tcell.Color220).Bold(true),
			Error:   apbase.Foreground(tcell.Color196).Bold(true),
		},
		Table: TableTheme{
			TH: LabelTheme{
				Normal:  apbase.Bold(true),
				Success: apbase.Foreground(tcell.Color42).Bold(true),
				Warn:    apbase.Foreground(tcell.Color220).Bold(true),
				Error:   apbase.Foreground(tcell.Color196).Bold(true),
			},
			TD: LabelTheme{
				Normal:  apbase,
				Success: apbase.Foreground(tcell.Color42),
				Warn:    apbase.Foreground(tcell.Color220),
				Error:   apbase.Foreground(tcell.Color196),
			},
			TDSelected: LabelTheme{
				Normal:  apbase.Background(tcell.Color241).Bold(true),
				Success: apbase.Foreground(tcell.Color42).Background(tcell.Color241).Bold(true),
				Warn:    apbase.Foreground(tcell.Color220).Background(tcell.Color241).Bold(true),
				Error:   apbase.Foreground(tcell.Color196).Background(tcell.Color241).Bold(true),
			},
		},
		Deflist: DeflistTheme{
			Term: LabelTheme{
				Normal:  apbase.Bold(true),
				Success: apbase.Foreground(tcell.Color42).Bold(true),
				Warn:    apbase.Foreground(tcell.Color220).Bold(true),
				Error:   apbase.Foreground(tcell.Color196).Bold(true),
			},
			Definition: LabelTheme{
				Normal:  apbase,
				Success: apbase.Foreground(tcell.Color42),
				Warn:    apbase.Foreground(tcell.Color220),
				Error:   apbase.Foreground(tcell.Color196),
			},
		},
	}

	ipbase        = Base.Foreground(tcell.Color242)
	ThemeInactive = Theme{
		Base:  ipbase,
		Title: ipbase.Background(tcell.Color18).Bold(true),
		Label: LabelTheme{
			Normal:  ipbase,
			Success: ipbase.Foreground(tcell.Color34),
			Warn:    ipbase.Foreground(tcell.Color221),
			Error:   ipbase.Foreground(tcell.Color166),
		},
		Table: TableTheme{
			TH: LabelTheme{
				Normal:  ipbase.Bold(true),
				Success: ipbase.Foreground(tcell.Color34).Bold(true),
				Warn:    ipbase.Foreground(tcell.Color221).Bold(true),
				Error:   ipbase.Foreground(tcell.Color166).Bold(true),
			},
			TD: LabelTheme{
				Normal:  ipbase,
				Success: ipbase.Foreground(tcell.Color34),
				Warn:    ipbase.Foreground(tcell.Color221),
				Error:   ipbase.Foreground(tcell.Color166),
			},
			TDSelected: LabelTheme{
				Normal:  ipbase.Background(tcell.Color235).Bold(true),
				Success: ipbase.Foreground(tcell.Color34).Background(tcell.Color235).Bold(true),
				Warn:    ipbase.Foreground(tcell.Color221).Background(tcell.Color235).Bold(true),
				Error:   ipbase.Foreground(tcell.Color166).Background(tcell.Color235).Bold(true),
			},
		},
		Deflist: DeflistTheme{
			Term: LabelTheme{
				Normal:  ipbase.Bold(true),
				Success: ipbase.Foreground(tcell.Color34).Bold(true),
				Warn:    ipbase.Foreground(tcell.Color221).Bold(true),
				Error:   ipbase.Foreground(tcell.Color166).Bold(true),
			},
			Definition: LabelTheme{
				Normal:  ipbase,
				Success: ipbase.Foreground(tcell.Color34),
				Warn:    ipbase.Foreground(tcell.Color221),
				Error:   ipbase.Foreground(tcell.Color166),
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
