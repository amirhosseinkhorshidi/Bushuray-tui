package mainmodel

import (
	profilelist "github.com/Keivan-sf/Bushuray-tui/components/ProfileList"
	servercmds "github.com/Keivan-sf/Bushuray-tui/lib/ServerCommands"
	sharedtypes "github.com/Keivan-sf/Bushuray-tui/shared_types"

	tea "github.com/charmbracelet/bubbletea"
)

func applyGroupAdded(msg sharedtypes.GroupAdded, m Model) (tea.Model, tea.Cmd) {
	m.ProfileList.Groups = append(m.ProfileList.Groups, profilelist.GroupItem{
		Id:   msg.Id,
		Name: msg.Name,
	})
	m.ProfileList.ActiveGroupIdx = len(m.ProfileList.Groups) - 1
	m.ProfileList.ResetCursor()
	if msg.SubscriptionUrl != "" {
		go servercmds.UpdateSubscription(msg.Id)
	}
	return m, nil
}
