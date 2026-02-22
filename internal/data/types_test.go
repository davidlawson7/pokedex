package data

import "testing"

// GenForVersion tests

func TestGenForVersion_Gen1(t *testing.T) {
	for _, v := range []GameVersion{GameRed, GameBlue, GameYellow} {
		if got := GenForVersion(v); got != 1 {
			t.Errorf("GenForVersion(%v) = %d, want 1", v, got)
		}
	}
}

func TestGenForVersion_Gen2(t *testing.T) {
	for _, v := range []GameVersion{GameGold, GameSilver, GameCrystal} {
		if got := GenForVersion(v); got != 2 {
			t.Errorf("GenForVersion(%v) = %d, want 2", v, got)
		}
	}
}

func TestGenForVersion_Gen3(t *testing.T) {
	for _, v := range []GameVersion{GameRuby, GameSapphire, GameEmerald, GameFireRed, GameLeafGreen} {
		if got := GenForVersion(v); got != 3 {
			t.Errorf("GenForVersion(%v) = %d, want 3", v, got)
		}
	}
}

// Move.TypeForGen tests

func TestMoveTypeForGen_StableMove(t *testing.T) {
	tackle := &Move{
		ID:   33,
		Name: "Tackle",
		Type: TypeNormal,
	}
	for _, gen := range []Generation{1, 2, 3} {
		if got := tackle.TypeForGen(gen); got != TypeNormal {
			t.Errorf("Tackle.TypeForGen(%d) = %v, want TypeNormal", gen, got)
		}
	}
}

func TestMoveTypeForGen_ChangedMove(t *testing.T) {
	bite := &Move{
		ID:        44,
		Name:      "Bite",
		Type:      TypeDark,
		PastTypes: []MoveTypePast{{UntilGen: 1, Type: TypeNormal}},
	}
	if got := bite.TypeForGen(1); got != TypeNormal {
		t.Errorf("Bite.TypeForGen(1) = %v, want TypeNormal", got)
	}
	if got := bite.TypeForGen(2); got != TypeDark {
		t.Errorf("Bite.TypeForGen(2) = %v, want TypeDark", got)
	}
	if got := bite.TypeForGen(3); got != TypeDark {
		t.Errorf("Bite.TypeForGen(3) = %v, want TypeDark", got)
	}
}

// Move.CategoryForGen tests

func TestMoveCategoryForGen_Gen3PerMove(t *testing.T) {
	surf := &Move{
		ID:       57,
		Name:     "Surf",
		Type:     TypeWater,
		Category: CategorySpecial,
	}
	if got := surf.CategoryForGen(3); got != CategorySpecial {
		t.Errorf("Surf.CategoryForGen(3) = %v, want CategorySpecial", got)
	}
}

func TestMoveCategoryForGen_Gen1TypeBased(t *testing.T) {
	surf := &Move{
		ID:       57,
		Name:     "Surf",
		Type:     TypeWater,
		Category: CategorySpecial,
	}
	// Water is Special in Gen 1 type-split
	if got := surf.CategoryForGen(1); got != CategorySpecial {
		t.Errorf("Surf.CategoryForGen(1) = %v, want CategorySpecial", got)
	}

	tackle := &Move{
		ID:       33,
		Name:     "Tackle",
		Type:     TypeNormal,
		Category: CategoryPhysical,
	}
	// Normal is Physical in Gen 1 type-split
	if got := tackle.CategoryForGen(1); got != CategoryPhysical {
		t.Errorf("Tackle.CategoryForGen(1) = %v, want CategoryPhysical", got)
	}
}

func TestMoveCategoryForGen_StatusAlwaysStatus(t *testing.T) {
	growl := &Move{
		ID:       45,
		Name:     "Growl",
		Type:     TypeNormal,
		Category: CategoryStatus,
	}
	for _, gen := range []Generation{1, 2, 3} {
		if got := growl.CategoryForGen(gen); got != CategoryStatus {
			t.Errorf("Growl.CategoryForGen(%d) = %v, want CategoryStatus", gen, got)
		}
	}
}

func TestMoveCategoryForGen_BiteGen1(t *testing.T) {
	// Bite was Normal in Gen 1 → Normal is Physical in Gen 1 type-split
	bite := &Move{
		ID:        44,
		Name:      "Bite",
		Type:      TypeDark,
		Category:  CategoryPhysical,
		PastTypes: []MoveTypePast{{UntilGen: 1, Type: TypeNormal}},
	}
	if got := bite.CategoryForGen(1); got != CategoryPhysical {
		t.Errorf("Bite.CategoryForGen(1) = %v, want CategoryPhysical", got)
	}
}

func TestMoveCategoryForGen_BiteGen2(t *testing.T) {
	// Bite is Dark in Gen 2; Dark is Physical in Gen 2 type-split
	bite := &Move{
		ID:        44,
		Name:      "Bite",
		Type:      TypeDark,
		Category:  CategoryPhysical,
		PastTypes: []MoveTypePast{{UntilGen: 1, Type: TypeNormal}},
	}
	if got := bite.CategoryForGen(2); got != CategoryPhysical {
		t.Errorf("Bite.CategoryForGen(2) = %v, want CategoryPhysical", got)
	}
}

// Pokemon.TypesForGen tests

func TestPokemonTypesForGen_Stable(t *testing.T) {
	bulbasaur := &Pokemon{
		ID:    1,
		Name:  "bulbasaur",
		Types: [2]PokeType{TypeGrass, TypePoison},
	}
	for _, gen := range []Generation{1, 2, 3} {
		got := bulbasaur.TypesForGen(gen)
		if got[0] != TypeGrass || got[1] != TypePoison {
			t.Errorf("Bulbasaur.TypesForGen(%d) = %v, want [Grass, Poison]", gen, got)
		}
	}
}

func TestPokemonTypesForGen_Changed(t *testing.T) {
	magnemite := &Pokemon{
		ID:        81,
		Name:      "magnemite",
		Types:     [2]PokeType{TypeElectric, TypeSteel},
		PastTypes: []PokemonTypePast{{UntilGen: 1, Types: [2]PokeType{TypeElectric, TypeNone}}},
	}
	gen1 := magnemite.TypesForGen(1)
	if gen1[0] != TypeElectric || gen1[1] != TypeNone {
		t.Errorf("Magnemite.TypesForGen(1) = %v, want [Electric, None]", gen1)
	}
	gen2 := magnemite.TypesForGen(2)
	if gen2[0] != TypeElectric || gen2[1] != TypeSteel {
		t.Errorf("Magnemite.TypesForGen(2) = %v, want [Electric, Steel]", gen2)
	}
	gen3 := magnemite.TypesForGen(3)
	if gen3[0] != TypeElectric || gen3[1] != TypeSteel {
		t.Errorf("Magnemite.TypesForGen(3) = %v, want [Electric, Steel]", gen3)
	}
}
