package mainmodel

import (
	sharedtypes "github.com/Keivan-sf/Bushuray-tui/shared_types"

	tea "github.com/charmbracelet/bubbletea"
)

func applyProfilesDeleted(msg sharedtypes.ProfilesDeleted, m Model) (tea.Model, tea.Cmd) {
	for _, profile := range msg.DeletedProfiles {
		if profile.Id == m.ProfileList.PrimaryId && profile.GroupId == m.ProfileList.PrimaryGroupId {
			m.ProfileList.PrimaryId = -1
			m.ProfileList.PrimaryGroupId = -1
		}
		idx := findProfile(profile.GroupId, profile.Id, m)
		if idx == -1 {
			continue
		}
		m.ProfileList.Items = append(m.ProfileList.Items[:idx], m.ProfileList.Items[idx+1:]...)
	}
	m.ProfileList.ResolveInvalidCursor()
	return m, nil
}
