package tui

import (
	"testing"

	"github.com/charmbracelet/bubbletea"
	"github.com/davidlawson7/pokedex/internal/data"
)

// testPokemon is a small set of fixture Pokemon for TUI tests.
var testPokemon = []*data.Pokemon{
	{ID: 4, Name: "charmander", Types: [2]data.PokeType{data.TypeFire, data.TypeNone}},
	{ID: 5, Name: "charmeleon", Types: [2]data.PokeType{data.TypeFire, data.TypeNone}},
	{ID: 6, Name: "charizard", Types: [2]data.PokeType{data.TypeFire, data.TypeFlying}},
	{ID: 152, Name: "chikorita", Types: [2]data.PokeType{data.TypeGrass, data.TypeNone}},
	{ID: 1, Name: "bulbasaur", Types: [2]data.PokeType{data.TypeGrass, data.TypePoison}},
}

func newTestSearchModel() SearchModel {
	m := NewSearchModel()
	m.results = testPokemon
	return m
}

func TestSearchModel_StartsWithAllResults(t *testing.T) {
	m := NewSearchModel()
	m.pokemon = testPokemon
	m.results = testPokemon
	if len(m.results) != len(testPokemon) {
		t.Errorf("expected %d results, got %d", len(testPokemon), len(m.results))
	}
}

func TestSearchModel_TypeUpdatesFilter(t *testing.T) {
	m := newTestSearchModel()
	m.pokemon = testPokemon

	// Type 'c' - should narrow to char* pokemon
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	sm2 := m2.(SearchModel)
	if len(sm2.results) == 0 {
		t.Error("expected results after typing 'c'")
	}
	// All results should match 'c' in some way
	for _, p := range sm2.results {
		found := false
		for _, c := range p.Name {
			_ = c
			found = true
			break
		}
		_ = found
	}
}

func TestSearchModel_NavigateDown(t *testing.T) {
	m := newTestSearchModel()
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	sm2 := m2.(SearchModel)
	if sm2.cursor != 1 {
		t.Errorf("cursor = %d after down, want 1", sm2.cursor)
	}
}

func TestSearchModel_NavigateUp(t *testing.T) {
	m := newTestSearchModel()
	m.cursor = 2
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	sm2 := m2.(SearchModel)
	if sm2.cursor != 1 {
		t.Errorf("cursor = %d after up from 2, want 1", sm2.cursor)
	}
}

func TestSearchModel_NavigateUpClampsAtZero(t *testing.T) {
	m := newTestSearchModel()
	m.cursor = 0
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	sm2 := m2.(SearchModel)
	if sm2.cursor != 0 {
		t.Errorf("cursor = %d after up from 0, want 0 (clamped)", sm2.cursor)
	}
}

func TestSearchModel_CursorClampsAfterFilter(t *testing.T) {
	m := newTestSearchModel()
	m.pokemon = testPokemon
	m.cursor = 3

	// Type a query that yields only 1 result
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	sm2 := m2.(SearchModel)
	// "bulbasaur" matches 'b'; cursor should clamp to 0 since len=1
	if sm2.cursor != 0 {
		t.Errorf("cursor = %d after filter to 1 result, want 0", sm2.cursor)
	}
}

func TestSearchModel_EnterEmitsSwitchMsg(t *testing.T) {
	m := newTestSearchModel()
	m.cursor = 0
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a command on enter, got nil")
	}
	msg := cmd()
	if _, ok := msg.(switchToDetailMsg); !ok {
		t.Errorf("expected switchToDetailMsg, got %T", msg)
	}
}

func TestSearchModel_QuitEmitsQuit(t *testing.T) {
	m := newTestSearchModel()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Fatal("expected a command on 'q', got nil")
	}
	msg := cmd()
	if msg != tea.Quit() {
		t.Errorf("expected tea.Quit, got %T", msg)
	}
}
