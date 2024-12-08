package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/iamhectorsosa/sqlite-shell/internal/database"
	"github.com/iamhectorsosa/sqlite-shell/internal/help"
	"github.com/iamhectorsosa/sqlite-shell/internal/helpers"
	"github.com/iamhectorsosa/sqlite-shell/internal/styles"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Required: <databasePath>")
	}

	databasePath := os.Args[1]

	p := tea.NewProgram(initialModel(databasePath), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type errMsg error

type keyMap struct {
	Tab  key.Binding
	Quit key.Binding
}

var keys = keyMap{
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "toggle input"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc", "quit"),
	),
}

type model struct {
	databasePath   string
	styles         *styles.Styles
	help           string
	err            error
	textInput      textinput.Model
	table          table.Model
	viewportWidth  int
	viewportHeight int
	ready          bool
}

func initialModel(databasePath string) model {
	styles := styles.New()
	textInput := textinput.New()
	textInput.Placeholder = "Write SQL..."
	textInput.PromptStyle = styles.Highlight
	textInput.Cursor.Style = styles.Highlight
	textInput.Focus()

	return model{
		databasePath: databasePath,
		styles:       styles,
		err:          nil,
		ready:        false,
		textInput:    textInput,
		help:         help.New(),
	}
}

func (m model) appBoundaryText(text string) string {
	return lipgloss.PlaceHorizontal(
		m.viewportWidth,
		lipgloss.Left,
		m.styles.BoundaryText.Render(text+" "),
		lipgloss.WithWhitespaceChars("•"),
		m.styles.WhitespaceStyle,
	)
}

func (m model) appErrorText(text string) string {
	return lipgloss.Place(m.viewportWidth, m.viewportHeight,
		lipgloss.Center, lipgloss.Center,
		m.styles.ErrorText.Render(lipgloss.NewStyle().
			Padding(0, 2).
			Width(50).
			Align(lipgloss.Center).
			Render(text),
		),
		lipgloss.WithWhitespaceChars("猫咪"),
		m.styles.WhitespaceBackgroundStyle,
	)
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.ready = true
		}

		m.viewportWidth = msg.Width - 2
		m.viewportHeight = msg.Height - 8

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			if m.table.Focused() {
				m.table.Blur()
				cmd := m.textInput.Focus()
				m.textInput.PromptStyle = m.styles.Highlight
				cmds = append(cmds, cmd)
			} else {
				m.table.Focus()
				m.textInput.Blur()
				m.textInput.PromptStyle = m.styles.Background
			}
		case tea.KeyEnter:
			query := strings.TrimSpace(m.textInput.Value())
			if query == "" {
				return m, nil
			}

			headers, rows, err := database.ExecCmd(m.databasePath, query)
			if err != nil {
				m.err = err
				m.styles.Error()
				m.textInput.PromptStyle = m.styles.Highlight
				m.textInput.Cursor.Style = m.styles.Highlight
				if len(m.table.Rows()) > 0 {
					s := table.DefaultStyles()
					s.Header = m.styles.TableHeader
					s.Selected = m.styles.TableRow
					m.table.SetStyles(s)
				}
				return m, nil
			} else {
				m.err = nil
				m.styles.Reset()
				m.textInput.PromptStyle = m.styles.Highlight
				m.textInput.Cursor.Style = m.styles.Highlight

			}

			m.textInput.Blur()
			m.textInput.PromptStyle = m.styles.Background

			height := len(rows) + 1
			if height > m.viewportHeight {
				height = m.viewportHeight
			}

			t := table.New(
				table.WithColumns(helpers.CreateColumns(headers, rows, m.viewportWidth)),
				table.WithRows(helpers.CreateRows(rows)),
				table.WithFocused(true),
				table.WithHeight(height),
			)

			s := table.DefaultStyles()
			s.Header = m.styles.TableHeader
			s.Selected = m.styles.TableRow
			t.SetStyles(s)

			m.table = t
			m.table.Focus()

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	var cmdTextInput tea.Cmd
	var cmdTable tea.Cmd
	m.textInput, cmdTextInput = m.textInput.Update(msg)
	m.table, cmdTable = m.table.Update(msg)
	cmds = append(cmds, cmdTextInput, cmdTable)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	header := m.appBoundaryText("SQLite Shell")
	if m.err != nil {
		header = m.appBoundaryText("An error has occured")
	}

	content := m.styles.TextBorder.Render(m.textInput.View())
	if len(m.table.Columns()) > 0 && m.err == nil {
		content = fmt.Sprintf("%s\n\n%s", content, m.table.View())
	}

	if m.err != nil {
		content = fmt.Sprintf("%s\n\n%s", content, m.appErrorText(m.err.Error()))
	}

	return m.styles.Base.Render(fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		header,
		content,
		m.appBoundaryText(m.help),
	))
}
