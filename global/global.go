package global

import (
	appconfig "github.com/Keivan-sf/Bushuray-tui/lib/AppConfig"

	"github.com/charmbracelet/lipgloss"
)

const Version = "1.0.5"

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
	// reds
	{Name: "Ember", Description: "burn it down", Primary: "#ef4444", ActiveBg: "#290a0a", CursorBg: "#1a0606"},
	{Name: "Wine", Description: "pour me another", Primary: "#be123c", ActiveBg: "#1f0611", CursorBg: "#14040b"},
	// warm
	{Name: "Copper", Description: "forge and fire", Primary: "#f97316", ActiveBg: "#291600", CursorBg: "#1a0e00"},
	{Name: "Sunset", Description: "golden hour forever", Primary: "#f59e0b", ActiveBg: "#29200a", CursorBg: "#1a1406"},
	// greens
	{Name: "Neon", Description: "too bright to ignore", Primary: "#a3e635", ActiveBg: "#182b00", CursorBg: "#101c00"},
	{Name: "Matrix", Description: "hack the planet", Primary: "#0d9e78", ActiveBg: "#0a2a1f", CursorBg: "#061a12"},
	{Name: "Teal", Description: "calm runs deep", Primary: "#14b8a6", ActiveBg: "#062420", CursorBg: "#031714"},
	// blues
	{Name: "Cyan", Description: "electric ice", Primary: "#06b6d4", ActiveBg: "#042129", CursorBg: "#02141a"},
	{Name: "Ocean", Description: "deep and cold", Primary: "#0369a1", ActiveBg: "#050f1a", CursorBg: "#030a12"},
	{Name: "Indigo", Description: "between worlds", Primary: "#6366f1", ActiveBg: "#0d0d2a", CursorBg: "#08081a"},
	// purples & pinks
	{Name: "Nebula", Description: "lost in space", Primary: "#8b5cf6", ActiveBg: "#160a29", CursorBg: "#0e061a"},
	{Name: "Sakura", Description: "soft but deadly", Primary: "#ec4899", ActiveBg: "#290a15", CursorBg: "#1a060d"},
	// neutral
	{Name: "Silver", Description: "cold steel", Primary: "#94a3b8", ActiveBg: "#141820", CursorBg: "#0d1016"},
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
