package mainmodel

import (
	"log"
	"time"

	notif_publisher "github.com/Keivan-sf/Bushuray-tui/lib/NotifPublisher"
	sharedtypes "github.com/Keivan-sf/Bushuray-tui/shared_types"

	tea "github.com/charmbracelet/bubbletea"
)

func applyWarning(msg sharedtypes.Warning, m Model) (tea.Model, tea.Cmd) {
	if msg.Key == "update-subscription-failed" {
		m.ProfileList.Warning = "Update failed"
		m.ProfileList.WarningMode = "warn"
		m.ProfileList.LastWarningTime = time.Now()
		go func() {
			time.Sleep(time.Second * 4)
			notif_publisher.ClearWarningsNotif(sharedtypes.ClearWarnings{})
		}()
		return m, nil
	}
	log.Println("[server warning]", msg.Key, msg.Content)
	return m, nil
}
