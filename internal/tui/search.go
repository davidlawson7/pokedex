package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/davidlawson7/pokedex/internal/data"
	"github.com/davidlawson7/pokedex/internal/search"
)

const maxVisible = 12

// SearchModel is the fuzzy search screen model.
type SearchModel struct {
	input   textinput.Model
	pokemon []*data.Pokemon
	results []*data.Pokemon
	cursor  int
	width   int
	height  int
}

// NewSearchModel creates a new search screen model.
func NewSearchModel() SearchModel {
	ti := textinput.New()
	ti.Placeholder = "Search Pokémon..."
	ti.Focus()

	all := search.Filter("")
	return SearchModel{
		input:   ti,
		pokemon: data.AllPokemon,
		results: all,
	}
}

func (m SearchModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch {
		case msg.Type == tea.KeyCtrlC || (msg.Type == tea.KeyRunes && string(msg.Runes) == "q"):
			return m, tea.Quit

		case msg.Type == tea.KeyUp || (msg.Type == tea.KeyCtrlP):
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case msg.Type == tea.KeyDown || (msg.Type == tea.KeyCtrlN):
			if m.cursor < len(m.results)-1 {
				m.cursor++
			}
			return m, nil

		case msg.Type == tea.KeyEnter:
			if len(m.results) > 0 && m.cursor < len(m.results) {
				id := m.results[m.cursor].ID
				return m, func() tea.Msg { return switchToDetailMsg{pokemonID: id} }
			}
			return m, nil
		}
	}

	// Pass all other messages to the text input
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	// Re-filter based on updated query
	query := m.input.Value()
	m.results = filterPokemon(m.pokemon, query)

	// Clamp cursor
	if m.cursor >= len(m.results) {
		if len(m.results) > 0 {
			m.cursor = len(m.results) - 1
		} else {
			m.cursor = 0
		}
	}

	return m, cmd
}

// filterPokemon filters the pokemon list using the search package.
func filterPokemon(pokemon []*data.Pokemon, query string) []*data.Pokemon {
	// Use the search package's filterOver via the exported Filter,
	// but we want to filter over the given slice (not data.AllPokemon).
	// Re-use the score logic inline for flexibility.
	return search.FilterOver(pokemon, query)
}

func (m SearchModel) View() string {
	var sb strings.Builder

	sb.WriteString("  Search: ")
	sb.WriteString(m.input.View())
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("─", max(m.width-2, 40)))
	sb.WriteString("\n")

	if len(m.results) == 0 {
		sb.WriteString(dimStyle.Render("  No results"))
		sb.WriteString("\n")
	} else {
		start := 0
		if m.cursor >= maxVisible {
			start = m.cursor - maxVisible + 1
		}
		end := start + maxVisible
		if end > len(m.results) {
			end = len(m.results)
		}
		for i := start; i < end; i++ {
			p := m.results[i]
			line := formatSearchResult(p)
			if i == m.cursor {
				sb.WriteString(selectedRowStyle.Render("  > " + line))
			} else {
				sb.WriteString("    " + line)
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n")
	sb.WriteString(footerStyle.Render("  [Enter] open   [↑↓] navigate   [q] quit"))
	return sb.String()
}

func formatSearchResult(p *data.Pokemon) string {
	name := strings.ToUpper(p.Name[:1]) + p.Name[1:]
	types := ""
	t1 := p.Types[0]
	t2 := p.Types[1]
	if t1 != data.TypeNone {
		types += "[" + t1.String() + "]"
	}
	if t2 != data.TypeNone {
		types += "[" + t2.String() + "]"
	}
	return fmt.Sprintf("#%03d %-12s %s", p.ID, name, types)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
