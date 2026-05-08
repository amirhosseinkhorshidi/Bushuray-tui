package cmds

import (
	t "github.com/Keivan-sf/Bushuray-tui/shared_types"

	tea "github.com/charmbracelet/bubbletea"
)

func ExitAddSelectorView() tea.Msg {
	return t.AddSelectorExit{}
}

func EnterAddSelectorView() tea.Msg {
	return t.AddSelectorEnter{}
}
