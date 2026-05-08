package list

import (
	"fmt"
	"strconv"

	"github.com/Keivan-sf/Bushuray-tui/global"
	"github.com/Keivan-sf/Bushuray-tui/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

type ListItem struct {
	ProfileId  int
	Name       string
	Protocol   string
	TestResult int
	Uri        string
}

type Model struct {
	Id      string
	GroupId int
	Items   []ListItem
	cursor  int
	Width   int
	Height  int
	offset  int
	Primary int
}

var protocol_w = 8
var test_result_w = 15

func primaryStyle() lipgloss.Style {
	return lipgloss.NewStyle().Background(global.GetTheme().ActiveBg).Foreground(global.GetPrimaryColor()).Bold(true)
}
func protocolPrimaryStyle() lipgloss.Style {
	return lipgloss.NewStyle().Background(global.GetPrimaryColor()).Foreground(lipgloss.Color("#FFF")).Width(protocol_w).Align(lipgloss.Center)
}
func underCursorStyle() lipgloss.Style {
	return lipgloss.NewStyle().Background(global.GetTheme().CursorBg).Foreground(lipgloss.Color("#FFF"))
}
func protocolUnderCursorStyle() lipgloss.Style {
	return lipgloss.NewStyle().Background(global.GetTheme().CursorBg).Foreground(lipgloss.Color("#94a3b8")).Width(protocol_w).Align(lipgloss.Center)
}

var protocol_style = lipgloss.NewStyle().Foreground(lipgloss.Color("#6b7280")).Width(protocol_w).Align(lipgloss.Center)
var item_style = lipgloss.NewStyle()

func (l Model) GetItemUnderCursor() (ListItem, error) {
	if len(l.Items) <= 0 {
		return ListItem{}, fmt.Errorf("no items in the list")
	}
	return l.Items[l.cursor], nil
}

func (l Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if l.cursor > 0 {
				l.cursor--
			}
			l.adjustOffsetForCursor()
		case "down":
			if l.cursor < len(l.Items)-1 {
				l.cursor++
			}
			l.adjustOffsetForCursor()
		case "T":
			l.testGroup()
		case "t":
			l.testProfile()
		case "delete", "d":
			l.deleteProfileUnderCursor()
		case "enter":
			l.connectToProfile()
		}
	case tea.MouseMsg:
		switch msg.Button {
		case tea.MouseButtonLeft:
			for i := range l.Items {
				if zone.Get(l.Id + strconv.Itoa(i)).InBounds(msg) {
					l.cursor = i
				}
			}
		case tea.MouseButtonWheelDown:
			if !zone.Get(l.Id).InBounds(msg) {
				break
			}
			if l.offset < len(l.Items)-l.Height {
				l.offset++
			}
			l.adjustCursorForOffset()
		case tea.MouseButtonWheelUp:
			if !zone.Get(l.Id).InBounds(msg) {
				break
			}
			if l.offset > 0 {
				l.offset--
			}
			l.adjustCursorForOffset()
		}
	}

	return l, nil
}

func (l Model) View() string {
	var s []string
	start := l.offset
	end := len(l.Items)
	w := l.Width - protocol_w - test_result_w
	for i := start; i < end; i++ {
		row := ""
		item := l.Items[i]
		item_str := fmt.Sprintf(" %s", utils.LimitStrLen(utils.SanitizeString(item.Name), w-2))
		if i == l.Primary {
			protocol := protocolPrimaryStyle().Render(item.Protocol)
			row += protocol
			row += primaryStyle().Width(w).MaxWidth(w).MaxHeight(1).Render(item_str) + styleTestPrimary(item.TestResult)
		} else if i == l.cursor {
			protocol := protocolUnderCursorStyle().Render(item.Protocol)
			row += protocol
			row += underCursorStyle().Width(w).MaxWidth(w).MaxHeight(1).Render(item_str) + styleTestUnderCursor(item.TestResult)
		} else {
			protocol := protocol_style.Background(global.GetBgColor()).Render(item.Protocol)
			row += protocol
			row += item_style.Width(w).MaxWidth(w).MaxHeight(1).Background(global.GetBgColor()).Render(item_str) + styleTestNormal(item.TestResult)
		}
		s = append(s, zone.Mark(l.Id+strconv.Itoa(i), row))
	}
	list_view := zone.Mark(l.Id, lipgloss.NewStyle().Height(l.Height).MaxHeight(l.Height).Render(lipgloss.JoinVertical(lipgloss.Top, s...)))
	return list_view
}
