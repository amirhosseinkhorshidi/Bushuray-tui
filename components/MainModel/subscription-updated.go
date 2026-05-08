package mainmodel

import (
	"time"

	profilelist "github.com/Keivan-sf/Bushuray-tui/components/ProfileList"
	notif_publisher "github.com/Keivan-sf/Bushuray-tui/lib/NotifPublisher"
	sharedtypes "github.com/Keivan-sf/Bushuray-tui/shared_types"

	tea "github.com/charmbracelet/bubbletea"
)

func applySubscriptionUpdated(msg sharedtypes.SubscriptionUpdated, m Model) (tea.Model, tea.Cmd) {
	primaryId := m.ProfileList.PrimaryId
	primaryGid := m.ProfileList.PrimaryGroupId

	var newItems []profilelist.ProfileItem
	for _, item := range m.ProfileList.Items {
		if item.GroupId != msg.GroupId {
			newItems = append(newItems, item)
		}
	}
	for _, profile := range msg.Profiles {
		newItems = append(newItems, makeProfileItem(profile))
	}
	m.ProfileList.Items = newItems

	// If primary was in the updated group, verify it still exists
	if primaryId != -1 && primaryGid == msg.GroupId {
		found := false
		for _, item := range newItems {
			if item.ProfileId == primaryId && item.GroupId == primaryGid {
				found = true
				break
			}
		}
		if !found {
			m.ProfileList.PrimaryId = -1
			m.ProfileList.PrimaryGroupId = -1
		}
	}

	m.ProfileList.ResolveInvalidCursor()
	m.ProfileList.Warning = "Sub updated"
	m.ProfileList.WarningMode = "success"
	m.ProfileList.LastWarningTime = time.Now()
	go func() {
		time.Sleep(time.Second * 4)
		notif_publisher.ClearWarningsNotif(sharedtypes.ClearWarnings{})
	}()
	return m, nil
}
