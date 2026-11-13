package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/jake-mckenzie/alt-tab/internal/audio"
	"github.com/jake-mckenzie/alt-tab/internal/chords"
	"github.com/jake-mckenzie/alt-tab/internal/tui"
)

// main builds the chord catalog and runs the terminal application.
func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "alt-tab: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	catalog := chords.NewCatalog()
	model := tui.New(catalog)
	player, audioErr := audio.Open()
	if audioErr == nil {
		defer player.Close()
		model = model.WithAudio(player)
	} else {
		fmt.Fprintf(os.Stderr, "alt-tab: audio unavailable: %v\n", audioErr)
	}
	program := tea.NewProgram(model)

	_, err := program.Run()
	return err
}
