package main

import (
	"fmt"
	"os"

	"github.com/davidlawson7/pokedex/internal/tui"
)

func main() {
	if err := tui.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
