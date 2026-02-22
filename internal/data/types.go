package data

// PokeType is a single byte; 17 types fit with room to spare.
type PokeType byte

const (
	TypeNone     PokeType = 0
	TypeNormal   PokeType = 1
	TypeFire     PokeType = 2
	TypeWater    PokeType = 3
	TypeGrass    PokeType = 4
	TypeElectric PokeType = 5
	TypeIce      PokeType = 6
	TypeFighting PokeType = 7
	TypePoison   PokeType = 8
	TypeGround   PokeType = 9
	TypeFlying   PokeType = 10
	TypePsychic  PokeType = 11
	TypeBug      PokeType = 12
	TypeRock     PokeType = 13
	TypeGhost    PokeType = 14
	TypeDragon   PokeType = 15
	TypeDark     PokeType = 16
	TypeSteel    PokeType = 17
)

var typeNames = [18]string{
	"", "Normal", "Fire", "Water", "Grass", "Electric", "Ice",
	"Fighting", "Poison", "Ground", "Flying", "Psychic", "Bug",
	"Rock", "Ghost", "Dragon", "Dark", "Steel",
}

func (t PokeType) String() string { return typeNames[t] }

// MoveCategory fits in 2 bits; using byte.
type MoveCategory byte

const (
	CategoryPhysical MoveCategory = 0
	CategorySpecial  MoveCategory = 1
	CategoryStatus   MoveCategory = 2
)

var categoryNames = [3]string{"Phys", "Spec", "Stat"}

func (c MoveCategory) String() string { return categoryNames[c] }

// LearnMethod fits in 2 bits; using byte.
type LearnMethod byte

const (
	LearnLevelUp LearnMethod = 0
	LearnMachine LearnMethod = 1
	LearnTutor   LearnMethod = 2
	LearnEgg     LearnMethod = 3
)

// GameVersion fits in 4 bits; using byte.
type GameVersion byte

const (
	GameRed       GameVersion = 1
	GameBlue      GameVersion = 2
	GameYellow    GameVersion = 3
	GameGold      GameVersion = 4
	GameSilver    GameVersion = 5
	GameCrystal   GameVersion = 6
	GameRuby      GameVersion = 7
	GameSapphire  GameVersion = 8
	GameEmerald   GameVersion = 9
	GameFireRed   GameVersion = 10
	GameLeafGreen GameVersion = 11
)

var versionNames = [12]string{
	"", "Red", "Blue", "Yellow", "Gold", "Silver", "Crystal",
	"Ruby", "Sapphire", "Emerald", "FireRed", "LeafGreen",
}

func (v GameVersion) String() string { return versionNames[v] }

// EncounterMethod fits in 3 bits; using byte.
type EncounterMethod byte

const (
	EncounterWalk      EncounterMethod = 0
	EncounterSurf      EncounterMethod = 1
	EncounterOldRod    EncounterMethod = 2
	EncounterGoodRod   EncounterMethod = 3
	EncounterSuperRod  EncounterMethod = 4
	EncounterRockSmash EncounterMethod = 5
	EncounterHeadbutt  EncounterMethod = 6
)

var encounterMethodNames = [7]string{
	"Walk", "Surf", "Old Rod", "Good Rod", "Super Rod", "Rock Smash", "Headbutt",
}

func (e EncounterMethod) String() string {
	if int(e) < len(encounterMethodNames) {
		return encounterMethodNames[e]
	}
	return "Unknown"
}

// Generation is 1, 2, or 3 — fits in 2 bits; using byte.
type Generation byte

// GenForVersion maps a GameVersion to its generation number.
func GenForVersion(v GameVersion) Generation {
	switch {
	case v <= GameYellow:
		return 1
	case v <= GameCrystal:
		return 2
	default:
		return 3
	}
}

// MoveID uniquely identifies a move. uint16 handles all current + future gens.
type MoveID uint16

// AbilityID uniquely identifies an ability.
type AbilityID uint16

// MoveTypePast records a move's type before it changed.
// UntilGen and Type are each 1 byte → struct is 2 bytes.
type MoveTypePast struct {
	UntilGen Generation
	Type     PokeType
}

// Move is stored once in AllMoves; LearnedMove references it by MoveID.
type Move struct {
	ID        MoveID
	Name      string
	Type      PokeType
	Category  MoveCategory
	Power     uint8
	Accuracy  uint8
	PP        uint8
	PastTypes []MoveTypePast
}

// TypeForGen returns the move's type for a given generation.
func (m *Move) TypeForGen(gen Generation) PokeType {
	for _, pt := range m.PastTypes {
		if gen <= pt.UntilGen {
			return pt.Type
		}
	}
	return m.Type
}

// preGen3PhysicalTypes maps PokeType to Physical for the Gen 1-2 type-based split.
// Physical: Normal, Fighting, Poison, Ground, Flying, Rock, Ghost, Bug, Dark, Steel.
// Dark and Steel were introduced in Gen 2 and are Physical in the type-split.
// Indexed by PokeType byte — O(1), zero allocation.
var preGen3PhysicalTypes = [18]bool{
	false, // TypeNone
	true,  // TypeNormal
	false, // TypeFire
	false, // TypeWater
	false, // TypeGrass
	false, // TypeElectric
	false, // TypeIce
	true,  // TypeFighting
	true,  // TypePoison
	true,  // TypeGround
	true,  // TypeFlying
	false, // TypePsychic
	true,  // TypeBug
	true,  // TypeRock
	true,  // TypeGhost
	false, // TypeDragon
	true,  // TypeDark  (Physical in Gen 2)
	true,  // TypeSteel (Physical in Gen 2)
}

// CategoryForGen returns the move's damage category for a given generation.
func (m *Move) CategoryForGen(gen Generation) MoveCategory {
	if m.Category == CategoryStatus {
		return CategoryStatus
	}
	if gen >= 3 {
		return m.Category
	}
	if preGen3PhysicalTypes[m.TypeForGen(gen)] {
		return CategoryPhysical
	}
	return CategorySpecial
}

// Ability is stored once in AllAbilities; Pokemon references it by AbilityID.
type Ability struct {
	ID        AbilityID
	Name      string
	ShortDesc string
}

// LearnedMove references the global AllMoves table by MoveID.
type LearnedMove struct {
	MoveID         MoveID
	Method         LearnMethod
	LevelLearnedAt uint8
	MachineNumber  uint8
}

// Move is a convenience helper for render time.
func (lm LearnedMove) Move() *Move { return AllMoves[lm.MoveID] }

// VersionedLearnset groups a move list by game version.
type VersionedLearnset struct {
	Version GameVersion
	Moves   []LearnedMove
}

// Location describes where a Pokemon can be encountered.
type Location struct {
	Game            GameVersion
	EncounterMethod EncounterMethod
	MinLevel        uint8
	MaxLevel        uint8
	Chance          uint8
	AreaName        string
}

// BaseStats holds the six base stats; all fit in uint8.
type BaseStats struct {
	HP, Attack, Defense, SpecialAttack, SpecialDefense, Speed uint8
}

// PokemonTypePast records a Pokemon's types before they changed.
type PokemonTypePast struct {
	UntilGen Generation
	Types    [2]PokeType
}

// Pokemon represents a single Pokémon entry.
type Pokemon struct {
	ID        uint16
	Name      string
	Types     [2]PokeType
	PastTypes []PokemonTypePast
	Stats     BaseStats
	Height    uint16
	Weight    uint16
	Abilities [2]AbilityID
	Moves     []VersionedLearnset
	Locations []Location
}

// TypesForGen returns the Pokemon's types for a given generation.
func (p *Pokemon) TypesForGen(gen Generation) [2]PokeType {
	for _, pt := range p.PastTypes {
		if gen <= pt.UntilGen {
			return pt.Types
		}
	}
	return p.Types
}
