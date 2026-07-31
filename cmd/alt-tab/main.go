package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/jake-mckenzie/alt-tab/internal/nativechords"
	"github.com/jake-mckenzie/alt-tab/internal/tui"
)

func main() {
	catalog := nativechords.NewNativeCatalog()
	program := tea.NewProgram(tui.New(catalog))

	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "alt-tab-tui: %v\n", err)
		os.Exit(1)
	}
}
