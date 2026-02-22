package tui

import (
	"fmt"
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davidlawson7/pokedex/internal/data"
)

// tabIndex names the three detail tabs.
type tabIndex int

const (
	tabStats     tabIndex = 0
	tabMoves     tabIndex = 1
	tabLocations tabIndex = 2
)

var tabNames = [3]string{"Stats", "Moves", "Locations"}

// DetailModel is the tabbed detail view for a single Pokémon.
type DetailModel struct {
	pokemon         *data.Pokemon
	activeTab       tabIndex
	selectedVersion data.GameVersion
	moveScroll      int
	locationScroll  int
	width           int
	height          int
}

// NewDetailModel creates a detail screen for the given pokemon ID.
// If the ID is not found, a zero Pokemon is shown.
func NewDetailModel(pokemonID uint16, width, height int) DetailModel {
	p := data.ByID[pokemonID]
	return DetailModel{
		pokemon:         p,
		activeTab:       tabStats,
		selectedVersion: data.GameRed,
		width:           width,
		height:          height,
	}
}

func (m DetailModel) Init() tea.Cmd { return nil }

func (m DetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch {
		case msg.Type == tea.KeyEsc:
			return m, func() tea.Msg { return switchToSearchMsg{} }

		case msg.Type == tea.KeyTab:
			m.activeTab = (m.activeTab + 1) % 3
			m.moveScroll = 0
			m.locationScroll = 0

		case msg.Type == tea.KeyShiftTab:
			m.activeTab = (m.activeTab + 2) % 3
			m.moveScroll = 0
			m.locationScroll = 0

		case msg.Type == tea.KeyUp:
			switch m.activeTab {
			case tabMoves:
				if m.moveScroll > 0 {
					m.moveScroll--
				}
			case tabLocations:
				if m.locationScroll > 0 {
					m.locationScroll--
				}
			}

		case msg.Type == tea.KeyDown:
			switch m.activeTab {
			case tabMoves:
				m.moveScroll++
			case tabLocations:
				m.locationScroll++
			}

		case msg.Type == tea.KeyRunes:
			// Version keys 1-9 map to GameVersion constants
			if len(msg.Runes) == 1 {
				r := msg.Runes[0]
				if r >= '1' && r <= '9' {
					v := data.GameVersion(r - '0')
					if v <= data.GameLeafGreen {
						m.selectedVersion = v
					}
				}
			}
		}
	}
	return m, nil
}

func (m DetailModel) View() string {
	if m.pokemon == nil {
		return "Pokemon not found."
	}

	gen := data.GenForVersion(m.selectedVersion)
	types := m.pokemon.TypesForGen(gen)

	var sb strings.Builder

	// Header
	name := capitalize(m.pokemon.Name)
	sb.WriteString(fmt.Sprintf("  #%03d %-14s ver: %s\n", m.pokemon.ID, name, m.selectedVersion))

	// Type badges
	t1 := types[0]
	t2 := types[1]
	typeStr := ""
	if t1 != data.TypeNone {
		typeStr += TypeBadge(t1.String()) + " "
	}
	if t2 != data.TypeNone {
		typeStr += TypeBadge(t2.String())
	}
	sb.WriteString("  Type: " + typeStr + "\n")
	sb.WriteString(strings.Repeat("─", max(m.width-2, 40)) + "\n")

	// Tabs
	tabBar := ""
	for i, name := range tabNames {
		if tabIndex(i) == m.activeTab {
			tabBar += activeTabStyle.Render("["+name+"]") + " "
		} else {
			tabBar += inactiveTabStyle.Render(name) + " "
		}
	}
	sb.WriteString("  " + tabBar + "\n")
	sb.WriteString(strings.Repeat("─", max(m.width-2, 40)) + "\n")

	// Tab content
	switch m.activeTab {
	case tabStats:
		sb.WriteString(m.renderStatsTab(gen, types))
	case tabMoves:
		sb.WriteString(m.renderMovesTab(gen))
	case tabLocations:
		sb.WriteString(m.renderLocationsTab())
	}

	// Footer
	sb.WriteString("\n")
	sb.WriteString(footerStyle.Render("  esc:back  tab:switch  1-9:version  ↑↓:scroll"))
	return sb.String()
}

func (m DetailModel) renderStatsTab(gen data.Generation, types [2]data.PokeType) string {
	var sb strings.Builder
	p := m.pokemon

	// Abilities (Gen 3 only)
	if gen >= 3 {
		ab1 := abilityName(p.Abilities[0])
		ab2 := abilityName(p.Abilities[1])
		if ab2 != "" {
			sb.WriteString(fmt.Sprintf("  Ability:  %s / %s\n", ab1, ab2))
		} else if ab1 != "" {
			sb.WriteString(fmt.Sprintf("  Ability:  %s\n", ab1))
		} else {
			sb.WriteString("  Ability:  —\n")
		}
	} else {
		sb.WriteString("  Ability:  (introduced in Gen 3)\n")
	}
	sb.WriteString("\n")

	// Base stats
	stats := []struct {
		label string
		val   uint8
	}{
		{"HP   ", p.Stats.HP},
		{"Atk  ", p.Stats.Attack},
		{"Def  ", p.Stats.Defense},
		{"SpAtk", p.Stats.SpecialAttack},
		{"SpDef", p.Stats.SpecialDefense},
		{"Speed", p.Stats.Speed},
	}
	for _, s := range stats {
		note := ""
		if gen < 2 && (s.label == "SpAtk" || s.label == "SpDef") {
			note = dimStyle.Render(" (= Spc in Gen 1)")
		}
		sb.WriteString(fmt.Sprintf("  %s  %s  %3d%s\n", s.label, StatBar(s.val), s.val, note))
	}
	_ = types
	return sb.String()
}

func (m DetailModel) renderMovesTab(gen data.Generation) string {
	var sb strings.Builder
	p := m.pokemon

	// Collect moves for selectedVersion
	var moves []data.LearnedMove
	noData := true
	for _, vls := range p.Moves {
		if vls.Version == m.selectedVersion {
			moves = vls.Moves
			noData = false
			break
		}
	}

	if noData || len(moves) == 0 {
		if noData {
			sb.WriteString(dimStyle.Render("  No data for this version"))
		} else {
			sb.WriteString(dimStyle.Render("  No moves for this version"))
		}
		sb.WriteString("\n")
		return sb.String()
	}

	// Header row
	sb.WriteString(headerStyle.Render(fmt.Sprintf("  %-14s %-8s %-5s %3s %3s %3s  %-6s\n",
		"Name", "Type", "Cat", "Pwr", "Acc", "PP", "Lv/TM")))

	// Apply scroll
	start := m.moveScroll
	if start > len(moves) {
		start = len(moves)
	}

	for _, lm := range moves[start:] {
		if data.AllMoves == nil || int(lm.MoveID) >= len(data.AllMoves) || data.AllMoves[lm.MoveID] == nil {
			continue
		}
		mv := lm.Move()
		moveType := mv.TypeForGen(gen)
		cat := mv.CategoryForGen(gen)

		power := "—"
		if mv.Power > 0 {
			power = fmt.Sprintf("%3d", mv.Power)
		}
		acc := "—"
		if mv.Accuracy > 0 {
			acc = fmt.Sprintf("%3d", mv.Accuracy)
		}

		lvTM := ""
		switch lm.Method {
		case data.LearnLevelUp:
			if lm.LevelLearnedAt > 0 {
				lvTM = fmt.Sprintf("Lv%3d", lm.LevelLearnedAt)
			} else {
				lvTM = "Lv  1"
			}
		case data.LearnMachine:
			if lm.MachineNumber > 0 {
				lvTM = fmt.Sprintf("TM%02d", lm.MachineNumber)
			} else {
				lvTM = "TM"
			}
		case data.LearnTutor:
			lvTM = "Tutor"
		case data.LearnEgg:
			lvTM = "Egg"
		}

		sb.WriteString(fmt.Sprintf("  %-14s %-8s %-5s %3s %3s %3d  %-6s\n",
			mv.Name, moveType.String(), cat.String(),
			power, acc, mv.PP, lvTM))
	}
	return sb.String()
}

func (m DetailModel) renderLocationsTab() string {
	var sb strings.Builder
	p := m.pokemon

	// Filter locations by selected version
	var locs []data.Location
	for _, loc := range p.Locations {
		if loc.Game == m.selectedVersion {
			locs = append(locs, loc)
		}
	}

	if len(locs) == 0 {
		sb.WriteString(dimStyle.Render("  Not found in the wild for this version"))
		sb.WriteString("\n")
		return sb.String()
	}

	start := m.locationScroll
	if start > len(locs) {
		start = len(locs)
	}

	sb.WriteString(headerStyle.Render(fmt.Sprintf("  %-30s %-12s %-8s %s\n",
		"Area", "Method", "Levels", "Chance")))
	for _, loc := range locs[start:] {
		levels := fmt.Sprintf("%d-%d", loc.MinLevel, loc.MaxLevel)
		sb.WriteString(fmt.Sprintf("  %-30s %-12s %-8s %d%%\n",
			loc.AreaName, loc.EncounterMethod.String(), levels, loc.Chance))
	}
	return sb.String()
}

// abilityName returns the display name for an ability ID, or "" if not found.
func abilityName(id data.AbilityID) string {
	if id == 0 || data.AllAbilities == nil || int(id) >= len(data.AllAbilities) {
		return ""
	}
	ab := data.AllAbilities[id]
	if ab == nil {
		return ""
	}
	return capitalize(strings.ReplaceAll(ab.Name, "-", " "))
}

// capitalize uppercases the first letter of each word.
func capitalize(s string) string {
	words := strings.Fields(s)
	for i, w := range words {
		if len(w) > 0 {
			runes := []rune(w)
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}
	return strings.Join(words, " ")
}

// EncounterMethod.String() for EncounterMethod display
func init() {
	// Ensure EncounterMethod has a String representation.
	// (Already defined via the type's String method below - this is a placeholder.)
}
