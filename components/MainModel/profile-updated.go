package mainmodel

import (
	"time"

	notif_publisher "github.com/Keivan-sf/Bushuray-tui/lib/NotifPublisher"
	sharedtypes "github.com/Keivan-sf/Bushuray-tui/shared_types"

	tea "github.com/charmbracelet/bubbletea"
)

func applyProfileUpdated(msg sharedtypes.ProfileUpdated, m Model) (tea.Model, tea.Cmd) {
	idx := findProfile(msg.Profile.GroupId, msg.Profile.Id, m)
	if idx != -1 {
		newItem := makeProfileItem(msg.Profile)
		if newItem.TestResult > 0 || newItem.TestResult == -1 {
			newItem.TestResultTime = time.Now()
			gid, pid := newItem.GroupId, newItem.ProfileId
			go func() {
				time.Sleep(10 * time.Second)
				notif_publisher.ClearTestResultNotif(sharedtypes.ClearTestResult{GroupId: gid, ProfileId: pid})
			}()
		}
		m.ProfileList.Items[idx] = newItem
	}
	return m, nil
}
