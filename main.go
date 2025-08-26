package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var helpStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("241")).
	Margin(1, 0)

type model struct {
	table    table.Model
	sepaData *SEPAData
	keys     keyMap
	err      error
	copied   bool
	copyMsg  string
}

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Copy   key.Binding
	Quit   key.Binding
	Help   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Copy, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Copy, k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Copy: key.NewBinding(
		key.WithKeys("c", "ctrl+c"),
		key.WithHelp("c", "copy value"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc"),
		key.WithHelp("q/esc", "quit"),
	),
}

func initialModel(filepath string) model {
	sepaData, err := ParseSEPAFile(filepath)
	if err != nil {
		return model{err: err, keys: keys}
	}

	columns := []table.Column{
		{Title: "Category", Width: 20},
		{Title: "Field", Width: 25},
		{Title: "Value", Width: 50},
	}

	rows := make([]table.Row, len(sepaData.Fields))
	for i, field := range sepaData.Fields {
		rows[i] = table.Row{field.Category, field.Field, field.Value}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return model{
		table:    t,
		sepaData: sepaData,
		keys:     keys,
	}
}

type copyMsg struct {
	value   string
	success bool
	err     error
}

type clearCopyMsg struct{}

func copyToClipboard(value string) tea.Cmd {
	return func() tea.Msg {
		err := clipboard.WriteAll(value)
		if err != nil {
			return copyMsg{value: value, success: false, err: err}
		}
		return copyMsg{value: value, success: true}
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Copy):
			if m.sepaData != nil && len(m.sepaData.Fields) > 0 {
				selectedRow := m.table.Cursor()
				if selectedRow < len(m.sepaData.Fields) {
					value := m.sepaData.Fields[selectedRow].Value
					return m, copyToClipboard(value)
				}
			}
		}

	case copyMsg:
		m.copied = true
		if msg.success {
			m.copyMsg = fmt.Sprintf("Copied: %s", truncateString(msg.value, 40))
		} else {
			m.copyMsg = fmt.Sprintf("Copy failed: %v", msg.err)
		}
		// Clear the copy message after a delay
		return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg {
			return clearCopyMsg{}
		})

	case clearCopyMsg:
		// Clear copy message
		m.copied = false
		m.copyMsg = ""
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress 'q' to quit.", m.err)
	}

	var b strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Margin(1, 0).
		Render("SEPA Payment Information")

	b.WriteString(title + "\n")
	b.WriteString(baseStyle.Render(m.table.View()) + "\n")

	// Copy status message
	if m.copied {
		copyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")).
			Margin(1, 0)
		b.WriteString(copyStyle.Render(m.copyMsg) + "\n")
	}

	// Help text
	help := helpStyle.Render("Navigation: ↑/↓ or j/k • Copy value: c • Quit: q/esc")
	b.WriteString(help)

	return b.String()
}

func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: sepa <filepath>")
		fmt.Println("Example: sepa example/SEPA_Example_2024.xml")
		os.Exit(1)
	}

	filepath := os.Args[1]
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' does not exist\n", filepath)
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel(filepath), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
