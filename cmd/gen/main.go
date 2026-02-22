package gen

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func Main() {
	dataDir := flag.String("data", "", "path to PokeAPI data dir (api/v2/)")
	outDir := flag.String("out", "", "output directory for generated Go files")
	flag.Parse()

	if *dataDir == "" || *outDir == "" {
		fmt.Fprintf(os.Stderr, "Usage: gen -data <path> -out <path>\n")
		os.Exit(1)
	}

	absData, err := filepath.Abs(*dataDir)
	if err != nil {
		log.Fatalf("resolving data path: %v", err)
	}
	absOut, err := filepath.Abs(*outDir)
	if err != nil {
		log.Fatalf("resolving out path: %v", err)
	}

	cfg := Config{
		DataDir: absData,
		OutDir:  absOut,
	}

	if err := Run(cfg); err != nil {
		log.Fatalf("codegen failed: %v", err)
	}

	fmt.Printf("Generated files written to %s\n", absOut)
}
