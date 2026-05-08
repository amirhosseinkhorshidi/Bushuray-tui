package cmds

import (
	t "github.com/Keivan-sf/Bushuray-tui/shared_types"

	tea "github.com/charmbracelet/bubbletea"
)

func ExitPasteProfileView() tea.Msg {
	return t.PasteProfileExit{}
}

func EnterPasteProfileView() tea.Msg {
	return t.PasteProfileEnter{}
}

func SubmitPasteProfileView() tea.Msg {
	return t.PasteProfileSubmitted{}
}
