package search

import (
	"sort"
	"strings"

	"github.com/davidlawson7/pokedex/internal/data"
)

// Score tiers
const (
	scoreExact      = 100
	scorePrefix     = 80
	scoreContains   = 60
	scoreSubseq     = 40
	scoreNoMatch    = 0
)

// filterOver is the pure, testable implementation. It accepts an injected slice.
func filterOver(pokemon []*data.Pokemon, query string) []*data.Pokemon {
	if query == "" {
		result := make([]*data.Pokemon, len(pokemon))
		copy(result, pokemon)
		return result
	}

	q := strings.ToLower(query)

	type scored struct {
		p     *data.Pokemon
		score int
		idx   int
	}

	var matches []scored
	for i, p := range pokemon {
		name := strings.ToLower(p.Name)
		s := scoreMatch(name, q)
		if s > scoreNoMatch {
			matches = append(matches, scored{p, s, i})
		}
	}

	// Stable sort: higher score first; original index (dex order) as tiebreaker.
	sort.SliceStable(matches, func(i, j int) bool {
		if matches[i].score != matches[j].score {
			return matches[i].score > matches[j].score
		}
		return matches[i].idx < matches[j].idx
	})

	result := make([]*data.Pokemon, len(matches))
	for i, m := range matches {
		result[i] = m.p
	}
	return result
}

// scoreMatch computes the match score for a single name against a query.
// Both name and query must already be lowercased.
func scoreMatch(name, query string) int {
	if name == query {
		return scoreExact
	}
	if strings.HasPrefix(name, query) {
		return scorePrefix
	}
	if strings.Contains(name, query) {
		return scoreContains
	}
	if isSubsequence(name, query) {
		return scoreSubseq
	}
	return scoreNoMatch
}

// isSubsequence returns true if query is a subsequence of name.
func isSubsequence(name, query string) bool {
	qi := 0
	for _, c := range name {
		if qi < len(query) && rune(query[qi]) == c {
			qi++
		}
	}
	return qi == len(query)
}

// Filter is the public API, wrapping filterOver with data.AllPokemon.
func Filter(query string) []*data.Pokemon {
	return filterOver(data.AllPokemon, query)
}

// FilterOver is the exported version of filterOver for use by the TUI package.
func FilterOver(pokemon []*data.Pokemon, query string) []*data.Pokemon {
	return filterOver(pokemon, query)
}
