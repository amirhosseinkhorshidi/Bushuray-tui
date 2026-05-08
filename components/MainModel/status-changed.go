package mainmodel

import (
	sharedtypes "github.com/Keivan-sf/Bushuray-tui/shared_types"

	tea "github.com/charmbracelet/bubbletea"
)

func applyStatusChanged(msg sharedtypes.ProxyStatus, m Model) (tea.Model, tea.Cmd) {
	m.ProfileList.IsConnected = msg.Connection == "connected"
	if m.ProfileList.IsConnected {
		m.ProfileList.PrimaryId = msg.Profile.Id
		m.ProfileList.PrimaryGroupId = msg.Profile.GroupId
	} else {
		m.ProfileList.PrimaryId = -1
		m.ProfileList.PrimaryGroupId = -1
	}
	return m, nil
}
