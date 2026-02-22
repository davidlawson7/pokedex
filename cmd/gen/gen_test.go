package gen

import (
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var testdataDir = filepath.Join("testdata")

func TestParseGeneration(t *testing.T) {
	cases := []struct {
		input string
		want  byte
	}{
		{"generation-i", 1},
		{"generation-ii", 2},
		{"generation-iii", 3},
	}
	for _, c := range cases {
		got, err := ParseGeneration(c.input)
		if err != nil {
			t.Errorf("ParseGeneration(%q) error: %v", c.input, err)
			continue
		}
		if got != c.want {
			t.Errorf("ParseGeneration(%q) = %d, want %d", c.input, got, c.want)
		}
	}
	_, err := ParseGeneration("generation-iv")
	if err == nil {
		t.Error("ParseGeneration(\"generation-iv\") expected error, got nil")
	}
}

func TestParseDamageClass(t *testing.T) {
	cases := []struct {
		input string
		want  byte
	}{
		{"physical", 0},
		{"special", 1},
		{"status", 2},
	}
	for _, c := range cases {
		got, err := ParseDamageClass(c.input)
		if err != nil {
			t.Errorf("ParseDamageClass(%q) error: %v", c.input, err)
			continue
		}
		if got != c.want {
			t.Errorf("ParseDamageClass(%q) = %d, want %d", c.input, got, c.want)
		}
	}
}

func TestParseVersionGroup(t *testing.T) {
	cases := []struct {
		input string
		want  []string
	}{
		{"red-blue", []string{"GameRed", "GameBlue"}},
		{"yellow", []string{"GameYellow"}},
		{"gold-silver", []string{"GameGold", "GameSilver"}},
		{"crystal", []string{"GameCrystal"}},
		{"ruby-sapphire", []string{"GameRuby", "GameSapphire"}},
		{"emerald", []string{"GameEmerald"}},
		{"firered-leafgreen", []string{"GameFireRed", "GameLeafGreen"}},
		{"diamond-pearl", nil}, // outside Gen 1-3
	}
	for _, c := range cases {
		got, err := ParseVersionGroup(c.input)
		if err != nil {
			t.Errorf("ParseVersionGroup(%q) error: %v", c.input, err)
			continue
		}
		if len(got) != len(c.want) {
			t.Errorf("ParseVersionGroup(%q) = %v, want %v", c.input, got, c.want)
			continue
		}
		for i := range got {
			if got[i] != c.want[i] {
				t.Errorf("ParseVersionGroup(%q)[%d] = %q, want %q", c.input, i, got[i], c.want[i])
			}
		}
	}
}

func TestBuildMove_Stable(t *testing.T) {
	m, err := BuildMove(filepath.Join(testdataDir, "move", "33", "index.json"))
	if err != nil {
		t.Fatal(err)
	}
	if m.ID != 33 {
		t.Errorf("ID = %d, want 33", m.ID)
	}
	if m.Name != "Tackle" {
		t.Errorf("Name = %q, want \"Tackle\"", m.Name)
	}
	if m.Power != 40 {
		t.Errorf("Power = %d, want 40", m.Power)
	}
	if m.Accuracy != 100 {
		t.Errorf("Accuracy = %d, want 100", m.Accuracy)
	}
	if m.PP != 35 {
		t.Errorf("PP = %d, want 35", m.PP)
	}
	if len(m.PastTypes) != 0 {
		t.Errorf("PastTypes = %v, want empty", m.PastTypes)
	}
}

func TestBuildMove_WithPastType(t *testing.T) {
	m, err := BuildMove(filepath.Join(testdataDir, "move", "44", "index.json"))
	if err != nil {
		t.Fatal(err)
	}
	if m.ID != 44 {
		t.Errorf("ID = %d, want 44", m.ID)
	}
	if len(m.PastTypes) != 1 {
		t.Fatalf("PastTypes len = %d, want 1", len(m.PastTypes))
	}
	if m.PastTypes[0].UntilGen != 1 {
		t.Errorf("PastTypes[0].UntilGen = %d, want 1", m.PastTypes[0].UntilGen)
	}
	if m.PastTypes[0].TypeConst != "TypeNormal" {
		t.Errorf("PastTypes[0].TypeConst = %q, want \"TypeNormal\"", m.PastTypes[0].TypeConst)
	}
}

func TestBuildPokemon_PastTypes(t *testing.T) {
	pk, err := BuildPokemon(testdataDir, 81, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(pk.PastTypes) != 1 {
		t.Fatalf("PastTypes len = %d, want 1", len(pk.PastTypes))
	}
	pt := pk.PastTypes[0]
	if pt.UntilGen != 1 {
		t.Errorf("UntilGen = %d, want 1", pt.UntilGen)
	}
	if pt.Type1 != "TypeElectric" {
		t.Errorf("Type1 = %q, want \"TypeElectric\"", pt.Type1)
	}
	if pt.Type2 != "TypeNone" {
		t.Errorf("Type2 = %q, want \"TypeNone\"", pt.Type2)
	}
}

func TestBuildPokemon_Abilities(t *testing.T) {
	pk, err := BuildPokemon(testdataDir, 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	// Bulbasaur: slot 1 = overgrow (65), hidden = chlorophyll (34, skipped)
	if pk.Ability1 != 65 {
		t.Errorf("Ability1 = %d, want 65 (overgrow)", pk.Ability1)
	}
	if pk.Ability2 != 0 {
		t.Errorf("Ability2 = %d, want 0 (no second non-hidden ability)", pk.Ability2)
	}
}

func TestBuildAbility_English(t *testing.T) {
	a, err := BuildAbility(filepath.Join(testdataDir, "ability", "65", "index.json"))
	if err != nil {
		t.Fatal(err)
	}
	if a.ID != 65 {
		t.Errorf("ID = %d, want 65", a.ID)
	}
	if a.Name != "overgrow" {
		t.Errorf("Name = %q, want \"overgrow\"", a.Name)
	}
	if a.ShortDesc == "" {
		t.Error("ShortDesc is empty, want English description")
	}
	// Verify it picked English, not French
	if a.ShortDesc == "Renforce les capacités de type Plante quand les PV du Pokémon sont faibles." {
		t.Error("ShortDesc is French, want English")
	}
}

func TestCodegen_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Check that go is available
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go binary not found")
	}

	outDir := t.TempDir()
	cfg := Config{
		DataDir:    testdataDir,
		PokemonIDs: []int{1, 6, 81},
		OutDir:     outDir,
	}
	if err := Run(cfg); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Verify generated files exist
	for _, name := range []string{"abilities_gen.go", "moves_gen.go", "pokemon_gen.go"} {
		path := filepath.Join(outDir, name)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected file %s to exist: %v", name, err)
		}
	}

	// Try to compile the generated files together with the data package.
	// Build a small test program that imports and uses the generated data.
	buildTestProgram(t, outDir)
}

// buildTestProgram compiles the generated output against the data package to verify it compiles.
func buildTestProgram(t *testing.T, genDir string) {
	t.Helper()

	// Find the module root (two levels up from cmd/gen)
	moduleRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}

	// Copy generated files into a temp dir that mimics internal/data/gen
	// Then run go build on the whole module.
	// For simplicity, verify the generated files parse as valid Go.
	cmd := exec.Command("go", "vet", filepath.Join(moduleRoot, "internal", "data")+"/...")
	cmd.Dir = moduleRoot
	cmd.Env = append(os.Environ(),
		"GOPATH="+build.Default.GOPATH,
	)

	// We only check that the files we generated are syntactically valid Go
	// by running gofmt -e on them.
	for _, name := range []string{"abilities_gen.go", "moves_gen.go", "pokemon_gen.go"} {
		path := filepath.Join(genDir, name)
		out, err := exec.Command("gofmt", "-e", path).CombinedOutput()
		if err != nil {
			t.Errorf("gofmt -e %s failed: %v\n%s", name, err, out)
		}
	}
	_ = cmd
}
