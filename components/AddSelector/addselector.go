package addselector

import (
	cmds "github.com/Keivan-sf/Bushuray-tui/commands"
	"github.com/Keivan-sf/Bushuray-tui/components/shared"
	"github.com/Keivan-sf/Bushuray-tui/global"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Width   int
	Height  int
	cursor  int
	options []option
}

type option struct {
	label string
	cmd   tea.Cmd
}

var blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#6b7280"))

func focusedStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(global.GetPrimaryColor())
}

func pointerStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(global.GetPrimaryColor())
}

func InitialModel() Model {
	return Model{
		cursor: 0,
		options: []option{
			{label: "Add Subscription", cmd: cmds.EnterAddProfileView},
			{label: "Add Proxy  (vmess / vless / ...)", cmd: cmds.EnterPasteProfileView},
		},
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, cmds.ExitAddSelectorView
		case "up", "shift+tab":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "tab":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "enter":
			return m, m.options[m.cursor].cmd
		}
	}
	return m, nil
}

func (m Model) View() string {
	bg := lipgloss.NewStyle().Background(global.GetBgColor())
	center := bg.Width(m.Width).Height(1).Align(lipgloss.Center)

	title := lipgloss.NewStyle().
		Foreground(global.GetPrimaryColor()).
		Background(global.GetBgColor()).
		Render("Add Sub / Proxy")
	titleRow := center.Render(title)

	var rows []string
	for i, opt := range m.options {
		var row string
		if i == m.cursor {
			pointer := pointerStyle().Background(global.GetBgColor()).Render("❯ ")
			label := focusedStyle().Background(global.GetBgColor()).Render(opt.label)
			row = center.Render(lipgloss.JoinHorizontal(lipgloss.Top, pointer, label))
		} else {
			pointer := blurredStyle.Background(global.GetBgColor()).Render("  ")
			label := blurredStyle.Background(global.GetBgColor()).Render(opt.label)
			row = center.Render(lipgloss.JoinHorizontal(lipgloss.Top, pointer, label))
		}
		rows = append(rows, row)
	}

	help := shared.GenHelp([]string{"↑↓", "enter", "esc"}, []string{"navigate", "select", "cancel"})
	helpRow := center.Render(help)

	views := []string{titleRow, ""}
	views = append(views, rows...)
	views = append(views, "", helpRow)

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
