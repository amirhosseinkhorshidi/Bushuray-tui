package mainmodel

import (
	"time"

	addprofile "github.com/Keivan-sf/Bushuray-tui/components/AddProfile"

	addselector "github.com/Keivan-sf/Bushuray-tui/components/AddSelector"
	helpview "github.com/Keivan-sf/Bushuray-tui/components/Help"
	pasteprofile "github.com/Keivan-sf/Bushuray-tui/components/PasteProfile"
	profilelist "github.com/Keivan-sf/Bushuray-tui/components/ProfileList"
	themeview "github.com/Keivan-sf/Bushuray-tui/components/ThemeView"
	updateprofile "github.com/Keivan-sf/Bushuray-tui/components/UpdateProfile"
	appconfig "github.com/Keivan-sf/Bushuray-tui/lib/AppConfig"
	sharedtypes "github.com/Keivan-sf/Bushuray-tui/shared_types"

	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

type Model struct {
	Width              int
	Height             int
	ProfileList        profilelist.Model
	AddProfile         addprofile.Model
	AddSelector        addselector.Model
	PasteProfile       pasteprofile.Model
	Help               helpview.Model
	UpdateProfile      updateprofile.Model
	ThemeView          themeview.Model
	ActiveSection      string
	expectingManualAdd bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case sharedtypes.ServerNotification:
		return HandleServerNotifs(msg, m)

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.ProfileList = m.ProfileList.SetWH(msg.Width, msg.Height)
		m.AddProfile = m.AddProfile.SetWH(msg.Width, msg.Height)
		m.AddSelector = m.AddSelector.SetWH(msg.Width, msg.Height)
		m.PasteProfile = m.PasteProfile.SetWH(msg.Width, msg.Height)
		m.Help = m.Help.SetWH(msg.Width, msg.Height)
		m.UpdateProfile = m.UpdateProfile.SetWH(msg.Width, msg.Height)
		m.ThemeView = m.ThemeView.SetWH(msg.Width, msg.Height)
		return m, nil

	case sharedtypes.ClearWarnings:
		if time.Now().Unix()-m.ProfileList.LastWarningTime.Unix() >= 4 {
			m.ProfileList.Warning = ""
		}
		return m, nil

	case sharedtypes.ClearTestResult:
		for i, item := range m.ProfileList.Items {
			if item.GroupId == msg.GroupId && item.ProfileId == msg.ProfileId {
				if time.Since(item.TestResultTime) >= 10*time.Second {
					m.ProfileList.Items[i].TestResult = 0
				}
				break
			}
		}
		return m, nil

	case sharedtypes.AddSelectorEnter:
		m.ActiveSection = "add-selector"
		return m, nil

	case sharedtypes.AddSelectorExit:
		m.ActiveSection = "profile-list"
		return m, nil

	case sharedtypes.AddProfileExit:
		m.ActiveSection = "profile-list"
		return m, nil

	case sharedtypes.AddProfileEnter:
		m.ActiveSection = "add-profile"
		return m, nil

	case sharedtypes.PasteProfileEnter:
		m.ActiveSection = "paste-profile"
		return m, nil

	case sharedtypes.PasteProfileSubmitted:
		m.ActiveSection = "profile-list"
		m.expectingManualAdd = true
		return m, nil

	case sharedtypes.PasteProfileExit:
		m.ActiveSection = "profile-list"
		return m, nil

	case sharedtypes.HelpViewEnter:
		m.ActiveSection = "helpview"
		return m, nil

	case sharedtypes.HelpViewExit:
		m.ActiveSection = "profile-list"
		return m, nil

	case sharedtypes.ThemeViewEnter:
		m.ThemeView = themeview.InitialModel().SetWH(m.Width, m.Height)
		m.ActiveSection = "themeview"
		return m, nil

	case sharedtypes.ThemeViewExit:
		m.ActiveSection = "profile-list"
		return m, nil

	case sharedtypes.UpdateProfileEnter:
		item := m.ProfileList.GetItemUnderCursor()
		if item == nil {
			return m, nil
		}
		m.UpdateProfile.SetProfile(updateprofile.UpdateProfileDetails{
			Id:      item.ProfileId,
			GroupId: item.GroupId,
			Name:    item.Name,
		})
		m.ActiveSection = "update-profile"
		return m, nil

	case sharedtypes.UpdateProfileExit:
		m.ActiveSection = "profile-list"
		return m, nil
	}

	if m.ActiveSection == "add-selector" {
		var cmd tea.Cmd
		m.AddSelector, cmd = m.AddSelector.Update(msg)
		return m, cmd
	}

	if m.ActiveSection == "add-profile" {
		var cmd tea.Cmd
		m.AddProfile, cmd = m.AddProfile.Update(msg)
		return m, cmd
	}

	if m.ActiveSection == "paste-profile" {
		var cmd tea.Cmd
		m.PasteProfile, cmd = m.PasteProfile.Update(msg)
		return m, cmd
	}

	if m.ActiveSection == "helpview" {
		var cmd tea.Cmd
		m.Help, cmd = m.Help.Update(msg)
		return m, cmd
	}

	if m.ActiveSection == "themeview" {
		var cmd tea.Cmd
		m.ThemeView, cmd = m.ThemeView.Update(msg)
		return m, cmd
	}

	if m.ActiveSection == "update-profile" {
		var cmd tea.Cmd
		m.UpdateProfile, cmd = m.UpdateProfile.Update(msg)
		return m, cmd
	}

	if m.ActiveSection == "profile-list" {
		var cmd tea.Cmd
		m.ProfileList, cmd = m.ProfileList.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	if m.ActiveSection == "themeview" {
		return m.ThemeView.View()
	}
	if m.ActiveSection == "helpview" {
		return m.Help.View()
	}
	if m.ActiveSection == "add-selector" {
		return m.AddSelector.View()
	}
	if m.ActiveSection == "add-profile" {
		return m.AddProfile.View()
	}
	if m.ActiveSection == "paste-profile" {
		return m.PasteProfile.View()
	}
	if m.ActiveSection == "update-profile" {
		return m.UpdateProfile.View()
	}
	return zone.Scan(m.ProfileList.View())
}

func InitModel() Model {
	config := appconfig.GetConfig()
	return Model{
		ActiveSection: "profile-list",
		AddProfile:    addprofile.InitialModel(),
		AddSelector:   addselector.InitialModel(),
		PasteProfile:  pasteprofile.InitialModel(),
		Help:          helpview.InitialModel(),
		UpdateProfile: updateprofile.InitialModel(),
		ThemeView:     themeview.InitialModel(),
		ProfileList: profilelist.Model{
			Id:             zone.NewPrefix(),
			PrimaryId:      -1,
			PrimaryGroupId: -1,
			Items:          []profilelist.ProfileItem{},
			Groups:         []profilelist.GroupItem{},
			SocksPort:      config.SocksPort,
			HttpPort:       config.HttpPort,
		},
	}
}
