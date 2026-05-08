package cmds

import (
	t "github.com/Keivan-sf/Bushuray-tui/shared_types"

	tea "github.com/charmbracelet/bubbletea"
)

func ExitAddProfileView() tea.Msg {
	return t.AddProfileExit{}
}

func EnterAddProfileView() tea.Msg {
	return t.AddProfileEnter{}
}
