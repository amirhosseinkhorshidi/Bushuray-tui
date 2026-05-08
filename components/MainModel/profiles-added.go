package mainmodel

import (
	"time"

	notif_publisher "github.com/Keivan-sf/Bushuray-tui/lib/NotifPublisher"
	sharedtypes "github.com/Keivan-sf/Bushuray-tui/shared_types"

	tea "github.com/charmbracelet/bubbletea"
)

func applyProfilesAdded(msg sharedtypes.ProfilesAdded, m Model) (tea.Model, tea.Cmd) {
	existing := make(map[[2]int]struct{}, len(m.ProfileList.Items))
	for _, item := range m.ProfileList.Items {
		existing[[2]int{item.GroupId, item.ProfileId}] = struct{}{}
	}
	for _, profile := range msg.Profiles {
		if _, dup := existing[[2]int{profile.GroupId, profile.Id}]; !dup {
			m.ProfileList.Items = append(m.ProfileList.Items, makeProfileItem(profile))
		}
	}

	if m.expectingManualAdd {
		m.expectingManualAdd = false
		m.ProfileList.LastWarningTime = time.Now()
		if len(msg.Profiles) > 0 {
			m.ProfileList.Warning = "Proxy added"
			m.ProfileList.WarningMode = "success"
		} else {
			m.ProfileList.Warning = "Invalid URI"
			m.ProfileList.WarningMode = "fatal"
		}
		go func() {
			time.Sleep(time.Second * 4)
			notif_publisher.ClearWarningsNotif(sharedtypes.ClearWarnings{})
		}()
	}

	return m, nil
}
