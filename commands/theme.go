package cmds

import (
	t "github.com/Keivan-sf/Bushuray-tui/shared_types"

	tea "github.com/charmbracelet/bubbletea"
)

func EnterThemeView() tea.Msg {
	return t.ThemeViewEnter{}
}

func ExitThemeView() tea.Msg {
	return t.ThemeViewExit{}
}
