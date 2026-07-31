package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"alt-tab/internal/tui"
)

func main() {
	program := tea.NewProgram(tui.New())

	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "alt-tab-tui: %v\n", err)
		os.Exit(1)
	}
}
