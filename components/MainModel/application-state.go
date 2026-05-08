package mainmodel

import (
	profilelist "github.com/Keivan-sf/Bushuray-tui/components/ProfileList"
	sharedtypes "github.com/Keivan-sf/Bushuray-tui/shared_types"

	tea "github.com/charmbracelet/bubbletea"
)

func applyApplicationState(msg sharedtypes.ApplicationState, m Model) (tea.Model, tea.Cmd) {
	m.ProfileList.IsConnected = msg.ConnectionStatus.Connection == "connected"

	var groups []profilelist.GroupItem
	for _, gwp := range msg.Groups {
		groups = append(groups, profilelist.GroupItem{
			Id:   gwp.Group.Id,
			Name: gwp.Group.Name,
		})
	}
	m.ProfileList.Groups = groups

	var items []profilelist.ProfileItem
	for _, gwp := range msg.Groups {
		for _, profile := range gwp.Profiles {
			item := makeProfileItem(profile)
			item.TestResult = 0
			items = append(items, item)
		}
	}
	m.ProfileList.Items = items

	m.ProfileList.PrimaryId = -1
	m.ProfileList.PrimaryGroupId = -1
	if m.ProfileList.IsConnected {
		m.ProfileList.PrimaryId = msg.ConnectionStatus.Profile.Id
		m.ProfileList.PrimaryGroupId = msg.ConnectionStatus.Profile.GroupId
	}

	m.ProfileList = m.ProfileList.SetWH(m.Width, m.Height)
	m.ProfileList.JumpToPrimary()
	return m, nil
}
