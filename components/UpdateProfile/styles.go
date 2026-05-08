package updateprofile

import (
	"github.com/Keivan-sf/Bushuray-tui/global"

	"github.com/charmbracelet/lipgloss"
)

var (
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#94a3b8"))
	noStyle             = lipgloss.NewStyle()
	grayStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("#6b7280"))
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
)

func placeHolderStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#374151"))
}

func focusedStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(global.GetPrimaryColor())
}

func cursorStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#6b7280"))
}

func focusedButton() string {
	return focusedStyle().Background(global.GetBgColor()).Render("[ Submit ]")
}

func blurredButton() string {
	return grayStyle.Background(global.GetBgColor()).Render("[ Submit ]")
}
