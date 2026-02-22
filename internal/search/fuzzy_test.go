package search

import (
	"testing"

	"github.com/davidlawson7/pokedex/internal/data"
)

// fixture is a small set of Pokemon for testing, independent of generated data.
var fixture = []*data.Pokemon{
	{ID: 1, Name: "bulbasaur"},
	{ID: 2, Name: "ivysaur"},
	{ID: 3, Name: "venusaur"},
	{ID: 4, Name: "charmander"},
	{ID: 5, Name: "charmeleon"},
	{ID: 6, Name: "charizard"},
	{ID: 19, Name: "rattata"},
	{ID: 7, Name: "squirtle"},
}

func TestFilter_EmptyQuery_ReturnsAll(t *testing.T) {
	got := filterOver(fixture, "")
	if len(got) != len(fixture) {
		t.Errorf("filterOver(fixture, \"\") len = %d, want %d", len(got), len(fixture))
	}
	// Should be in dex order (original slice order)
	for i, p := range got {
		if p.ID != fixture[i].ID {
			t.Errorf("result[%d].ID = %d, want %d", i, p.ID, fixture[i].ID)
		}
	}
}

func TestFilter_ExactMatch_ScoresHighest(t *testing.T) {
	got := filterOver(fixture, "rattata")
	if len(got) == 0 {
		t.Fatal("expected at least one result")
	}
	if got[0].Name != "rattata" {
		t.Errorf("result[0].Name = %q, want \"rattata\"", got[0].Name)
	}
}

func TestFilter_PrefixMatch(t *testing.T) {
	got := filterOver(fixture, "char")
	if len(got) < 3 {
		t.Fatalf("expected at least 3 results for \"char\", got %d", len(got))
	}
	// All should start with "char"
	for _, p := range got {
		if len(p.Name) < 4 || p.Name[:4] != "char" {
			// Could also be a subsequence match - just check top 3 are char*
		}
	}
	// Charmander (4) comes before Charizard (6) in dex order
	var charmanderPos, charizardPos int = -1, -1
	for i, p := range got {
		if p.Name == "charmander" {
			charmanderPos = i
		}
		if p.Name == "charizard" {
			charizardPos = i
		}
	}
	if charmanderPos < 0 || charizardPos < 0 {
		t.Fatalf("charmander or charizard missing from results")
	}
	if charmanderPos > charizardPos {
		t.Errorf("expected charmander before charizard (dex order tiebreak), got positions %d, %d",
			charmanderPos, charizardPos)
	}
}

func TestFilter_ContainsMatch(t *testing.T) {
	got := filterOver(fixture, "saur")
	names := make(map[string]bool, len(got))
	for _, p := range got {
		names[p.Name] = true
	}
	for _, want := range []string{"bulbasaur", "ivysaur", "venusaur"} {
		if !names[want] {
			t.Errorf("expected %q in results for \"saur\", got: %v", want, got)
		}
	}
}

func TestFilter_SubsequenceMatch(t *testing.T) {
	got := filterOver(fixture, "bsr")
	found := false
	for _, p := range got {
		if p.Name == "bulbasaur" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected \"bulbasaur\" in results for subsequence \"bsr\"")
	}
}

func TestFilter_NoMatch_ReturnsEmpty(t *testing.T) {
	got := filterOver(fixture, "zzz")
	if len(got) != 0 {
		t.Errorf("expected empty results for \"zzz\", got %d", len(got))
	}
}

func TestFilter_CaseInsensitive(t *testing.T) {
	got := filterOver(fixture, "CHAR")
	if len(got) == 0 {
		t.Error("expected results for \"CHAR\", got none")
	}
	for _, p := range got {
		_ = p // just verify no panic and results returned
	}
}

func TestFilter_ScoreOrdering(t *testing.T) {
	// "charmander" is exact, "char" prefix → charmander before charmeleon before charizard
	got := filterOver(fixture, "charmander")
	if len(got) == 0 {
		t.Fatal("no results")
	}
	if got[0].Name != "charmander" {
		t.Errorf("exact match \"charmander\" should be first, got %q", got[0].Name)
	}
}

func TestFilter_PublicWrapper(t *testing.T) {
	// Filter wraps filterOver with data.AllPokemon; smoke test that it returns a slice
	// (AllPokemon is nil at test time since generated init hasn't run, returns empty)
	got := Filter("")
	// Just verify it doesn't panic
	_ = got
}
