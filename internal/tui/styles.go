package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Layout
	borderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))

	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("255"))

	selectedRowStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("255"))

	dimStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	footerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	// Tab styles
	activeTabStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("212")).
		Underline(true)

	inactiveTabStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	// Type badge colors (approximating game palette)
	typeBadgeColors = map[string]lipgloss.Color{
		"Normal":   lipgloss.Color("250"),
		"Fire":     lipgloss.Color("202"),
		"Water":    lipgloss.Color("33"),
		"Grass":    lipgloss.Color("70"),
		"Electric": lipgloss.Color("220"),
		"Ice":      lipgloss.Color("153"),
		"Fighting": lipgloss.Color("124"),
		"Poison":   lipgloss.Color("129"),
		"Ground":   lipgloss.Color("178"),
		"Flying":   lipgloss.Color("105"),
		"Psychic":  lipgloss.Color("205"),
		"Bug":      lipgloss.Color("106"),
		"Rock":     lipgloss.Color("143"),
		"Ghost":    lipgloss.Color("60"),
		"Dragon":   lipgloss.Color("62"),
		"Dark":     lipgloss.Color("95"),
		"Steel":    lipgloss.Color("103"),
	}
)

// TypeBadge renders a colored type badge string.
func TypeBadge(typeName string) string {
	if typeName == "" {
		return ""
	}
	color, ok := typeBadgeColors[typeName]
	if !ok {
		color = lipgloss.Color("250")
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(color).
		Padding(0, 1).
		Render(typeName)
}

// StatBar renders a simple ASCII bar for a base stat value (0-255).
func StatBar(value uint8) string {
	const maxBars = 10
	filled := int(value) * maxBars / 255
	if filled > maxBars {
		filled = maxBars
	}
	bar := ""
	for i := 0; i < maxBars; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	return bar
}
