package mainmodel

import (
	profilelist "github.com/Keivan-sf/Bushuray-tui/components/ProfileList"
	sharedtypes "github.com/Keivan-sf/Bushuray-tui/shared_types"

	tea "github.com/charmbracelet/bubbletea"
)

func applyGroupDeleted(msg sharedtypes.GroupDeleted, m Model) (tea.Model, tea.Cmd) {
	if m.ProfileList.PrimaryGroupId == msg.Id {
		m.ProfileList.PrimaryId = -1
		m.ProfileList.PrimaryGroupId = -1
	}

	deletedIdx := -1
	var newGroups []profilelist.GroupItem
	for i, grp := range m.ProfileList.Groups {
		if grp.Id == msg.Id {
			deletedIdx = i
		} else {
			newGroups = append(newGroups, grp)
		}
	}
	m.ProfileList.Groups = newGroups

	if deletedIdx != -1 && m.ProfileList.ActiveGroupIdx >= deletedIdx {
		m.ProfileList.ActiveGroupIdx--
	}
	if m.ProfileList.ActiveGroupIdx < 0 {
		m.ProfileList.ActiveGroupIdx = 0
	}

	var newItems []profilelist.ProfileItem
	for _, item := range m.ProfileList.Items {
		if item.GroupId != msg.Id {
			newItems = append(newItems, item)
		}
	}
	m.ProfileList.Items = newItems

	m.ProfileList.ResolveInvalidCursor()
	return m, nil
}
