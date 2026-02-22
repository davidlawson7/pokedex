package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap holds all key bindings for the application.
type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Enter    key.Binding
	Escape   key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
	Quit     key.Binding
}

// DefaultKeyMap returns the default key bindings.
var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "ctrl+p"),
		key.WithHelp("↑/ctrl+p", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "ctrl+n"),
		key.WithHelp("↓/ctrl+n", "down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "open"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next tab"),
	),
	ShiftTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "prev tab"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
