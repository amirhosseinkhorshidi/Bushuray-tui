package pasteprofile

import (
	cmds "github.com/Keivan-sf/Bushuray-tui/commands"
	"github.com/Keivan-sf/Bushuray-tui/components/shared"
	"github.com/Keivan-sf/Bushuray-tui/global"
	servercmds "github.com/Keivan-sf/Bushuray-tui/lib/ServerCommands"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const manualGroupId = 0

type Model struct {
	Width      int
	Height     int
	input      textinput.Model
	focusOnBtn bool
}

var blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#6b7280"))

func focusedStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(global.GetPrimaryColor())
}

func placeHolderStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#374151"))
}

func InitialModel() Model {
	t := textinput.New()
	t.Placeholder = "vmess://... or vless://..."
	t.Focus()
	t.TextStyle = focusedStyle().Background(global.GetBgColor())
	t.PromptStyle = focusedStyle().Background(global.GetBgColor())
	t.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#6b7280")).Background(global.GetBgColor())
	t.Cursor.TextStyle = focusedStyle().Background(global.GetBgColor())
	t.PlaceholderStyle = placeHolderStyle().Background(global.GetBgColor())
	t.Width = 50
	t.CharLimit = 2000

	return Model{input: t}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.reset()
			return m, cmds.ExitPasteProfileView
		case "down":
			if !m.focusOnBtn {
				m.focusOnBtn = true
				m.input.Blur()
				m.input.PromptStyle = blurredStyle.Background(global.GetBgColor())
				m.input.TextStyle = blurredStyle.Background(global.GetBgColor())
				return m, nil
			}
		case "up":
			if m.focusOnBtn {
				m.focusOnBtn = false
				m.input.Focus()
				m.input.PromptStyle = focusedStyle().Background(global.GetBgColor())
				m.input.TextStyle = focusedStyle().Background(global.GetBgColor())
				return m, nil
			}
		case "enter":
			uri := m.input.Value()
			if uri != "" {
				servercmds.AddProfiles(uri, manualGroupId)
				m.reset()
				return m, cmds.SubmitPasteProfileView
			}
		}
	}

	if !m.focusOnBtn {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *Model) reset() {
	m.input.Reset()
	m.focusOnBtn = false
	m.input.Focus()
	m.input.PromptStyle = focusedStyle().Background(global.GetBgColor())
	m.input.TextStyle = focusedStyle().Background(global.GetBgColor())
}

func (m Model) View() string {
	bg := lipgloss.NewStyle().Background(global.GetBgColor())
	center := bg.Width(m.Width).Height(1).Align(lipgloss.Center)

	title := lipgloss.NewStyle().
		Foreground(global.GetPrimaryColor()).
		Background(global.GetBgColor()).
		Render("Add Proxy")
	titleRow := center.Render(title)

	inputRow := center.Render(m.input.View())

	var btnStr string
	if m.focusOnBtn {
		btnStr = focusedStyle().Background(global.GetBgColor()).Render("[ Submit ]")
	} else {
		btnStr = blurredStyle.Background(global.GetBgColor()).Render("[ Submit ]")
	}
	btnRow := bg.Width(m.Width).Height(1).Align(lipgloss.Center).Render(
		bg.Width(54).Render(btnStr),
	)

	help := shared.GenHelp([]string{"↑↓", "enter", "esc"}, []string{"navigate", "submit", "cancel"})
	helpRow := center.Render(bg.Width(54).Render(help))

	views := []string{titleRow, "", inputRow, "", btnRow, "", helpRow}
	content := lipgloss.JoinVertical(lipgloss.Top, views...)
	container := bg.Render(content)
	placed := lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, container)
	return lipgloss.NewStyle().Background(global.GetBgColor()).Render(placed)
}

func (m Model) SetWH(width, height int) Model {
	m.Width = width
	m.Height = height
	return m
}
