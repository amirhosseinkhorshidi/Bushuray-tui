package global

import (
	appconfig "github.com/Keivan-sf/Bushuray-tui/lib/AppConfig"

	"github.com/charmbracelet/lipgloss"
)

var NoColor = lipgloss.NoColor{}
var BgColor = lipgloss.Color("#020308")

type ThemeColors struct {
	Name        string
	Description string
	Primary     lipgloss.Color
	ActiveBg    lipgloss.Color // dark bg for connected proxy row
	CursorBg    lipgloss.Color // dark bg for row under cursor
}

var Themes = []ThemeColors{
	{Name: "Matrix", Description: "hack the planet", Primary: "#0d9e78", ActiveBg: "#0a2a1f", CursorBg: "#061a12"},
	{Name: "Ocean", Description: "deep and cold", Primary: "#3b82f6", ActiveBg: "#0a1229", CursorBg: "#06091a"},
	{Name: "Nebula", Description: "lost in space", Primary: "#8b5cf6", ActiveBg: "#160a29", CursorBg: "#0e061a"},
	{Name: "Ember", Description: "burn it down", Primary: "#ef4444", ActiveBg: "#290a0a", CursorBg: "#1a0606"},
	{Name: "Sakura", Description: "soft but deadly", Primary: "#ec4899", ActiveBg: "#290a15", CursorBg: "#1a060d"},
	{Name: "Sunset", Description: "golden hour forever", Primary: "#f59e0b", ActiveBg: "#29200a", CursorBg: "#1a1406"},
}

var currentThemeIdx = 0

func LoadTheme(themeName string) {
	for i, t := range Themes {
		if t.Name == themeName {
			currentThemeIdx = i
			return
		}
	}
	currentThemeIdx = 0
}

func GetTheme() ThemeColors {
	return Themes[currentThemeIdx]
}

func GetPrimaryColor() lipgloss.Color {
	return Themes[currentThemeIdx].Primary
}

func SetTheme(name string) {
	for i, t := range Themes {
		if t.Name == name {
			currentThemeIdx = i
			return
		}
	}
}

func GetBgColor() lipgloss.TerminalColor {
	if appconfig.GetConfig().NoBackground {
		return NoColor
	} else {
		return BgColor
	}
}
