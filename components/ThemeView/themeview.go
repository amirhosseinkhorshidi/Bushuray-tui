package themeview

import (
	cmds "github.com/Keivan-sf/Bushuray-tui/commands"
	appconfig "github.com/Keivan-sf/Bushuray-tui/lib/AppConfig"
	"github.com/Keivan-sf/Bushuray-tui/components/shared"
	"github.com/Keivan-sf/Bushuray-tui/global"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Width  int
	Height int
	cursor int
}

func InitialModel() Model {
	idx := 0
	current := global.GetTheme()
	for i, t := range global.Themes {
		if t.Name == current.Name {
			idx = i
			break
		}
	}
	return Model{cursor: idx}
}

func (m Model) SetWH(width, height int) Model {
	m.Width = width
	m.Height = height
	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(global.Themes)-1 {
				m.cursor++
			}
		case "enter":
			chosen := global.Themes[m.cursor].Name
			global.SetTheme(chosen)
			go appconfig.SaveTheme(chosen)
			return m, cmds.ExitThemeView
		case "esc", "q", "C":
			return m, cmds.ExitThemeView
		}
	}
	return m, nil
}

func (m Model) View() string {
	bg := global.GetBgColor()
	bgStyle := lipgloss.NewStyle().Background(bg)
	// every row is full terminal width, content centered inside
	fullRow := bgStyle.Width(m.Width).Height(1).Align(lipgloss.Center)

	// fixed inner widths so all rows align as a column
	maxNameW := 0
	maxDescW := 0
	for _, t := range global.Themes {
		if w := lipgloss.Width(t.Name); w > maxNameW {
			maxNameW = w
		}
		if w := lipgloss.Width(t.Description); w > maxDescW {
			maxDescW = w
		}
	}

	title := lipgloss.NewStyle().
		Foreground(global.GetPrimaryColor()).
		Background(bg).
		Render("Color Theme")

	var views []string
	views = append(views, fullRow.Render(title))
	views = append(views, fullRow.Render(""))

	for i, theme := range global.Themes {
		isActive := theme.Name == global.GetTheme().Name
		isSelected := i == m.cursor

		check := bgStyle.Width(2).Render(" ")
		if isActive {
			check = lipgloss.NewStyle().Foreground(global.GetPrimaryColor()).Background(bg).Width(2).Render("✓")
		}

		pointer := bgStyle.Width(2).Render(" ")
		if isSelected {
			pointer = lipgloss.NewStyle().Foreground(theme.Primary).Background(bg).Width(2).Render("❯")
		}

		swatch := lipgloss.NewStyle().Foreground(theme.Primary).Background(bg).Width(2).Render("█")

		var nameStyle, descStyle lipgloss.Style
		if isSelected {
			nameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Background(bg).Bold(true)
			descStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#94a3b8")).Background(bg)
		} else {
			nameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#6b7280")).Background(bg)
			descStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#374151")).Background(bg)
		}
		name := nameStyle.Width(maxNameW + 2).Render(theme.Name)
		desc := descStyle.Width(maxDescW).Render(theme.Description)

		inner := lipgloss.JoinHorizontal(lipgloss.Top, check, pointer, swatch, name, desc)
		views = append(views, fullRow.Render(inner))
	}

	views = append(views, fullRow.Render(""))
	help := shared.GenHelp([]string{"↑↓", "enter", "esc"}, []string{"navigate", "select", "cancel"})
	views = append(views, fullRow.Render(help))

	content := lipgloss.JoinVertical(lipgloss.Top, views...)
	container := bgStyle.Render(content)
	placed := lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, container)
	return lipgloss.NewStyle().Background(bg).Render(placed)
}
