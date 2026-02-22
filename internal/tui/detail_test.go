package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davidlawson7/pokedex/internal/data"
)

// buildDetailModel creates a DetailModel with a hand-crafted Pokemon for testing.
func buildDetailModel(p *data.Pokemon) DetailModel {
	// Register the pokemon in the data.ByID map for the test.
	if data.ByID == nil {
		data.ByID = make(map[uint16]*data.Pokemon)
		data.ByName = make(map[string]*data.Pokemon)
	}
	data.ByID[p.ID] = p
	data.ByName[p.Name] = p
	m := NewDetailModel(p.ID, 80, 24)
	return m
}

var detailTestMagnemite = &data.Pokemon{
	ID:    81,
	Name:  "magnemite",
	Types: [2]data.PokeType{data.TypeElectric, data.TypeSteel},
	PastTypes: []data.PokemonTypePast{
		{UntilGen: 1, Types: [2]data.PokeType{data.TypeElectric, data.TypeNone}},
	},
	Stats:     data.BaseStats{HP: 25, Attack: 35, Defense: 70, SpecialAttack: 95, SpecialDefense: 55, Speed: 45},
	Abilities: [2]data.AbilityID{42, 5},
}

var detailTestBulbasaur = &data.Pokemon{
	ID:        1,
	Name:      "bulbasaur",
	Types:     [2]data.PokeType{data.TypeGrass, data.TypePoison},
	Stats:     data.BaseStats{HP: 45, Attack: 49, Defense: 49, SpecialAttack: 65, SpecialDefense: 65, Speed: 45},
	Abilities: [2]data.AbilityID{65, 0},
}

var detailTestCharizardWithBite = &data.Pokemon{
	ID:    6,
	Name:  "charizard",
	Types: [2]data.PokeType{data.TypeFire, data.TypeFlying},
	Stats: data.BaseStats{HP: 78, Attack: 84, Defense: 78, SpecialAttack: 109, SpecialDefense: 85, Speed: 100},
	Moves: []data.VersionedLearnset{
		{
			Version: data.GameRed,
			Moves: []data.LearnedMove{
				{MoveID: 44, Method: data.LearnLevelUp, LevelLearnedAt: 33},
			},
		},
	},
}

func setupMovesForTest() {
	if data.AllMoves == nil {
		data.AllMoves = make([]*data.Move, 50)
	}
	if len(data.AllMoves) <= 44 {
		newMoves := make([]*data.Move, 50)
		copy(newMoves, data.AllMoves)
		data.AllMoves = newMoves
	}
	data.AllMoves[44] = &data.Move{
		ID:        44,
		Name:      "Bite",
		Type:      data.TypeDark,
		Category:  data.CategoryPhysical,
		Power:     60,
		Accuracy:  100,
		PP:        25,
		PastTypes: []data.MoveTypePast{{UntilGen: 1, Type: data.TypeNormal}},
	}
}

func setupAbilitiesForTest() {
	if data.AllAbilities == nil {
		data.AllAbilities = make([]*data.Ability, 100)
	}
	if len(data.AllAbilities) <= 65 {
		newAbs := make([]*data.Ability, 100)
		copy(newAbs, data.AllAbilities)
		data.AllAbilities = newAbs
	}
	data.AllAbilities[65] = &data.Ability{ID: 65, Name: "overgrow", ShortDesc: "Powers up Grass-type moves when the Pokémon's HP is low."}
	data.AllAbilities[42] = &data.Ability{ID: 42, Name: "magnet-pull", ShortDesc: "Prevents Steel-type Pokémon from fleeing."}
	data.AllAbilities[5] = &data.Ability{ID: 5, Name: "sturdy", ShortDesc: "The Pokémon is unaffected by one-hit KO moves."}
}

func TestDetailModel_TabCycles(t *testing.T) {
	m := buildDetailModel(detailTestBulbasaur)
	if m.activeTab != tabStats {
		t.Errorf("initial tab = %d, want tabStats (0)", m.activeTab)
	}
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	dm2 := m2.(DetailModel)
	if dm2.activeTab != tabMoves {
		t.Errorf("after tab: activeTab = %d, want tabMoves (1)", dm2.activeTab)
	}
	m3, _ := dm2.Update(tea.KeyMsg{Type: tea.KeyTab})
	dm3 := m3.(DetailModel)
	if dm3.activeTab != tabLocations {
		t.Errorf("after tab: activeTab = %d, want tabLocations (2)", dm3.activeTab)
	}
	m4, _ := dm3.Update(tea.KeyMsg{Type: tea.KeyTab})
	dm4 := m4.(DetailModel)
	if dm4.activeTab != tabStats {
		t.Errorf("after tab: activeTab = %d, want tabStats (0) (wrapped)", dm4.activeTab)
	}
}

func TestDetailModel_ShiftTabCyclesBack(t *testing.T) {
	m := buildDetailModel(detailTestBulbasaur)
	// 0 → shift+tab → 2 → shift+tab → 1 → shift+tab → 0
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	dm2 := m2.(DetailModel)
	if dm2.activeTab != tabLocations {
		t.Errorf("after shift+tab from 0: activeTab = %d, want tabLocations (2)", dm2.activeTab)
	}
	m3, _ := dm2.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	dm3 := m3.(DetailModel)
	if dm3.activeTab != tabMoves {
		t.Errorf("after shift+tab from 2: activeTab = %d, want tabMoves (1)", dm3.activeTab)
	}
	m4, _ := dm3.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	dm4 := m4.(DetailModel)
	if dm4.activeTab != tabStats {
		t.Errorf("after shift+tab from 1: activeTab = %d, want tabStats (0)", dm4.activeTab)
	}
}

func TestDetailModel_EscapeEmitsSwitchMsg(t *testing.T) {
	m := buildDetailModel(detailTestBulbasaur)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected cmd on esc, got nil")
	}
	msg := cmd()
	if _, ok := msg.(switchToSearchMsg); !ok {
		t.Errorf("expected switchToSearchMsg, got %T", msg)
	}
}

func TestDetailModel_VersionKey1SetsRed(t *testing.T) {
	m := buildDetailModel(detailTestBulbasaur)
	m.selectedVersion = data.GameEmerald // start somewhere else
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	dm2 := m2.(DetailModel)
	if dm2.selectedVersion != data.GameRed {
		t.Errorf("selectedVersion = %v, want GameRed", dm2.selectedVersion)
	}
}

func TestDetailModel_VersionKey5SetsSilver(t *testing.T) {
	m := buildDetailModel(detailTestBulbasaur)
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}})
	dm2 := m2.(DetailModel)
	if dm2.selectedVersion != data.GameSilver {
		t.Errorf("selectedVersion = %v, want GameSilver", dm2.selectedVersion)
	}
}

func TestDetailModel_VersionKey6SetsCrystal(t *testing.T) {
	m := buildDetailModel(detailTestBulbasaur)
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'6'}})
	dm2 := m2.(DetailModel)
	if dm2.selectedVersion != data.GameCrystal {
		t.Errorf("selectedVersion = %v, want GameCrystal", dm2.selectedVersion)
	}
}

func TestDetailModel_TypesVersionAware_Magnemite(t *testing.T) {
	m := buildDetailModel(detailTestMagnemite)

	// Set to Red (Gen 1) → should show Electric only
	m.selectedVersion = data.GameRed
	view := m.View()
	if !strings.Contains(view, "Electric") {
		t.Error("expected Electric in view for Gen 1 Magnemite")
	}
	// Steel should not appear as a type badge in Gen 1
	// (it appears in the header area; verify it's NOT there in type line)
	// We check the type row specifically
	gen1 := detailTestMagnemite.TypesForGen(1)
	if gen1[1] != data.TypeNone {
		t.Errorf("Gen 1 Magnemite type2 = %v, want TypeNone", gen1[1])
	}

	// Set to Emerald (Gen 3) → should show Electric/Steel
	m.selectedVersion = data.GameEmerald
	gen3 := detailTestMagnemite.TypesForGen(3)
	if gen3[0] != data.TypeElectric || gen3[1] != data.TypeSteel {
		t.Errorf("Gen 3 Magnemite types = %v/%v, want Electric/Steel", gen3[0], gen3[1])
	}
}

func TestDetailModel_MoveCategoryVersionAware(t *testing.T) {
	setupMovesForTest()
	m := buildDetailModel(detailTestCharizardWithBite)

	// Gen 1 (Red): Bite is Normal/Physical (type-based split, Normal = Physical)
	m.selectedVersion = data.GameRed
	m.activeTab = tabMoves
	bite := data.AllMoves[44]
	gen1Type := bite.TypeForGen(1)
	if gen1Type != data.TypeNormal {
		t.Errorf("Bite type in Gen 1 = %v, want TypeNormal", gen1Type)
	}
	gen1Cat := bite.CategoryForGen(1)
	if gen1Cat != data.CategoryPhysical {
		t.Errorf("Bite category in Gen 1 = %v, want Physical", gen1Cat)
	}

	// Gen 2 (Gold): Bite is Dark/Physical (Dark = Physical in type split)
	gen2Type := bite.TypeForGen(2)
	if gen2Type != data.TypeDark {
		t.Errorf("Bite type in Gen 2 = %v, want TypeDark", gen2Type)
	}
	gen2Cat := bite.CategoryForGen(2)
	if gen2Cat != data.CategoryPhysical {
		t.Errorf("Bite category in Gen 2 = %v, want Physical", gen2Cat)
	}
}

func TestDetailModel_AbilitiesHiddenInGen1(t *testing.T) {
	setupAbilitiesForTest()
	m := buildDetailModel(detailTestMagnemite)
	m.selectedVersion = data.GameRed
	view := m.View()
	if strings.Contains(view, "Magnet Pull") || strings.Contains(view, "Sturdy") {
		t.Error("abilities should not be shown in Gen 1 view")
	}
	if !strings.Contains(view, "Gen 3") {
		t.Error("view should indicate abilities introduced in Gen 3")
	}
}

func TestDetailModel_AbilitiesShownInGen3(t *testing.T) {
	setupAbilitiesForTest()
	m := buildDetailModel(detailTestMagnemite)
	m.selectedVersion = data.GameEmerald
	view := m.View()
	if !strings.Contains(view, "Magnet") {
		t.Error("expected magnet-pull ability in Gen 3 view")
	}
}

func TestDetailModel_ScrollMoves(t *testing.T) {
	setupMovesForTest()
	m := buildDetailModel(detailTestCharizardWithBite)
	m.activeTab = tabMoves

	// scroll down
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	dm2 := m2.(DetailModel)
	if dm2.moveScroll != 1 {
		t.Errorf("moveScroll = %d after down, want 1", dm2.moveScroll)
	}

	// scroll up
	m3, _ := dm2.Update(tea.KeyMsg{Type: tea.KeyUp})
	dm3 := m3.(DetailModel)
	if dm3.moveScroll != 0 {
		t.Errorf("moveScroll = %d after up, want 0", dm3.moveScroll)
	}

	// scroll up from 0 clamps
	m4, _ := dm3.Update(tea.KeyMsg{Type: tea.KeyUp})
	dm4 := m4.(DetailModel)
	if dm4.moveScroll != 0 {
		t.Errorf("moveScroll = %d after up from 0, want 0 (clamped)", dm4.moveScroll)
	}
}
