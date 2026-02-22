package gen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

// --- JSON shape structs ---

type apiNamedResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type apiPokemon struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
	Types  []struct {
		Slot int              `json:"slot"`
		Type apiNamedResource `json:"type"`
	} `json:"types"`
	PastTypes []struct {
		Generation apiNamedResource `json:"generation"`
		Types      []struct {
			Slot int              `json:"slot"`
			Type apiNamedResource `json:"type"`
		} `json:"types"`
	} `json:"past_types"`
	Stats []struct {
		BaseStat int              `json:"base_stat"`
		Stat     apiNamedResource `json:"stat"`
	} `json:"stats"`
	Abilities []struct {
		Slot     int              `json:"slot"`
		IsHidden bool             `json:"is_hidden"`
		Ability  apiNamedResource `json:"ability"`
	} `json:"abilities"`
	Moves []struct {
		Move                apiNamedResource `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int              `json:"level_learned_at"`
			MoveLearnMethod apiNamedResource `json:"move_learn_method"`
			VersionGroup    apiNamedResource `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
}

type apiMove struct {
	ID          int              `json:"id"`
	Name        string           `json:"name"`
	Power       *int             `json:"power"`
	Accuracy    *int             `json:"accuracy"`
	PP          *int             `json:"pp"`
	DamageClass apiNamedResource `json:"damage_class"`
	Type        apiNamedResource `json:"type"`
	PastValues  []struct {
		Type         *apiNamedResource `json:"type"`
		VersionGroup apiNamedResource  `json:"version_group"`
	} `json:"past_values"`
}

type apiAbility struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	EffectEntries []struct {
		Language    apiNamedResource `json:"language"`
		ShortEffect string           `json:"short_effect"`
	} `json:"effect_entries"`
}

type apiEncounterEntry struct {
	LocationArea apiNamedResource `json:"location_area"`
	VersionDetails []struct {
		EncounterDetails []struct {
			Chance   int              `json:"chance"`
			MaxLevel int              `json:"max_level"`
			MinLevel int              `json:"min_level"`
			Method   apiNamedResource `json:"method"`
		} `json:"encounter_details"`
		MaxChance int              `json:"max_chance"`
		Version   apiNamedResource `json:"version"`
	} `json:"version_details"`
}

// --- Enum parsing ---

// ParseGeneration converts "generation-i" → 1, "generation-ii" → 2, "generation-iii" → 3.
// Returns an error for any generation outside Gen 1-3.
func ParseGeneration(name string) (byte, error) {
	switch name {
	case "generation-i":
		return 1, nil
	case "generation-ii":
		return 2, nil
	case "generation-iii":
		return 3, nil
	}
	return 0, fmt.Errorf("unknown generation: %q", name)
}

// parseGenerationFull converts any generation name to its number (1-9+).
// Used for past_types entries that may span gen 4-9 (e.g. Fairy changes in Gen 6).
func parseGenerationFull(name string) (byte, bool) {
	switch name {
	case "generation-i":
		return 1, true
	case "generation-ii":
		return 2, true
	case "generation-iii":
		return 3, true
	case "generation-iv":
		return 4, true
	case "generation-v":
		return 5, true
	case "generation-vi":
		return 6, true
	case "generation-vii":
		return 7, true
	case "generation-viii":
		return 8, true
	case "generation-ix":
		return 9, true
	}
	return 0, false
}

// parsePokeTypeLenient converts a type name to its byte constant.
// Returns (TypeNone=0, nil) for unknown types (post-Gen 3 types like Fairy).
func parsePokeTypeLenient(name string) byte {
	v, _ := ParsePokeType(name)
	return v // returns 0 (TypeNone) on error, which is intentional
}

// ParseDamageClass converts "physical"→0, "special"→1, "status"→2.
func ParseDamageClass(name string) (byte, error) {
	switch name {
	case "physical":
		return 0, nil
	case "special":
		return 1, nil
	case "status":
		return 2, nil
	}
	return 0, fmt.Errorf("unknown damage class: %q", name)
}

// ParsePokeType converts a type name string to its byte constant value.
func ParsePokeType(name string) (byte, error) {
	switch name {
	case "normal":
		return 1, nil
	case "fire":
		return 2, nil
	case "water":
		return 3, nil
	case "grass":
		return 4, nil
	case "electric":
		return 5, nil
	case "ice":
		return 6, nil
	case "fighting":
		return 7, nil
	case "poison":
		return 8, nil
	case "ground":
		return 9, nil
	case "flying":
		return 10, nil
	case "psychic":
		return 11, nil
	case "bug":
		return 12, nil
	case "rock":
		return 13, nil
	case "ghost":
		return 14, nil
	case "dragon":
		return 15, nil
	case "dark":
		return 16, nil
	case "steel":
		return 17, nil
	}
	return 0, fmt.Errorf("unknown type: %q", name)
}

// ParseVersionGroup maps a version group name to the list of GameVersion constants.
// Returns nil for version groups outside Gen 1-3.
func ParseVersionGroup(name string) ([]string, error) {
	switch name {
	case "red-blue":
		return []string{"GameRed", "GameBlue"}, nil
	case "yellow":
		return []string{"GameYellow"}, nil
	case "gold-silver":
		return []string{"GameGold", "GameSilver"}, nil
	case "crystal":
		return []string{"GameCrystal"}, nil
	case "ruby-sapphire":
		return []string{"GameRuby", "GameSapphire"}, nil
	case "emerald":
		return []string{"GameEmerald"}, nil
	case "firered-leafgreen":
		return []string{"GameFireRed", "GameLeafGreen"}, nil
	}
	return nil, nil // outside Gen 1-3, silently ignored
}

// prevGenForVersionGroup returns the generation just before this version group.
// Used for move past_values: the version_group is where the change happened,
// so the old value applied until the end of the previous generation.
// e.g. Bite: {type: normal, version_group: gold-silver} → change happened in Gen 2
//      → old type (Normal) was valid until end of Gen 1 → UntilGen = 1.
func prevGenForVersionGroup(name string) byte {
	switch name {
	case "red-blue", "yellow":
		return 0 // change was in Gen 1 (or before) — old type never applies in Gen 1-3
	case "gold-silver", "crystal":
		return 1 // change was in Gen 2 → old type valid until end of Gen 1
	case "ruby-sapphire", "emerald", "firered-leafgreen":
		return 2 // change was in Gen 3 → old type valid until end of Gen 2
	default:
		return 3 // change was in Gen 4+ → old type valid through all of Gen 1-3
	}
}

// versionNameToGameVersion converts an individual version name (as seen in encounters)
// to its GameVersion constant name.
func versionNameToGameVersion(name string) string {
	switch name {
	case "red":
		return "GameRed"
	case "blue":
		return "GameBlue"
	case "yellow":
		return "GameYellow"
	case "gold":
		return "GameGold"
	case "silver":
		return "GameSilver"
	case "crystal":
		return "GameCrystal"
	case "ruby":
		return "GameRuby"
	case "sapphire":
		return "GameSapphire"
	case "emerald":
		return "GameEmerald"
	case "firered":
		return "GameFireRed"
	case "leafgreen":
		return "GameLeafGreen"
	}
	return ""
}

// encounterMethodConstant converts an encounter method name to its constant name.
func encounterMethodConstant(name string) string {
	switch name {
	case "walk", "grass", "tall-grass":
		return "EncounterWalk"
	case "surf", "water":
		return "EncounterSurf"
	case "old-rod":
		return "EncounterOldRod"
	case "good-rod":
		return "EncounterGoodRod"
	case "super-rod":
		return "EncounterSuperRod"
	case "rock-smash":
		return "EncounterRockSmash"
	case "headbutt", "headbutt-normal", "headbutt-special":
		return "EncounterHeadbutt"
	}
	return "EncounterWalk"
}

// idFromURL extracts the numeric ID from a PokeAPI URL like ".../ability/65/".
func idFromURL(url string) (int, error) {
	url = strings.TrimRight(url, "/")
	idx := strings.LastIndex(url, "/")
	if idx < 0 {
		return 0, fmt.Errorf("invalid URL: %q", url)
	}
	return strconv.Atoi(url[idx+1:])
}

// --- Data structures for codegen ---

// MoveData is the parsed representation of a move.
type MoveData struct {
	ID       int
	Name     string
	Type     byte
	Category byte
	Power    uint8
	Accuracy uint8
	PP       uint8
	// PastTypes: each entry is {UntilGen, TypeConst}
	PastTypes []MovePastTypeData
}

// MovePastTypeData records a move's past type for codegen.
type MovePastTypeData struct {
	UntilGen byte
	TypeConst string // e.g. "TypeNormal"
}

// AbilityData is the parsed representation of an ability.
type AbilityData struct {
	ID        int
	Name      string
	ShortDesc string
}

// VersionedMoveEntry is a move learned by a specific method in a version.
type VersionedMoveEntry struct {
	GameVersion    string // e.g. "GameRed"
	MoveID         int
	Method         string // "LearnLevelUp", "LearnMachine", "LearnTutor", "LearnEgg"
	LevelLearnedAt uint8
	MachineNumber  uint8
}

// PastTypeData records a pokemon's past types for a generation.
type PastTypeData struct {
	UntilGen byte
	Type1    string // e.g. "TypeElectric"
	Type2    string // e.g. "TypeNone"
}

// LocationData records a single encounter for codegen.
type LocationData struct {
	GameVersion     string
	EncounterMethod string
	MinLevel        uint8
	MaxLevel        uint8
	Chance          uint8
	AreaName        string
}

// PokemonData is the parsed representation of a pokemon.
type PokemonData struct {
	ID        int
	Name      string
	Type1     string
	Type2     string
	PastTypes []PastTypeData
	HP        uint8
	Attack    uint8
	Defense   uint8
	SpAtk     uint8
	SpDef     uint8
	Speed     uint8
	Height    uint16
	Weight    uint16
	Ability1  int // 0 = none
	Ability2  int // 0 = none
	// VersionedMoves: grouped by game version constant
	VersionedMoves map[string][]VersionedMoveEntry
	Locations      []LocationData
}

// typeConstant converts a byte type value to its Go constant name.
func typeConstant(t byte) string {
	names := [18]string{
		"TypeNone", "TypeNormal", "TypeFire", "TypeWater", "TypeGrass",
		"TypeElectric", "TypeIce", "TypeFighting", "TypePoison", "TypeGround",
		"TypeFlying", "TypePsychic", "TypeBug", "TypeRock", "TypeGhost",
		"TypeDragon", "TypeDark", "TypeSteel",
	}
	if int(t) < len(names) {
		return names[t]
	}
	return "TypeNone"
}

// --- File reading helpers ---

func readJSON(path string, v interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}

// --- Build functions ---

// BuildMove parses a move JSON file path and returns a MoveData.
func BuildMove(path string) (MoveData, error) {
	var m apiMove
	if err := readJSON(path, &m); err != nil {
		return MoveData{}, err
	}

	typeByte, err := ParsePokeType(m.Type.Name)
	if err != nil {
		return MoveData{}, err
	}
	catByte, err := ParseDamageClass(m.DamageClass.Name)
	if err != nil {
		return MoveData{}, err
	}

	var power, accuracy, pp uint8
	if m.Power != nil {
		power = uint8(*m.Power)
	}
	if m.Accuracy != nil {
		accuracy = uint8(*m.Accuracy)
	}
	if m.PP != nil {
		pp = uint8(*m.PP)
	}

	// Build past types: the version_group in past_values marks when the change HAPPENED.
	// UntilGen = prevGen(version_group) = the last generation where the old type applied.
	var pastTypes []MovePastTypeData
	for _, pv := range m.PastValues {
		if pv.Type == nil {
			continue
		}
		untilGen := prevGenForVersionGroup(pv.VersionGroup.Name)
		if untilGen == 0 {
			continue // change was in Gen 1 or before — irrelevant for our Gen 1-3 scope
		}
		pastTypeByte, err := ParsePokeType(pv.Type.Name)
		if err != nil {
			continue
		}
		pastTypes = append(pastTypes, MovePastTypeData{
			UntilGen:  untilGen,
			TypeConst: typeConstant(pastTypeByte),
		})
	}

	return MoveData{
		ID:        m.ID,
		Name:      strings.Title(m.Name), // display name: capitalize first letter
		Type:      typeByte,
		Category:  catByte,
		Power:     power,
		Accuracy:  accuracy,
		PP:        pp,
		PastTypes: pastTypes,
	}, nil
}

// BuildAbility parses an ability JSON file path and returns an AbilityData.
func BuildAbility(path string) (AbilityData, error) {
	var a apiAbility
	if err := readJSON(path, &a); err != nil {
		return AbilityData{}, err
	}
	var shortDesc string
	for _, e := range a.EffectEntries {
		if e.Language.Name == "en" {
			shortDesc = e.ShortEffect
			break
		}
	}
	return AbilityData{
		ID:        a.ID,
		Name:      a.Name,
		ShortDesc: shortDesc,
	}, nil
}

// CollectMoves reads all move files referenced by Pokemon in dataDir and returns a map[id]MoveData.
func CollectMoves(dataDir string, pokemonIDs []int) (map[int]MoveData, error) {
	// First pass: collect move IDs from all pokemon files.
	moveIDs := make(map[int]bool)
	for _, id := range pokemonIDs {
		path := filepath.Join(dataDir, "pokemon", strconv.Itoa(id), "index.json")
		var p apiPokemon
		if err := readJSON(path, &p); err != nil {
			return nil, fmt.Errorf("reading pokemon %d: %w", id, err)
		}
		for _, m := range p.Moves {
			mid, err := idFromURL(m.Move.URL)
			if err != nil {
				return nil, err
			}
			moveIDs[mid] = true
		}
	}

	moves := make(map[int]MoveData, len(moveIDs))
	for id := range moveIDs {
		path := filepath.Join(dataDir, "move", strconv.Itoa(id), "index.json")
		m, err := BuildMove(path)
		if err != nil {
			// Skip moves with unknown types (post-Gen 3 types like Fairy, Shadow, etc.)
			continue
		}
		moves[id] = m
	}
	return moves, nil
}

// CollectAbilities reads all ability files referenced by Pokemon in dataDir and returns a map[id]AbilityData.
func CollectAbilities(dataDir string, pokemonIDs []int) (map[int]AbilityData, error) {
	abilityIDs := make(map[int]bool)
	for _, id := range pokemonIDs {
		path := filepath.Join(dataDir, "pokemon", strconv.Itoa(id), "index.json")
		var p apiPokemon
		if err := readJSON(path, &p); err != nil {
			return nil, fmt.Errorf("reading pokemon %d: %w", id, err)
		}
		for _, a := range p.Abilities {
			if a.IsHidden {
				continue
			}
			aid, err := idFromURL(a.Ability.URL)
			if err != nil {
				return nil, err
			}
			abilityIDs[aid] = true
		}
	}

	abilities := make(map[int]AbilityData, len(abilityIDs))
	for id := range abilityIDs {
		path := filepath.Join(dataDir, "ability", strconv.Itoa(id), "index.json")
		a, err := BuildAbility(path)
		if err != nil {
			return nil, fmt.Errorf("building ability %d: %w", id, err)
		}
		abilities[id] = a
	}
	return abilities, nil
}

// learnMethodConstant maps a move_learn_method name to its Go constant.
func learnMethodConstant(name string) string {
	switch name {
	case "level-up":
		return "LearnLevelUp"
	case "machine":
		return "LearnMachine"
	case "tutor":
		return "LearnTutor"
	case "egg":
		return "LearnEgg"
	}
	return "LearnLevelUp"
}

// BuildPokemon parses a single pokemon JSON and returns a PokemonData.
func BuildPokemon(dataDir string, id int, abilities map[int]AbilityData) (PokemonData, error) {
	path := filepath.Join(dataDir, "pokemon", strconv.Itoa(id), "index.json")
	var p apiPokemon
	if err := readJSON(path, &p); err != nil {
		return PokemonData{}, err
	}

	// Types (canonical — use lenient parser; post-Gen 3 types like Fairy become TypeNone)
	var type1, type2 string = "TypeNone", "TypeNone"
	for _, t := range p.Types {
		tb := parsePokeTypeLenient(t.Type.Name)
		if t.Slot == 1 {
			type1 = typeConstant(tb)
		} else if t.Slot == 2 {
			type2 = typeConstant(tb)
		}
	}

	// Past types — use full generation parser to capture Gen 4-5 entries.
	// These are needed for Pokemon whose types changed in Gen 6 (e.g. Fairy retrochange):
	// their Gen 5 past_type entry records what they were in Gen 1-3.
	var pastTypes []PastTypeData
	for _, pt := range p.PastTypes {
		gen, ok := parseGenerationFull(pt.Generation.Name)
		if !ok {
			continue
		}
		var pt1, pt2 string = "TypeNone", "TypeNone"
		for _, t := range pt.Types {
			tb := parsePokeTypeLenient(t.Type.Name)
			if t.Slot == 1 {
				pt1 = typeConstant(tb)
			} else if t.Slot == 2 {
				pt2 = typeConstant(tb)
			}
		}
		pastTypes = append(pastTypes, PastTypeData{UntilGen: gen, Type1: pt1, Type2: pt2})
	}

	// Stats
	statsMap := make(map[string]uint8)
	for _, s := range p.Stats {
		statsMap[s.Stat.Name] = uint8(s.BaseStat)
	}

	// Abilities (non-hidden only, slots 1 and 2)
	var ab1, ab2 int
	for _, a := range p.Abilities {
		if a.IsHidden {
			continue
		}
		aid, err := idFromURL(a.Ability.URL)
		if err != nil {
			continue
		}
		if a.Slot == 1 {
			ab1 = aid
		} else if a.Slot == 2 {
			ab2 = aid
		}
	}

	// Moves: group by version
	versionedMoves := make(map[string][]VersionedMoveEntry)
	for _, m := range p.Moves {
		mid, err := idFromURL(m.Move.URL)
		if err != nil {
			continue
		}
		for _, vgd := range m.VersionGroupDetails {
			versions, _ := ParseVersionGroup(vgd.VersionGroup.Name)
			if versions == nil {
				continue
			}
			method := learnMethodConstant(vgd.MoveLearnMethod.Name)
			for _, ver := range versions {
				entry := VersionedMoveEntry{
					GameVersion:    ver,
					MoveID:         mid,
					Method:         method,
					LevelLearnedAt: uint8(vgd.LevelLearnedAt),
				}
				versionedMoves[ver] = append(versionedMoves[ver], entry)
			}
		}
	}
	// Sort each version's moves by level then move ID for deterministic output
	for ver := range versionedMoves {
		moves := versionedMoves[ver]
		sort.Slice(moves, func(i, j int) bool {
			if moves[i].LevelLearnedAt != moves[j].LevelLearnedAt {
				return moves[i].LevelLearnedAt < moves[j].LevelLearnedAt
			}
			return moves[i].MoveID < moves[j].MoveID
		})
		versionedMoves[ver] = moves
	}

	// Locations
	encPath := filepath.Join(dataDir, "pokemon", strconv.Itoa(id), "encounters", "index.json")
	var encounters []apiEncounterEntry
	if err := readJSON(encPath, &encounters); err != nil {
		// encounters file might not exist; treat as no encounters
		encounters = nil
	}

	var locations []LocationData
	for _, enc := range encounters {
		areaName := enc.LocationArea.Name
		for _, vd := range enc.VersionDetails {
			gameVer := versionNameToGameVersion(vd.Version.Name)
			if gameVer == "" {
				continue // outside Gen 1-3
			}
			for _, ed := range vd.EncounterDetails {
				locations = append(locations, LocationData{
					GameVersion:     gameVer,
					EncounterMethod: encounterMethodConstant(ed.Method.Name),
					MinLevel:        uint8(ed.MinLevel),
					MaxLevel:        uint8(ed.MaxLevel),
					Chance:          uint8(ed.Chance),
					AreaName:        areaName,
				})
			}
		}
	}

	return PokemonData{
		ID:             p.ID,
		Name:           p.Name,
		Type1:          type1,
		Type2:          type2,
		PastTypes:      pastTypes,
		HP:             statsMap["hp"],
		Attack:         statsMap["attack"],
		Defense:        statsMap["defense"],
		SpAtk:          statsMap["special-attack"],
		SpDef:          statsMap["special-defense"],
		Speed:          statsMap["speed"],
		Height:         uint16(p.Height),
		Weight:         uint16(p.Weight),
		Ability1:       ab1,
		Ability2:       ab2,
		VersionedMoves: versionedMoves,
		Locations:      locations,
	}, nil
}

// --- Code generation templates ---

var abilitiesTemplate = template.Must(template.New("abilities").Parse(`// Code generated by cmd/gen/main.go. DO NOT EDIT.
package data

func init() {
	AllAbilities = make([]*Ability, {{.Size}})
{{- range .Abilities}}
	AllAbilities[{{.ID}}] = &Ability{ID: {{.ID}}, Name: {{printf "%q" .Name}}, ShortDesc: {{printf "%q" .ShortDesc}}}
{{- end}}
}
`))

var movesTemplate = template.Must(template.New("moves").Parse(`// Code generated by cmd/gen/main.go. DO NOT EDIT.
package data

func init() {
	AllMoves = make([]*Move, {{.Size}})
{{- range .Moves}}
	AllMoves[{{.ID}}] = &Move{ID: {{.ID}}, Name: {{printf "%q" .Name}}, Type: {{.TypeConst}}, Category: {{.CategoryConst}}, Power: {{.Power}}, Accuracy: {{.Accuracy}}, PP: {{.PP}}{{if .PastTypes}}, PastTypes: []MoveTypePast{ {{- range .PastTypes}}{UntilGen: {{.UntilGen}}, Type: {{.TypeConst}}}, {{end}}}{{end}}}
{{- end}}
}
`))

var pokemonTemplate = template.Must(template.New("pokemon").Funcs(template.FuncMap{
	"versionOrder": versionOrder,
	"sortedVersions": func(m map[string][]VersionedMoveEntry) []string {
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			return versionOrder(keys[i]) < versionOrder(keys[j])
		})
		return keys
	},
}).Parse(`// Code generated by cmd/gen/main.go. DO NOT EDIT.
package data

func init() {
	AllPokemon = append(AllPokemon,
{{- range .Pokemon}}
		&Pokemon{
			ID:    {{.ID}},
			Name:  {{printf "%q" .Name}},
			Types: [2]PokeType{ {{.Type1}}, {{.Type2}} },
			{{- if .PastTypes}}
			PastTypes: []PokemonTypePast{
				{{- range .PastTypes}}
				{UntilGen: {{.UntilGen}}, Types: [2]PokeType{ {{.Type1}}, {{.Type2}} }},
				{{- end}}
			},
			{{- end}}
			Stats:     BaseStats{HP: {{.HP}}, Attack: {{.Attack}}, Defense: {{.Defense}}, SpecialAttack: {{.SpAtk}}, SpecialDefense: {{.SpDef}}, Speed: {{.Speed}}},
			Height:    {{.Height}},
			Weight:    {{.Weight}},
			Abilities: [2]AbilityID{ {{.Ability1}}, {{.Ability2}} },
			{{- if .VersionedMoves}}
			Moves: []VersionedLearnset{
				{{- range (sortedVersions .VersionedMoves)}}
				{Version: {{.}}, Moves: []LearnedMove{
					{{- range (index $.VersionedMovesByPokemon $.PokemonIdx .)}}{MoveID: {{.MoveID}}, Method: {{.Method}}, LevelLearnedAt: {{.LevelLearnedAt}}, MachineNumber: {{.MachineNumber}}},
					{{- end}}
				}},
				{{- end}}
			},
			{{- end}}
			{{- if .Locations}}
			Locations: []Location{
				{{- range .Locations}}
				{Game: {{.GameVersion}}, EncounterMethod: {{.EncounterMethod}}, MinLevel: {{.MinLevel}}, MaxLevel: {{.MaxLevel}}, Chance: {{.Chance}}, AreaName: {{printf "%q" .AreaName}}},
				{{- end}}
			},
			{{- end}}
		},
{{- end}}
	)
}
`))

// versionOrder returns a sort key for GameVersion constants.
func versionOrder(name string) int {
	order := map[string]int{
		"GameRed": 1, "GameBlue": 2, "GameYellow": 3,
		"GameGold": 4, "GameSilver": 5, "GameCrystal": 6,
		"GameRuby": 7, "GameSapphire": 8, "GameEmerald": 9,
		"GameFireRed": 10, "GameLeafGreen": 11,
	}
	if v, ok := order[name]; ok {
		return v
	}
	return 99
}

// --- Template data helpers ---

type moveTplEntry struct {
	ID            int
	Name          string
	TypeConst     string
	CategoryConst string
	Power         uint8
	Accuracy      uint8
	PP            uint8
	PastTypes     []struct {
		UntilGen  byte
		TypeConst string
	}
}

func categoryConstant(c byte) string {
	switch c {
	case 0:
		return "CategoryPhysical"
	case 1:
		return "CategorySpecial"
	case 2:
		return "CategoryStatus"
	}
	return "CategoryPhysical"
}

// --- Run: full codegen pipeline ---

// Config holds codegen parameters.
type Config struct {
	DataDir    string // path to api/v2/
	OutDir     string // output directory for generated Go files
	PokemonIDs []int  // IDs to process; nil means scan 1-386
}

// Run executes the full codegen pipeline.
func Run(cfg Config) error {
	ids := cfg.PokemonIDs
	if ids == nil {
		ids = make([]int, 386)
		for i := range ids {
			ids[i] = i + 1
		}
	}

	// Filter to IDs that actually have files
	var validIDs []int
	for _, id := range ids {
		p := filepath.Join(cfg.DataDir, "pokemon", strconv.Itoa(id), "index.json")
		if _, err := os.Stat(p); err == nil {
			validIDs = append(validIDs, id)
		}
	}
	ids = validIDs

	// Collect moves
	moves, err := CollectMoves(cfg.DataDir, ids)
	if err != nil {
		return fmt.Errorf("collecting moves: %w", err)
	}

	// Collect abilities
	abilities, err := CollectAbilities(cfg.DataDir, ids)
	if err != nil {
		return fmt.Errorf("collecting abilities: %w", err)
	}

	// Collect pokemon
	var allPokemon []PokemonData
	for _, id := range ids {
		pk, err := BuildPokemon(cfg.DataDir, id, abilities)
		if err != nil {
			return fmt.Errorf("building pokemon %d: %w", id, err)
		}
		allPokemon = append(allPokemon, pk)
	}

	if err := os.MkdirAll(cfg.OutDir, 0755); err != nil {
		return err
	}

	// Emit abilities_gen.go
	if err := emitAbilities(cfg.OutDir, abilities); err != nil {
		return err
	}

	// Emit moves_gen.go
	if err := emitMoves(cfg.OutDir, moves); err != nil {
		return err
	}

	// Emit pokemon_gen.go
	if err := emitPokemon(cfg.OutDir, allPokemon); err != nil {
		return err
	}

	return nil
}

func emitAbilities(outDir string, abilities map[int]AbilityData) error {
	// Sort for deterministic output
	ids := make([]int, 0, len(abilities))
	for id := range abilities {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	maxID := 0
	for _, id := range ids {
		if id > maxID {
			maxID = id
		}
	}

	type tplData struct {
		Size      int
		Abilities []AbilityData
	}
	sorted := make([]AbilityData, 0, len(ids))
	for _, id := range ids {
		sorted = append(sorted, abilities[id])
	}

	f, err := os.Create(filepath.Join(outDir, "abilities_gen.go"))
	if err != nil {
		return err
	}
	defer f.Close()
	return abilitiesTemplate.Execute(f, tplData{Size: maxID + 10, Abilities: sorted})
}

func emitMoves(outDir string, moves map[int]MoveData) error {
	ids := make([]int, 0, len(moves))
	for id := range moves {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	maxID := 0
	for _, id := range ids {
		if id > maxID {
			maxID = id
		}
	}

	type moveTpl struct {
		ID            int
		Name          string
		TypeConst     string
		CategoryConst string
		Power         uint8
		Accuracy      uint8
		PP            uint8
		PastTypes     []struct {
			UntilGen  byte
			TypeConst string
		}
	}

	entries := make([]moveTpl, 0, len(ids))
	for _, id := range ids {
		m := moves[id]
		var pastTypes []struct {
			UntilGen  byte
			TypeConst string
		}
		for _, pt := range m.PastTypes {
			pastTypes = append(pastTypes, struct {
				UntilGen  byte
				TypeConst string
			}{pt.UntilGen, pt.TypeConst})
		}
		entries = append(entries, moveTpl{
			ID:            m.ID,
			Name:          m.Name,
			TypeConst:     typeConstant(m.Type),
			CategoryConst: categoryConstant(m.Category),
			Power:         m.Power,
			Accuracy:      m.Accuracy,
			PP:            m.PP,
			PastTypes:     pastTypes,
		})
	}

	type tplData struct {
		Size  int
		Moves []moveTpl
	}

	f, err := os.Create(filepath.Join(outDir, "moves_gen.go"))
	if err != nil {
		return err
	}
	defer f.Close()
	return movesTemplate.Execute(f, tplData{Size: maxID + 10, Moves: entries})
}

// pokemonTplEntry is a flat struct for template use (avoids map access in template).
type pokemonTplEntry struct {
	ID             int
	Name           string
	Type1          string
	Type2          string
	PastTypes      []PastTypeData
	HP             uint8
	Attack         uint8
	Defense        uint8
	SpAtk          uint8
	SpDef          uint8
	Speed          uint8
	Height         uint16
	Weight         uint16
	Ability1       int
	Ability2       int
	VersionedMoves map[string][]VersionedMoveEntry
	Locations      []LocationData
	PokemonIdx     int
}

func emitPokemon(outDir string, pokemon []PokemonData) error {
	// Build flat template entries
	entries := make([]pokemonTplEntry, len(pokemon))
	for i, p := range pokemon {
		entries[i] = pokemonTplEntry{
			ID:             p.ID,
			Name:           p.Name,
			Type1:          p.Type1,
			Type2:          p.Type2,
			PastTypes:      p.PastTypes,
			HP:             p.HP,
			Attack:         p.Attack,
			Defense:        p.Defense,
			SpAtk:          p.SpAtk,
			SpDef:          p.SpDef,
			Speed:          p.Speed,
			Height:         p.Height,
			Weight:         p.Weight,
			Ability1:       p.Ability1,
			Ability2:       p.Ability2,
			VersionedMoves: p.VersionedMoves,
			Locations:      p.Locations,
			PokemonIdx:     i,
		}
	}

	// Build VersionedMovesByPokemon: a 2D slice [pokemonIdx][versionKey] → []VersionedMoveEntry
	// Used by the template to iterate moves without map access issues.
	// We'll use a simpler approach: inline the versioned moves in the template.

	// Actually, the template approach needs restructuring since Go templates can't easily
	// do map[string][]T lookups by string key from a range variable.
	// Instead, emit pokemon_gen.go using direct Go code generation (fmt.Fprintf).
	return emitPokemonDirect(outDir, entries)
}

func emitPokemonDirect(outDir string, entries []pokemonTplEntry) error {
	f, err := os.Create(filepath.Join(outDir, "pokemon_gen.go"))
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "// Code generated by cmd/gen/main.go. DO NOT EDIT.\npackage data\n\nfunc init() {\n\tAllPokemon = append(AllPokemon,\n")

	for _, p := range entries {
		fmt.Fprintf(f, "\t\t&Pokemon{\n")
		fmt.Fprintf(f, "\t\t\tID:    %d,\n", p.ID)
		fmt.Fprintf(f, "\t\t\tName:  %q,\n", p.Name)
		fmt.Fprintf(f, "\t\t\tTypes: [2]PokeType{%s, %s},\n", p.Type1, p.Type2)

		if len(p.PastTypes) > 0 {
			fmt.Fprintf(f, "\t\t\tPastTypes: []PokemonTypePast{\n")
			for _, pt := range p.PastTypes {
				fmt.Fprintf(f, "\t\t\t\t{UntilGen: %d, Types: [2]PokeType{%s, %s}},\n", pt.UntilGen, pt.Type1, pt.Type2)
			}
			fmt.Fprintf(f, "\t\t\t},\n")
		}

		fmt.Fprintf(f, "\t\t\tStats:     BaseStats{HP: %d, Attack: %d, Defense: %d, SpecialAttack: %d, SpecialDefense: %d, Speed: %d},\n",
			p.HP, p.Attack, p.Defense, p.SpAtk, p.SpDef, p.Speed)
		fmt.Fprintf(f, "\t\t\tHeight:    %d,\n", p.Height)
		fmt.Fprintf(f, "\t\t\tWeight:    %d,\n", p.Weight)
		fmt.Fprintf(f, "\t\t\tAbilities: [2]AbilityID{%d, %d},\n", p.Ability1, p.Ability2)

		if len(p.VersionedMoves) > 0 {
			// Sort versions for deterministic output
			versions := make([]string, 0, len(p.VersionedMoves))
			for v := range p.VersionedMoves {
				versions = append(versions, v)
			}
			sort.Slice(versions, func(i, j int) bool {
				return versionOrder(versions[i]) < versionOrder(versions[j])
			})
			fmt.Fprintf(f, "\t\t\tMoves: []VersionedLearnset{\n")
			for _, ver := range versions {
				fmt.Fprintf(f, "\t\t\t\t{Version: %s, Moves: []LearnedMove{\n", ver)
				for _, m := range p.VersionedMoves[ver] {
					fmt.Fprintf(f, "\t\t\t\t\t{MoveID: %d, Method: %s, LevelLearnedAt: %d, MachineNumber: %d},\n",
						m.MoveID, m.Method, m.LevelLearnedAt, m.MachineNumber)
				}
				fmt.Fprintf(f, "\t\t\t\t}},\n")
			}
			fmt.Fprintf(f, "\t\t\t},\n")
		}

		if len(p.Locations) > 0 {
			fmt.Fprintf(f, "\t\t\tLocations: []Location{\n")
			for _, loc := range p.Locations {
				fmt.Fprintf(f, "\t\t\t\t{Game: %s, EncounterMethod: %s, MinLevel: %d, MaxLevel: %d, Chance: %d, AreaName: %q},\n",
					loc.GameVersion, loc.EncounterMethod, loc.MinLevel, loc.MaxLevel, loc.Chance, loc.AreaName)
			}
			fmt.Fprintf(f, "\t\t\t},\n")
		}

		fmt.Fprintf(f, "\t\t},\n")
	}

	fmt.Fprintf(f, "\t)\n}\n")
	return nil
}
