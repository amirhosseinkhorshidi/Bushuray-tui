package mainmodel

import (
	profilelist "github.com/Keivan-sf/Bushuray-tui/components/ProfileList"
	sharedtypes "github.com/Keivan-sf/Bushuray-tui/shared_types"
	"github.com/Keivan-sf/Bushuray-tui/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func HandleServerNotifs(msg sharedtypes.ServerNotification, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case sharedtypes.Warning:
		return applyWarning(msg, m)
	case sharedtypes.ApplicationState:
		return applyApplicationState(msg, m)
	case sharedtypes.ProfileUpdated:
		return applyProfileUpdated(msg, m)
	case sharedtypes.ProxyStatus:
		return applyStatusChanged(msg, m)
	case sharedtypes.ProfilesAdded:
		return applyProfilesAdded(msg, m)
	case sharedtypes.ProfilesDeleted:
		return applyProfilesDeleted(msg, m)
	case sharedtypes.GroupAdded:
		return applyGroupAdded(msg, m)
	case sharedtypes.GroupDeleted:
		return applyGroupDeleted(msg, m)
	case sharedtypes.SubscriptionUpdated:
		return applySubscriptionUpdated(msg, m)
	}
	return m, nil
}

func findProfile(gid int, id int, m Model) int {
	for i, item := range m.ProfileList.Items {
		if item.GroupId == gid && item.ProfileId == id {
			return i
		}
	}
	return -1
}

func convertProtocolForDisplay(name string) string {
	switch name {
	case "vless":
		return "V-LESS"
	case "vmess":
		return "V-MESS"
	case "socks":
		return "SOCKS5"
	case "shadowsocks":
		return "SHADOW"
	case "trojan":
		return "TROJAN"
	}
	return name
}

func makeProfileItem(p sharedtypes.Profile) profilelist.ProfileItem {
	address, port, transport := utils.ParseURIInfo(p.Uri)
	return profilelist.ProfileItem{
		ProfileId:  p.Id,
		GroupId:    p.GroupId,
		Name:       p.Name,
		Protocol:   convertProtocolForDisplay(p.Protocol),
		TestResult: p.TestResult,
		Uri:        p.Uri,
		Address:    address,
		Port:       port,
		Transport:  transport,
	}
}
