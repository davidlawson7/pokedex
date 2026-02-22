package data

//go:generate go run ../../cmd/gencmd/main.go -data ../../_data/api-data/data/api/v2 -out .

// AllPokemon is the ordered dex slice; appended to by generated init() in pokemon_gen.go.
var AllPokemon []*Pokemon

// AllMoves is indexed by MoveID; slot 0 unused. Populated by moves_gen.go init().
var AllMoves []*Move

// AllAbilities is indexed by AbilityID; slot 0 unused. Populated by abilities_gen.go init().
var AllAbilities []*Ability

// ByID and ByName are built after all generated init() blocks have run.
var ByID map[uint16]*Pokemon
var ByName map[string]*Pokemon

func init() {
	ByID = make(map[uint16]*Pokemon, len(AllPokemon))
	ByName = make(map[string]*Pokemon, len(AllPokemon))
	for _, p := range AllPokemon {
		ByID[p.ID] = p
		ByName[p.Name] = p
	}
}
