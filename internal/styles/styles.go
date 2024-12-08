package styles

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	aquamarine    = lipgloss.Color("122")
	outer_space   = lipgloss.Color("238")
	dark_charcoal = lipgloss.Color("236")
	strong_red    = lipgloss.Color("161")
	white         = lipgloss.Color("231")
)

type Styles struct {
	Base                      lipgloss.Style
	BoundaryText              lipgloss.Style
	ErrorText                 lipgloss.Style
	Highlight                 lipgloss.Style
	Background                lipgloss.Style
	TableHeader               lipgloss.Style
	TableRow                  lipgloss.Style
	TextBorder                lipgloss.Style
	WhitespaceStyle           lipgloss.WhitespaceOption
	WhitespaceBackgroundStyle lipgloss.WhitespaceOption
	accent                    lipgloss.Color
	foreground                lipgloss.Color
}

func setStyles(s *Styles) {
	s.Base = lipgloss.NewStyle().
		Padding(0, 1)
	s.BoundaryText = lipgloss.NewStyle().
		Foreground(s.accent).
		Bold(true)
	s.ErrorText = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.accent).
		Padding(1, 0).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)
	s.Highlight = lipgloss.NewStyle().
		Foreground(s.accent)
	s.Background = lipgloss.NewStyle().Foreground(outer_space)
	s.TableHeader = table.DefaultStyles().Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(outer_space).
		BorderBottom(true).
		Bold(false)
	s.TableRow = table.DefaultStyles().Selected.
		Foreground(s.foreground).
		Background(s.accent).
		Bold(false)
	s.TextBorder = lipgloss.NewStyle().
		PaddingLeft(1).
		BorderStyle(lipgloss.ThickBorder()).
		BorderLeft(true).
		BorderForeground(outer_space)
	s.WhitespaceStyle = lipgloss.WithWhitespaceForeground(s.accent)
	s.WhitespaceBackgroundStyle = lipgloss.WithWhitespaceForeground(dark_charcoal)
}

func New() *Styles {
	s := Styles{}

	s.accent = aquamarine
	s.foreground = outer_space
	setStyles(&s)

	return &s
}

func (s *Styles) Error() {
	s.accent = strong_red
	s.foreground = white

	setStyles(s)
}

func (s *Styles) Reset() {
	s.accent = aquamarine
	s.foreground = outer_space

	setStyles(s)
}
