package profilelist

import (
	"fmt"
	"strconv"
	"time"

	cmds "github.com/Keivan-sf/Bushuray-tui/commands"
	"github.com/Keivan-sf/Bushuray-tui/components/shared"
	"github.com/Keivan-sf/Bushuray-tui/global"
	notif_publisher "github.com/Keivan-sf/Bushuray-tui/lib/NotifPublisher"
	servercmds "github.com/Keivan-sf/Bushuray-tui/lib/ServerCommands"
	sharedtypes "github.com/Keivan-sf/Bushuray-tui/shared_types"
	"github.com/Keivan-sf/Bushuray-tui/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

type GroupItem struct {
	Id   int
	Name string
}

type ProfileItem struct {
	ProfileId      int
	GroupId        int
	Name           string
	Protocol       string
	TestResult     int
	TestResultTime time.Time
	Uri            string
	Address        string
	Port           string
	Transport      string
}

type Model struct {
	Id             string
	Groups         []GroupItem
	ActiveGroupIdx int
	Items          []ProfileItem
	cursor         int
	offset         int
	Width          int
	Height         int
	PrimaryId      int
	PrimaryGroupId int
	IsConnected    bool
	SocksPort      int
	HttpPort       int
	Warning        string
	WarningMode    string
	LastWarningTime time.Time
}

var (
	protocolW   = 8
	testResultW = 15

	cursorRowStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	protocolCursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Align(lipgloss.Center)
	protocolNormalStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#6b7280")).Align(lipgloss.Center)
	itemNormalStyle     = lipgloss.NewStyle()

	notifWarnColor  = lipgloss.Color("#f59e0b")
	notifFatalColor = lipgloss.Color("#f43f5e")
)

func renderFadedName(text string, w int, fg, f1, f2, f3, bg lipgloss.TerminalColor) string {
	runes := []rune(" " + utils.SanitizeString(text))
	if len(runes) <= w {
		return lipgloss.NewStyle().Foreground(fg).Background(bg).Width(w).MaxWidth(w).MaxHeight(1).Render(string(runes))
	}
	if w < 6 {
		return lipgloss.NewStyle().Foreground(fg).Background(bg).Width(w).MaxWidth(w).MaxHeight(1).Render(string(runes[:w]))
	}
	runes = runes[:w]
	result := lipgloss.NewStyle().Foreground(fg).Background(bg).MaxHeight(1).Render(string(runes[:w-5]))
	result += lipgloss.NewStyle().Foreground(f1).Background(bg).MaxHeight(1).Render(string(runes[w-5]))
	result += lipgloss.NewStyle().Foreground(f2).Background(bg).MaxHeight(1).Render(string(runes[w-4]))
	result += lipgloss.NewStyle().Foreground(f2).Background(bg).MaxHeight(1).Render(string(runes[w-3]))
	result += lipgloss.NewStyle().Foreground(f3).Background(bg).MaxHeight(1).Render(string(runes[w-2]))
	result += lipgloss.NewStyle().Foreground(bg).Background(bg).MaxHeight(1).Render(string(runes[w-1]))
	return result
}

func formatMetaStr(address, port, transport string) string {
	meta := address
	if port != "" {
		if meta != "" {
			meta += "  :" + port
		} else {
			meta = ":" + port
		}
	}
	if transport != "" {
		if meta != "" {
			meta += "  " + transport
		} else {
			meta = transport
		}
	}
	return meta
}

func InitialModel() Model {
	return Model{
		Id:             zone.NewPrefix(),
		PrimaryId:      -1,
		PrimaryGroupId: -1,
		Items:          []ProfileItem{},
		Groups:         []GroupItem{},
	}
}

func (m Model) SetWH(width, height int) Model {
	m.Width = width
	m.Height = height
	return m
}

func (m Model) listHeight() int {
	h := m.Height - 11
	if h < 0 {
		return 0
	}
	return h
}

func (m *Model) activeGroupItems() []ProfileItem {
	if len(m.Groups) == 0 || m.ActiveGroupIdx < 0 || m.ActiveGroupIdx >= len(m.Groups) {
		return nil
	}
	gid := m.Groups[m.ActiveGroupIdx].Id
	var result []ProfileItem
	for _, item := range m.Items {
		if item.GroupId == gid {
			result = append(result, item)
		}
	}
	return result
}

func (m *Model) GetCurrentGroupId() int {
	if len(m.Groups) == 0 || m.ActiveGroupIdx < 0 || m.ActiveGroupIdx >= len(m.Groups) {
		return -1
	}
	return m.Groups[m.ActiveGroupIdx].Id
}

func (m *Model) GetItemUnderCursor() *ProfileItem {
	items := m.activeGroupItems()
	if len(items) == 0 || m.cursor < 0 || m.cursor >= len(items) {
		return nil
	}
	item := items[m.cursor]
	return &item
}

func (m *Model) ResetCursor() {
	m.cursor = 0
	m.offset = 0
}

func (m *Model) ResolveInvalidCursor() {
	items := m.activeGroupItems()
	n := len(items)
	if n == 0 {
		m.cursor = 0
		m.offset = 0
		return
	}
	if m.cursor >= n {
		m.cursor = n - 1
	} else if m.cursor < 0 {
		m.cursor = 0
	}
	m.adjustOffsetForCursor()
}

func (m *Model) JumpToPrimary() {
	if m.PrimaryId == -1 {
		return
	}
	for i, grp := range m.Groups {
		if grp.Id == m.PrimaryGroupId {
			m.ActiveGroupIdx = i
			break
		}
	}
	items := m.activeGroupItems()
	for i, item := range items {
		if item.ProfileId == m.PrimaryId && item.GroupId == m.PrimaryGroupId {
			m.cursor = i
			m.adjustOffsetForCursor()
			return
		}
	}
}

func (m *Model) adjustOffsetForCursor() {
	h := m.listHeight()
	if h <= 0 {
		return
	}
	if m.cursor < m.offset {
		m.offset = m.cursor
	} else if m.cursor >= m.offset+h {
		m.offset = m.cursor - h + 1
	}
}

func (m *Model) adjustCursorForOffset() {
	h := m.listHeight()
	if h <= 0 {
		return
	}
	if m.offset > m.cursor {
		m.cursor = m.offset
	} else if m.offset+h-1 < m.cursor {
		m.cursor = m.offset + h - 1
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left":
			if m.ActiveGroupIdx > 0 {
				m.ActiveGroupIdx--
				m.cursor = 0
				m.offset = 0
			}
		case "right":
			if m.ActiveGroupIdx < len(m.Groups)-1 {
				m.ActiveGroupIdx++
				m.cursor = 0
				m.offset = 0
			}
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
			m.adjustOffsetForCursor()
		case "down":
			items := m.activeGroupItems()
			if m.cursor < len(items)-1 {
				m.cursor++
			}
			m.adjustOffsetForCursor()
		case "enter":
			m.connectToProfile()
		case "d":
			m.deleteProfileUnderCursor()
		case "D":
			if m.GetCurrentGroupId() != 0 {
				m.deleteGroupUnderCursor()
			}
		case "t":
			m.testProfileUnderCursor()
		case "T":
			m.testAll()
		case "a":
			return m, cmds.EnterAddSelectorView
		case "?":
			return m, cmds.EnterHelpView
		case "C":
			return m, cmds.EnterThemeView
		case "U":
			if m.IsConnected {
				break
			}
			gid := m.GetCurrentGroupId()
			if gid == -1 {
				break
			}
			m.Warning = "Updating..."
			m.WarningMode = "warn"
			m.LastWarningTime = time.Now()
			go func() {
				time.Sleep(time.Second * 4)
				notif_publisher.ClearWarningsNotif(sharedtypes.ClearWarnings{})
			}()
			servercmds.UpdateSubscription(gid)
		case "q", "esc":
			return m, tea.Quit
		}
	case tea.MouseMsg:
		switch msg.Button {
		case tea.MouseButtonLeft:
			if msg.Action != tea.MouseActionRelease {
				break
			}
			for i := range m.Groups {
				if zone.Get(m.Id+"tab"+strconv.Itoa(i)).InBounds(msg) {
					m.ActiveGroupIdx = i
					m.cursor = 0
					m.offset = 0
					return m, nil
				}
			}
			items := m.activeGroupItems()
			h := m.listHeight()
			for i := m.offset; i < len(items) && i < m.offset+h; i++ {
				if zone.Get(m.Id+strconv.Itoa(i)).InBounds(msg) {
					m.cursor = i
					break
				}
			}
		case tea.MouseButtonWheelDown:
			if !zone.Get(m.Id).InBounds(msg) {
				break
			}
			items := m.activeGroupItems()
			h := m.listHeight()
			if m.offset < len(items)-h {
				m.offset++
			}
			m.adjustCursorForOffset()
		case tea.MouseButtonWheelUp:
			if !zone.Get(m.Id).InBounds(msg) {
				break
			}
			if m.offset > 0 {
				m.offset--
			}
			m.adjustCursorForOffset()
		}
	}
	return m, nil
}

func (m *Model) connectToProfile() {
	item := m.GetItemUnderCursor()
	if item == nil {
		return
	}
	isPrimary := item.ProfileId == m.PrimaryId && item.GroupId == m.PrimaryGroupId
	if isPrimary {
		servercmds.Disconnect()
	} else {
		m.PrimaryId = item.ProfileId
		m.PrimaryGroupId = item.GroupId
		servercmds.Connect(item.GroupId, item.ProfileId)
	}
}

func (m *Model) deleteProfileUnderCursor() {
	item := m.GetItemUnderCursor()
	if item == nil {
		return
	}
	if item.ProfileId == m.PrimaryId && item.GroupId == m.PrimaryGroupId {
		return
	}
	servercmds.DeleteProfiles([]sharedtypes.ProfileID{{Id: item.ProfileId, GroupId: item.GroupId}})
}

func (m *Model) deleteGroupUnderCursor() {
	gid := m.GetCurrentGroupId()
	if gid == -1 {
		return
	}
	servercmds.DeleteGroup(gid)
}

func (m *Model) testProfileUnderCursor() {
	item := m.GetItemUnderCursor()
	if item == nil {
		return
	}
	for i, flatItem := range m.Items {
		if flatItem.ProfileId == item.ProfileId && flatItem.GroupId == item.GroupId {
			m.Items[i].TestResult = -2
			break
		}
	}
	servercmds.Test(item.GroupId, item.ProfileId)
}

func (m *Model) testAll() {
	gid := m.GetCurrentGroupId()
	if gid == -1 {
		return
	}
	for i, item := range m.Items {
		if item.GroupId == gid {
			m.Items[i].TestResult = -2
			servercmds.Test(item.GroupId, item.ProfileId)
		}
	}
}

func (m Model) View() string {
	bg := lipgloss.NewStyle().Background(global.GetBgColor())
	content := lipgloss.JoinVertical(lipgloss.Top,
		m.renderAppTitle(),
		m.renderHelp(),
		m.renderTabs(),
		m.renderProfiles(),
		m.renderStatusBar(),
	)
	return bg.Render(content)
}

func (m Model) renderAppTitle() string {
	primary := global.GetPrimaryColor()
	title := lipgloss.NewStyle().Foreground(primary).Background(global.GetBgColor()).Render("Bushuray-tui")
	sep := lipgloss.NewStyle().Foreground(primary).Background(global.GetBgColor()).Render(" • ")
	ver := lipgloss.NewStyle().Foreground(primary).Background(global.GetBgColor()).Render("v" + global.Version)
	content := lipgloss.JoinHorizontal(lipgloss.Top, title, sep, ver)
	return lipgloss.NewStyle().Width(m.Width).Height(3).MaxHeight(3).Padding(1, 0).Background(global.GetBgColor()).Align(lipgloss.Center).Render(content)
}

func (m Model) renderTabs() string {
	bg := lipgloss.NewStyle().Background(global.GetBgColor())
	spacer := bg.Width(1).Height(3).Render("")

	if len(m.Groups) == 0 {
		return renderTabFill(m.Width)
	}

	parts := []string{spacer}
	totalWidth := 1

	for i, grp := range m.Groups {
		name := grp.Name
		if grp.Id == 0 {
			name = "Manual"
		}
		var tab string
		if i == m.ActiveGroupIdx {
			tab = zone.Mark(m.Id+"tab"+strconv.Itoa(i), renderActiveTabTitle(name))
		} else {
			tab = zone.Mark(m.Id+"tab"+strconv.Itoa(i), renderTabTitle(name))
		}
		totalWidth += lipgloss.Width(tab)
		parts = append(parts, tab)
	}

	remaining := m.Width - totalWidth
	if remaining > 0 {
		parts = append(parts, renderTabFill(remaining))
	}

	return bg.Render(
		lipgloss.JoinHorizontal(lipgloss.Bottom, parts...),
	)
}

func renderTabTitle(title string) string {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true).
		Foreground(lipgloss.Color("#6b7280")).
		Background(global.GetBgColor()).
		BorderBackground(global.GetBgColor()).
		BorderForeground(lipgloss.Color("#6b7280")).
		Padding(0, 1).Render(title)
}

func renderActiveTabTitle(title string) string {
	primary := global.GetPrimaryColor()
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true).
		Foreground(primary).
		Background(global.GetBgColor()).
		BorderBackground(global.GetBgColor()).
		BorderForeground(primary).
		Padding(0, 1).Render(title)
}

func renderTabFill(width int) string {
	if width <= 0 {
		return ""
	}
	return lipgloss.NewStyle().Background(global.GetBgColor()).Width(width).Height(3).Render("")
}

func (m Model) renderHelp() string {
	var keys, helps []string
	if len(m.Groups) == 0 {
		keys = []string{"a", "?"}
		helps = []string{"add subscription", "help"}
	} else {
		keys = []string{"enter", "a", "←→", "?"}
		helps = []string{"connect / disconnect", "add sub / proxy", "switch", "help"}
	}
	help := shared.GenHelp(keys, helps)
	return lipgloss.NewStyle().Width(m.Width).Height(2).MaxHeight(2).Align(lipgloss.Center).Background(global.GetBgColor()).Render(help)
}

func (m Model) renderProfiles() string {
	h := m.listHeight()
	items := m.activeGroupItems()
	innerWidth := m.Width - 2

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true).
		BorderForeground(global.GetPrimaryColor()).
		BorderBackground(global.GetBgColor()).
		Background(global.GetBgColor()).
		Width(innerWidth)

	if len(m.Groups) == 0 {
		msg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6b7280")).
			Background(global.GetBgColor()).
			Render("No subscriptions yet — press 'a' to add one")
		return borderStyle.Height(h).Align(lipgloss.Center, lipgloss.Center).Render(msg)
	}

	if len(items) == 0 {
		emptyText := "No proxies added yet"
		activeGroup := m.Groups[m.ActiveGroupIdx]
		if activeGroup.Id != 0 {
			emptyText = "No proxy in this subscription"
		}
		msg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6b7280")).
			Background(global.GetBgColor()).
			Render(emptyText)
		return borderStyle.Height(h).Align(lipgloss.Center, lipgloss.Center).Render(msg)
	}

	available := innerWidth - protocolW - testResultW
	nameW := available * 2 / 5
	if nameW < 18 {
		nameW = 18
	}
	if nameW > 40 {
		nameW = 40
	}
	metaW := available - nameW
	if metaW < 0 {
		metaW = 0
	}

	var rows []string
	for i := m.offset; i < len(items) && i < m.offset+h; i++ {
		item := items[i]
		isPrimary := item.ProfileId == m.PrimaryId && item.GroupId == m.PrimaryGroupId
		bg := global.GetBgColor()
		primary := global.GetPrimaryColor()
		metaLimit := metaW - 2
		if metaLimit < 0 {
			metaLimit = 0
		}
		metaStr := "  " + utils.LimitStrLen(formatMetaStr(item.Address, item.Port, item.Transport), metaLimit)
		protocolPrimaryStyle := lipgloss.NewStyle().Background(primary).Foreground(lipgloss.Color("#FFF")).Align(lipgloss.Center)
		var row string
		if isPrimary {
			row = protocolPrimaryStyle.Width(protocolW).Render(item.Protocol) +
				renderFadedName(item.Name, nameW, primary, lipgloss.Color("#6b7280"), lipgloss.Color("#374151"), lipgloss.Color("#1f2937"), bg) +
				itemNormalStyle.Width(metaW).MaxWidth(metaW).MaxHeight(1).Background(bg).Foreground(primary).Bold(true).Render(metaStr) +
				styleTestResult(item.TestResult, "normal")
		} else if i == m.cursor {
			row = protocolCursorStyle.Width(protocolW).Background(bg).Render(item.Protocol) +
				renderFadedName(item.Name, nameW, lipgloss.Color("#ffffff"), lipgloss.Color("#9ca3af"), lipgloss.Color("#6b7280"), lipgloss.Color("#374151"), bg) +
				itemNormalStyle.Width(metaW).MaxWidth(metaW).MaxHeight(1).Background(bg).Foreground(lipgloss.Color("#9ca3af")).Render(metaStr) +
				styleTestResult(item.TestResult, "cursor")
		} else {
			row = protocolNormalStyle.Width(protocolW).Background(bg).Render(item.Protocol) +
				renderFadedName(item.Name, nameW, lipgloss.Color("#6b7280"), lipgloss.Color("#4b5563"), lipgloss.Color("#374151"), lipgloss.Color("#1f2937"), bg) +
				itemNormalStyle.Width(metaW).MaxWidth(metaW).MaxHeight(1).Background(bg).Foreground(lipgloss.Color("#4b5563")).Render(metaStr) +
				styleTestResult(item.TestResult, "normal")
		}
		rows = append(rows, zone.Mark(m.Id+strconv.Itoa(i), row))
	}

	listContent := lipgloss.JoinVertical(lipgloss.Top, rows...)
	innerList := zone.Mark(m.Id, lipgloss.NewStyle().Height(h).MaxHeight(h).Width(innerWidth).Background(global.GetBgColor()).Render(listContent))
	return borderStyle.Render(innerList)
}

func (m Model) renderWarning() string {
	color := notifWarnColor
	if m.WarningMode == "success" {
		color = global.GetPrimaryColor()
	} else if m.WarningMode == "fatal" {
		color = notifFatalColor
	}
	return lipgloss.NewStyle().Height(1).Width(m.Width).Background(global.GetBgColor()).Foreground(color).Render(m.Warning)
}

func (m Model) renderStatusBar() string {
	var dot, label string
	if m.IsConnected {
		primary := global.GetPrimaryColor()
		dot = lipgloss.NewStyle().Foreground(primary).Background(global.GetBgColor()).Render("●")
		label = lipgloss.NewStyle().Foreground(primary).Background(global.GetBgColor()).Render(" CONNECTED")
	} else {
		dot = lipgloss.NewStyle().Foreground(lipgloss.Color("#6b7280")).Background(global.GetBgColor()).Render("●")
		label = lipgloss.NewStyle().Foreground(lipgloss.Color("#6b7280")).Background(global.GetBgColor()).Render(" NOT CONNECTED")
	}

	ports := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6b7280")).
		Background(global.GetBgColor()).
		Render(fmt.Sprintf("   Server: 127.0.0.1  SOCKS5: %d  HTTP: %d", m.SocksPort, m.HttpPort))

	parts := []string{dot, label, ports}

	if m.Warning != "" {
		warnColor := notifWarnColor
		if m.WarningMode == "success" {
			warnColor = global.GetPrimaryColor()
		} else if m.WarningMode == "fatal" {
			warnColor = notifFatalColor
		}
		warn := lipgloss.NewStyle().Foreground(warnColor).Background(global.GetBgColor()).Render("   " + m.Warning)
		parts = append(parts, warn)
	}

	content := lipgloss.JoinHorizontal(lipgloss.Top, parts...)
	return lipgloss.NewStyle().Width(m.Width).Height(1).Background(global.GetBgColor()).Padding(0, 1).Render(content)
}

func styleTestResult(ping int, mode string) string {
	var base lipgloss.Style
	switch mode {
	case "primary":
		base = lipgloss.NewStyle().Background(global.GetPrimaryColor()).Foreground(lipgloss.Color("#FFF"))
	case "cursor":
		base = cursorRowStyle.Background(global.GetBgColor())
	default:
		base = lipgloss.NewStyle().Background(global.GetBgColor())
	}
	w := testResultW
	if ping > 0 {
		return base.Foreground(global.GetPrimaryColor()).Width(w).MaxWidth(w).Render("OK " + strconv.Itoa(ping) + "ms")
	} else if ping == -1 {
		return base.Foreground(lipgloss.Color("#f87171")).Width(w).MaxWidth(w).Render("FAILED")
	} else if ping == -2 {
		return base.Foreground(lipgloss.Color("#94a3b8")).Width(w).MaxWidth(w).Render("TESTING...")
	}
	return base.Width(w).MaxWidth(w).Render("")
}
