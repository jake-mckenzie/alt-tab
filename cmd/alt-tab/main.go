package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/jake-mckenzie/alt-tab/internal/chords"
	"github.com/jake-mckenzie/alt-tab/internal/tui"
)

// main builds the chord catalog and runs the terminal application.
func main() {
	catalog := chords.NewCatalog()
	program := tea.NewProgram(tui.New(catalog))

	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "alt-tab: %v\n", err)
		os.Exit(1)
	}
}
