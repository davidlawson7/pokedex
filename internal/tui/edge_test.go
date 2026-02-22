package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davidlawson7/pokedex/internal/data"
)

func TestSearchModel_NoResultsView(t *testing.T) {
	m := newTestSearchModel()
	m.pokemon = testPokemon
	// Send a query that matches nothing
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'z'}})
	sm2 := m2.(SearchModel)
	m3, _ := sm2.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'z'}})
	sm3 := m3.(SearchModel)
	m4, _ := sm3.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'z'}})
	sm4 := m4.(SearchModel)

	view := sm4.View()
	if !strings.Contains(view, "No results") {
		t.Errorf("expected 'No results' in view, got: %q", view)
	}
}

func TestDetailModel_NoLocationsMsg(t *testing.T) {
	// Bulbasaur is a starter — no locations
	p := &data.Pokemon{
		ID:    1,
		Name:  "bulbasaur",
		Types: [2]data.PokeType{data.TypeGrass, data.TypePoison},
		Stats: data.BaseStats{HP: 45, Attack: 49, Defense: 49, SpecialAttack: 65, SpecialDefense: 65, Speed: 45},
		// Locations is nil
	}
	m := buildDetailModel(p)
	m.activeTab = tabLocations
	m.selectedVersion = data.GameRed

	view := m.View()
	if !strings.Contains(view, "Not found in the wild") {
		t.Errorf("expected 'Not found in the wild' in locations tab view, got: %q", view)
	}
}

func TestDetailModel_NoMovesForVersion(t *testing.T) {
	p := &data.Pokemon{
		ID:    99,
		Name:  "testpoke",
		Types: [2]data.PokeType{data.TypeNormal, data.TypeNone},
		Stats: data.BaseStats{},
		Moves: []data.VersionedLearnset{
			// Only has Gold moves, not Red
			{Version: data.GameGold, Moves: []data.LearnedMove{
				{MoveID: 33, Method: data.LearnLevelUp, LevelLearnedAt: 1},
			}},
		},
	}
	if data.ByID == nil {
		data.ByID = make(map[uint16]*data.Pokemon)
		data.ByName = make(map[string]*data.Pokemon)
	}
	data.ByID[99] = p
	data.ByName["testpoke"] = p

	m := NewDetailModel(99, 80, 24)
	m.activeTab = tabMoves
	m.selectedVersion = data.GameRed // Red has no moves for this pokemon

	view := m.View()
	if !strings.Contains(view, "No data for this version") {
		t.Errorf("expected 'No data for this version' in moves tab view, got: %q", view)
	}
}

func TestDetailModel_WindowResize(t *testing.T) {
	m := buildDetailModel(detailTestBulbasaur)

	// Should not panic on any window size including small ones
	for _, size := range [][2]int{{0, 0}, {10, 5}, {200, 50}, {80, 24}} {
		m2, _ := m.Update(tea.WindowSizeMsg{Width: size[0], Height: size[1]})
		dm2 := m2.(DetailModel)
		_ = dm2.View() // must not panic
	}
}

func TestSearchModel_WindowResize(t *testing.T) {
	m := newTestSearchModel()
	// Should not panic
	m2, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	sm2 := m2.(SearchModel)
	_ = sm2.View()
}
